package xmlutils

// XmlContype categorizes the XML file. See variable "XmlContypes".
type XmlContype string

// XmlDoctype is just a DOCTYPE string, for example: <!DOCTYPE html>
type XmlDoctype string

// XmlContypes note: maybe DTDmod should be DTDelms.
var XmlContypes = []XmlContype{"Unknown", "DTD", "DTDmod", "DTDent",
	"RootTagData", "RootTagMixedContent", "MultipleRootTags", "INVALID"}
