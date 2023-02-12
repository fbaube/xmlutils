package xmlutils

// This file contains LwDITA-specific stuff, but it is hard-coded
// and does not pull in other packages, so we leave it alone for now.

var knownRootTags = []string{"html", "map", "topic", "task", "concept", "reference"}

// DoctypeMType maps a DOCTYPE string to an MType string and a bool, Is it LwDITA?
type DoctypeMType struct {
	ToMatch       string
	DoctypesMType string
	RootElm       string
	IsLwDITA      bool
	// LwDITA, HTML5, and not much more (if any)
	IsInScope bool
}

// DTMTmap maps DOCTYPEs to MTypes (and: Is it LwDITA ?). This list
// should suffice for all ordinary XML files (except of course Docbook).
var DTMTmap = []DoctypeMType{
	// This will require special handling
	{"html", "html/cnt/html5", "html", false, true},
	// uri="dtd/lw-topic.dtd"
	{"//DTD LIGHTWEIGHT DITA Topic//", "xml/cnt/topic", "topic", true, true},
	{"//DTD LW DITA Topic//", "xml/cnt/topic", "topic", true, true},
	{"//DTD XDITA Topic//", "html/cnt/topic", "topic", true, true},
	// uri="dtd/lw-map.dtd"
	{"//DTD LIGHTWEIGHT DITA Map//", "xml/map/---", "map", true, true},
	{"//DTD LW DITA Map//", "xml/map/---", "map", true, true},
	{"//DTD XDITA Map//", "html/map/---", "map", true, true},
	// DITA 1.3
	{"//DTD DITA Concept//", "xml/cnt/concept", "concept", false, false},
	{"//DTD DITA Topic//", "xml/cnt/topic", "topic", false, false},
	{"//DTD DITA Task//", "xml/cnt/task", "task", false, false},
	//
	// https://www.w3.org/QA/2002/04/valid-dtd-list.html"
	// NOTE: The root element "html" of the document must contain an
	// xmlns declaration for the XHTML namespace [XMLNS]. The namespace
	// for XHTML is defined to be http://www.w3.org/1999/xhtml
	//
	{"//DTD HTML 4.", "html/cnt/html4", "html", false, false},
	{"//DTD XHTML 1.0 ", "html/cnt/xhtml1.0", "html", false, false},
	{"//DTD XHTML 1.1//", "html/cnt/xhtml1.1", "html", false, false},
	{"//DTD MathML 2.0//", "html/cnt/mathml", "", false, false},
	{"//DTD SVG 1.0//", "xml/img/svg1.0", "svg", false, false},
	{"//DTD SVG 1.1", "xml/img/svg", "svg", false, false},
	{"//DTD XHTML Basic 1.1//", "html/cnt/topic", "html", false, false},
	{"//DTD XHTML 1.1 plus MathML 2.0 plus SVG 1.1//", "html/cnt/blarg", "html", false, false},
}
