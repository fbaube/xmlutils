package xmlutils

// This file: funcs to make working with an [aml.Token] less annoying.

import (
	"encoding/xml"
	// "fmt"
	SU "github.com/fbaube/stringutils"
	S "strings"
)

// Redefine some standard library XML types
// (for simplicity and convenience) so that
// we can (a) attach methods to them and
// (b) use them for other types of markup
// too (such as Markdown!).

type XName xml.Name
type XAtt xml.Attr
type XAtts []XAtt

func (x1 XAtts) AsStdLibXml() []xml.Attr {
	var x2 []XAtt
	var x3 []xml.Attr
	x2 = x1
	// x3 = []xml.Attr(x2)
	for _, A := range x2 {
		x3 = append(x3, xml.Attr(A))
	}
	return x3
}

type XToken struct {
	// Token is the [xml.Token] from [xml.Decoder],
	// kept around "just in case". It has to be a
	// clone, gotten by calling [xml.CopyToken].
	xml.Token
	// TDType enumerates (a) the types of [xml.Token],
	// plus (b) the (sub)types of XML directives.
	// Note that [TD_type_ENDLM] ("EndElement") is
	// superfluous when token depth info is available.
	TDType
	// TagOrDirective (ex-"TagOrPrcsrDrctv", ex-"Keyword")
	// is a convenient one-word summary. It enables quick
	// checks; it is NOT any kind of complete description.
	//  - For a tag: a simple string (minus any tag namespace)
	//  - For a PI: the processor name (i.e the first string)
	//  - An XML directive (in upper case: "DOCTYPE", etc.)
	TagOrDirective string
	// DirectiveText is for directives ONLY, and
	// not for [TD_type_ELMNT] and [TD_type_ENDLM].
	DirectiveText string
	// XName is ONLY for [TD_type_ELMNT] and [TD_type_ENDLM].
	XName
	// XAtts is ONLY for [TD_type_ELMNT].
	XAtts
}

// NewXToken returns a unified token tyoe, to replace the
// unwieldy multi-typed mess of the standard library. It
// returns a ptr so that ignorable, skippable tokens (like,
// all-whitespace) can be marked as such, by returning nil.
func NewXToken(XT xml.Token) *XToken {
	xtkn := new(XToken)
	xtkn.Token = XT
	switch XT.(type) {
	case xml.StartElement:
		xtkn.TDType = TD_type_ELMNT
		// type xml.StartElement struct {
		//     Name Name ; Attr []Attr }
		var xSE xml.StartElement
		xSE = xml.CopyToken(XT).(xml.StartElement)
		xtkn.TagOrDirective = xSE.Name.Local
		xtkn.XName = XName(xSE.Name)
		xtkn.XName.FixNS()
		// println("Elm:", xtkn.XName.String())

		// Is this the place check for any of the other
		// "standard" XML namespaces that we might encounter ?
		if xtkn.XName.Space == NS_XML {
			xtkn.XName.Space = "xml:"
		}
		for _, xA := range xSE.Attr {
			if xA.Name.Space == NS_XML {
				// println("TODO check name.local:
				// newgtoken xml:" + A.Name.Local)
				xA.Name.Space = "xml:"
			}
			gA := XAtt(xA)
			xtkn.XAtts = append(xtkn.XAtts, gA)
		}

	case xml.EndElement:
		// An EndElement has a Name (XName).
		xtkn.TDType = TD_type_ENDLM
		// type xml.EndElement struct { Name Name }
		var xEE xml.EndElement
		xEE = xml.CopyToken(XT).(xml.EndElement)
		xtkn.TagOrDirective = xEE.Name.Local
		xtkn.XName = XName(xEE.Name)
		if xtkn.XName.Space == NS_XML {
			xtkn.XName.Space = "xml:"
		}
		// fmt.Printf("<!--End-Tagnt--> %s \n", outGT.Echo())

	case xml.Comment:
		// type Comment []byte
		xtkn.TDType = TD_type_COMNT
		xtkn.TagOrDirective = "//" // TD_type_COMNT
		xtkn.DirectiveText = S.TrimSpace(
			string([]byte(XT.(xml.Comment))))
		// fmt.Printf("<!-- Comment --> <!-- %s --> \n", outGT.DirectiveText)

	case xml.ProcInst:
		xtkn.TDType = TD_type_PINST
		// type xml.ProcInst struct { Target string ; Inst []byte }
		xPI := XT.(xml.ProcInst)
		xtkn.TagOrDirective = S.TrimSpace(xPI.Target)
		xtkn.DirectiveText = S.TrimSpace(string(xPI.Inst))
		// 2023.04 This works OK :-D
		if xtkn.TagOrDirective == "xml" {
			// fmt.Printf("XML!! %s \n", xtkn.DirectiveText)
		}
		// fmt.Printf("<!--ProcInstr--> <?%s %s?> \n",
		// 	outGT.Keyword, outGT.DirectiveText)

	case xml.Directive: // type Directive []byte
		xtkn.TDType = TD_type_DRCTV
		s := S.TrimSpace(string([]byte(XT.(xml.Directive))))
		xtkn.TagOrDirective, xtkn.DirectiveText = SU.SplitOffFirstWord(s)
		// TODO: Assign TagOrDirective to xtkn.TDType
		// 2023.04 This works OK :-D
		// fmt.Printf("DRCTV||%s||%s||%s||\n",
		//	xtkn.TDType, xtkn.TagOrDirective, xtkn.DirectiveText)
		// fmt.Printf("<!--Directive--> <!%s %s> \n",
		// 	outGT.Keyword, outGT.Otherwo rds)

	case xml.CharData:
		// type CharData []byte
		xtkn.TDType = TD_type_CDATA
		xtkn.TagOrDirective = "\"\""
		bb := []byte(xml.CopyToken(XT).(xml.CharData))
		s := S.TrimSpace(string(bb))
		// This might cause problems in a scanario
		// where we actually have to worry about
		// the finer points of whitespace handing.
		// But ignore it for now, to preserve sanity.
		if s == "" {
			return nil
		}
		xtkn.DirectiveText = s
		// fmt.Printf("<!--Char-Data--> %s \n", outGT.DirectiveText)

	default:
		xtkn.TDType = TD_type_ERROR
		// L.L.Error("Unrecognized xml.Token type<%T> for: %+v", XT, XT)
		// continue
	}
	return xtkn
}

func (xt XToken) IsNonElement() bool {
	switch xt.TDType {
	case TD_type_DOCMT, TD_type_ELMNT, TD_type_ENDLM, TD_type_VOIDD:
		return false
	case TD_type_CDATA, TD_type_PINST, TD_type_COMNT, TD_type_DRCTV,
		// DIRECTIVE SUBTYPES
		TD_type_Doctype, TD_type_Element, TD_type_Attlist,
		TD_type_Entity, TD_type_Notation,
		// TBD/experimental
		TD_type_ID, TD_type_IDREF, TD_type_Enum:
		return true
	case TD_type_ERROR:
		panic("XU.IsNonElement")
	}
	return true
}

/* TMP print stuff to restore

if xtkn.TDType != TD_type_ENDLM {
	fmt.Printf("[%s] %s (%s) %s%s%s %s \n",
		pCPR.AsString(i), S.Repeat("  ", prDpth),
		xtkn.TDType, quote, xtkn.Echo(), quote, sCS)
}
*/
