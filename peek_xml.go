package xmlutils

// Problem:
// ! hasRootTag
// pPeek.RawDoctype is ""
// pPeek.RawPreamble is ""

import (
	"github.com/nbio/xml"
	"fmt"
	S "strings"

	CT "github.com/fbaube/ctoken"
	L "github.com/fbaube/mlog"
	_ "github.com/fbaube/fileutils" // for docu 
)

// XmlPeek is used by [fileutils.AnalyseFile] 
// when preparing a [fileutils.AnalysisRecord].
// Note that ContentityBasics has chunks of 
// Raw but not the full "Raw" string.
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
// It is called by [fileutils.AnalyzeFile].
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
			L.L.Debug("Directive: TorD<%s> T.TDType<%s>",
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
			L.L.Debug("XU.Peek: %s", TorD)
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
				L.L.Debug("XML-PrI <%s> <%s>",
					T.Text, T.ControlStrings[0])
				if (pPeek.PreambleRaw == "") && !didFirstPass {
					pPeek.PreambleRaw = CT.Raw("<?xml " + sInst + "?>")
					L.L.Debug("GOT Raw.PREAMBLE: %s", pPeek.PreambleRaw)
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
				L.L.Debug("End xml root at: " + LAT.FilePosition.Info())
			case pPeek.Meta.TagName:
				pPeek.Meta.End = LAT.FilePosition
				pPeek.Meta.End.Pos += len(localName) + 3
				L.L.Debug("End meta elm at: " + LAT.FilePosition.Info())
			case pPeek.Text.TagName:
				pPeek.Text.End = LAT.FilePosition
				pPeek.Text.End.Pos += len(localName) + 3
				L.L.Debug("End text elm at: " + LAT.FilePosition.Info())
			}

		case CT.TD_type_ELMNT:
			// type xml.StartElement struct { Name Name ; Attr []Attr }
			localName := T.CName.Local
			if !foundRootElm {
				L.L.Info("Found root tag: " + localName)
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
					L.L.Debug("Got key elm.beg <%s>:%s => meta<%s> text<%s>",
						localName, pPeek.XmlRoot.Beg.Info(),
						metaTagToFind, textTagToFind)
				}
			} else {
				if localName == metaTagToFind {
					pPeek.Meta.TagName = localName
					pPeek.Meta.Atts = T.CAtts.AsStdLibXml() // tok.Attr
					pPeek.Meta.Beg = LAT.FilePosition
					L.L.Debug("Got meta elm <%s> at %s",
						metaTagToFind, pPeek.Meta.Beg.Info())
				}
				if localName == textTagToFind {
					pPeek.Text.TagName = localName
					pPeek.Text.Atts = T.CAtts.AsStdLibXml() //tok.Attr
					pPeek.Text.Beg = LAT.FilePosition
					L.L.Debug("Got text elm <%s> at %s \n",
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
