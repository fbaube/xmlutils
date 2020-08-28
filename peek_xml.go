package xmlmodels

import (
	"encoding/xml"
	"fmt"
	S "strings"

	SU "github.com/fbaube/stringutils"
)

// XmlStructurePeek is created by FU.AnalyseFile(..)
// when preparing an FU.AnalysisRecord.
type XmlStructurePeek struct {
	Preamble    string
	Doctype     string
	RootTag     string // xml.StartElement
	HasDTDstuff bool
	error
	// scratch variable
	// KeyElms   map[string]*FilePosition
	// These next two are set ONLY if only a single key elm from the list is found.
	KeyElmTag string
	KeyElmPos FilePosition
}

// PeekAtStructure_xml takes a string and does the bare minimum to find XML
// preamble, DOCTYPE, root element, whether DTD stuff was encountered, and
// elements that surround metadata and body text.
func PeekAtStructure_xml(content string, keyElms []string) *XmlStructurePeek {
	var e error
	var s string

	r := S.NewReader(content)
	var parser = xml.NewDecoder(r)
	parser.Strict = false
	parser.AutoClose = xml.HTMLAutoClose
	parser.Entity = xml.HTMLEntity

	var didFirstPass bool
	var foundRootElm bool
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
		case xml.StartElement:
			// type xml.StartElement struct { Name Name ; Attr []Attr }
			var tok xml.StartElement
			tok = xml.CopyToken(T).(xml.StartElement)
			// Found the XML root tag ?
			localName := tok.Name.Local
			if !foundRootElm {
				pXSP.RootTag = localName
				foundRootElm = true
			}
			_, bb := SU.IsInSlice(localName, keyElms)
			if bb {
				if pXSP.KeyElmTag == "" {
					pXSP.KeyElmTag = localName
					pXSP.KeyElmPos = LAT.FilePosition
					fmt.Printf("--> Found key elm <%s> at %s (%d) \n",
						localName, pXSP.KeyElmPos, pXSP.KeyElmPos.Pos)
				} else {
					println("Mult. key elms:", pXSP.KeyElmTag, localName)
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
