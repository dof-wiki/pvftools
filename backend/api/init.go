package api

import (
	"context"
	"github.com/dof-wiki/godof/parser"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"pvftools/backend/common/ctx"
	"pvftools/backend/internal/data_source"
	"sync"
)

// App struct
type App struct {
	ctx context.Context

	cacheMu     sync.Mutex
	parserCache map[string]*parser.Parser
	lstParser   map[string]*parser.LstParser
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		parserCache: make(map[string]*parser.Parser),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(_ctx context.Context) {
	a.ctx = _ctx
	ctx.Ctx = &_ctx
}

func (a *App) getParser(path string) *parser.Parser {
	a.cacheMu.Lock()
	defer a.cacheMu.Unlock()
	if _, ok := a.parserCache[path]; !ok {
		c, err := data_source.GetDataSource().GetFileContent(path)
		if err != nil {
			runtime.LogError(a.ctx, err.Error())
			return parser.NewParser("")
		}
		a.parserCache[path] = parser.NewParser(c)
	}
	return a.parserCache[path]
}

func (a *App) getLstParser(path string) *parser.LstParser {
	a.cacheMu.Lock()
	defer a.cacheMu.Unlock()
	if _, ok := a.lstParser[path]; !ok {
		c, err := data_source.GetDataSource().GetFileContent(path)
		if err != nil {
			runtime.LogError(a.ctx, err.Error())
			return parser.NewLstParser("")
		}
		a.lstParser[path] = parser.NewLstParser(c)
	}
	return a.lstParser[path]
}

func (a *App) saveData(path string) error {
	a.cacheMu.Lock()
	defer a.cacheMu.Unlock()
	if _, ok := a.parserCache[path]; !ok {
		return nil
	}
	return data_source.GetDataSource().SaveFileContent(path, a.getParser(path).Render())
}
