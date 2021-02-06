package xmlmodels

import (
	"encoding/xml"
	"fmt"

	SU "github.com/fbaube/stringutils"
)

// ContentitySections is embedded in FU.AnalysisRecord
type ContentitySections struct {
	Raw string // The entire input file
	// Text_raw + Meta_raw = Raw (maybe plus surrounding tags)
	Text_raw   string
	Meta_raw   string
	MetaFormat string
	MetaProps  SU.PropSet
}

type KeyElmInfo struct {
	Name string
	Meta string
	Text string
}

var KeyElmInfos = []*KeyElmInfo{
	{"html", "head", "body"},
	{"topic", "prolog", "body"},
	{"map", "topicmeta", ""},
	{"reference", "", ""},
	{"task", "", ""},
	{"bookmap", "", ""},
	{"glossentry", "", ""},
	{"glossgroup", "", ""},
}

func GetKeyElm(localName string) *KeyElmInfo {
	for _, ke := range KeyElmInfos {
		if ke.Name == localName {
			return ke
		}
	}
	return nil
}

type ElmExtent struct {
	Name   string
	Atts   []xml.Attr
	BegPos FilePosition
	EndPos FilePosition
}

func (p *ElmExtent) String() string {
	return fmt.Sprintf("%s,%da,%d:%d",
		p.Name, len(p.Atts), p.BegPos.Pos, p.EndPos.Pos)
}

// KeyElms is embedded in XmlStructurePeek
type KeyElms struct {
	RootElm ElmExtent
	MetaElm ElmExtent
	TextElm ElmExtent
}

func (p *KeyElms) HasNone() bool {
	return p.RootElm.Name == "" &&
		p.MetaElm.Name == "" &&
		p.TextElm.Name == ""
}

func (p *KeyElms) SetToAllText() bool {
	return p.RootElm.Name == "" &&
		p.MetaElm.Name == "" &&
		p.TextElm.Name == ""
}

func (p *KeyElms) CheckXml() bool {
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
		println("--> Key elm RootElm has no closing tag")
	}
	if p.MetaElm.BegPos.Pos != 0 && p.MetaElm.EndPos.Pos == 0 {
		println("--> Key elm MetaElm has no closing tag")
	}
	if p.TextElm.BegPos.Pos != 0 && p.TextElm.EndPos.Pos == 0 {
		println("--> Key elm TextElm has no closing tag")
	}
	return true
}

func (p *KeyElms) IsSplittable() bool {
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

func (p *AnalysisRecord) MakeContentitySections(sCont string) {
	pCS := new(ContentitySections)
	// Fields to set:
	// Raw      string
	// Text_raw string
	// Meta_raw string
	pCS.Raw = sCont
	fmt.Printf("xm.nuCS: key<%s> meta<%s> text<%s> \n",
		p.RootElm.String(), p.MetaElm.String(), p.TextElm.String())
	if p.MetaElm.BegPos.Pos != 0 {
		pCS.Meta_raw = sCont[p.MetaElm.BegPos.Pos:p.MetaElm.EndPos.Pos]
		fmt.Printf("D=> xm.KE: set Meta_raw <%d:%d> %s \n",
			p.MetaElm.BegPos.Pos, p.MetaElm.EndPos.Pos, pCS.Meta_raw)
		println("D=> xm.nuCS: Meta_raw:", pCS.Meta_raw)
	}
	if p.TextElm.BegPos.Pos != 0 {
		pCS.Text_raw = sCont[p.TextElm.BegPos.Pos:p.TextElm.EndPos.Pos]
		fmt.Printf("D=> xm.KE: set Text_raw <%d:%d> %s \n",
			p.TextElm.BegPos.Pos, p.TextElm.EndPos.Pos, pCS.Text_raw)
		println("D=> xm.nuCS: Text_raw:", pCS.Text_raw)
	}
	// If nothing found, assume it is entirely Text.
	if p.KeyElms.HasNone() {
		println("--> No meta/text division detected")
		p.Text_raw = p.Raw
	}
	p.ContentitySections = *pCS
}
