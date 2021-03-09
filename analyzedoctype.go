package xmlmodels

import (
	"errors"
	"fmt"
	S "strings"

	// FU "github.com/fbaube/fileutils"
	L "github.com/fbaube/mlog"
	SU "github.com/fbaube/stringutils"
)

// This file contains LwDITA-specific stuff, but it is hard-coded
// and does not pull in other packages, so we leave it alone for now.

var knownRootTags = []string{"html", "map", "topic", "task", "concept", "reference"}

// ContypingInfo has fields related to typing content (i.e. determining its type).
type ContypingInfo struct {
	FileExt  string
	MimeType string
	MType    string
	Doctype  string
	IsLwDita bool
	// IsProcbl means, is it processable (by us) ?
	// i.e. CAN we process it ? (Even if it might not be LwDITA.)
	IsProcbl bool
}

// DoctypeMType maps a DOCTYPE string to an MType string and a bool, Is it LwDITA?
type DoctypeMType struct {
	ToMatch       string
	DoctypesMType string
	RootElm       string
	IsLwDITA      bool
}

// DTMTmap maps DOCTYPEs to MTypes (and: Is it LwDITA ?). This list
// should suffice for all ordinary XML files (except of course Docbook).
var DTMTmap = []DoctypeMType{
	// This will require special handling
	{"html", "html/cnt/html5", "html", false},
	// uri="dtd/lw-topic.dtd"
	{"//DTD LIGHTWEIGHT DITA Topic//", "xml/cnt/topic", "topic", true},
	{"//DTD LW DITA Topic//", "xml/cnt/topic", "topic", true},
	{"//DTD XDITA Topic//", "html/cnt/topic", "topic", true},
	// uri="dtd/lw-map.dtd"
	{"//DTD LIGHTWEIGHT DITA Map//", "xml/map/---", "map", true},
	{"//DTD LW DITA Map//", "xml/map/---", "map", true},
	{"//DTD XDITA Map//", "html/map/---", "map", true},
	// DITA 1.3
	{"//DTD DITA Concept//", "xml/cnt/concept", "concept", false},
	{"//DTD DITA Topic//", "xml/cnt/topic", "topic", false},
	{"//DTD DITA Task//", "xml/cnt/task", "task", false},
	//
	// https://www.w3.org/QA/2002/04/valid-dtd-list.html"
	// NOTE: The root element "html" of the document must contain an xmlns
	// declaration for the XHTML namespace [XMLNS]. The namespace for XHTML
	// is defined to be http://www.w3.org/1999/xhtml
	//
	{"//DTD HTML 4.", "html/cnt/html4", "html", false},
	{"//DTD XHTML 1.0 ", "html/cnt/xhtml1.0", "html", false},
	{"//DTD XHTML 1.1//", "html/cnt/xhtml1.1", "html", false},
	{"//DTD MathML 2.0//", "html/cnt/mathml", "", false},
	{"//DTD SVG 1.0//", "xml/img/svg1.0", "svg", false},
	{"//DTD SVG 1.1", "xml/img/svg", "svg", false},
	{"//DTD XHTML Basic 1.1//", "html/cnt/topic", "html", false},
	{"//DTD XHTML 1.1 plus MathML 2.0 plus SVG 1.1//", "html/cnt/blarg", "html", false},
}

func (p ContypingInfo) String() (s string) {
	return fmt.Sprintf("filext:%s mtype:%s mimetype:%s", //  isLwdita:%s isProcbl:%s",
		p.FileExt, p.MType, p.MimeType) // , SU.Yn(p.IsLwDita), SU.Yn(p.IsProcbl))
}

