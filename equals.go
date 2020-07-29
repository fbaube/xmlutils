package xmlmodels

import "encoding/xml"

// Structure details of `xml.Name`:
//   type Name struct { Space, Local string }//
// Structure details of `xml.Attr`:
//   type Attr struct {
//     // xml.Name :: Space, Local string
//     Name  Name
//     Value string }

func XmlStartElmsEqual(arg1, arg2 xml.StartElement) bool {
	return XmlNamesEqual(arg1.Name, arg2.Name) &&
		XmlAttSlicesEqual(arg1.Attr, arg2.Attr)
}

func XmlNamesEqual(arg1, arg2 xml.Name) bool {
	return arg1.Space == arg2.Space && arg1.Local == arg2.Local
}

func XmlAttsEqual(arg1, arg2 xml.Attr) bool {
	return XmlNamesEqual(arg1.Name, arg2.Name) && arg1.Value == arg2.Value
}

func XmlAttSlicesEqual(arg1, arg2 []xml.Attr) bool {
	if (arg1 == nil || len(arg1) == 0) && (arg2 == nil || len(arg2) == 0) {
		return true
	}
	if arg1 == nil || arg2 == nil || len(arg1) == 0 || len(arg2) == 0 {
		return false
	}
	if len(arg1) != len(arg2) {
		return false
	}
	var i int
	var A xml.Attr
	for i, A = range arg1 {
		if !XmlAttsEqual(A, arg2[i]) {
			return false
		}
	}
	return true
}
