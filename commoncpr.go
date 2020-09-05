package xmlmodels

import (
	"fmt"
	"io"
)

type CommonCPR struct {
	NodeDepths []int
	FilePosns  []*FilePosition
	CPR_raw    string
	DumpDest   io.Writer
}

func NewCommonCPR() *CommonCPR {
	p := new(CommonCPR)
	p.NodeDepths = make([]int, 0)
	p.FilePosns = make([]*FilePosition, 0)
	return p
}

func (p *CommonCPR) AsString(i int) string {
	if p.NodeDepths == nil {
		println("OOPS NodeDepths")
	}
	if p.FilePosns == nil {
		println("OOPS FilePosns")
	}
	fmt.Printf("## CmnCPR: i %d nd %d fp %d \n", i, len(p.NodeDepths), len(p.FilePosns))
	return fmt.Sprintf("i%02d,Lv%02d,%s", i, p.NodeDepths[i], p.FilePosns[i])
}
