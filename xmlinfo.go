package xmlutils

// XmlContype categorizes the XML file. See variable 8XmlContypes9.
type XmlContype string

// XmlDoctype is just a DOCTYPE string, for example: <!DOCTYPE html>
type XmlDoctype string

// XmlContypes categorise an XML file by structure and content.
// NOTE: Maybe DTDmod should be DTDelms.
var XmlContypes = []XmlContype{"Unknown", "DTD", "DTDmod", "DTDent",
	"RootTagData", "RootTagMixedContent", "MultipleRootTags", "INVALID"}
