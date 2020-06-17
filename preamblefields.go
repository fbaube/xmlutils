package xmlmodels

import (
	"fmt"
	S "strings"
	SU "github.com/fbaube/stringutils"
)

// XmlPreambleFields is a parse of the optional PI "<?xml ... ?>"  at the
// start of the file. XML has a version 1.1 but nobody uses it. So here,
// the major version MUST be 1, and minor version 0, because that's what
// the Go stdlib XML parser understands. Note that strictly speaking,
// the preamble is NOT a PI.
//
//  <?xml version="version_number"         <= required, 1.0
//       encoding="encoding_declaration"   <= optional, assume "UTF-8"
//     standalone="standalone_status" ?>   <= optional, can be "yes", assume "no"
//
type XmlPreambleFields struct {
	Preamble_raw string
	MinorVersion string // e.g. expect "0" for XML 1.0
	// Encoding: not sure what the valid values are or what form they are.
	Encoding     string
	IsStandalone bool // "yes"  or "no"
}

// NewXmlPreambleFields parses an XML preamble,
// which MUST be the first line in a file.
//
//  - Example: <?xml version="1.0" encoding='UTF-8' standalone="yes"?>
//  - Also OK:   xml version="1.0" encoding='UTF-8' standalone="yes"
//  - Also OK:       version="1.0" encoding='UTF-8' standalone="yes"
//
func NewXmlPreambleFields(s string) (*XmlPreambleFields, error) {
	if s == "" {
		return nil, nil
	}
	// println("Doing preamble:", s)
	// Be sure to trim a trailing newline.
	s = S.TrimSpace(s)
	// if matching outer brackets "<?xml .. ?>", remove them.
	if S.HasPrefix(s, "<?xml ") {
		if !S.HasSuffix(s, "?>") {
			return nil, fmt.Errorf("Malformed preamble: " + s)
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
	p := new(XmlPreambleFields)
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
			return p, fmt.Errorf("xml.preamble.new.badquotes<%s>", prop)
		}
		value = SU.MustXmlUnquote(sides[1])

		switch varbl {
		case "encoding":
			p.Encoding = value
		case "version":
			p.MinorVersion = S.TrimPrefix(value, "1.")
			if "0" != p.MinorVersion {
				return p, fmt.Errorf("xml.preamble.new.: bad XML minor version number<" + p.MinorVersion + ">")
			}
		case "standalone":
			// We may safely ignore bogus values.
			p.IsStandalone = (value == "yes")
		}
	}
	return p, nil
}

func (xp XmlPreambleFields) Echo() string { return xp.Preamble_raw + "\n" }

// XmlPreambleFields includes a terminating newline.
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
