package xmlmodels

import(
	"fmt"
)

// DitaMarkupLg is a [Lw]DITA flavor. See enumeration "DitaMLs".
type DitaMarkupLg string
// DitaContype is a [Lw]DITA Topic, Map, etc. See enumeration "DitaContypes".
type DitaContype  string

// DitaInfo is two enumerations (so far): Markup language and Content type.
// They are both "" IFF the file is not DITA/LwDITA.
//  - ML: "1.2", "1.3", "XDITA", "HDITA", "MDATA".
//  - CT: "Map", "Bookmap", "Topic", "Task", "Concept", "Reference",
//        "Dita", "Glossary", "Conrefs", "LwMap", "LwTopic"
type DitaInfo struct {
	DitaMarkupLg
	DitaContype
}

// DitaMLs - see "type DitaMarkupLg".
var DitaMLs      = []DitaMarkupLg { "1.2", "1.3", "XDITA", "HDITA", "MDATA"}
// DitaContypes - see "type DitaContype".
var DitaContypes = []DitaContype  { "Map", "Bookmap", "Topic", "Task", "Concept",
	 "Reference", "Dita", "Glossary", "Conrefs", "LwMap", "LwTopic"}

func (di DitaInfo) String() string {
	return fmt.Sprintf("ML<%s> CT<%s>", di.DitaMarkupLg, di.DitaContype)
}

func (di DitaInfo) DString() string {
	return "<-- DITA " + string(di.DitaMarkupLg) +
		" " + string(di.DitaContype) + " -->\n"
}
