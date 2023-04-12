package xmlutils

import (
	// "encoding/xml"
	S "strings"
)

// XName is a generic XML name.
//
// type Name struct { Space, Local string }

func (p1 *XName) Equals(p2 *XName) bool {
	return p1.Space == p2.Space && p1.Local == p2.Local
}

func (p *XName) FixNS() {
	if p.Space != "" && !S.HasSuffix(p.Space, ":") {
		p.Space = p.Space + ":"
	}
}

// NewXName adds a colon to a non-empty namespace if it is not there already.
func NewXName(ns, local string) *XName {
	p := new(XName)
	if ns != "" && !S.HasSuffix(ns, ":") {
		ns += ":"
	}
	p.Space = ns
	p.Local = local
	return p
}
