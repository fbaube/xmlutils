package xmlmodels

import (
	"fmt"
	S "strings"
	SU "github.com/fbaube/stringutils"
)

// This file contains LwDITA-specific stuff, but it is hard-coded
// and does not pull in other packages, so we leave it alone for now.

var knownRootTags = []string{"html", "map", "topic", "task", "concept", "reference"}

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
// For [Lw]DITA, what interests us is
// PUBLIC "-//OASIS//DTD (PublicTextDesc)//EN" or sometimes
// PUBLIC "-//OASIS//ELEMENTS (PublicTextDesc)//EN" and
// maybe followed by SYSTEM...
//  * PUBLIC | SYSTEM = Availability
//  * - = Reg'n = Org'zn & DTD are not reg'd with ISO.
//  * OASIS = Org'zn
//  * DTD = Public Text Class (CAPACITY | CHARSET | DOCUMENT |
//    DTD | ELEMENTS | ENTITIES | LPD | NONSGML | NOTATION |
//    SHORTREF | SUBDOC | SYNTAX | TEXT )
//  * (*) = Public Text Description, incl. any version number
//  * EN = Public Text Language
//  * URL = optional, explicit
//
// We don't include the RAW Doctype here cos this field can
// be nil but the Doctype string needs to be in the DB as a
// separate column, even if it is empty.
//
type XmlDoctypeFields struct {
	Availability string // "PUBLIC" or "SYSTEM"
	FPIfields
	XmlPublicIDcatalogRecord
	// TopTag is the tag declared in the DOCTYPE
	TopTag string
	// MType is here because a DOCTYPE does indeed let us create one.
	DoctypeMType string
}

