package xmlmodels

import (
	"encoding/xml"
	"fmt"
	S "strings"
)

// XmlStructurePeek is created by FU.AnalyseFile(..)
// when preparing an FU.AnalysisRecord.
type XmlStructurePeek struct {
	Preamble    string
	Doctype     string
	HasDTDstuff bool
	KeyElms
	error
}

// PeekAtStructure_xml takes a string and does the bare minimum to find XML
// preamble, DOCTYPE, root element, whether DTD stuff was encountered, and
// elements that surround metadata and body text.
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
			case pXSP.RootElm.Name:
				pXSP.RootElm.EndPos = LAT.FilePosition
				pXSP.RootElm.EndPos.Pos += len(localName) + 3
				println("--> End root elm at", LAT.FilePosition.String())
			case pXSP.MetaElm.Name:
				pXSP.MetaElm.EndPos = LAT.FilePosition
				pXSP.MetaElm.EndPos.Pos += len(localName) + 3
				println("--> End meta elm at", LAT.FilePosition.String())
			case pXSP.TextElm.Name:
				pXSP.TextElm.EndPos = LAT.FilePosition
				pXSP.TextElm.EndPos.Pos += len(localName) + 3
				println("--> End text elm at", LAT.FilePosition.String())
			}

		case xml.StartElement:
			// type xml.StartElement struct { Name Name ; Attr []Attr }
			var tok xml.StartElement
			tok = xml.CopyToken(T).(xml.StartElement)
			localName := tok.Name.Local

			if !foundRootElm {
				pXSP.RootElm.Name = localName
				pXSP.RootElm.Atts = tok.Attr
				pXSP.RootElm.BegPos = LAT.FilePosition
				foundRootElm = true

				var pKeyElm *KeyElmInfo
				pKeyElm = GetKeyElm(localName)
				if pKeyElm == nil {
					println("==> Can't find info for key elm:", localName)
				} else {
					metaTagToFind = pKeyElm.Meta
					textTagToFind = pKeyElm.Text
					fmt.Printf("--> Got key elm beg <%s> at %s (%d), needs meta<%s> text<%s> \n",
						localName, pXSP.RootElm.BegPos.String(), pXSP.RootElm.BegPos.Pos,
						metaTagToFind, textTagToFind)
				}
			} else {
				if localName == metaTagToFind {
					pXSP.MetaElm.Name = localName
					pXSP.MetaElm.Atts = tok.Attr
					pXSP.MetaElm.BegPos = LAT.FilePosition
					fmt.Printf("--> Got meta elm <%s> at %s (%d) \n",
						metaTagToFind, pXSP.MetaElm.BegPos.String(), pXSP.MetaElm.BegPos.Pos)
				}
				if localName == textTagToFind {
					pXSP.TextElm.Name = localName
					pXSP.TextElm.Atts = tok.Attr
					pXSP.TextElm.BegPos = LAT.FilePosition
					fmt.Printf("--> Got text elm <%s> at %s (%d) \n",
						textTagToFind, pXSP.TextElm.BegPos.String(), pXSP.TextElm.BegPos.Pos)
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
