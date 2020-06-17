package xmlmodels

import(
	"fmt"
)

type DitaMarkupLg string
type DitaContype  string

// DitaInfo is two enums (so far): Markup language & Content type.
//  - ML: "1.2", "1.3", "XDITA", "HDITA", "MDATA".
//  - CT: "Map", "Bookmap", "Topic", "Task", "Concept", "Reference",
//        "Dita", "Glossary", "Conrefs", "LwMap", "LwTopic"
type DitaInfo struct {
	// These next two are "" IFF the file is not DITA/LwDITA.
	DitaMarkupLg
	DitaContype
}

var DitaMLs      = []DitaMarkupLg { "1.2", "1.3", "XDITA", "HDITA", "MDATA"}
var DitaContypes = []DitaContype  { "Map", "Bookmap", "Topic", "Task", "Concept",
	 "Reference", "Dita", "Glossary", "Conrefs", "LwMap", "LwTopic"}

func (di DitaInfo) String() string {
	return fmt.Sprintf("ML<%s> CT<%s>", di.DitaMarkupLg, di.DitaContype)
}

func (di DitaInfo) DString() string {
	return "<-- DITA " + string(di.DitaMarkupLg) +
		" " + string(di.DitaContype) + " -->\n"
}
