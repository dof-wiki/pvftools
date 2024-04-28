package script

import (
	"github.com/dof-wiki/godof/parser"
	"pvftools/backend/common/consts"
	"pvftools/backend/common/log"
	"pvftools/backend/common/utils"
)

const (
	ActionSet    = "set"
	ActionAdd    = "add"
	ActionDel    = "del"
	ActionBreath = "breath"
)

type Action struct {
	Type string   `json:"type"`
	Args []string `json:"args"`
}

func (a *Action) Run(p *parser.Parser) {
	defer func() {
		if err := recover(); err != nil {
			log.LogError("error: %v", err)
		}
	}()
	switch a.Type {
	case ActionSet:
		if len(a.Args) < 2 {
			return
		}
		key := a.Args[0]
		isClosed := a.Args[1] == "true"
		values := ""
		if len(a.Args) > 2 {
			values = a.Args[2]
		}
		tokens := []*parser.Token{
			parser.NewRawContentToken(values),
		}
		p.SetAny(key, tokens, isClosed)
	case ActionAdd:
		if len(a.Args) < 2 {
			return
		}
		key := a.Args[0]
		isClosed := a.Args[1] == "true"
		values := ""
		if len(a.Args) > 2 {
			values = a.Args[2]
		}
		tokens := []*parser.Token{
			parser.NewRawContentToken(values),
		}
		p.AddAny(key, tokens, isClosed)
	case ActionDel:
		if len(a.Args) < 1 {
			return
		}
		key := a.Args[0]
		p.DelKey(key)
	case ActionBreath:
		tv := p.GetAny(consts.LabelEquipmentType)
		equTp := utils.RemoveWrapper(tv.GetCleanTokens()[0].RawContent())
		tokens := make([]*parser.Token, 0)
		for job := 0; job <= 10; job++ {
			for part := 1; part <= 4; part++ {
				tokens = append(tokens, parser.GenTokenList(job, part, equTp)...)
			}
		}
		p.SetAny(consts.LabelCharacterItemCheck, tokens, true)
	}
}
