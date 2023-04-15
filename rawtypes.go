package xmlutils

import SU "github.com/fbaube/stringutils"

type Raw string

type TypedRaw struct {
	Raw
	SU.MarkupType
}

func (p *TypedRaw) S() string {
	return string(p.Raw)
}

// RawType is a convenience function so that
// if (i.e. when) it becomes convenient, the
// elements of [TypedRaw] can be unexported.
// .
func (p *TypedRaw) RawType() SU.MarkupType {
	return p.MarkupType
}
