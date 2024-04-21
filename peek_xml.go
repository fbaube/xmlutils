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
	
	// Structs are provided at the end of this file, for reference

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
			TorD = "drctv:" + T.ControlStrings[0] +
			     "," + T.ControlStrings[1] + "," + T.Text
			L.L.Dbg("Directive: TorD<%s> T.TDType<%s>",
				TorD, string(T.TDType)) 
		case CT.TD_type_PINST:
			TorD = "pr.i.: " + T.Text + "," + T.ControlStrings[0]
		case CT.TD_type_COMNT:
			TorD = "comnt:" + T.Text
		default:
			L.L.Panic("OOPS, bad TDT in LAToken: %s", T.TDType)
		}
		if !skippable {
			// fmt.Printf("XTKN: %s :: %s \n", TorD, T.DirectiveText)
			L.L.Dbg("XU.Peek: %s", TorD)
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
		     	// T is a CToken, so dump it when needed 
			// L.L.Warning("Directive! CToken: %+v", T)

		     	// html5: T.ControlStrings[0] is "DOCTYPE"
			// html5: T.ControlStrings[1] is "html"
			// html5: T.Text is "" 
			// Found the DOCTYPE ?
			// type Directive []byte
			
			// CS[0] is "DOCTYPE" etc.
			sDrctvSubtype := S.ToUpper(T.ControlStrings[0]) 
			// sRE = Root Element
			sRE := T.ControlStrings[1]
			// CS[1] is "HTML", "MAP", other root element 
			// if "DOCTYPE", else is WHOLE REMAINING TEXT
			// Text is the rest of a DOCTYPE declaration 
			
			// ============================
			//  Special handling for html5
			// ============================
			if sDrctvSubtype == "DOCTYPE" &&
			   S.EqualFold(T.ControlStrings[1], "HTML") &&
			   len(T.ControlStrings) == 2 {
			   	pPeek.DoctypeRaw = CT.Raw("<!DOCTYPE html>")
				L.L.Info("peek: Raw html DOCTYPE: %s",
					pPeek.DoctypeRaw)
				L.L.Info("peek: Got html5; TODO set bounds")
				// TODO: set beg+end bounds: for tag? more?
				pPeek.XmlRoot.TagName = "html" 
				didFirstPass = true
				return pPeek, nil
			}
			switch sDrctvSubtype {
			case "ELEMENT", "ATTLIST", "ENTITY", "NOTATION":
				pPeek.HasDTDstuff = true
				// TODO: Do something with sRE
			case "DOCTYPE":
				if pPeek.DoctypeRaw != "" {
				   L.L.Warning("XU.Peek: Got second " +
				   	"DOCTYPE; ignoring it")
				} else {
				   // len(T.ControlStrings) should be 2 
				   var theRest string 
				   theRest = T.Text // CtlStrings[2] // "etc etc"
				   pPeek.DoctypeRaw = CT.Raw("<!" +
					sDrctvSubtype + " " + sRE +
					" " + theRest + ">")
				   pPeek.XmlRoot.TagName = sRE
				   L.L.Info("XU.Peek: Raw DOCTYPE (recon" +
				   	 "stituted): %s", pPeek.DoctypeRaw)
				}
			}
			didFirstPass = true
		default:
			didFirstPass = true
		}
	}
	return pPeek, nil
}

