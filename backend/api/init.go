package api

import (
	"context"
	"github.com/dof-wiki/godof/parser"
	"github.com/dof-wiki/godof/parser/tree_parser"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"pvftools/backend/common/consts"
	"pvftools/backend/common/ctx"
	"pvftools/backend/internal/data_source"
	"pvftools/backend/internal/updater"
	"sync"
)

// App struct
type App struct {
	ctx context.Context

	cacheMu     sync.RWMutex
	parserCache map[string]*parser.Parser
	lstParser   map[string]*parser.LstParser
	treeParser  map[string]*tree_parser.TreeParser

	updater *updater.Updater
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		parserCache: make(map[string]*parser.Parser),
		lstParser:   make(map[string]*parser.LstParser),
		treeParser:  make(map[string]*tree_parser.TreeParser),
		updater:     updater.NewUpdater(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(_ctx context.Context) {
	a.ctx = _ctx
	ctx.Ctx = &_ctx
	go a.checkUpdate()
}

func (a *App) checkUpdate() bool {
	if !a.updater.CheckUpdate() {
		return false
	}
	newVer, msg := a.updater.GetUpdateInfo()
	runtime.EventsEmit(a.ctx, consts.EventHasSelfUpdate, newVer, msg)
	return true
}

func (a *App) getParserInner(path string) *parser.Parser {
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

func (a *App) getParser(path string) *parser.Parser {
	a.cacheMu.RLock()
	defer a.cacheMu.RUnlock()
	return a.getParserInner(path)
}

func (a *App) getLstParser(path string) *parser.LstParser {
	a.cacheMu.RLock()
	defer a.cacheMu.RUnlock()
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
	return data_source.GetDataSource().SaveFileContent(path, a.getParserInner(path).Render())
}

func (a *App) getTreeParserInner(path string) *tree_parser.TreeParser {
	if _, ok := a.treeParser[path]; !ok {
		c, err := data_source.GetDataSource().GetFileContent(path)
		if err != nil {
			runtime.LogError(a.ctx, err.Error())
			return tree_parser.NewTreeParser("")
		}
		a.treeParser[path] = tree_parser.NewTreeParser(c)
	}
	return a.treeParser[path]
}

func (a *App) getTreeParser(path string) *tree_parser.TreeParser {
	a.cacheMu.RLock()
	defer a.cacheMu.RUnlock()
	return a.getTreeParserInner(path)
}

func (a *App) delTreeParser(path string) {
	a.cacheMu.RLock()
	defer a.cacheMu.RUnlock()
	if _, ok := a.treeParser[path]; ok {
		delete(a.treeParser, path)
	}
}
