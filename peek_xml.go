package xmlmodels

import (
	"encoding/xml"
	"fmt"
	S "strings"
)

type XmlStructure struct {
	Preamble    string
	Doctype     string
	RootTag     xml.StartElement
	HasDTDstuff bool
	error
	KeyElms map[string]*FilePosition
}

// PeekAtStructure_xml takes a string and does the bare minimum to find XML
// preamble, DOCTYPE, root element, whether DTD stuff was encountered, and
// elements that surround metadata and body text.
func PeekAtStructure_xml(content string, keyElms map[string]*FilePosition) *XmlStructure {
	var e error
	var s string

	r := S.NewReader(content)
	var parser = xml.NewDecoder(r)
	parser.Strict = false
	parser.AutoClose = xml.HTMLAutoClose
	parser.Entity = xml.HTMLEntity

	var didFirstPass bool
	var foundRootElm bool
	var pXS *XmlStructure
	pXS = new(XmlStructure)
	pXS.KeyElms = keyElms

	// DoParse_xml_locationAware(s string) (xtokens []LAToken, err error) {
	var latokens []LAToken
	var LAT LAToken
	var T xml.Token

	latokens, e = DoParse_xml_locationAware(content)
	if e != nil {
		println("Peek: Error:", e.Error())
		pXS.SetError(fmt.Errorf("peek: parser error: %w", e))
		return pXS
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
				if (pXS.Preamble == "") && !didFirstPass {
					pXS.Preamble = "<?xml " + s + "?>"
				} else {
					println("xm.peek: Got xml PI as non-first / repeated token ?!:", s)
				}
			}
			didFirstPass = true
		case xml.StartElement:
			// type xml.StartElement struct { Name Name ; Attr []Attr }
			var tok xml.StartElement
			tok = xml.CopyToken(T).(xml.StartElement)
			// Found the XML root tag ?
			if !foundRootElm {
				pXS.RootTag = tok
				foundRootElm = true
			}
			localName := tok.Name.Local
			v, ok := keyElms[localName]
			if ok {
				if v == nil {
					keyElms[localName] = &LAT.FilePosition
					fmt.Printf("--> Found <%s> at %s \n", localName, LAT.FilePosition)
				} else {
					println("DUPE!:", localName)
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
				pXS.HasDTDstuff = true
				continue
			}
			if S.HasPrefix(s, "DOCTYPE ") {
				if pXS.Doctype != "" {
					println("xm.peek: Second DOCTYPE ?!")
				} else {
					pXS.Doctype = "<!" + s + ">"
				}
			}
			didFirstPass = true
		default:
			didFirstPass = true
		}
	}
	return pXS
}

// === Implement interface Errable

func (p *XmlStructure) HasError() bool {
	return p.error != nil && p.error.Error() != ""
}

// GetError is necessary cos "Error()"" dusnt tell you whether "error"
// is "nil", which is the indication of no error. Therefore we need
// this function, which can actually return the telltale "nil".
func (p *XmlStructure) GetError() error {
	return p.error
}

// Error satisfies interface "error", but the
// weird thing is that "error" can be nil.
func (p *XmlStructure) Error() string {
	if p.error != nil {
		return p.error.Error()
	}
	return ""
}

func (p *XmlStructure) SetError(e error) {
	p.error = e
}
