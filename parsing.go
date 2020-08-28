package xmlmodels

import (
	"encoding/xml"
	"fmt"
	"io"
	S "strings"

	SU "github.com/fbaube/stringutils"
)

type ConcreteParseResults_xml struct {
	// ParseTree ??
	NodeList []xml.Token
	CommonCPR
}

func GetConcreteParseResults_xml(s string) (*ConcreteParseResults_xml, error) {
	var nl []xml.Token
	var e error
	nl, e = DoParse_xml(s)
	if e != nil {
		return nil, fmt.Errorf("pu.xml.parseResults: %w", e)
	}
	p := new(ConcreteParseResults_xml)
	p.CommonCPR = *NewCommonCPR()
	p.NodeList = nl
	p.CPR_raw = s
	return p, nil
}

// DoParse_xml takes a string, so we can assume that we can
// discard it after use cos the caller has another copy of it.
// To be safe, it copies every token using `xml.CopyToken(T)`.
func DoParse_xml(s string) (xtokens []xml.Token, err error) {
	return doParse_xml_maybeRaw(s, false)
}

func DoParseRaw_xml(s string) (xtokens []xml.Token, err error) {
	return doParse_xml_maybeRaw(s, true)
}

func doParse_xml_maybeRaw(s string, doRaw bool) (xtokens []xml.Token, err error) {
	var e error
	var T, TT xml.Token
	xtokens = make([]xml.Token, 0, 100)
	// println("(DD) XmlTokenizeBuffer:", s)

	r := S.NewReader(s)
	var parser *xml.Decoder
	parser = NewConfiguredDecoder(r)

	for {
		if doRaw {
			T, e = parser.RawToken()
		} else {
			T, e = parser.Token()
		}
		if e == io.EOF {
			break
		}
		if e != nil {
			return xtokens, fmt.Errorf("pu.xml.doParse: %w", e)
		}
		TT = xml.CopyToken(T)
		xtokens = append(xtokens, TT)
	}
	return xtokens, nil
}

func DoParse_xml_locationAware(s string) (xtokens []LAToken, err error) {
	var e error
	var T, TT xml.Token
	var LAT LAToken
	xtokens = make([]LAToken, 0, 100)
	// println("(DD) XmlTokenizeBuffer:", s)
	var idcs []int
	idcs = SU.AllIndices(s, "\n")

	r := S.NewReader(s)
	var parser *xml.Decoder
	parser = NewConfiguredDecoder(r)
	var pos int

	for {
		pos = int(parser.InputOffset())
		T, e = parser.RawToken()
		if e == io.EOF {
			break
		}
		if e != nil {
			return xtokens, fmt.Errorf("pu.xml.doParse: %w", e)
		}
		TT = xml.CopyToken(T)
		LAT = *new(LAToken)
		LAT.Token = TT
		LAT.Pos = pos
		LAT.Lnr, LAT.Col = LnrAndColFromPos(pos, idcs)
		xtokens = append(xtokens, LAT)
	}
	return xtokens, nil
}

func LnrAndColFromPos(pos int, idcs []int) (int, int) {
	if pos < idcs[0] {
		return 1, pos + 1
	}
	for i, v := range idcs {
		if pos < v {
			return i + 1, pos - idcs[i-1] + 1
		}
	}
	return -1, -1
}

func NewConfiguredDecoder(r io.Reader) *xml.Decoder {
	var parser *xml.Decoder
	parser = xml.NewDecoder(r)
	// Strict mode does not enforce XML namespace requirements. In parti-
	// cular it does not reject namespace tags that use undefined prefixes.
	// Such tags are recorded with the unknown prefix as the namespace URL.
	// We will use this because we want details of namespace parsing.
	parser.Strict = false
	// When Strict == false, AutoClose is a set of elements to consider
	// closed immediately after they are opened, regardless of whether
	// an end element is present. For example, <br/>.
	// TODO Add anything for LwDITA ?
	parser.AutoClose = xml.HTMLAutoClose
	// Entity can map non-standard entity names to string replacements.
	// The parser is preloaded with the following standard XML mappings,
	// whether or not they are also provided in the actual map content:
	//	"lt": "<", "gt": ">", "amp": "&", "apos": "'", "quot": `"`
	// NOTE It doesn't do parameter entities, and we havnt necessarily
	// parsed any entities at all yet, so don't bother trying to use this.
	// NOTE If you dump all these, you find that there's a zillion of'em.
	parser.Entity = xml.HTMLEntity
	// Done!
	return parser
}

/*

// SetRootTagAndDoctype_xml takes a string.
// func GetPreambleDoctypeRootTag_xml(s string) (PI xml.ProcInst, DT string, RT xml.StartElement, err error) {
func GetPreambleDoctypeRootTag_xml(content string) (preambleTail string, doctypeTail string, rootTag xml.StartElement, err error) {
	var e error
	var s string
	var T xml.Token

	r := S.NewReader(content)
	var parser = xml.NewDecoder(r)
	parser.Strict    = false
	parser.AutoClose = xml.HTMLAutoClose
	parser.Entity    = xml.HTMLEntity

	var firstpass = true

	for {
		T, e = parser.Token()
		if e == io.EOF { break }
		if e != nil {
			return preambleTail, doctypeTail, rootTag,
				fmt.Errorf("pu.xml.getPRDTRT: %w", e)
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
				if firstpass && (preambleTail == "") {
					preambleTail = s
				} else {
					println("Got xml PI as non-first / repeated token ?!:", s)
				}
			}
			firstpass = false
		case xml.StartElement:
			// Found the XML root tag ?
			// type xml.StartElement struct { Name Name ; Attr []Attr }
			var tok xml.StartElement
			tok = xml.CopyToken(T).(xml.StartElement)
			rootTag = tok
			// println("XML-RT:", fmt.Sprintf("StElm:<%v>", tok))
			return preambleTail, doctypeTail, rootTag, err
		case xml.Directive:
			// Found the DOCTYPE ?
			// type Directive []byte
			var tok xml.Directive
			tok = xml.CopyToken(T).(xml.Directive)
			s = S.TrimSpace(string(tok))
			// s1, s2 := SU.SplitOffFirstWord(s)
			// println("XML-DT:", s1, "//", s2)
			if S.HasPrefix(s, "DOCTYPE ") {
				if doctypeTail != "" {
					println ("Second DOCTYPE ?!")
				} else {
					doctypeTail = S.TrimPrefix(s, "DOCTYPE ")
				}
			}
			firstpass = false
		default:
			firstpass = false
		}
	}
	return preambleTail, doctypeTail, rootTag, nil
}

*/
