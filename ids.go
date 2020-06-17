package xmlmodels

import (
	"encoding/xml"
	"fmt"
	S "strings"
	// FU "github.com/fbaube/fileutils"
	SU "github.com/fbaube/stringutils"
)

// DTDtypeFileExtensions are for content guessing.
var DTDtypeFileExtensions = []string{".dtd", ".mod", ".ent"}

type XmlSystemID string
type XmlPublicID string

// FPIfields holds the parsed results of (for example)
// "-//OASIS//DTD LIGHTWEIGHT DITA Topic//EN"
type FPIfields struct {
	// Registration is "+" or "-"
	Registration string
	// IsOasis but if not, then could be any of many others
	IsOasis bool
	// Organization is "OASIS" or maybe something else
	Organization string
	// PublicTextClass is typically "DTD" (filename.dtd)
	// or "ELEMENTS" (filename.mod)
	PublicTextClass string
	// PublicTextDesc is the distinguishing string,
	// e.g. PUBLIC "-//OASIS//DTD (_PublicTextDesc_)//EN".
	// It can end with the root tag of the document
	// (e.g. "Topic"). It can have an optional
	// embedded version number, such as "DITA 1.3".
	PublicTextDesc string
}

// XmlPublicIDcatalogRecord representa a line item from a parsed XML catalog file.
// One with a simple structure, such as the catalog file for LwDITA.
type XmlPublicIDcatalogRecord struct {
	XMLName  xml.Name `xml:"public"`
	// XmlPublicID is the DOCTYPE string
	XmlPublicID `xml:"publicId,attr"`
	FPIfields // PublicID 
	// XmlSystemID is the path to the file. Tipicly a relative filepath.
	XmlSystemID `xml:"uri,attr"`
	// The filepath long form, as resolved.
	// Note that we must use a string in order to avoid an import cycle.
	AbsFilePath string // FU.AbsFilePath
	Err error // in case an entry barfs
	// DoctypeIdentifierFields
}

func NewXmlPublicIDcatalogRecord(s string) (*XmlPublicIDcatalogRecord, error) {
	// println("NewXmlPublicID:", s)
	if s == "" {
		return nil, nil
	}
	if SU.IsXmlQuoted(s) {
		s = SU.MustXmlUnquote(s)
	}
	// -//OASIS//DTD LIGHTWEIGHT DITA Topic//EN
	var ss []string
	ss = S.Split(s, "/")
	// fmt.Printf("(DD) (%d) %#v \n", len(ss), ss)
	ss = SU.DeleteEmptyStrings(ss)
	// {"-", "OASIS", "DTD LIGHTWEIGHT DITA Topic", "EN"}
	// fmt.Printf("(DD) (%d) %#v \n", len(ss), ss)
	if len(ss) != 4 || ss[0] != "-" || ss[3] != "EN" {
		return nil, fmt.Errorf("Malformed Public ID<" + s + ">")
	}
	p := new(XmlPublicIDcatalogRecord)
	p.Registration = ss[0]
	p.Organization = ss[1]
	p.IsOasis = ("OASIS" == p.Organization)
	p.PublicTextClass, p.PublicTextDesc = SU.SplitOffFirstWord(ss[2])
	// ilog.Printf("PubID|%s|%s|%s|\n", p.Organization, p.PTClass, p.PTDesc)
	// fmt.Printf("(DD:pPID) PubID<%s|%s|%s>\n", p.Organization, p.PTClass, p.PTDesc)
	return p, nil
}

// Echo returns the public ID _unquoted_.
// <!DOCTYPE topic "-//OASIS//DTD LIGHTWEIGHT DITA Topic//EN">
func (p XmlPublicIDcatalogRecord) Echo() string {
	return fmt.Sprintf("%s//%s//%s %s//EN",
		p.Registration, p.Organization, p.PublicTextClass, p.PublicTextDesc)
}

// String returns the juicy part. For example,
// <!DOCTYPE topic "-//OASIS//DTD LIGHTWEIGHT DITA Topic//EN">
// maps to "DTD LIGHTWEIGHT DITA Topic".
func (p XmlPublicIDcatalogRecord) String() string {
	// return fmt.Sprintf("%s//%s//%s %s//EN",
	// p.Registration, p.Organization, p.PublicTextClass, p.PublicTextDesc)
	return fmt.Sprintf("%s %s", p.PublicTextClass, p.PublicTextDesc)
}
