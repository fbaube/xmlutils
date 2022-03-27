package xmlutils

import (
	"encoding/xml"
	"fmt"
	S "strings"

	L "github.com/fbaube/mlog"
)

// XmlStructurePeek is called by FU.AnalyseFile(..)
// when preparing an FU.AnalysisRecord .
type XmlStructurePeek struct {
	RawPreamble string
	RawDoctype  string
	HasDTDstuff bool
	ContentityStructure
	// error
}

// PeekAtStructure_xml takes a string and does the bare minimum to find XML
// preamble, DOCTYPE, root element, whether DTD stuff was encountered, and
// the locations of outer elements containing metadata and body text.
//
// It uses the Go stdlib parser, so success in finding a root element in
// this function all but guarantees that the string is valid XML.
//
// It is called by FU.AnalyzeFile
func PeekAtStructure_xml(content string) (*XmlStructurePeek, error) {
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
	var pXSP *XmlStructurePeek
	pXSP = new(XmlStructurePeek)

	// DoParse_xml_locationAware(s string) (xtokens []LAToken, err error) {
	var latokens []LAToken
	var LAT LAToken
	var T xml.Token

	latokens, e = DoParse_xml_locationAware(content)
	if e != nil {
		L.L.Error("xm.peek: " + e.Error())
		return pXSP, fmt.Errorf("xm.peek: parser error: %w", e)
	}
	for _, LAT = range latokens {
		T = LAT.Token
		switch T.(type) {
		case xml.ProcInst:
			// Found the XML preamble ?
			// type xml.ProcInst struct { Target string ; Inst []byte }
			var tok xml.ProcInst
			tok = xml.CopyToken(T).(xml.ProcInst)
			if S.TrimSpace(tok.Target) == "xml" {
				s = S.TrimSpace(string(tok.Inst))
				// println("XML-PR:", tok.Target, tok.Inst)
				if (pXSP.RawPreamble == "") && !didFirstPass {
					pXSP.RawPreamble = "<?xml " + s + "?>"
				} else {
					// Not fatal
					L.L.Error("xm.peek: Got xml PI as non-first / repeated token")
				}
			}
			didFirstPass = true
		case xml.EndElement:
			// type xml.EndElement struct { Name Name ; Attr []Attr }
			var tok xml.EndElement
			tok = xml.CopyToken(T).(xml.EndElement)
			// Found the XML root tag ?
			localName := tok.Name.Local
			switch localName {
			case pXSP.Root.TagName:
				pXSP.Root.End = LAT.FilePosition
				pXSP.Root.End.Pos += len(localName) + 3
				L.L.Dbg("End root elm at: " + LAT.FilePosition.String())
			case pXSP.Meta.TagName:
				pXSP.Meta.End = LAT.FilePosition
				pXSP.Meta.End.Pos += len(localName) + 3
				L.L.Dbg("End meta elm at: " + LAT.FilePosition.String())
			case pXSP.Text.TagName:
				pXSP.Text.End = LAT.FilePosition
				pXSP.Text.End.Pos += len(localName) + 3
				L.L.Dbg("End text elm at: " + LAT.FilePosition.String())
			}

		case xml.StartElement:
			// type xml.StartElement struct { Name Name ; Attr []Attr }
			var tok xml.StartElement
			tok = xml.CopyToken(T).(xml.StartElement)
			localName := tok.Name.Local

			if !foundRootElm {
				pXSP.Root.TagName = localName
				pXSP.Root.Atts = tok.Attr
				pXSP.Root.Beg = LAT.FilePosition
				foundRootElm = true

				var pKeyElmTriplet *KeyElmTriplet
				pKeyElmTriplet = GetKeyElmTriplet(localName)
				if pKeyElmTriplet == nil {
					L.L.Warning("No info for key elm: " + localName)
				} else {
					metaTagToFind = pKeyElmTriplet.Meta
					textTagToFind = pKeyElmTriplet.Text
					L.L.Progress("Got key elm.beg <%s>:%s => meta<%s> text<%s>",
						localName, pXSP.Root.Beg.String(),
						metaTagToFind, textTagToFind)
				}
			} else {
				if localName == metaTagToFind {
					pXSP.Meta.TagName = localName
					pXSP.Meta.Atts = tok.Attr
					pXSP.Meta.Beg = LAT.FilePosition
					L.L.Dbg("Got meta elm <%s> at %s",
						metaTagToFind, pXSP.Meta.Beg.String())
				}
				if localName == textTagToFind {
					pXSP.Text.TagName = localName
					pXSP.Text.Atts = tok.Attr
					pXSP.Text.Beg = LAT.FilePosition
					L.L.Dbg("Got text elm <%s> at %s \n",
						textTagToFind, pXSP.Text.Beg.String())
				}
			}
			didFirstPass = true
		case xml.Directive:
			// Found the DOCTYPE ?
			// type Directive []byte
			var tok xml.Directive
			tok = xml.CopyToken(T).(xml.Directive)
			s = S.TrimSpace(string(tok))
			if S.HasPrefix(s, "ELEMENT ") || S.HasPrefix(s, "ATTLIST ") ||
				S.HasPrefix(s, "ENTITY ") || S.HasPrefix(s, "NOTATION ") {
				pXSP.HasDTDstuff = true
				continue
			}
			if S.HasPrefix(s, "DOCTYPE ") {
				if pXSP.RawDoctype != "" {
					L.L.Warning("xm.peek: Got second DOCTYPE")
				} else {
					pXSP.RawDoctype = "<!" + s + ">"
				}
			}
			didFirstPass = true
		default:
			didFirstPass = true
		}
	}
	return pXSP, nil
}
