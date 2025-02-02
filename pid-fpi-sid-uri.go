package xmlutils

import (
	"github.com/nbio/xml"
	"fmt"
	S "strings"

	SU "github.com/fbaube/stringutils"
)

// XmlPublicID = PID = Public ID = FPI = Formal Public Identifier
type XmlPublicID string

// XmlSystemID = SID = System ID = URI (Universal Resource Identifier)
// (can be a filepath or an HTTP address)
type XmlSystemID string

// PIDFPIfields holds the parsed results of a PID (PublicID) a.k.a. Formal
// Public Identifier, for example "-//OASIS//DTD LIGHTWEIGHT DITA Topic//EN"
type PIDFPIfields struct {
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

func (p PIDFPIfields) String() string {
	return fmt.Sprintf("%s,%s,%s,%s,%s",
		p.Registration, SU.Yn(p.IsOasis), p.Organization,
		p.PublicTextClass, p.PublicTextDesc)
}

// PIDSIDcatalogFileRecord representa a line item from a parsed XML catalog file.
// One with a simple structure, such as the catalog file for LwDITA. This same
// struct is also used to record the PID and/or SID of a DOCTYPE declaration.
type PIDSIDcatalogFileRecord struct {
	// XMLName probably does not ever need to be printed.
	XMLName xml.Name `xml:"public"`
	// XmlPublicID (PID) (FPI) is the DOCTYPE string
	XmlPublicID  `xml:"publicId,attr"`
	PIDFPIfields // PublicID
	// XmlSystemID is the path to the file. Tipicly a relative filepath.
	XmlSystemID `xml:"uri,attr"`
	// The filepath long form, as resolved.
	// Note that we must use a string in order to avoid an import cycle.
	AbsFilePath string // FU.AbsFilePath
	HttpPath    string
	Err         error // in case an entry barfs
	// DoctypeIdentifierFields
}

func (p *PIDSIDcatalogFileRecord) HasPID() bool {
	return p.XmlPublicID != ""
}

func (p *PIDSIDcatalogFileRecord) HasSID() bool {
	return p.XmlSystemID != ""
}

// NewPIDSIDcatalogFileRecord is pretty self-explanatory.
func NewPIDSIDcatalogFileRecord(pid string, sid string) (*PIDSIDcatalogFileRecord, error) {
	// println("NewXmlPublicID:", s)
	if pid == "" && sid == "" {
		return nil, nil
	}
	if SU.IsXmlQuoted(pid) {
		pid = SU.MustXmlUnquote(pid)
	}
	if SU.IsXmlQuoted(sid) {
		sid = SU.MustXmlUnquote(sid)
	}
	p := new(PIDSIDcatalogFileRecord)

	if pid != "" {
		// -//OASIS//DTD LIGHTWEIGHT DITA Topic//EN
		var ss []string
		ss = S.Split(pid, "/")
		// fmt.Printf("(DD) (%d) %#v \n", len(ss), ss)
		ss = SU.DeleteEmptyStrings(ss)
		// {"-", "OASIS", "DTD LIGHTWEIGHT DITA Topic", "EN"}
		// fmt.Printf("(DD) (%d) %#v \n", len(ss), ss)
		if len(ss) != 4 || ss[0] != "-" || ss[3] != "EN" {
			return nil, fmt.Errorf("Malformed Public ID: " + pid)
		}
		p.XmlPublicID = XmlPublicID(pid)
		p.Registration = ss[0]
		p.Organization = ss[1]
		p.IsOasis = (p.Organization == "OASIS")
		p.PublicTextClass, p.PublicTextDesc = SU.SplitOffFirstWord(ss[2])
		// ilog.Printf("PubID|%s|%s|%s|\n", p.Organization, p.PTClass, p.PTDesc)
		// fmt.Printf("(DD:pPID) PubID<%s|%s|%s>\n", p.Organization, p.PTClass, p.PTDesc)
	}
	if sid != "" {
		p.XmlSystemID = XmlSystemID(sid)
	}
	return p, nil
}

// Echo returns the public ID _unquoted_.
// <!DOCTYPE topic "-//OASIS//DTD LIGHTWEIGHT DITA Topic//EN">
func (p PIDSIDcatalogFileRecord) Echo() string {
	return fmt.Sprintf("%s//%s//%s %s//EN",
		p.Registration, p.Organization, p.PublicTextClass, p.PublicTextDesc)
}

// String returns the juicy part. For example,
// <!DOCTYPE topic "-//OASIS//DTD LIGHTWEIGHT DITA Topic//EN">
// maps to "DTD LIGHTWEIGHT DITA Topic".
func (p PIDSIDcatalogFileRecord) String() string {
	// return fmt.Sprintf("%s//%s//%s %s//EN",
	// p.Registration, p.Organization, p.PublicTextClass, p.PublicTextDesc)
	return fmt.Sprintf("%s %s", p.PublicTextClass, p.PublicTextDesc)
}

// DString returns a comprehensive dump.
func (p PIDSIDcatalogFileRecord) DString() string {
	var s string
	if p.XmlPublicID != "" {
		s += fmt.Sprintf("PID<%s:%s> ", p.XmlPublicID, p.PIDFPIfields)
	}
	if p.XmlSystemID != "" {
		s += fmt.Sprintf("SID<%s:%s> ", p.XmlSystemID, p.AbsFilePath)
	}
	if p.Err != nil {
		s += fmt.Sprintf("ERR<%s> ", p.Err.Error())
	}
	return s
}
