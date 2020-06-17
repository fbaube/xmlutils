package xmlmodels

import (
	"encoding/xml"
)

// STD_PREAMBLE is "<?xml version="1.0" encoding="UTF-8"?>" + "\n"
var STD_PREAMBLE string = xml.Header
var STD_PreambleFields XmlPreambleFields

func init() {
	pf, e := NewXmlPreambleFields(STD_PREAMBLE)
	if e != nil {
		panic("xm.xmlinfo.stdpreamble: " + e.Error())
	}
	STD_PreambleFields = *pf
}

type XmlInfo struct {
	XmlContype
	// Defaults to xmlmodels.STD_PreambleFields
 	XmlPreambleFields
	XmlDoctype
 *XmlDoctypeFields

 	// TagDefCt is for DTD-type files (.dtd, .mod, .ent)
 	// // TagDefCt int // Nr of <!ELEMENT ...>
 	// RootTagIndex int  // Or some sort of pointer into the tree.
 	// RootTagCt is >1 means mark the content as a Fragment.
 	// // RootTagCt int

	// (Obs.cmt) XML items are
	//  - (DOCS) IDs & IDREFs
	//  - (DTDs) Elm defs (incl. Att defs) & Ent defs.

	// It is not precisely defined how to handle relative paths in external
	// IDs and entity substitutions, so we need to maintain this list.
	// EntSearchDirs []string // TODO

	// GEnts is "ENTITY"" directives (both with "%" and without).
	// GEnts map[string]*gparse.GEnt
	// DElms is "ELEMENT" directives.
	// DElms map[string]*gtree.GTag
	// TODO Maybe also add maps for NOTs (Notations)
}

type XmlContype string
type XmlDoctype string

// XmlContypes, maybe DTDmod should be DTDelms.
var XmlContypes = []XmlContype{"Unknown", "DTD", "DTDmod", "DTDent",
	"RootTagData", "RootTagMixedContent", "MultipleRootTags", "INVALID"}

func (p *XmlInfo) String() string {
	return "XmlInfo:meh"
}
