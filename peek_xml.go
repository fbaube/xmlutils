package xmlutils

// Problem:
// ! hasRootTag
// pPeek.RawDoctype is ""
// pPeek.RawPreamble is ""

import (
	"encoding/xml"
	"fmt"
	S "strings"

	L "github.com/fbaube/mlog"
)

// XmlPeek is called by FU.AnalyseFile(..)
// when preparing an FU.AnalysisRecord .
// ContentityBasics has chunks of Raw
// but not the full "Raw" string.
// .
type XmlPeek struct { // has has Raw
	RawPreamble string
	RawDoctype  string
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
	var s string

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
	var latokens []LAToken
	var LAT LAToken
	// var T xml.Token
	var T XToken

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
	for _, LAT = range latokens {
		T = LAT.XToken
		var TorD string
		var skippable = false
		if T.TDType == TD_type_ELMNT {
			TorD = "+" + T.TagOrDirective + "+"
		} else if T.TDType == TD_type_ENDLM {
			TorD = "-" + T.TagOrDirective + "-"
		} else if T.TDType == TD_type_CDATA {
			TorD = "\"" + T.TagOrDirective + "\""
			if T.DirectiveText == "" {
				skippable = true
				TorD = "(nil.str)"
				panic("OOPS, SLIPT THRU")
			}
		} else {
			TorD = "DRCTV! " + T.TagOrDirective
		}
		if !skippable {
			// fmt.Printf("XTKN: %s :: %s \n", TorD, T.DirectiveText)
			fmt.Printf("((%s::%s)) ", TorD, T.DirectiveText)
		}
		// switch T.(type) {
		switch T.TDType {

		case TD_type_PINST: // xml.ProcInst:
			// Found the XML preamble ?
			// type xml.ProcInst struct { Target string ; Inst []byte }

			/* OBS
			var tok xml.ProcInst
			tok = xml.CopyToken(T).(xml.ProcInst)
			if S.TrimSpace(tok.Target) == "xml" {
			*/
			if T.TagOrDirective == "xml" {

				s = T.DirectiveText // S.TrimSpace(string(tok.Inst))
				// println("XML-PR:", tok.Target, tok.Inst)
				if (pPeek.RawPreamble == "") && !didFirstPass {
					pPeek.RawPreamble = "<?xml " + s + "?>"
					// fmt.Printf("GOT Raw.PREAMBLE: %s \n", s)
				} else {
					// Not fatal
					L.L.Error("xm.peek: Got \"<?xml ...>\" prolog PI " +
						"as non-first / repeated token")
				}
			}
			didFirstPass = true
		case TD_type_ENDLM: // xml.EndElement:
			// type xml.EndElement struct { Name Name ; Attr []Attr }

			/* OBS
			var tok xml.EndElement
			tok = xml.CopyToken(T).(xml.EndElement)
			// Found the XML root tag ?
			localName := tok.Name.Local
			*/
			localName := T.XName.Local

			switch localName {
			case pPeek.Root.TagName:
				pPeek.Root.End = LAT.FilePosition
				pPeek.Root.End.Pos += len(localName) + 3
				L.L.Dbg("End root elm at: " + LAT.FilePosition.String())
			case pPeek.Meta.TagName:
				pPeek.Meta.End = LAT.FilePosition
				pPeek.Meta.End.Pos += len(localName) + 3
				L.L.Dbg("End meta elm at: " + LAT.FilePosition.String())
			case pPeek.Text.TagName:
				pPeek.Text.End = LAT.FilePosition
				pPeek.Text.End.Pos += len(localName) + 3
				L.L.Dbg("End text elm at: " + LAT.FilePosition.String())
			}

		case TD_type_ELMNT: // xml.StartElement:
			// type xml.StartElement struct { Name Name ; Attr []Attr }

			/* OBS
			var tok xml.StartElement
			tok = xml.CopyToken(T).(xml.StartElement)
			localName := tok.Name.Local
			*/
			localName := T.XName.Local

			if !foundRootElm {
				fmt.Printf("FOUND FIRST (ROOT) TAG: %s \n", localName)
				pPeek.Root.TagName = localName
				pPeek.Root.Atts = T.XAtts.AsStdLibXml() // tok.Attr
				pPeek.Root.Beg = LAT.FilePosition
				foundRootElm = true

				var pKeyElmTriplet *KeyElmTriplet
				pKeyElmTriplet = GetKeyElmTriplet(localName)
				if pKeyElmTriplet == nil {
					L.L.Warning("No info for key elm: " + localName)
				} else {
					metaTagToFind = pKeyElmTriplet.Meta
					textTagToFind = pKeyElmTriplet.Text
					L.L.Progress("Got key elm.beg <%s>:%s => meta<%s> text<%s>",
						localName, pPeek.Root.Beg.String(),
						metaTagToFind, textTagToFind)
				}
			} else {
				if localName == metaTagToFind {
					pPeek.Meta.TagName = localName
					pPeek.Meta.Atts = T.XAtts.AsStdLibXml() // tok.Attr
					pPeek.Meta.Beg = LAT.FilePosition
					L.L.Dbg("Got meta elm <%s> at %s",
						metaTagToFind, pPeek.Meta.Beg.String())
				}
				if localName == textTagToFind {
					pPeek.Text.TagName = localName
					pPeek.Text.Atts = T.XAtts.AsStdLibXml() //tok.Attr
					pPeek.Text.Beg = LAT.FilePosition
					L.L.Dbg("Got text elm <%s> at %s \n",
						textTagToFind, pPeek.Text.Beg.String())
				}
			}
			didFirstPass = true

		case TD_type_DRCTV: // xml.Directive:
			// Found the DOCTYPE ?
			// type Directive []byte
			/* OBS
			var tok xml.Directive
			tok = xml.CopyToken(T).(xml.Directive)
			s = S.TrimSpace(string(tok))
			if S.HasPrefix(s, "ELEMENT ") || S.HasPrefix(s, "ATTLIST ") ||
				S.HasPrefix(s, "ENTITY ") || S.HasPrefix(s, "NOTATION ") {
				pPeek.HasDTDstuff = true
				continue
			}
			*/
			s = T.TagOrDirective
			switch s {
			case "ELEMENT", "ATTLIST", "ENTITY", "NOTATION":
				pPeek.HasDTDstuff = true
			}
			if s == "DOCTYPE" {
				if pPeek.RawDoctype != "" {
					L.L.Warning("xm.peek: Got second DOCTYPE")
				} else {
					pPeek.RawDoctype = "<!" + s + ">"
					// fmt.Printf("GOT Raw.DOCTYPE: %s \n", s)
				}
			}
			didFirstPass = true
		default:
			didFirstPass = true
		}
	}
	return pPeek, nil
}