// NewXmlDoctypeFieldsInclMType parses an XML DOCTYPE declaration.
// (However it does not process internal DTD subsets.)
//
//  <!DOCTYPE topic PUBLIC "-//OASIS//DTD LIGHTWEIGHT DITA Topic//EN">
//  <!DOCTYPE topic PUBLIC "-//OASIS//DTD LWDITA Topic//EN" "./foo.dtd">
//  <!DOCTYPE topic SYSTEM "./foo.dtd">
//    DOCTYPE topic PUBLIC "-//OASIS//DTD LWDITA Topic//EN" (etc. etc.)
//            topic PUBLIC "-//OASIS//DTD LWDITA Topic//EN" (etc. etc.)
//                  PUBLIC "-//OASIS//DTD LWDITA Topic//EN" (etc. etc.)
//
// It should also work on a DOCTYPE reference plucked out of a DTD file,
// one that tells the user what DOCTYPE declaration will reference the DTD.
// In other words, the XML Catalog reference. Therefore, let this function
// parse a string that begins as minimally as a PUBLIC or SYSTEM (see the
// last example above), and maybe don't worry about how the string ends.
//
// Target strings of great interest:
//  DOCTYPE topic PUBLIC "-//OASIS//DTD LIGHTWEIGHT DITA Topic//EN"
//  DOCTYPE map   PUBLIC "-//OASIS//DTD LIGHTWEIGHT DITA Map//EN"
//  DOCTYPE html     (i.e. HTML5)
//  DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.0 Transitional//EN" (MAYBE!)
//
func NewXmlDoctypeFieldsInclMType(s string) (*XmlDoctypeFields, error) {

	if s == "" {
		return nil, nil
	}
	pDTF := new(XmlDoctypeFields)
	// println("Doing doctype:", s)

	// Be sure to trim any trailing space.
	s = S.TrimSpace(s)
	// if brackets, remove them.
	if S.HasPrefix(s, "<!") {
		s = S.TrimPrefix(s, "<!")
		s = S.TrimSuffix(s, ">") // does not need to succeed
		s = S.TrimSpace(s)
	}
	// if leading "DOCTYPE ", remove it.
	if S.HasPrefix(s, "DOCTYPE ") {
		s = S.TrimPrefix(s, "DOCTYPE ")
		s = S.TrimSpace(s)
	}
	pDTF.DoctypeMType = "-/-/-"

	// println("--> Parsing doctype:\n    ", s)

	// Quick exit: HTML5
	if S.EqualFold(s, "html") || S.EqualFold(s, "html>") {
		pDTF.TopTag = "html"
		pDTF.DoctypeMType = "html/cnt/html5"
		pDTF.PublicTextClass = "(HTML5)"
		return pDTF, nil
	}

	// Possible here:
	// [topic] PUBLIC "-//OASIS//DTD LIGHTWEIGHT DITA Topic//EN" "lw-topic.dtd"

	// If we split off the first word, it should be PUBLIC, SYSTEM, or a root tag.
	var unk string
	unk, _ = SU.SplitOffFirstWord(s)
	if unk != "PUBLIC" && unk != "SYSTEM" {
		if !SU.IsInSliceIgnoreCase(unk, knownRootTags) {
			return nil, fmt.Errorf("<%s>: Unrecognized DOCTYPE root element or " +
				"bad DOCTYPE availability (neither PUBLIC nor SYSTEM)", unk)
		}
		pDTF.TopTag, s = SU.SplitOffFirstWord(s)
	}
	pDTF.Availability, s = SU.SplitOffFirstWord(s)
	if pDTF.Availability != "PUBLIC" && pDTF.Availability != "SYSTEM" {
		panic("Lost the PUBLIC/SYSTEM")
		// return nil, fmt.Errorf("Bad DOCTYPE availability<" +
		// 	p.Availability + "> (neither PUBLIC nor SYSTEM)")
	}

	// ===============================
	//  Let's use our DOCTYPE matcher
	// ===============================
	// var isLwdita bool
	pDTF.DoctypeMType, _ = GetMTypeByDoctype(s)
	// println("-->", "DT/MType search results (#2):", mt)

	// Possible here:
	// "-//OASIS//DTD LIGHTWEIGHT DITA Topic//EN" "lw-topic.dtd"

	// =========================================================
	//  Now srsly parse the entire DOCTYPE. Do this even tho it
	//  is optional - we can't get any usable info that we did
	//  not get from the preceding call to GetMTypeByDoctype.
	// =========================================================
	// The next item(s) ("FPI" Public ID and/or "URI" System ID)
	// have to be a quoted strings. The spec says they use only
	// double quotes, not single, but let's be open-minded.
	// FIXME Handle cases of bad quoting.
	qtd1, qtd2, e := SU.SplitOffQuotedToken(s)
	if e != nil {
		return pDTF, fmt.Errorf("xm.dtflds.SplitOffQuotedToken(1)<%s>", s)
	}
	qtd2 = S.TrimSpace(qtd2)
	if qtd2 != "" {
		if !SU.IsXmlQuoted(qtd2) {
			return pDTF, fmt.Errorf("xm.dtflds.SplitOffQuotedToken(2)<%s>", s)
		}
		qtd2 = SU.MustXmlUnquote(qtd2)
	}
	if pDTF.Availability == "SYSTEM" {
		if qtd2 != "" {
			return pDTF, fmt.Errorf("xm.dtflds.SecondSYSTEMargument<%s>", qtd2)
		}
		pDTF.XmlSystemID = XmlSystemID(qtd1)
	} else if pDTF.Availability == "PUBLIC" {
		ppid, e := NewXmlPublicIDcatalogRecord(qtd1)
		if e != nil {
			return nil, fmt.Errorf("xm.dtflds.NewXmlPublicID<%s>: %w", qtd1, e)
		}
		pDTF.XmlPublicIDcatalogRecord = *ppid
		pDTF.XmlSystemID = XmlSystemID(qtd2)
	} else {
		panic("Unkwnown availability: " + pDTF.Availability)
	}

	sd := pDTF.XmlPublicIDcatalogRecord.FPIfields.PublicTextClass
	if sd == "" {
		return pDTF, nil
	}
	// Now let's set the MType using some intelligent guesses,
	// to compare to the results of GetMTypeByDocType(..)
	if S.Contains(sd, "DITA") {
		 pDTF.DoctypeMType = "dita/[TBS]/" + pDTF.TopTag
	}
	if S.Contains(sd, "XDITA") ||
		 S.Contains(sd, "LW DITA") ||
		 S.Contains(sd, "LIGHTWEIGHT DITA") {
		 pDTF.DoctypeMType = "lwdita/xdita/" + pDTF.TopTag
	}
	return pDTF, nil
}

func (xdf XmlDoctypeFields) Echo() string {
	return "OOPS:TBS"
	} // xd.raw + "\n" }

func (xdf XmlDoctypeFields) String() string {
	if "" == xdf.TopTag {
		panic("xdf.TopTag")
	}
	// "-//OASIS//DTD LIGHTWEIGHT DITA Topic//EN" "lw-topic.dtd"
	return fmt.Sprintf("(%s,%s,%s) %s", xdf.Availability, xdf.TopTag,
		xdf.DoctypeMType, xdf.XmlPublicIDcatalogRecord)
}

func (xdf XmlDoctypeFields) DString() string {
	return fmt.Sprintf("xm.xdf.DS: %+v", xdf)
}
