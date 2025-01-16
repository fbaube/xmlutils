package xmlutils

import (
	"encoding/xml"
	CT "github.com/fbaube/ctoken"
)

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
	"topic", "concept", "reference", "task", "bookmap",
	"map", "glossentry", "glossgroup"}

// MiscFileExtensions are all the other file extensions that
// we want to process.
var MiscFileExtensions = []string{".sqlar"}

// STD_PREAMBLE is "<?xml version="1.0" encoding="UTF-8"?>" + "\n"
var STD_PREAMBLE CT.Raw = xml.Header

// STD_PreambleFields is our parse of variable "STD_PREAMBLE".
var STD_PreambleParsed ParsedPreamble

func init() {
	pf, e := ParsePreamble(STD_PREAMBLE)
	if e != nil {
		panic("xm.xmlinfo.stdpreamble: " + e.Error())
	}
	STD_PreambleParsed = *pf
}
