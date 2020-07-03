package xmlmodels

import (
	"fmt"
	"errors"
	S "strings"
	SU "github.com/fbaube/stringutils"
)

// XmlPreambleFields is a parse of an optional PI (processing instruction) at
// the start of an XML file. The most typical form is defined in the stdlib:
//
//  "<?xml version="1.0" encoding="UTF-8"?>" + "\n"
//
// Here the major version MUST be 1. XML has a version 1.1 but nobody uses it,
// so also the minor version MUST be 0, because that's what the Go stdlib XML
// parser understands, and anything else is gonna cause crazy breakage. Fields:
//
//  <?xml version="version_number"         <= required, "1.0"
//       encoding="encoding_declaration"   <= optional, assume "UTF-8"
//     standalone="standalone_status" ?>   <= optional, can be "yes", assume "no"
//
// Probably any errors returned by this function should be panicked on, because
// any such error is pretty fundamental and also ridiculous. Note also that
// strictly speaking, an XML preamble is NOT a PI.
type XmlPreambleFields struct {
	// Do not include a trailing newline.
	Preamble_raw string
	// e.g. "0" means XML 1.0
	MinorVersion string
	// Valid values and forms are TBS.
	Encoding     string
	// "yes"  or "no"
	IsStandalone bool
}

// NewXmlPreambleFields parses an XML preamble, which (BTW) MUST be the first
// line in a file. XML version MUST be "1.0". Encoding handling is incomplete.
//
//  - Example: <?xml version="1.0" encoding='UTF-8' standalone="yes"?>
//  - Also OK:   xml version="1.0" encoding='UTF-8' standalone="yes"
//  - Also OK:       version="1.0" encoding='UTF-8' standalone="yes"
//  - Also OK:   fields as documented for struct "XmlPreambleFields".
//
func NewXmlPreambleFields(s string) (*XmlPreambleFields, error) {
	if s == "" {
		return nil, nil
	}
	// println("Doing preamble:", s)
	// Be sure to trim a trailing newline.
	s = S.TrimSpace(s)
	p := new(XmlPreambleFields)

	// if matching outer brackets "<?xml .. ?>", remove them.
	if S.HasPrefix(s, "<?xml ") {
		if !S.HasSuffix(s, "?>") {
			return nil, fmt.Errorf("XML preamble is malformed: " + s)
		}
		s = S.TrimPrefix(s, "<?xml ")
		s = S.TrimSuffix(s, "?>")
		s = S.TrimSpace(s)
	} else if
		// if leading "xml ", remove it.
		S.HasPrefix(s, "xml ") {
		s = S.TrimPrefix(s, "xml ")
		s = S.TrimSpace(s)
	}
	p.Preamble_raw = "<?xml " + s + "?>"

	var props, sides []string
	var prop, varbl, value string
	// Break at spaces to get one to three properties.
	// println("Getting props from:", s)
	props = S.Split(s, " ")
	for _, prop = range props {
		sides = S.Split(prop, "=")
		varbl = sides[0]
		// println("Splitting prop:", prop)
		if !SU.IsXmlQuoted(sides[1]) {
			return p, errors.New("XML preamble property has bad quoting: " + prop)
		}
		value = SU.MustXmlUnquote(sides[1])

		switch varbl {
		case "encoding":
			p.Encoding = value
		case "version":
			if !S.HasPrefix(value, "1.") {
				return p, errors.New("XML preamble has bad XML major version number: " + value[:2])
			}
			p.MinorVersion = S.TrimPrefix(value, "1.")
			if "0" != p.MinorVersion {
				return p, errors.New("XML preamble has bad XML minor version number: " + p.MinorVersion)
			}
		case "standalone":
			// Let's say we may safely ignore bogus values.
			p.IsStandalone = (value == "yes")
			if value != "yes" && value != "no" {
				return p, errors.New("XML preamble has bad standalone property: " + value)
			}
		}
	}
	return p, nil
}

// Echo returns the raw preamble that was parsed, with a terminating newline.
func (xp XmlPreambleFields) Echo() string { return xp.Preamble_raw + "\n" }

// String includes a terminating newline.
func (xp XmlPreambleFields) String() string {
	var xmlver, encodg, stdaln string
	xmlver = fmt.Sprintf("<?xml version=\"1.%s\"", xp.MinorVersion)
	if xp.Encoding != "" {
		encodg = fmt.Sprintf(" encoding=\"%s\"", xp.Encoding)
	}
	if xp.IsStandalone {
		stdaln = " standalone=\"yes\""
	}
	return xmlver + encodg + stdaln + "?>\n"
}
