package xmlutils

import (
	// "encoding/xml"
	CT "github.com/fbaube/ctoken"
	L "github.com/fbaube/mlog"
	SU "github.com/fbaube/stringutils"
)

// ContentityBasics has Raw,Root,Text,Meta,MetaProps
// and is embedded in XU.AnalysisRecord.
// .
type ContentityBasics struct {
	// XmlRoot is not meaningful for non-XML
	XmlRoot CT.Span
	Text    CT.Span
	Meta    CT.Span
	// MetaFormat is? "YAML","XML"
	MetaFormat string
	// MetaProps uses dot separators if hierarchy is needed
	MetaProps SU.PropSet
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

func (p *ContentityBasics) HasNone() bool {
	return p.XmlRoot.TagName == "" && p.Meta.TagName == "" && p.Text.TagName == ""
}

// SetToNonXml just needs the length of the content.
// .
func (p *ContentityBasics) SetToNonXml(L int) {
	p.XmlRoot.Beg.Pos = 0
	p.XmlRoot.End.Pos = L
	p.Meta.Beg.Pos = 0
	p.Meta.End.Pos = 0
	p.Text.Beg.Pos = 0
	p.Text.End.Pos = L
}

// HasRootTag returns true is a root element was found,
// and a message about missing top-level constructs,
// and can write warnings.
// .
func (p *ContentityBasics) HasRootTag() (bool, string) {
	if p.XmlRoot.TagName == "" {
		return false, "No XML root element found"
	}
	var s string
	if p.Meta.TagName == "" {
		s = "No top-level metadata header element found "
	}
	if p.Text.TagName == "" {
		s += "No top-level content body text element found"
	}
	if p.XmlRoot.Beg.Pos != 0 && p.XmlRoot.End.Pos == 0 {
		L.L.Warning("Key elm root has no closing tag")
	}
	if p.Meta.Beg.Pos != 0 && p.Meta.End.Pos == 0 {
		L.L.Warning("Key elm for metadata header has no closing tag")
	}
	if p.Text.Beg.Pos != 0 && p.Text.End.Pos == 0 {
		L.L.Warning("Key elm for body text has no closing tag")
	}
	return true, s
}
