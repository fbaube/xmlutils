package xmlmodels

import "encoding/xml"

// GetFirstStartElmByTag checks the start-element's tag's local name only,
// not any namespace. If no match, it returns an empty xml.StartElement.
// This func returns only the naked Token, without context, so it is meant
// only for processing XML catalog files. General XML processing should use
// the GToken version, which returns a GToken in the context of a tree structure.
func GetFirstStartElmByTag(tkzn []xml.Token, s string) xml.StartElement {
	if s == "" {
		return xml.StartElement{}
	}
	for _, p := range tkzn {
		if SE, OK := p.(xml.StartElement); OK {
			if SE.Name.Local == s {
				return SE
			}
		}
	}
	return xml.StartElement{}
}

// GetAllStartElmsByTag returns a new GTokenization.
// It checks the basic tag only, not any namespace.
func GetAllStartElmsByTag(tkzn []xml.Token, s string) []xml.StartElement {
	if s == "" {
		return nil
	}
	var ret []xml.StartElement
	ret = make([]xml.StartElement, 0)
	for _, p := range tkzn {
		if SE, OK := p.(xml.StartElement); OK {
			if SE.Name.Local == s {
				// fmt.Printf("found a match [%d] %s (NS:%s)\n", i, p.GName.Local, p.GName.Space)
				ret = append(ret, SE)
			}
		}
	}
	return ret
}

// GetAttVal returns the attribute's string value, or "" if not found.
func GetAttVal(se xml.StartElement, att string) string {
	for _, A := range se.Attr {
		if A.Name.Local == att {
			return A.Value
		}
	}
	return ""
}
