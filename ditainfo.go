package xmlutils

// DitaFlavor is a [Lw]DITA flavor. See enumeration "DitaFlavors".
type DitaFlavor string

// DitaContype is a [Lw]DITA Topic, Map, etc. See enumeration "DitaContypes".
type DitaContype string

// DitaInfo is two enumerations (so far): Markup flavor and Content type.
// They are both "" IFF the file is not DITA/LwDITA.
//  - MF: "1.2", "1.3", "XDITA", "HDITA", "MDATA".
//  - CT: "Map", "Bookmap", "Topic", "Task", "Concept", "Reference",
//        "Dita", "Glossary", "Conrefs", "LwMap", "LwTopic"
/*
type DitaInfo struct {
	DitaFlavor
	DitaContype
}
*/

// DitaFlavors - see "type DitaFlavor".
var DitaFlavors = []DitaFlavor{"1.2", "1.3", "XDITA", "HDITA", "MDATA"}

// DitaContypes - see "type DitaContype".
var DitaContypes = []DitaContype{"Map", "Bookmap", "Topic", "Task", "Concept",
	"Reference", "Dita", "Glossary", "Conrefs", "LwMap", "LwTopic"}

/*
func (di DitaInfo) String() string {
	return fmt.Sprintf("Flav<%s> CT<%s>", di.DitaFlavor, di.DitaContype)
}

func (di DitaInfo) DString() string {
	return "<-- DITA " + string(di.DitaFlavor) +
		" " + string(di.DitaContype) + " -->\n"
}
*/
