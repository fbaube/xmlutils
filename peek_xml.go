package xmlmodels

import (
	"encoding/xml"
	"fmt"
	"io"
	S "strings"
)

// Peek_xml takes a string.
func Peek_xml(content string) (preamble string, doctype string, rootTag xml.StartElement, gotDTDstuff bool, err error) {
	var e error
	var s string
	var T xml.Token

	r := S.NewReader(content)
	var parser = xml.NewDecoder(r)
	parser.Strict    = false
	parser.AutoClose = xml.HTMLAutoClose
	parser.Entity    = xml.HTMLEntity

	var didFirstPass bool
	var foundRootElm bool

	for {
		T, e = parser.Token()
		if e == io.EOF { break }
		if e != nil {
			println("Peek: Error:", e.Error())
			return preamble, doctype, rootTag, gotDTDstuff, fmt.Errorf("xm.peek: ERROR: %w", e)
		}
		switch T.(type) {
		case xml.ProcInst:
			// Found the XML preamble ?
			// type xml.ProcInst struct { Target string ; Inst []byte }
			var tok xml.ProcInst
			tok = xml.CopyToken(T).(xml.ProcInst)
			if "xml" == S.TrimSpace(tok.Target) {
				s = S.TrimSpace(string(tok.Inst))
				// println("XML-PR:", tok.Target, tok.Inst)
				if (preamble == "") && !didFirstPass {
					preamble = "<?xml " + s + "?>"
				} else {
					println("xm.peek: Got xml PI as non-first / repeated token ?!:", s)
				}
			}
			didFirstPass = true
		case xml.StartElement:
			if foundRootElm { continue }
			// Found the XML root tag ?
			// type xml.StartElement struct { Name Name ; Attr []Attr }
			var tok xml.StartElement
			tok = xml.CopyToken(T).(xml.StartElement)
			rootTag = tok
			foundRootElm = true
		case xml.Directive:
			// Found the DOCTYPE ?
			// type Directive []byte
			var tok xml.Directive
			tok = xml.CopyToken(T).(xml.Directive)
			s = S.TrimSpace(string(tok))
			if S.HasPrefix(s, "ELEMENT ") || S.HasPrefix(s, "ATTLIST ") ||
			   S.HasPrefix(s, "ENTITY ")  || S.HasPrefix(s, "NOTATION ") {
				gotDTDstuff = true
				continue
			}
			if S.HasPrefix(s, "DOCTYPE ") {
				if doctype != "" {
					println ("xm.peek: Second DOCTYPE ?!")
				} else {
					doctype = "<!" + s + ">"
				}
			}
			didFirstPass = true
		default:
			didFirstPass = true
		}
	}
	return preamble, doctype, rootTag, gotDTDstuff, nil
}
