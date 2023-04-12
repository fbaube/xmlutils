package xmlutils

import _ "encoding/xml" // For documentation

// TDType specifies the type of a markup tag (assumed to be
// XML) or an XML directive. Values are based on the tokens
// output'd by the stdlib [xml.Decoder], with some additions
// to accommodate DIRECTIVE subtypes, IDs, and ENUM.
// .
type TDType string

const (
	TD_type_ERROR TDType = "ERR" // ERROR

	TD_type_DOCMT = "Docmt"
	TD_type_ELMNT = "Elmnt"
	TD_type_ENDLM = "endlm"
	TD_type_VOIDD = "Voidd" // A void tag is one that needs/takes no closing tag
	TD_type_CDATA = "CData"
	TD_type_PINST = "PInst"
	TD_type_COMNT = "Comnt"
	TD_type_DRCTV = "Drctv"
	// The following are actually DIRECTIVE SUBTYPES, but they
	// are put in this list so that they can be assigned freely.
	TD_type_Doctype  = "Doctype"
	TD_type_Element  = "Element"
	TD_type_Attlist  = "Attlist"
	TD_type_Entity   = "Entitty"
	TD_type_Notation = "Notat:n"
	// The following are TBD/experimental.
	TD_type_ID    = "ID"
	TD_type_IDREF = "IDREF"
	TD_type_Enum  = "ENUM"
)

func (tdt TDType) LongForm() string {
	switch tdt {
	case TD_type_ELMNT:
		return "Start-Tag"
	case TD_type_ENDLM:
		return "End'g-Tag"
	case TD_type_CDATA:
		return "Char-Data"
	case TD_type_COMNT:
		return "_Comment_"
	case TD_type_PINST:
		return "ProcInstr"
	case TD_type_DRCTV:
		return "Directive"
	case TD_type_VOIDD:
		return "Void--Tag"
	case TD_type_DOCMT:
		return "DocuStart"
	}
	return string(tdt)
}
