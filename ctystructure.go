package xmlmodels

import (
	"encoding/xml"
	"fmt"

	SU "github.com/fbaube/stringutils"
)

type KeyElmTriplet struct {
	Name string
	Meta string
	Text string
}

// ContentityStructure is embedded in FU.AnalysisRecord
type ContentityStructure struct {
	Raw string // The entire input file
	// Text_raw + Meta_raw = Raw (maybe plus surrounding tags)
	Root Span // not meaningful for non-XML
	Text Span
	Meta Span
	// MetaFormat is? YAML,XML
	MetaFormat string
	// MetaProps uses dot separators if hierarchy is needed
	MetaProps SU.PropSet
}

/* KeyElmsWithRanges is embedded in XmlStructurePeek
type Spans struct {
	RootElm Span
	MetaElm Span
	TextElm Span
}
*/

type Span struct {
	// FIXME Make this a ptr to a ContentityNode
	Name string
	Atts []xml.Attr
	FileRange
}

type FileRange struct {
	Beg FilePosition
	End FilePosition
}

/* type FilePosition struct {
	Pos int // Position, from xml.Decoder.InputOffset()
	Lnr int // Line number
	Col int // Column [number]
} */

var KeyElmTriplets = []*KeyElmTriplet{
	{"html", "head", "body"},
	{"topic", "prolog", "body"},
	{"map", "topicmeta", ""},
	{"reference", "", ""},
	{"task", "", ""},
	{"bookmap", "", ""},
	{"glossentry", "", ""},
	{"glossgroup", "", ""},
}

func GetKeyElmTriplet(localName string) *KeyElmTriplet {
	for _, ke := range KeyElmTriplets {
		if ke.Name == localName {
			return ke
		}
	}
	return nil
}

func (p *Span) String() string {
	return fmt.Sprintf("%s(%d:%d)", p.Name, p.Beg.Pos, p.End.Pos)
}

func (p *ContentityStructure) HasNone() bool {
	return p.Root.Name == "" && p.Meta.Name == "" && p.Text.Name == ""
}

func (p *ContentityStructure) SetToAllText() {
	println("!!> alltext: no-op")
	return
}

// CheckXmlSections returns true is a root element was found,
// and writes messages about other findings.
func (p *ContentityStructure) CheckXmlSections() bool {
	if p.Root.Name == "" {
		// println("--> Key elm RootElm not found")
		return false
	}
	if p.Meta.Name == "" {
		println("--> Metadata header element not found")
	}
	if p.Text.Name == "" {
		println("--> Content body text element not found")
	}
	if p.Root.Beg.Pos != 0 && p.Root.End.Pos == 0 {
		println("--> Key elm root has no closing tag")
	}
	if p.Meta.Beg.Pos != 0 && p.Meta.End.Pos == 0 {
		println("--> Key elm for metadata header has no closing tag")
	}
	if p.Text.Beg.Pos != 0 && p.Text.End.Pos == 0 {
		println("--> Key elm for body text has no closing tag")
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

func (p *ContentityStructure) MetaRaw() string {
	return p.Raw[p.Meta.FileRange.Beg.Pos:p.Meta.FileRange.End.Pos]
}

func (p *ContentityStructure) TextRaw() string {
	return p.Raw[p.Text.FileRange.Beg.Pos:p.Text.FileRange.End.Pos]
}

func (p *AnalysisRecord) MakeXmlContentitySections(sCont string) bool {
	// If nothing found, assume it is entirely Text.
	if p.ContentityStructure.HasNone() {
		println("--> No meta/text division detected")
		// p.Text_raw = p.Raw
		p.ContentityStructure.Text.FileRange.Beg.Pos = 0
		p.ContentityStructure.Text.FileRange.End.Pos = len(p.Raw)
		return false
	}
	p.Raw = sCont
	// BEFOR?: root<html(146:306)> meta<(0:0)> text<body(155:298)>
	// AFTER?: (mmm:root (mmm:meta/:nnn) (146:html/:306) /root:nnn)
	//    OR?: (mmm:root (146:html/:306) /root:nnn)
	fmt.Printf("xm.nuCS: root<%s> meta<%s> text<%s> \n",
		p.Root.String(), p.Meta.String(), p.Text.String())
	if p.Meta.Beg.Pos != 0 {
		fmt.Printf("D=> xm.KE: MetaRaw: <%d:%d> |%s| \n",
			p.Meta.Beg.Pos, p.Meta.End.Pos, p.MetaRaw())
		// println("D=> xm.nuCS: MetaRaaw:", p.MetaRaw())
	}
	if p.Text.Beg.Pos != 0 {
		fmt.Printf("D=> xm.KE: TextRaw: <%d:%d> |%s| \n",
			p.Text.Beg.Pos, p.Text.End.Pos, p.TextRaw())
		// println("D=> xm.nuCS: TextRaw:", p.TextRaw())
	}
	return true
}
