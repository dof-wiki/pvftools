package progress

import (
	"pvftools/backend/common/consts"
	"pvftools/backend/common/ctx"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ProgressController struct {
	current int
	total   int
}

func NewProgressController(total int) *ProgressController {
	runtime.EventsEmit(*ctx.Ctx, consts.EventProgressStart, total)
	return &ProgressController{
		current: 0,
		total:   total,
	}
}

func (p *ProgressController) Increase(c int) {
	if p.current >= p.total {
		return
	}
	p.current += c
	if p.current >= p.total {
		p.current = p.total
	}
	runtime.EventsEmit(*ctx.Ctx, consts.EventProgressUpdate, p.current)
}

func (p *ProgressController) GetProgress() float64 {
	return float64(p.current) / float64(p.total)
}

func (p *ProgressController) IsEnd() bool {
	return p.current >= p.total
}
