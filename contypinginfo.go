package xmlutils

import (
	"errors"
	"fmt"
	S "strings"

	// FU "github.com/fbaube/fileutils"
	L "github.com/fbaube/mlog"
	SU "github.com/fbaube/stringutils"
)

// ContypingInfo has simple fields related to
// typing content (i.e. determining its type).
// .
type ContypingInfo struct {
	FileExt         string
	MimeType        string
	MimeTypeAsSnift string
	MType           string
	// !! // !! // IsLwDita        bool
	// ?? FIXME ?? Doctype         string
	// IsProcbl means, is it processable (by us) ?
	// i.e. CAN we process it ? (Even if it might not be LwDITA.)
	// IsProcbl bool
}

func (p ContypingInfo) String() (s string) {
	return fmt.Sprintf("<%s> MimeType<%s>(stdlib:%s) MType<%s>",
		p.FileExt, p.MimeTypeAsSnift, p.MimeType, p.MType)
}

func (p ContypingInfo) MultilineString() (s string) {
	var mismatch string
	if p.MimeTypeAsSnift != p.MimeType {
		mismatch = fmt.Sprintf("(not: %s)", p.MimeType)
	}
	return fmt.Sprintf("<%s> \n\t MimeTp: %s %s \n\t M-Type: %s",
		p.FileExt, p.MimeTypeAsSnift, mismatch, p.MType)
}

// ParseDoctype should probably NOT be a method on ContypingInfo !!
//
// AnalyzeDoctype expects to receive a file extension plus a content
// type as determined by the HTTP stdlib. However a DOCTYPE is always
// considered authoritative, so this func can ignore things like the
// file extension, and overwrite or set any field it wants to.
//
// It works by first trying to match the DOCTYPE against a list.
// If that fails, stronger measures are called for.
//
// Note two things about this function:
//
// Firstly, it can handle PID, SID, or both:
//
//	<!DOCTYPE topic PUBLIC "-//OASIS//DTD LWDITA Topic//EN">
//	<!DOCTYPE topic PUBLIC "-//OASIS//DTD LWDITA Topic//EN" "./foo.dtd">
//	<!DOCTYPE topic SYSTEM "./foo.dtd">
//
// Secondly, it can handle a less-than-complete declaration:
//
//	DOCTYPE topic PUBLIC "-//OASIS//DTD LWDITA Topic//EN" (and variations)
//	        topic PUBLIC "-//OASIS//DTD LWDITA Topic//EN" (and variations)
//	              PUBLIC "-//OASIS//DTD LWDITA Topic//EN" (and variations)
//
// The last one is quite important because it is the format that appears
// in XML catalog files.
func (pC *ContypingInfo) ParseDoctype(aDoctype string) (*ParsedDoctype, error) {

	var pPDT *ParsedDoctype
	L.L.Dbg("xm.adt: inDoctp?<%s> inCntpg: %s", SU.Yn(aDoctype == ""), pC.String())
	// pC.IsLwDita = false
	// pC.IsProcbl = false
	pPDT = new(ParsedDoctype)
	aDoctype = S.TrimSpace(aDoctype)
	pPDT.RawDoctype = aDoctype

	// First, try to match the DOCTYPE. This is the former func
	// func GetMTypeByDoctype(dt string) (mtype string, isLwdita bool)

	// A quick win ?
	if aDoctype == "<!DOCTYPE html>" || aDoctype == "html" {
		pPDT.DTrootElm = "html"
		pC.MType = "html/cnt/html5"
		// Not sure about this next line
		pPDT.PublicTextClass = "(HTML5)"
		L.L.Dbg("xm.adt: Got HTML5")
		return pPDT, nil
	}
	for _, p := range DTMTmap {
		// println("(1) " + aDoctype)
		// println("(2) " + p.ToMatch)
		if S.Contains(aDoctype, p.ToMatch) {
			pPDT.DTrootElm = p.RootElm
			pC.MType = p.DoctypesMType
			// !! // !! // pC.IsLwDita = p.IsLwDITA
			// pC.IsProcbl = p.IsLwDITA
			L.L.Info("DOCTYPE matches: " + pC.MType)
			return pPDT, nil
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
		pPDT.DTrootElm = "html"
		pC.MType = "html/cnt/html5"
		pPDT.PublicTextClass = "(HTML5)"
		return pPDT, nil
	}

	// Possible here:
	// [topic] PUBLIC "-//OASIS//DTD LIGHTWEIGHT DITA Topic//EN" "lw-topic.dtd"

	// If we split off the first word, it should be PUBLIC, SYSTEM, or a root tag.
	var unk string
	unk, _ = SU.SplitOffFirstWord(aDoctype)
	if unk != "PUBLIC" && unk != "SYSTEM" {
		if !SU.IsInSliceIgnoreCase(unk, knownRootTags) {
			return pPDT, errors.New(
				"Unrecognized DOCTYPE root element or " +
					"bad DOCTYPE availability (neither PUBLIC nor SYSTEM): " + unk)
		}
		pPDT.DTrootElm, aDoctype = SU.SplitOffFirstWord(aDoctype)
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
		return pPDT, errors.New(
			"xm.adt.SplitOffQuotedToken(1): " + aDoctype)
	}
	qtd2 = S.TrimSpace(qtd2)
	if qtd2 != "" {
		if !SU.IsXmlQuoted(qtd2) {
			return pPDT, errors.New(
				"xm.adt.SplitOffQuotedToken(2):" + aDoctype)
		}
		qtd2 = SU.MustXmlUnquote(qtd2)
	}

	// If both qtd1 and qtd2 are set then they must be FPI and URI.
	// If only qtd1 is set, it can be either FPI (PUBLIC) or URI (SYSTEM).
	var pPidSid *PIDSIDcatalogFileRecord

	if PubOrSys == "SYSTEM" {
		if qtd2 != "" {
			return pPDT, errors.New(
				"xm.adt.SecondArgumentForSYSTEM: %s" + qtd2)
		}
		pPidSid, e = NewPIDSIDcatalogFileRecord("", qtd1)
		if e != nil {
			return pPDT, errors.New(fmt.Sprintf(
				"xm.adt.NewPIDSIDcatalogFileRecord<%s|%s>: %w", qtd1, qtd2, e))
		}
		// pDTF.PIDSIDcatalogFileRecord =
	} else if PubOrSys == "PUBLIC" {
		pPidSid, e = NewPIDSIDcatalogFileRecord(qtd1, qtd2)
		if e != nil {
			return pPDT, errors.New(fmt.Sprintf(
				"xm.adt.NewXmlPublicID<%s|%s>: %w", qtd1, qtd2, e))
		}
	} else {
		panic("Unkwnown availability: " + PubOrSys)
	}
	pPDT.PIDSIDcatalogFileRecord = *pPidSid

	sd := pPDT.PIDSIDcatalogFileRecord.PIDFPIfields.PublicTextClass
	if sd == "" {
		println("!!> Odd exit")
		return pPDT, nil
	}
	// Now let's set the MType using some intelligent guesses,
	// to compare to the results of GetMTypeByDocType(..)
	if S.Contains(sd, "DITA") {
		pC.MType = "dita/[TBS]/" + pPDT.DTrootElm
	}
	if S.Contains(sd, "XDITA") ||
		S.Contains(sd, "LW DITA") ||
		S.Contains(sd, "LIGHTWEIGHT DITA") {
		pC.MType = "lwdita/xdita/" + pPDT.DTrootElm
	}
	if pC.MType == "" {
		println("!!> No MType in AR!")
	}
	return pPDT, nil
}
