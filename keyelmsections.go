package xmlmodels

import (
	"encoding/xml"
	"fmt"

	SU "github.com/fbaube/stringutils"
)

// ContentityRawSections is embedded in FU.AnalysisRecord
type ContentityRawSections struct {
	Raw string // The entire input file
	// Text_raw + Meta_raw = Raw (maybe plus surrounding tags)
	Text_raw   string
	Meta_raw   string
	MetaFormat string
	MetaProps  SU.PropSet
}

type KeyElmTriplet struct {
	Name string
	Meta string
	Text string
}

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

type ElmWithRange struct {
	Name   string
	Atts   []xml.Attr
	BegPos FilePosition
	EndPos FilePosition
}

func (p *ElmWithRange) String() string {
	return fmt.Sprintf("%s,atts[%d],%d:%d",
		p.Name, len(p.Atts), p.BegPos.Pos, p.EndPos.Pos)
}

// KeyElmsWithRanges is embedded in XmlStructurePeek
type KeyElmsWithRanges struct {
	RootElm ElmWithRange
	MetaElm ElmWithRange
	TextElm ElmWithRange
}

func (p *KeyElmsWithRanges) HasNone() bool {
	return p.RootElm.Name == "" &&
		p.MetaElm.Name == "" &&
		p.TextElm.Name == ""
}

func (p *KeyElmsWithRanges) SetToAllText() bool {
	return p.RootElm.Name == "" &&
		p.MetaElm.Name == "" &&
		p.TextElm.Name == ""
}

// CheckXmlSections returns true is a root element was found,
// and writes messages about other findings.
func (p *KeyElmsWithRanges) CheckXmlSections() bool {
	if p.RootElm.Name == "" {
		// println("--> Key elm RootElm not found")
		return false
	}
	if p.MetaElm.Name == "" {
		println("--> Metadata header element not found")
	}
	if p.TextElm.Name == "" {
		println("--> Content body text element not found")
	}
	if p.RootElm.BegPos.Pos != 0 && p.RootElm.EndPos.Pos == 0 {
		println("--> Key elm root has no closing tag")
	}
	if p.MetaElm.BegPos.Pos != 0 && p.MetaElm.EndPos.Pos == 0 {
		println("--> Key elm for metadata header has no closing tag")
	}
	if p.TextElm.BegPos.Pos != 0 && p.TextElm.EndPos.Pos == 0 {
		println("--> Key elm for body text has no closing tag")
	}
	return true
}

func (p *KeyElmsWithRanges) IsSplittable() bool {
	/* fmt.Printf("--> IsSplittable: %d,%d,%d,%d,%d,%d \n",
	p.RootElm.BegPos.Pos, p.RootElm.EndPos.Pos, p.MetaElm.BegPos.Pos,
	p.MetaElm.EndPos.Pos, p.TextElm.BegPos.Pos, p.TextElm.EndPos.Pos) */
	return p.RootElm.BegPos.Pos != 0 &&
		p.RootElm.EndPos.Pos != 0 &&
		p.MetaElm.BegPos.Pos != 0 &&
		p.MetaElm.EndPos.Pos != 0 &&
		p.TextElm.BegPos.Pos != 0 &&
		p.TextElm.EndPos.Pos != 0
}

func (p *AnalysisRecord) MakeXmlContentitySections(sCont string) bool {
	// If nothing found, assume it is entirely Text.
	if p.KeyElmsWithRanges.HasNone() {
		println("--> No meta/text division detected")
		p.Text_raw = p.Raw
		return false
	}
	// Fields to set:
	// Raw      string
	// Text_raw string
	// Meta_raw string
	p.Raw = sCont
	fmt.Printf("xm.nuCS: key<%s> meta<%s> text<%s> \n",
		p.RootElm.String(), p.MetaElm.String(), p.TextElm.String())
	if p.MetaElm.BegPos.Pos != 0 {
		p.Meta_raw = sCont[p.MetaElm.BegPos.Pos:p.MetaElm.EndPos.Pos]
		fmt.Printf("D=> xm.KE: set Meta_raw <%d:%d> %s \n",
			p.MetaElm.BegPos.Pos, p.MetaElm.EndPos.Pos, p.Meta_raw)
		println("D=> xm.nuCS: Meta_raw:", p.Meta_raw)
	}
	if p.TextElm.BegPos.Pos != 0 {
		p.Text_raw = sCont[p.TextElm.BegPos.Pos:p.TextElm.EndPos.Pos]
		fmt.Printf("D=> xm.KE: set Text_raw <%d:%d> %s \n",
			p.TextElm.BegPos.Pos, p.TextElm.EndPos.Pos, p.Text_raw)
		println("D=> xm.nuCS: Text_raw:", p.Text_raw)
	}
	return true
}
