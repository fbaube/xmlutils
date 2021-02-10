package xmlmodels

/*
type XmlInfo struct {
	XmlContype
	// nil if no preamble - defaults to xmlmodels.STD_PreambleFields
	*XmlPreambleFields
	// XmlDoctype
	// XmlDoctypeFields is a ptr - nil if there is no DOCTYPE declaration.
	*XmlDoctypeFields

	// TagDefCt is for DTD-type files (.dtd, .mod, .ent)
	// // TagDefCt int // Nr of <!ELEMENT ...>
	// RootTagIndex int  // Or some sort of pointer into the tree.
	// RootTagCt is >1 means mark the content as a Fragment.
	// // RootTagCt int

	// (Obs.cmt) XML items are
	//  - (DOCS) IDs & IDREFs
	//  - (DTDs) Elm defs (incl. Att defs) & Ent defs.

	// It is not precisely defined how to handle relative paths in external
	// IDs and entity substitutions, so we need to maintain this list.
	// EntSearchDirs []string // TODO

	// GEnts is "ENTITY"" directives (both with "%" and without).
	// GEnts map[string]*gparse.GEnt
	// DElms is "ELEMENT" directives.
	// DElms map[string]*gtree.GTag
	// TODO Maybe also add maps for NOTs (Notations)
}
*/

// XmlContype categorizes the XML file. See variable "XmlContypes".
type XmlContype string

// XmlDoctype is just a DOCTYPE string, for example: <!DOCTYPE html>
type XmlDoctype string

// XmlContypes note: maybe DTDmod should be DTDelms.
var XmlContypes = []XmlContype{"Unknown", "DTD", "DTDmod", "DTDent",
	"RootTagData", "RootTagMixedContent", "MultipleRootTags", "INVALID"}
