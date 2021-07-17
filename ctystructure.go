package xmlutils

import (
	"encoding/xml"
	"fmt"

	L "github.com/fbaube/mlog"
	SU "github.com/fbaube/stringutils"
)

// ContentityStructure is embedded in XM.AnalysisRecord
type ContentityStructure struct {
	// Text_raw + Meta_raw = Raw (maybe plus surrounding tags)
	Raw string // The entire input file
	// Root is not meaningful for non-XML
	Root Span
	Text Span
	Meta Span
	// MetaFormat is? "YAML","XML"
	MetaFormat string
	// MetaProps uses dot separators if hierarchy is needed
	MetaProps SU.PropSet
}

func (p *ContentityStructure) GetSpan(sp Span) string {
	if sp.End.Pos == 0 {
		return ""
	}
	if sp.End.Pos == -1 && sp.Beg.Pos == 0 {
		return p.Raw
	}
	if sp.Beg.Pos > sp.End.Pos {
		panic(fmt.Sprintf("BEG %d END %d", sp.Beg.Pos, sp.End.Pos))
	}
	if len(p.Raw) == 0 {
		panic("Zero-len Raw")
	}
	return p.Raw[sp.Beg.Pos:sp.End.Pos]
}

/* KeyElmsWithRanges is embedded in XmlStructurePeek
type Spans struct {
	RootElm Span
	MetaElm Span
	TextElm Span
}
*/

// Span FIXME Make this a ptr to a ContentityNode
type Span struct {
	TagName string
	Atts    []xml.Attr
	// SliceBounds
	FileRange
}

type FileRange struct {
	Beg FilePosition
	End FilePosition
}

type SliceBounds struct {
	BegIdx, EndIdx int
}

type KeyElmTriplet struct {
	Name string
	Meta string
	Text string
}

var KeyElmTriplets = []*KeyElmTriplet{
	// WHATWG: "The head element of a document is the first head element that
	// is a child of the html element, if there is one, or null otherwise.
	// The body element of a document is the first of the html element's
	// children that is either a body element or a frameset element, or
	// null if there is no such element.
	{"html", "head", "body"},
	{"topic", "prolog", "body"},
	{"map", "topicmeta", ""},
	{"reference", "", ""},
	{"task", "", ""},
	{"bookmap", "", ""},
	{"glossentry", "", ""},
	{"glossgroup", "", ""},
}

// HtmlKeyContentElms is elements that often surround the actual page content.
var HtmlKeyContentElms = []string{"main", "content"}

// HtmlSectioningContentElms have internal sections and subsections.
var HtmlSectioningContentElms = []string{"article", "aside", "nav", "section"}

// HtmlSectioningRootElms have their OWN outlines, separate from the
// outlines of their ancestors, i.e. self-contained hierarchies.
var HtmlSectioningRootElms = []string{
	"blockquote", "body", "details", "dialog", "fieldset", "figure", "td"}

func GetKeyElmTriplet(localName string) *KeyElmTriplet {
	for _, ke := range KeyElmTriplets {
		if ke.Name == localName {
			return ke
		}
	}
	return nil
}

func (sp Span) String() string {
	return fmt.Sprintf("%s[%d:%d]", sp.TagName, sp.Beg.Pos, sp.End.Pos)
}

func (p *ContentityStructure) HasNone() bool {
	return p.Root.TagName == "" && p.Meta.TagName == "" && p.Text.TagName == ""
}

func (p *ContentityStructure) SetToAllText() {
	p.Root.Beg.Pos = 0
	p.Root.End.Pos = len(p.Raw)
	p.Meta.Beg.Pos = 0
	p.Meta.End.Pos = 0
	p.Text.Beg.Pos = 0
	p.Text.End.Pos = len(p.Raw)
}

// CheckXmlSections returns true is a root element was found,
// and writes messages about other findings.
func (p *ContentityStructure) CheckXmlSections() bool {
	if p.Root.TagName == "" {
		L.L.Info("No XML root element found")
		return false
	}
	if p.Meta.TagName == "" {
		L.L.Info("No top-level metadata header element found")
	}
	if p.Text.TagName == "" {
		L.L.Info("No top-level content body text element found")
	}
	if p.Root.Beg.Pos != 0 && p.Root.End.Pos == 0 {
		L.L.Warning("Key elm root has no closing tag")
	}
	if p.Meta.Beg.Pos != 0 && p.Meta.End.Pos == 0 {
		L.L.Warning("Key elm for metadata header has no closing tag")
	}
	if p.Text.Beg.Pos != 0 && p.Text.End.Pos == 0 {
		L.L.Warning("Key elm for body text has no closing tag")
	}
	return true
}

/*
func (p *KeyElmsWithRanges) IsSplittable() bool {
	/* fmt.Printf("--> IsSplittable: %d,%d,%d,%d,%d,%d \n",
	p.RootElm.BegPos.Pos, p.RootElm.EndPos.Pos, p.MetaElm.BegPos.Pos,
	p.MetaElm.EndPos.Pos, p.TextElm.BegPos.Pos, p.TextElm.EndPos.Pos) * /
	return p.RootElm.BegPos.Pos != 0 &&
		p.RootElm.EndPos.Pos != 0 &&
		p.MetaElm.BegPos.Pos != 0 &&
		p.MetaElm.EndPos.Pos != 0 &&
		p.TextElm.BegPos.Pos != 0 &&
		p.TextElm.EndPos.Pos != 0
}
*/
/*
func (p *ContentityStructure) MetaRaw() string {
	return p.Raw[p.Meta.FileRange.Beg.Pos:p.Meta.FileRange.End.Pos]
}

func (p *ContentityStructure) TextRaw() string {
	return p.Raw[p.Text.FileRange.Beg.Pos:p.Text.FileRange.End.Pos]
}
*/

func (p *AnalysisRecord) MakeXmlContentitySections(sCont string) bool {
	// If nothing found, assume it is entirely Text.
	if p.ContentityStructure.HasNone() {
		println("--> No meta/text division detected")
		// p.Text_raw = p.Raw
		p.ContentityStructure.Text.FileRange.Beg.Pos = 0
		p.ContentityStructure.Text.FileRange.End.Pos = len(p.Raw)
		return false
	}
	if p.Raw == "" {
		L.L.Dbg("MakeXmlContentitySections: no Raw")
		p.Raw = sCont
	}
	// BEFOR?: root<html(146:306)> meta<(0:0)> text<body(155:298)>
	// AFTER?: (mmm:root (mmm:meta/:nnn) (146:html/:306) /root:nnn)
	//    OR?: (mmm:root (146:html/:306) /root:nnn)
	L.L.Info("Key elm ranges: root<%s> meta<%s> text<%s>",
		p.Root.String(), p.Meta.String(), p.Text.String())
	if p.Meta.Beg.Pos != 0 {
		L.L.Dbg("xm.KE: MetaRaw: <%d:%d> |%s|",
			p.Meta.Beg.Pos, p.Meta.End.Pos, p.GetSpan(p.Meta)) // p.MetaRaw())
		// println("D=> xm.nuCS: MetaRaaw:", p.MetaRaw())
	}
	if p.Text.Beg.Pos != 0 {
		L.L.Dbg("xm.KE: TextRaw: <%d:%d> |%s|",
			p.Text.Beg.Pos, p.Text.End.Pos,
			SU.NormalizeWhitespace(p.GetSpan(p.Text)))
		// println("D=> xm.nuCS: TextRaw:", p.TextRaw())
	}
	return true
}
