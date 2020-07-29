package xmlmodels

import "encoding/xml"

// LAToken is a location-aware token.
type LAToken struct {
	xml.Token
	Pos int // Position, from xml.Decoder
	Lnr int // Line number
	Col int // Column [number]
}
