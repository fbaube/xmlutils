package xmlmodels

import (
	S "strings"
)

// DoctypeMType maps a DOCTYPE string to an MType string and a bool, Is it LwDITA?
type DoctypeMType struct {
	ToMatch  string
	MType    string
	IsLwDITA bool
}

// DTMTmap maps DOCTYPEs to MTypes (and: Is it LwDITA ?). This list
// should suffice for all ordinary XML files (except of course Docbook).
var DTMTmap = []DoctypeMType {
	   // This will require special handling
   { "html", "html/cnt/html5", false },
	   // uri="dtd/lw-topic.dtd"
	{ "//DTD LIGHTWEIGHT DITA Topic//", "xml/cnt/topic", true },
	{ "//DTD LW DITA Topic//", "xml/cnt/topic", true },
	{ "//DTD XDITA Topic//",  "html/cnt/topic", true },
		 // uri="dtd/lw-map.dtd"
	{ "//DTD LIGHTWEIGHT DITA Map//", "xml/map/---", true },
	{ "//DTD LW DITA Map//", "xml/map/---",  true },
	{ "//DTD XDITA Map//", "html/map/---", true },
	   // DITA 1.3
	{ "//DTD DITA Concept//", "xml/cnt/concept", false },
	{ "//DTD DITA Topic//", "xml/cnt/topic", false },
	{ "//DTD DITA Task//", "xml/cnt/task", false },
		 //
		 // https://www.w3.org/QA/2002/04/valid-dtd-list.html
		 //
	{ "//DTD HTML 4.",      "html/cnt/html4",    false },
	{ "//DTD XHTML 1.0 ",   "html/cnt/xhtml1.0", false },
	{ "//DTD XHTML 1.1//",  "html/cnt/xhtml1.1", false },
	{ "//DTD MathML 2.0//", "html/cnt/mathml",   false },
	{ "//DTD SVG 1.0//",    "xml/img/svg1.0",    false },
	{ "//DTD SVG 1.1",      "xml/img/svg",       false },
	{ "//DTD XHTML Basic 1.1//","html/cnt/topic",false },
	{ "//DTD XHTML 1.1 plus MathML 2.0 plus SVG 1.1//", "html/cnt/blarg", false },
}

// GetMTypeByDoctype returns MType "" if no match.
func GetMTypeByDoctype(dt string) (mtype string, isLwdita bool) {
	if dt == "<!DOCTYPE html>" {
		return "html/cont/html", false
	}
	for _,p := range DTMTmap {
		if S.Contains(dt, p.ToMatch) { return p.MType, p.IsLwDITA }
	}
	return "", false
}
