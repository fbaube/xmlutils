package xmlmodels

import (
	"encoding/xml"
	"fmt"
)

type FilePosition struct {
	Pos int // Position, from xml.Decoder
	Lnr int // Line number
	Col int // Column [number]
}

func (fp FilePosition) String() string {
	return fmt.Sprintf("L%02dc%02d", fp.Lnr, fp.Col)
}

// LAToken is a location-aware token.
type LAToken struct {
	xml.Token
	FilePosition
}
