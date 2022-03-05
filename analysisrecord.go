package xmlutils

import (
	"strings"
	S "strings"
)

type MType string
type Doctype string
type MimeType string

// AnalysisRecord is the results of content analysis. It is named
// "Record" because it is meant to be persisted to the database.
// It is embedded in db.ContentityRecord
type AnalysisRecord struct {
	// ContypingInfo is simple fields:
	// FileExt MimeType MType Doctype IsLwDita 
	ContypingInfo
	MarkdownFlavor string
	// ContentityStructure includes Raw (the entire input content)
	ContentityStructure
	// KeyElms is: (Root,Meta,Text)ElmExtent
	// KeyElmsWithRanges
	// ContentitySections is: Text_raw, Meta_raw, MetaFormat; MetaProps SU.PropSet
	// ContentityRawSections
	// XmlInfo is: XmlPreambleFields, XmlDoctype, XmlDoctypeFields, ENTITY stuff
	/* XmlInfo */
	// XmlContype is an enum: "Unknown", "DTD", "DTDmod", "DTDent",
	// "RootTagData", "RootTagMixedContent", "MultipleRootTags", "INVALID"}
	XmlContype string
	// XmlPreambleFields is nil if no preamble - it can always
	// default to xmlutils.STD_PreambleFields (from stdlib)
	*ParsedPreamble
	// XmlDoctypeFields is a ptr - nil if ContypingInfo.Doctype
	// is "", i.e. if there is no DOCTYPE declaration
	*ParsedDoctype
	// DitaInfo
	DitaFlavor  string
	DitaContype string
}

// IsXML is true for all XML, including all HTML.
func (p AnalysisRecord) IsXML() bool {
	s := p.FileType()
	return s == "XML" || s == "HTML"
}

func (p *AnalysisRecord) String() string {
	/*
		// ContypingInfo is simple fields:
		// FileExt MimeType MType Doctype IsLwDita 
		ContypingInfo
		MarkdownFlavor string
		// ContentityStructure includes Raw (the entire input content)
		ContentityStructure
		// KeyElms is: (Root,Meta,Text)ElmExtent
		// KeyElmsWithRanges
		// ContentitySections is: Text_raw, Meta_raw, MetaFormat; MetaProps SU.PropSet
		// ContentityRawSections
		// XmlInfo is: XmlPreambleFields, XmlDoctype, XmlDoctypeFields, ENTITY stuff
		/* XmlInfo * /
		// XmlContype is an enum: "Unknown", "DTD", "DTDmod", "DTDent",
		// "RootTagData", "RootTagMixedContent", "MultipleRootTags", "INVALID"}
		XmlContype string
		// XmlPreambleFields is nil if no preamble - it can always
		// default to xmlutils.STD_PreambleFields (from stdlib)
		*XmlPreambleFields
		// XmlDoctypeFields is a ptr - nil if ContypingInfo.Doctype
		// is "", i.e. if there is no DOCTYPE declaration
		*XmlDoctypeFields
		// DitaInfo
		DitaFlavor  string
		DitaContype string
	*/
	var sb strings.Builder
	sb.WriteString("AnalysisRecord: ")
	sb.WriteString("CntpgInfo: \n\t" + p.ContypingInfo.String() + "\n\t")
	sb.WriteString("XmlCntp<" + p.XmlContype + "> ")
	sb.WriteString("XmlDctp<" + p.ParsedDoctype.String() + "> ")
	return sb.String()
}

// FileType returns "XML", "MKDN", "HTML", or future stuff TBD.
func (p AnalysisRecord) FileType() string {
	// HTML is an exceptional case
	if S.HasPrefix(p.MType, "xml/html/") {
		return "HTML"
	}
	if S.HasPrefix(p.MimeType, "text/html") {
		return "HTML"
	}
	// Normal case
	// return S.ToUpper(MTypeSub(p.MType, 0))
	// Cut & Paste
	if p.MType == "" {
		return "OOPS_NO_MType"
	}
	i := S.Index(p.MType, "/")
	if i == -1 {
		return "OOPS_NO_/_IN_MType:" + p.MType
	}
	return S.ToUpper(p.MType[:i])
}
