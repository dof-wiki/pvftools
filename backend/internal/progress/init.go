package progress

import (
	"pvftools/backend/common/consts"
	"pvftools/backend/common/ctx"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ProgressController struct {
	name    string
	current int
	total   int
}

func NewProgressController(name string, total int) *ProgressController {
	runtime.EventsEmit(*ctx.Ctx, consts.EventProgressStart, name, total)
	return &ProgressController{
		current: 0,
		total:   total,
		name:    name,
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
	runtime.EventsEmit(*ctx.Ctx, consts.EventProgressUpdate, p.name, p.current)
}

func (p *ProgressController) End() {
	runtime.EventsEmit(*ctx.Ctx, consts.EventProgressEnd, p.name)
}

func (p *ProgressController) GetProgress() float64 {
	return float64(p.current) / float64(p.total)
}

func (p *ProgressController) IsEnd() bool {
	return p.current >= p.total
}

func (p *ProgressController) Write(buf []byte) (int, error) {
	n := len(buf)
	p.Increase(n)
	return n, nil
}
