package xmlutils

import (
	"fmt"
	CT "github.com/fbaube/ctoken"
	"io"
)

type CommonCPR struct {
	NodeDepths []int
	FilePosns  []*CT.FilePosition
	CPR_raw    string
	// Writer is usually the GTokens Writer
	io.Writer
}

func NewCommonCPR() *CommonCPR {
	p := new(CommonCPR)
	p.NodeDepths = make([]int, 0)
	p.FilePosns = make([]*CT.FilePosition, 0)
	return p
}

func (p *CommonCPR) AsString(i int) string {
	var sND, sFP = "?", "?"

	if p.NodeDepths != nil && len(p.NodeDepths) > 0 {
		sND = fmt.Sprintf("%d", p.NodeDepths[i])
	}
	if p.FilePosns != nil && len(p.FilePosns) > 0 && p.FilePosns[i].Pos > 0 {
		sFP = fmt.Sprintf("%03d", p.FilePosns[i].Pos)
	}
	return fmt.Sprintf("i%02d,L%s,c%s", i, sND, sFP)
}
