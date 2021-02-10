package xmlmodels

import (
	S "strings"
)

// AnalysisRecord is the results of content analysis. It is named
// "Record" because it is meant to be persisted to the database.
// It is embedded in db.ContentRecord
type AnalysisRecord struct {
	// ContypingInfo is simple fields:
	// FileExt MimeType MType Doctype IsLwDita IsProcbl
	ContypingInfo
	MarkdownFlavor string
	// KeyElms is: (Root,Meta,Text)ElmExtent
	KeyElms
	// ContentitySections is: Text_raw, Meta_raw, MetaFormat; MetaProps SU.PropSet
	ContentitySections
	// XmlInfo is: XmlPreambleFields, XmlDoctype, XmlDoctypeFields, ENTITY stuff
	/* XmlInfo */
	// XmlContype is: "Unknown", "DTD", "DTDmod", "DTDent", "RootTagData",
	// "RootTagMixedContent", "MultipleRootTags", "INVALID"}
	XmlContype
	// XmlPreambleFields is nil if no preamble - it can always
	// default to xmlmodels.STD_PreambleFields (from stdlib)
	*XmlPreambleFields
	// XmlDoctypeFields is a ptr - nil if ContypingInfo.Doctype
	// is "", i.e. if there is no DOCTYPE declaration
	*XmlDoctypeFields
	// DitaInfo
	DitaMarkupLg
	DitaContype
}

// IsXML is true for all XML, including all HTML.
func (p AnalysisRecord) IsXML() bool {
	s := p.FileType()
	return s == "XML" || s == "HTML"
}

func (p *AnalysisRecord) String() string {
	return ("AR!")
}

// FileType returns "XML", "MKDN", "HTML", or future stuff TBD.
func (p AnalysisRecord) FileType() string {
	// Exceptional case
	if S.HasPrefix(p.MType, "xml/html/") {
		return "HTML"
	}
	if S.HasPrefix(p.MimeType, "text/html") {
		return "HTML"
	}
	// Normal case
	// return S.ToUpper(MTypeSub(p.MType, 0))
	// Cut & Paste
	i := S.Index(p.MType, "/")
	return S.ToUpper(p.MType[:i])
}
