package xmlutils

import (
	"github.com/nbio/xml"
	"fmt"
	CT "github.com/fbaube/ctoken"
	"io"
	S "strings"
)

// ParserResults_xml is TBS.
type ParserResults_xml struct {
	// NodeSlice is driven by the stdlib XML parser in [encoding/xml].
	NodeSlice []CT.CToken // []xml.Token
	CommonCPR
}

// GenerateParserResults_xml is TBS.
func GenerateParserResults_xml(s string) (*ParserResults_xml, error) {
	var nl []CT.CToken // []xml.Token
	var e error
	nl, e = DoParse_xml(s)
	if e != nil {
		return nil, fmt.Errorf("pu.xml.parseResults: %w", e)
	}
	p := new(ParserResults_xml)
	p.CommonCPR = *NewCommonCPR()
	p.NodeSlice = nl
	p.CPR_raw = s
	return p, nil
}

// DoParse_xml takes a string, so we can assume that we can
// discard it after use cos the caller has another copy of it.
// To be safe, it copies every token using `xml.CopyToken(T)`.
func DoParse_xml(s string) (xtokens []CT.CToken, err error) {
	return doParse_xml_maybeRaw(s, false)
}

// DoParseRaw_xml( is TBS.
func DoParseRaw_xml(s string) (xtokens []CT.CToken, err error) {
	return doParse_xml_maybeRaw(s, true)
}

// doParse_xml_maybeRaw is TBS.
func doParse_xml_maybeRaw(s string, doRaw bool) (xtokens []CT.CToken, err error) {
	var e error
	var TT *CT.CToken // xml.Token
	var ttt xml.Token
	xtokens = make([]CT.CToken, 0, 100)
	// println("(DD) XmlTokenizeBuffer:", s)

	r := S.NewReader(s)
	var parser *xml.Decoder
	parser = NewConfiguredDecoder(r)

	for {
		if doRaw {
			// func (d *Decoder) RawToken() (Token, error) API:
			// RawToken is like func [xml.Token] but (1) does 
			// not verify that start & end elements match, and
			// (2) does not translate namespace prefixes to
			// their corresponding URLs.
			ttt, e = parser.RawToken()
		} else {
			// func (d *Decoder) Token() (Token, error) API:
			// Token returns the next XML token in the input
			// stream. At the end of the input stream, Token
			// returns nil, io.EOF.
			// Token expands self-closing elements such as
			// <br> into separate start and end elements
			// returned by successive calls.
			// Token guarantees that the StartElement and
			// EndElement tokens it returns are properly
			// nested and matched: if Token encounters an
			// unexpected end element or EOF before all
			// expected end elements, it returns an error.
			// NAMESPACES: Token implements XML name spaces as
			// described by https://www.w3.org/TR/REC-xml-names/ .
			// Each of the Name structures contained in the
			// Token has the Space set to the URL identifying
			// its name space when known. If Token encounters
			// an unrecognized name space prefix, it uses the
			// prefix as the Space rather than report an error.
			ttt, e = parser.Token()
		}
		if e == io.EOF {
			break
		}
		if e != nil {
			return xtokens, fmt.Errorf("xu.xml.doParse.1: %w", e)
		}
		// TT = xml.CopyToken(T)
		TT = CT.NewCTokenFromXmlToken(ttt)
		if TT == nil {
		   // fmt.Printf("Got NIL CToken for whitespace token \n")
		   continue
		   }
		xtokens = append(xtokens, *TT)
	}
	return xtokens, nil
}

// DoParse_xml_locationAware is TBS.
func DoParse_xml_locationAware(s string) (xtokens []CT.LAToken, err error) {
	var e error
	var T xml.Token
	var pXT *CT.CToken
	var XT CT.CToken // xml.Token
	var LAT CT.LAToken
	xtokens = make([]CT.LAToken, 0, 100)
	// println("(DD) XmlTokenizeBuffer:", s)
	// var idcs []int
	// idcs = SU.AllIndices(s, "\n")

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
			return xtokens, fmt.Errorf("pu.xml.doParse.2: %w", e)
		}
		pXT = CT.NewCTokenFromXmlToken(T)
		if pXT == nil {
		   // fmt.Printf("Got NIL CToken for whitespace token \n")
			continue
		}
		XT = *pXT

		LAT = *new(CT.LAToken)
		LAT.CToken = XT
		LAT.Pos = pos
		// InputPos returns the line of the current decoder
		// position and the 1 based input position of the line.
		// The position gives the location of the end of the
		// most recently returned token.
		LAT.Lnr, LAT.Col = // LnrAndColFromPos(pos, idcs)
			parser.InputPos()
		// ll, cc := LnrAndColFromPos(pos, idcs)
		// fmt.Printf("OLD L%d C%d NEW L%d C%d \n",
		//     ll, cc, LAT.Lnr, LAT.Col)
		xtokens = append(xtokens, LAT)
	}
	return xtokens, nil
}

// NewConfiguredDecoder returns a new [xml.Decoder] that has been
// confgured with non-strict namespace parsing, HTML auto-closing
// tags, and HTML entities. 
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

func (p *ParserResults_xml) NodeCount() int {
        return len(p.NodeSlice)
}

func (p *ParserResults_xml) NodeDebug(i int) string {
        if i >= len(p.NodeSlice) {
                return "(indexOverrun)"
	}
        h := p.NodeSlice[i]
        return fmt.Sprintf("|%+v|", h)
        /* return fmt.Sprintf("|tp:%d:%s,data:%s,ns:%s,kids:%s,atts:%v|",
		h.Type, NodeTypeString[h.Type], DataOfHtmlNode(&h),
                h.Namespace, SU.Yn(h.FirstChild != nil), h.Attr) */
}

func (p *ParserResults_xml) NodeEcho(i int) string {
        if i >= len(p.NodeSlice) {
                return "(indexOverrun)"
        }
	return "(XML NodeEcho)"
	/*
        // var pBB *bytes.Buffer
        var pSB = new(S.Builder)
        // FIXME this call is recursive!
        // html.Render(pSB, p.NodeSlice[i])
        n := p.NodeSlice[i]
        FC, LC := n.FirstChild, n.LastChild
        n.FirstChild, n.LastChild = nil, nil
        html.Render(pSB, n)
        n.FirstChild, n.LastChild = FC, LC
        return pSB.String()
	*/
}

func (p *ParserResults_xml) NodeInfo(i int) string {
        if i >= len(p.NodeSlice) {
                return "(indexOverrun)"
        }
        // return fmt.Sprintf("<h[%d] lv%d,ch%d,%s>",
        //        i, p.NodeDepths[i], p.FilePosns[i].Pos, p.NodeDebug(i))
        return fmt.Sprintf("<h[%d] %s>", i, p.NodeDebug(i))
}

