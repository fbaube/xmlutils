package xmlutils

import (
	// "encoding/xml"
	"fmt"
)

// LAToken is a location-aware XML token.
type LAToken struct {
	// xml.Token
	XToken
	FilePosition
}

// FilePosition is a char.position (from xml.Decoder)
// plus line nr & column nr (both calculated).
type FilePosition struct {
	Pos int // byte Position in file, from xml.Decoder.InputOffset()
	Lnr int // Line number
	Col int // Column [number]
}

func (fp FilePosition) String() string {
	if fp.Lnr == 0 && fp.Col == 0 {
		return fmt.Sprintf("[%03d]", fp.Pos)
	}
	return fmt.Sprintf("[%d](L%02dc%02d)", fp.Pos, fp.Lnr, fp.Col)
}

func NewFilePosition(i int) *FilePosition {
	p := new(FilePosition)
	p.Pos = i
	return p
}
