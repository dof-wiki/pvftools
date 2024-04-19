package script

import "github.com/dof-wiki/godof/parser"

const (
	ActionSet = "set"
	ActionAdd = "add"
	ActionDel = "del"
)

type Action struct {
	Type string `json:"type"`
	Args []any  `json:"args"`
}

func (a *Action) Run(p *parser.Parser) {
	switch a.Type {
	case ActionSet:
		if len(a.Args) < 2 {
			return
		}
		key := a.Args[0].(string)
		isClosed := a.Args[1].(bool)
		values := a.Args[2:]
		tokens := parser.GenTokenList(values...)
		p.SetAny(key, tokens, isClosed)
	case ActionAdd:
		if len(a.Args) < 2 {
			return
		}
		key := a.Args[0].(string)
		isClosed := a.Args[1].(bool)
		values := a.Args[2:]
		tokens := parser.GenTokenList(values...)
		p.SetAny(key, tokens, isClosed)
	case ActionDel:
		if len(a.Args) < 1 {
			return
		}
		key := a.Args[0].(string)
		p.SetAny(key, nil)
	}
}