/* REFERENCE MATERIAL

// LAToken is a location-aware XML token.
type LAToken struct {
        CToken
        FilePosition
}

// CToken is the lowest common denominator of tokens parsed
// from XML mixed content and other content-oriented markup.
// It has [stringutils.MarkupType].
//
// CToken:
//   - Common Token
//   - Content Token
//   - Combined Token
//   - Canonical Token
//   - Consolidated Token
//   - ConMuchoGusto Token :-P
//
// A CToken contains all that can be parsed from a token that
// is considered in isolation, as-is, without the context of
// surrounding markup. It should record/reflect/reproduce any
// XML (or HTML) token faithfully, and also accommodate any
// token from Markdown or (in the future) related markup
// such as Docbook or Asciidoc or RST (restructured text).
//
// The use of an XML-like data structure as the lingua franca
// is also meant to make XML-style automated processing simpler.
//
// The use of a single unified token representation is intended
// most of all to simplify & unify tokenisation across LwDITA's
// three supported input formats: XDITA XML, HDITA HTML5, and
// MDITA-XP Markdown. It also serves to represent all the
// various kinds of XML directives, including DTDs(!).
//
// Creation of a new CToken from an [encoding/xml.Token] is
// by design very straightforward, but creation from other
// types of token, such as HTML or Markdown, must be done
// in their other packages in order to prevent circular
// dependencies.
//
// For convenience & simplicity, some items in the struct
// are simply aliases for Go's XML structs, but then these
// must also be adaptable for Markdown. For example, when
// Pandoc-style attributes are used.
//
// 2024.04 Declarations like <!DOCTYPE html> are causing a
// lot of toruble, so they will be discussed in the code.
//
// CToken implements interface [stringutils.Stringser].
// .
type CToken struct {
	// ==================================
	// The original ("source code") token,
	// and other information about it
	// ==================================
	// SourceToken is the original token.
	// Keep it around "just in case".
	// TODO: Make this an Echoer !
	// Types:
	//  - XML: [xml.Token] from [xml.Decoder]
	//  - HTML: TBS
	//  - Markdown: TBS
	// Note that an XML Token is transitory,
	// so every Token has to be cloned, by
	// calling [xml.CopyToken].
	SourceToken interface{}
	// MarkupType of the original token; the value is
	// one of MU_type_(XML/HTML/MKDN/BIN/SQL/DIRLIKE). 
	// It is particularly helpful to have this info at the
	// token level when we consider that for example, we can
	// embed HTML tags in Markdown. Note that in the future,
	// each value could actually be a namespace declaration.
	SU.MarkupType
	// FilePosition is char position, and line nr & column nr.
	FilePosition

	// TDType comprises (a) the types of [xml.Token]
	// (they are all different struct's, actually),
	// plus (b) the (sub)types of [xml.Directive].
	// Note that [TD_type_ENDLM] ("EndElement") is
	// superfluous when token depth info is available.
	TDType
	// CName is ONLY for elements
	// (i.e. [TD_type_ELMNT] and [TD_type_ENDLM]).
	CName
	// CAtts is ONLY for [TD_type_ELMNT].
	CAtts
	// Text holds CDATA, and a PI's Instruction,
	// and a DOCTYPE's root element declaration,
	// and
	Text string
	// ControlStrings is tipicly XML PI & Directive stuff.
	// When it is used, its length is 1 or 2.
	//  - XML PI: the Target field
	//  - XML directive: the directive subtype
	// But this field also available for other data that
	// is not classifiable as source text.
	ControlStrings []string
}

func NewCTokenFromXmlToken(XT xml.Token) *CToken {

	case xml.Directive: // type Directive []byte
		ctkn.TDType = TD_type_DRCTV
		var fullDrctv, string0, tmp string
		fullDrctv = S.TrimSpace(string([]byte(XT.(xml.Directive))))
		// ctkn.Strings = make([]string, 3)
		string0, tmp = SU.SplitOffFirstWord(fullDrctv)
		// TODO: Assign TagOrDirective to ctkn.TDType
		// 2023.04 This works OK :-D
		if string0 != "DOCTYPE" {
			ctkn.ControlStrings = make([]string, 1)
			ctkn.ControlStrings[0] = string0
			ctkn.Text = tmp
			fmt.Printf("newCtkn L212 (!Drctv): %s: %s||%s \n",
				ctkn.TDType, string0, tmp)
		} else {
			ctkn.ControlStrings = make([]string, 2)
			ctkn.ControlStrings[0] = string0
			ctkn.ControlStrings[1], ctkn.Text =
				SU.SplitOffFirstWord(tmp)
			L.L.Okay("NewCtoken.Directive: TDType<%s> // " +
				"%s // %s\n\t%s",
				ctkn.TDType, string0,
				ctkn.ControlStrings[1], tmp)
		}
		// fmt.Printf("NewCtkn: Drctv: [0] %s [1] %s [2] %s",
		//	ctkn.Strings[0], ctkn.Strings[1], ctkn.Strings[2])
		// fmt.Printf("<!--Directive--> <!%s %s> \n",
		// 	outGT.Keyword, outGT.Otherwo rds)

*/