// AnalyzeDoctype expects to receive a file extension plus a content type
// as determined by the HTTP stdlib. However a DOCTYPE is always considered
// authoritative, so this func can ignore things like the file extension,
// and overwrite or set any field it wants to.
//
// It works by first trying to match the DOCTYPE against a list. If that fails,
// stronger measures are called for.
//
// Note two things about this function:
//
// Firstly, it can handle PID, SID, or both:
//  <!DOCTYPE topic PUBLIC "-//OASIS//DTD LWDITA Topic//EN">
//  <!DOCTYPE topic PUBLIC "-//OASIS//DTD LWDITA Topic//EN" "./foo.dtd">
//  <!DOCTYPE topic SYSTEM "./foo.dtd">
//
// Secondly, it can handle a less-than-complete declaration:
//  DOCTYPE topic PUBLIC "-//OASIS//DTD LWDITA Topic//EN" (and variations)
//          topic PUBLIC "-//OASIS//DTD LWDITA Topic//EN" (and variations)
//                PUBLIC "-//OASIS//DTD LWDITA Topic//EN" (and variations)
//
// The last one is quite important because it is the format that appears
// in XML catalog files.
//
func (pC *ContypingInfo) AnalyzeXmlDoctype(aDoctype string) *XmlDoctypeFields {

	L.L.Dbg("xm.adt: inDoctp?<%s> inCntpg: %s", SU.Yn("" == aDoctype), pC.String())
	pC.IsLwDita = false
	pC.IsProcbl = false
	pDF := new(XmlDoctypeFields)
	pDF.ContypingInfo = *pC

	aDoctype = S.TrimSpace(aDoctype)

	// First, try to match the DOCTYPE. This is the former func
	// func GetMTypeByDoctype(dt string) (mtype string, isLwdita bool)

	// A quick win ?
	if aDoctype == "<!DOCTYPE html>" || aDoctype == "html" {
		pDF.DTrootElm = "html"
		pDF.MType = "html/cnt/html5"
		// Not sure about this next line
		pDF.PublicTextClass = "(HTML5)"
		L.L.Dbg("xm.adt: Got HTML5")
		return pDF
	}
	for _, p := range DTMTmap {
		if S.Contains(aDoctype, p.ToMatch) {
			pDF.MType = p.DoctypesMType
			pDF.DTrootElm = p.RootElm
			pDF.IsLwDita = p.IsLwDITA
			pDF.IsProcbl = p.IsLwDITA
			L.L.Dbg("xm.adt: Got a match on: " + pDF.MType)
			return pDF
		}
	}

	// OK so we did not match the DOCTYPE. So now let's analyze it in
	// excruciating detail.

	// Copied from mcfile.go:
	// [0] XML, BIN, TXT, MD
	// [1] IMG, CNT (Content), TOC (Map), SCH(ema)
	// [2] XML: per-DTD; BIN: fmt/filext; MD: flavor; SCH: fmt/filext

	// type XmlDoctypeFamily string
	//      XmlDoctypeFamilies are the broad groups of DOCTYPES.
	//  var XmlDoctypeFamilies = []XmlDoctypeFamily {
	//	"lwdita",
	//	"dita",
	//	"html5",
	//	"html",
	//	"other",
	// }

	// XmlDoctype is a parse of a complete DOCTYPE declaration.
	// For [Lw]DITA, what interests us is something like
	//  PUBLIC "-//OASIS//DTD (PublicTextDesc)//EN" or sometimes
	//  PUBLIC "-//OASIS//ELEMENTS (PublicTextDesc)//EN" and
	//  maybe followed by SYSTEM...
	//
	// The structure of a DOCTYPE is like so:
	//  * PUBLIC | SYSTEM = Availability
	//  * - = Registration = Organization & DTD are not registeredd with ISO.
	//  * OASIS = Organization
	//  * DTD = Public Text Class (CAPACITY | CHARSET | DOCUMENT |
	//      DTD | ELEMENTS | ENTITIES | LPD | NONSGML | NOTATION |
	//      SHORTREF | SUBDOC | SYNTAX | TEXT )
	//  * (*) = Public Text Description, incl. any version number
	//  * EN = Public Text Language
	//  * URL = optional, explicit
	//
	// We don't include the raw DOCTYPE here because this structure can be optional
	// but we still need to have the Doctype string in the DB as a separate column,
	// even if it is empty (i.e. "").

	// NewXmlDoctypeFieldsInclMType parses an XML DOCTYPE declaration.
	// (Note that it does not however process internal DTD subsets.)
	// Valid input forms:
	//
	//  <!DOCTYPE topic PUBLIC "-//OASIS//DTD LWDITA Topic//EN">
	//  <!DOCTYPE topic PUBLIC "-//OASIS//DTD LWDITA Topic//EN" "./foo.dtd">
	//  <!DOCTYPE topic SYSTEM "./foo.dtd">
	//    DOCTYPE topic PUBLIC "-//OASIS//DTD LWDITA Topic//EN" (and variations)
	//            topic PUBLIC "-//OASIS//DTD LWDITA Topic//EN" (and variations)
	//                  PUBLIC "-//OASIS//DTD LWDITA Topic//EN" (and variations)
	//
	// It should also work on a DOCTYPE reference plucked out of a DTD file,
	// one that tells the user what DOCTYPE declaration will reference the DTD.
	// In other words, the XML Catalog reference. Therefore, this function should
	// parse an input string that begins as minimally as a PUBLIC or SYSTEM (see
	// the last example above), and maybe don't worry about how the input string ends.
	//
	// Some input strings of great interest:
	//  DOCTYPE topic PUBLIC "-//OASIS//DTD LIGHTWEIGHT DITA Topic//EN"
	//  DOCTYPE map   PUBLIC "-//OASIS//DTD LIGHTWEIGHT DITA Map//EN"
	//  DOCTYPE html       (i.e. HTML5)
	//  DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.0 Transitional//EN" (MAYBE!)
	//

	// Second, analyze the DOCTYPE in detail. This is the former func
	// func NewXmlDoctypeFieldsInclMType(s string) (*XmlDoctypeFields, error)

	// But first we want it in a normalized form.

	// if brackets, remove them.
	if S.HasPrefix(aDoctype, "<!") {
		aDoctype = S.TrimPrefix(aDoctype, "<!")
		aDoctype = S.TrimSuffix(aDoctype, ">") // does not need to succeed
		aDoctype = S.TrimSpace(aDoctype)
	}
	// if leading "DOCTYPE ", remove it.
	if S.HasPrefix(aDoctype, "DOCTYPE ") {
		aDoctype = S.TrimPrefix(aDoctype, "DOCTYPE ")
		aDoctype = S.TrimSpace(aDoctype)
	}

	// Quick exit: HTML5
	if S.EqualFold(aDoctype, "html") || S.EqualFold(aDoctype, "html>") {
		println("==> Caught HTML5 doctype later rather than sooner ?!")
		pDF.DTrootElm = "html"
		pDF.MType = "html/cnt/html5"
		pDF.PublicTextClass = "(HTML5)"
		return pDF
	}

	// Possible here:
	// [topic] PUBLIC "-//OASIS//DTD LIGHTWEIGHT DITA Topic//EN" "lw-topic.dtd"

	// If we split off the first word, it should be PUBLIC, SYSTEM, or a root tag.
	var unk string
	unk, _ = SU.SplitOffFirstWord(aDoctype)
	if unk != "PUBLIC" && unk != "SYSTEM" {
		if !SU.IsInSliceIgnoreCase(unk, knownRootTags) {
			pDF.SetError(errors.New("Unrecognized DOCTYPE root element or " +
				"bad DOCTYPE availability (neither PUBLIC nor SYSTEM): " + unk))
		}
		pDF.DTrootElm, aDoctype = SU.SplitOffFirstWord(aDoctype)
	}
	var PubOrSys string
	PubOrSys, aDoctype = SU.SplitOffFirstWord(aDoctype)
	if PubOrSys != "PUBLIC" && PubOrSys != "SYSTEM" {
		panic("Lost the PUBLIC/SYSTEM")
		// return nil, fmt.Errorf("Bad DOCTYPE availability<" +
		// 	p.Availability + "> (neither PUBLIC nor SYSTEM)")
	}

	// Possible here:
	// "-//OASIS//DTD LIGHTWEIGHT DITA Topic//EN" "lw-topic.dtd"

	// =========================================================
	//  Now srsly parse the entire DOCTYPE. Do this even tho it
	//  probably won't bring any new information.
	// =========================================================
	// The next item(s) ("FPI" Public ID and/or "URI" System ID)
	// have to be a quoted strings. The spec says they use only
	// double quotes, not single, but let's be open-minded.
	// FIXME Handle cases of bad quoting.
	qtd1, qtd2, e := SU.SplitOffQuotedToken(aDoctype)
	if e != nil {
		pDF.SetError(fmt.Errorf("xm.adt.SplitOffQuotedToken(1)<%s>", aDoctype))
		return pDF
	}
	qtd2 = S.TrimSpace(qtd2)
	if qtd2 != "" {
		if !SU.IsXmlQuoted(qtd2) {
			pDF.SetError(fmt.Errorf("xm.adt.SplitOffQuotedToken(2)<%s>", aDoctype))
			return pDF
		}
		qtd2 = SU.MustXmlUnquote(qtd2)
	}

	// If both qtd1 and qtd2 are set then they must be FPI and URI.
	// If only qtd1 is set, it can be either FPI (PUBLIC) or URI (SYSTEM).
	var pPidSid *PIDSIDcatalogFileRecord

	if PubOrSys == "SYSTEM" {
		if qtd2 != "" {
			pDF.SetError(fmt.Errorf("xm.adt.SecondArgumentForSYSTEM: %s", qtd2))
			return pDF
		}
		pPidSid, e = NewPIDSIDcatalogFileRecord("", qtd1)
		// pDTF.PIDSIDcatalogFileRecord =
	} else if PubOrSys == "PUBLIC" {
		pPidSid, e = NewPIDSIDcatalogFileRecord(qtd1, qtd2)
		if e != nil {
			pDF.SetError(fmt.Errorf("xm.adt.NewXmlPublicID<%s|%s>: %w", qtd1, qtd2, e))
			return pDF
		}
	} else {
		panic("Unkwnown availability: " + PubOrSys)
	}
	pDF.PIDSIDcatalogFileRecord = *pPidSid

	sd := pDF.PIDSIDcatalogFileRecord.PIDFPIfields.PublicTextClass
	if sd == "" {
		println("!!> Odd exit")
		return pDF
	}
	// Now let's set the MType using some intelligent guesses,
	// to compare to the results of GetMTypeByDocType(..)
	if S.Contains(sd, "DITA") {
		pDF.MType = "dita/[TBS]/" + pDF.DTrootElm
	}
	if S.Contains(sd, "XDITA") ||
		S.Contains(sd, "LW DITA") ||
		S.Contains(sd, "LIGHTWEIGHT DITA") {
		pDF.MType = "lwdita/xdita/" + pDF.DTrootElm
	}
	/*
		if pDTF.TopTag != "" && Peek.RootTag != "" &&
			S.ToLower(pDTF.TopTag) != S.ToLower(Peek.RootTag) {
			fmt.Printf("--> RootTag MISMATCH: doctype<%s> bodytext<%s> \n",
				pDTF.TopTag, Peek.RootTag)
			panic("ROOT TAG MISMATCH")
		}
	*/
	if pDF.MType == "" {
		println("!!> No MType in AR!")
	}
	return pDF
}
