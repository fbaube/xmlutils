package xmlutils

// Problem:
// ! hasRootTag
// pPeek.RawDoctype is ""
// pPeek.RawPreamble is ""

import (
	"encoding/xml"
	"fmt"
	S "strings"

	CT "github.com/fbaube/ctoken"
	L "github.com/fbaube/mlog"
)

// XmlPeek is called by FU.AnalyseFile(..)
// when preparing an FU.AnalysisRecord .
// ContentityBasics has chunks of Raw
// but not the full "Raw" string.
// .
type XmlPeek struct { // has Raw
	PreambleRaw CT.Raw // string
	DoctypeRaw  CT.Raw // string
	HasDTDstuff bool
	ContentityBasics
	// error
}

// Peek_xml takes a string and does the minimum to find XML preamble,
// DOCTYPE, root element, whether DTD stuff was encountered, and the
// locations of outer elements containing metadata and body text.
//
// It uses the Go stdlib parser, so success in finding a root element
// in this function all but guarantees that the string is valid XML.
//
// It is called by FU.AnalyzeFile
// .
func Peek_xml(content string) (*XmlPeek, error) {

	// The return value !
	var pPeek *XmlPeek
	pPeek = new(XmlPeek)

	var e error

	r := S.NewReader(content)
	var parser = xml.NewDecoder(r)
	parser.Strict = false
	parser.AutoClose = xml.HTMLAutoClose
	parser.Entity = xml.HTMLEntity

	var didFirstPass bool
	var foundRootElm bool
	var metaTagToFind string
	var textTagToFind string
	// Obsolete ?
	// pPeek.Raw = content

	// DoParse_xml_locationAware(s string) (xtokens []LAToken, err error) {
	var latokens []CT.LAToken
	var LAT CT.LAToken
	// var T xml.Token
	var T CT.CToken

	latokens, e = DoParse_xml_locationAware(content)
	if e != nil {
		L.L.Error("xm.peek: " + e.Error())
		return pPeek, fmt.Errorf("xm.peek: parser error: %w", e)
	}
	/* REF
	type XToken struct {
		xml.Token
		TDType
		// TagOrDirective is a convenient one-word summary.
		TagOrDirective string
		XName
		XAtts
		// DirectiveText is for directives ONLY, and
		// not for [TD_type_ELMNT] and [TD_type_ENDLM].
		DirectiveText string
	*/
	// This loop just creates a printable string.
	for _, LAT = range latokens {
		T = LAT.CToken
		var TorD string
		var skippable = false
		switch T.TDType {
		case CT.TD_type_ELMNT:
			TorD = "+" + T.Text + "+"
			skippable = true
		case CT.TD_type_ENDLM:
			TorD = "-" + T.Text + "-"
			skippable = true
		case CT.TD_type_CDATA:
			TorD = "\"" + T.Text + "\""
			if T.Text == "" {
				skippable = true
				TorD = "(nil.str)"
				panic("OOPS, SLIPT THRU")
			}
			skippable = true
		case CT.TD_type_DRCTV:
			TorD = "drctv:" + T.ControlStrings[0] + "," + T.ControlStrings[1] + "," + T.Text
		case CT.TD_type_PINST:
			TorD = "pr.i.: " + T.Text + "," + T.ControlStrings[0]
		case CT.TD_type_COMNT:
			TorD = "comnt:" + T.Text
		default:
			L.L.Panic("OOPS, bad TDT in LAToken: %s", T.TDType)
		}
		if !skippable {
			// fmt.Printf("XTKN: %s :: %s \n", TorD, T.DirectiveText)
			L.L.Dbg("peek: %s", TorD)
		}

		// -------------------------------------------
		// This loop is where real processing is done.
		// -------------------------------------------
		switch T.TDType {
		// TT := T.SourceToken
		// switch TT.(type) {

		// case xml.ProcInst:
		case CT.TD_type_PINST:
			// Found the XML preamble ?
			// xml.ProcInst struct { Target string ; Inst []byte }
			// if TT.Target == "xml" {
			if T.Text == "xml" {
				sInst := T.ControlStrings[0]
				L.L.Dbg("XML-PrI <%s> <%s>",
					T.Text, T.ControlStrings[0])
				if (pPeek.PreambleRaw == "") && !didFirstPass {
					pPeek.PreambleRaw = CT.Raw("<?xml " + sInst + "?>")
					L.L.Dbg("GOT Raw.PREAMBLE: %s", pPeek.PreambleRaw)
				} else {
					// Not fatal
					L.L.Error("xm.peek: Got \"<?xml ...>\" prolog PI " +
						"as non-first / repeated token")
				}
			}
			didFirstPass = true

		case CT.TD_type_ENDLM: // xml.EndElement:
			// type xml.EndElement struct { Name Name ; Attr []Attr }
			localName := T.CName.Local
			switch localName {
			case pPeek.XmlRoot.TagName:
				pPeek.XmlRoot.End = LAT.FilePosition
				pPeek.XmlRoot.End.Pos += len(localName) + 3
				L.L.Dbg("End xml root at: " + LAT.FilePosition.Info())
			case pPeek.Meta.TagName:
				pPeek.Meta.End = LAT.FilePosition
				pPeek.Meta.End.Pos += len(localName) + 3
				L.L.Dbg("End meta elm at: " + LAT.FilePosition.Info())
			case pPeek.Text.TagName:
				pPeek.Text.End = LAT.FilePosition
				pPeek.Text.End.Pos += len(localName) + 3
				L.L.Dbg("End text elm at: " + LAT.FilePosition.Info())
			}

		case CT.TD_type_ELMNT:
			// type xml.StartElement struct { Name Name ; Attr []Attr }
			localName := T.CName.Local
			if !foundRootElm {
				L.L.Progress("Found root tag: " + localName)
				pPeek.XmlRoot.TagName = localName
				pPeek.XmlRoot.Atts = T.CAtts.AsStdLibXml() // tok.Attr
				pPeek.XmlRoot.Beg = LAT.FilePosition
				foundRootElm = true

				var pKeyElmTriplet *KeyElmTriplet
				pKeyElmTriplet = GetKeyElmTriplet(localName)
				if pKeyElmTriplet == nil {
					L.L.Warning("No info for key " +
						"(root?) elm: " + localName)
				} else {
					metaTagToFind = pKeyElmTriplet.Meta
					textTagToFind = pKeyElmTriplet.Text
					L.L.Progress("Got key elm.beg <%s>:%s => meta<%s> text<%s>",
						localName, pPeek.XmlRoot.Beg.Info(),
						metaTagToFind, textTagToFind)
				}
			} else {
				if localName == metaTagToFind {
					pPeek.Meta.TagName = localName
					pPeek.Meta.Atts = T.CAtts.AsStdLibXml() // tok.Attr
					pPeek.Meta.Beg = LAT.FilePosition
					L.L.Dbg("Got meta elm <%s> at %s",
						metaTagToFind, pPeek.Meta.Beg.Info())
				}
				if localName == textTagToFind {
					pPeek.Text.TagName = localName
					pPeek.Text.Atts = T.CAtts.AsStdLibXml() //tok.Attr
					pPeek.Text.Beg = LAT.FilePosition
					L.L.Dbg("Got text elm <%s> at %s \n",
						textTagToFind, pPeek.Text.Beg.Info())
				}
			}
			didFirstPass = true

		case CT.TD_type_DRCTV: // xml.Directive:
			// Found the DOCTYPE ?
			// type Directive []byte
			sDT := T.Text // "DOCTYPE"
			// sRE = Root Element
			sRE := T.ControlStrings[0] // "HTML" (etc.) IF "DOCTYPE", else WHOLE REM. TEXT
			sXX := T.ControlStrings[1] // "etc etc etc"
			switch sDT {
			case "ELEMENT", "ATTLIST", "ENTITY", "NOTATION":
				pPeek.HasDTDstuff = true
				// TODO: Do something with sRE
			case "DOCTYPE":
				if pPeek.DoctypeRaw != "" {
					L.L.Warning("xm.peek: Got second DOCTYPE")
				} else {
					pPeek.DoctypeRaw = CT.Raw("<!" +
						sDT + " " + sRE + " " + sXX + ">")
					fmt.Printf("peek: Raw.DOCTYPE: %s \n", pPeek.DoctypeRaw)
					L.L.Dbg("peek: Raw.DOCTYPE: %s", pPeek.DoctypeRaw)
				}
			}
			didFirstPass = true
		default:
			didFirstPass = true
		}
	}
	return pPeek, nil
}
