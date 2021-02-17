package xmlmodels

import (
	"encoding/xml"
	"fmt"
	S "strings"
)

// XmlStructurePeek is called by FU.AnalyseFile(..)
// when preparing an FU.AnalysisRecord .
type XmlStructurePeek struct {
	Preamble    string
	Doctype     string
	HasDTDstuff bool
	ContentityStructure
	error
}

// PeekAtStructure_xml takes a string and does the bare minimum to find XML
// preamble, DOCTYPE, root element, whether DTD stuff was encountered, and
// the locations of outer elements containing metadata and body text.
//
// It uses the Go stdlib parser, so success in finding a root element in
// this function all but guarantees that the string is valid XML.
//
// It is called by FU.AnalyzeFile
func PeekAtStructure_xml(content string) *XmlStructurePeek {
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
		println("Peek: Error:", e.Error())
		pXSP.SetError(fmt.Errorf("peek: parser error: %w", e))
		return pXSP
	}
	for _, LAT = range latokens {
		T = LAT.Token
		switch T.(type) {
		case xml.ProcInst:
			// Found the XML preamble ?
			// type xml.ProcInst struct { Target string ; Inst []byte }
			var tok xml.ProcInst
			tok = xml.CopyToken(T).(xml.ProcInst)
			if "xml" == S.TrimSpace(tok.Target) {
				s = S.TrimSpace(string(tok.Inst))
				// println("XML-PR:", tok.Target, tok.Inst)
				if (pXSP.Preamble == "") && !didFirstPass {
					pXSP.Preamble = "<?xml " + s + "?>"
				} else {
					println("xm.peek: Got xml PI as non-first / repeated token ?!:", s)
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
			case pXSP.Root.Name:
				pXSP.Root.End = LAT.FilePosition
				pXSP.Root.End.Pos += len(localName) + 3
				println("--> End root elm at", LAT.FilePosition.String())
			case pXSP.Meta.Name:
				pXSP.Meta.End = LAT.FilePosition
				pXSP.Meta.End.Pos += len(localName) + 3
				println("--> End meta elm at", LAT.FilePosition.String())
			case pXSP.Text.Name:
				pXSP.Text.End = LAT.FilePosition
				pXSP.Text.End.Pos += len(localName) + 3
				println("--> End text elm at", LAT.FilePosition.String())
			}

		case xml.StartElement:
			// type xml.StartElement struct { Name Name ; Attr []Attr }
			var tok xml.StartElement
			tok = xml.CopyToken(T).(xml.StartElement)
			localName := tok.Name.Local

			if !foundRootElm {
				pXSP.Root.Name = localName
				pXSP.Root.Atts = tok.Attr
				pXSP.Root.Beg = LAT.FilePosition
				foundRootElm = true

				var pKeyElmTriplet *KeyElmTriplet
				pKeyElmTriplet = GetKeyElmTriplet(localName)
				if pKeyElmTriplet == nil {
					println("==> Can't find info for key elm:", localName)
				} else {
					metaTagToFind = pKeyElmTriplet.Meta
					textTagToFind = pKeyElmTriplet.Text
					fmt.Printf("--> Got key elm.beg <%s> at %s, find meta<%s> text<%s> \n",
						localName, pXSP.Root.Beg.String(),
						metaTagToFind, textTagToFind)
				}
			} else {
				if localName == metaTagToFind {
					pXSP.Meta.Name = localName
					pXSP.Meta.Atts = tok.Attr
					pXSP.Meta.Beg = LAT.FilePosition
					fmt.Printf("--> Got meta elm <%s> at %s \n",
						metaTagToFind, pXSP.Meta.Beg.String())
				}
				if localName == textTagToFind {
					pXSP.Text.Name = localName
					pXSP.Text.Atts = tok.Attr
					pXSP.Text.Beg = LAT.FilePosition
					fmt.Printf("--> Got text elm <%s> at %s \n",
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
				if pXSP.Doctype != "" {
					println("xm.peek: Second DOCTYPE ?!")
				} else {
					pXSP.Doctype = "<!" + s + ">"
				}
			}
			didFirstPass = true
		default:
			didFirstPass = true
		}
	}
	return pXSP
}

// === Implement interface Errable

func (p *XmlStructurePeek) HasError() bool {
	return p.error != nil && p.error.Error() != ""
}

// GetError is necessary cos "Error()"" dusnt tell you whether "error"
// is "nil", which is the indication of no error. Therefore we need
// this function, which can actually return the telltale "nil".
func (p *XmlStructurePeek) GetError() error {
	return p.error
}

// Error satisfies interface "error", but the
// weird thing is that "error" can be nil.
func (p *XmlStructurePeek) Error() string {
	if p.error != nil {
		return p.error.Error()
	}
	return ""
}

func (p *XmlStructurePeek) SetError(e error) {
	p.error = e
}
