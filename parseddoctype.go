package xmlutils

import (
// !!!! "fmt"
)

// This file contains LwDITA-specific stuff, but it is hard-coded
// and does not pull in other packages, so we leave it alone for now.

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

// ParsedDoctype is a parse of a complete DOCTYPE declaration.
// For [Lw]DITA, what interests us is something like
//
//	PUBLIC "-//OASIS//DTD (PublicTextDesc)//EN" or sometimes
//	PUBLIC "-//OASIS//ELEMENTS (PublicTextDesc)//EN" and
//	maybe followed by SYSTEM...
//
// The structure of a DOCTYPE is like so:
//   - PUBLIC | SYSTEM = Availability
//   - - = Registration = Organization & DTD are not registeredd with ISO.
//   - OASIS = Organization
//   - DTD = Public Text Class (CAPACITY | CHARSET | DOCUMENT |
//     DTD | ELEMENTS | ENTITIES | LPD | NONSGML | NOTATION |
//     SHORTREF | SUBDOC | SYNTAX | TEXT )
//   - (*) = Public Text Description, incl. any version number
//   - EN = Public Text Language
//   - URL = optional, explicit
//
// We don't include the raw DOCTYPE here because this structure can be optional
// but we still need to have the Doctype string in the DB as a separate column,
// even if it is empty (i.e. "").
type ParsedDoctype struct {
	RawDoctype string
	// PIDSIDcatalogFileRecord is the PID + SID.
	PIDSIDcatalogFileRecord
	// DTrootElm is the tag declared in the DOCTYPE, which
	// should match the root tag in the text of the file.
	DTrootElm string
	// MType is here because a DOCTYPE does indeed give
	// us enough information to create one.
	// DoctypeMType string
	error
}

// NewXmlDoctypeFieldsInclMType parses an XML DOCTYPE declaration.
// (Note that it does not however process internal DTD subsets.)
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

/* OBS print stuff

func (xdf ParsedDoctype) Echo() string {
	return "OOPS:TBS"
} // xd.raw + "\n" }

func (xdf ParsedDoctype) String() string {
	TT := xdf.DTrootElm
	if TT == "" {
		TT = "(no rootElm)"
	}
	// "-//OASIS//DTD LIGHTWEIGHT DITA Topic//EN" "lw-topic.dtd"
	return fmt.Sprintf("rootElm:%s,PIDSIDrec <|> %s <|>",
		TT /* dtmt, xdf.ContypingInfo, * / !!!!, xdf.PIDSIDcatalogFileRecord.DString())
}

func (xdf ParsedDoctype) DString() string {
	return xdf.String() // fmt.Sprintf("xm.xdf.DS: %+v", xdf)
}

*/
