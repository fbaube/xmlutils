package xmlmodels

import "encoding/xml"

// DTDtypeFileExtensions are all the file extensions that are
// automatically classified as being DTD-type.
var DTDtypeFileExtensions = []string{".dtd", ".mod", ".ent"}

// MarkdownFileExtensions are all the file extensions that are automatically
// classified as being Markdown-type, even tho we generally use a regex instead.
var MarkdownFileExtensions = []string{".md", ".mdown", ".markdown", ".mkdn"}

// DITAtypeFileExtensions are all the file extensions that are
// automatically classified as being DITA-type.
var DITAtypeFileExtensions = []string{".dita", ".ditamap", ".ditaval"}

// DITArootElms are all the XML root elements that can be
// classified as DITA-type. Note that LwDITA uses only "topic".
var DITArootElms = []string{
	"topic", "concept", "reference", "task", "bookmap", "glossentry", "glossgroup"}

// STD_PREAMBLE is "<?xml version="1.0" encoding="UTF-8"?>" + "\n"
var STD_PREAMBLE string = xml.Header

// STD_PreambleFields is our parse of variable "STD_PREAMBLE".
var STD_PreambleFields XmlPreambleFields

func init() {
	pf, e := NewXmlPreambleFields(STD_PREAMBLE)
	if e != nil {
		panic("xm.xmlinfo.stdpreamble: " + e.Error())
	}
	STD_PreambleFields = *pf
}
