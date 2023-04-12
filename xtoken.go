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

func NewXToken(XT xml.Token) XToken {
	xtkn := XToken{}
	switch XT.(type) {
	case xml.StartElement:
		xtkn.TDType = TD_type_ELMNT
		// type xml.StartElement struct {
		//     Name Name ; Attr []Attr }
		var xSE xml.StartElement
		xSE = xml.CopyToken(XT).(xml.StartElement)
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
		xtkn.XName = XName(xEE.Name)
		if xtkn.XName.Space == NS_XML {
			xtkn.XName.Space = "xml:"
		}
		// fmt.Printf("<!--End-Tagnt--> %s \n", outGT.Echo())

	case xml.Comment:
		// type Comment []byte
		xtkn.TDType = TD_type_COMNT
		xtkn.DirectiveText = S.TrimSpace(
			string([]byte(XT.(xml.Comment))))
		// fmt.Printf("<!-- Comment --> <!-- %s --> \n", outGT.DirectiveText)

	case xml.ProcInst:
		xtkn.TDType = TD_type_PINST
		// type xml.ProcInst struct { Target string ; Inst []byte }
		xPI := XT.(xml.ProcInst)
		xtkn.TagOrDirective = S.TrimSpace(xPI.Target)
		xtkn.DirectiveText = S.TrimSpace(string(xPI.Inst))
		// fmt.Printf("<!--ProcInstr--> <?%s %s?> \n",
		// 	outGT.Keyword, outGT.DirectiveText)

	case xml.Directive: // type Directive []byte
		xtkn.TDType = TD_type_DRCTV
		s := S.TrimSpace(string([]byte(XT.(xml.Directive))))
		xtkn.TagOrDirective, xtkn.DirectiveText = SU.SplitOffFirstWord(s)
		// fmt.Printf("<!--Directive--> <!%s %s> \n",
		// 	outGT.Keyword, outGT.Otherwo rds)

	case xml.CharData:
		// type CharData []byte
		xtkn.TDType = TD_type_CDATA
		bb := []byte(xml.CopyToken(XT).(xml.CharData))
		s := S.TrimSpace(string(bb))
		// xtkn.Keyword remains ""
		xtkn.DirectiveText = s
		// fmt.Printf("<!--Char-Data--> %s \n", outGT.DirectiveText)

	default:
		xtkn.TDType = TD_type_ERROR
		// L.L.Error("Unrecognized xml.Token type<%T> for: %+v", XT, XT)
		// continue
	}
	/* OBS
	if xtkn.TDType != TD_type_ENDLM {
		fmt.Printf("[%s] %s (%s) %s%s%s %s \n",
			pCPR.AsString(i), S.Repeat("  ", prDpth),
			xtkn.TDType, quote, xtkn.Echo(), quote, sCS)
	}
	*/
	return xtkn
}
