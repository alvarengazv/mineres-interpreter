package parser

import (
	"mineres-interpreter/src/lexer"
	"strconv"
)

func (p *Parser) newLabel(prefix string) *lexer.TuplaLex {
	label := &lexer.TuplaLex{
		Token:  lexer.Identifier,
		Lexema: prefix + strconv.Itoa(p.labelCount),
	}
	p.labelCount++
	return label
}

func (p *Parser) newLabelTrue() *lexer.TuplaLex {
	return p.newLabel("V")
}

func (p *Parser) newLabelFalse() *lexer.TuplaLex {
	return p.newLabel("F")
}

func (p *Parser) newLabelEndIf() *lexer.TuplaLex {
	return p.newLabel("endIf")
}

func (p *Parser) newLabelLoopInit() *lexer.TuplaLex {
	return p.newLabel("loopInit")
}

func (p *Parser) newLabelForInc() *lexer.TuplaLex {
	return p.newLabel("forInc")
}

func (p *Parser) newLabelLoopEnd() *lexer.TuplaLex {
	return p.newLabel("loopEnd")
}

func (p *Parser) newLabelCaseIfT() *lexer.TuplaLex {
	return p.newLabel("caseIfT")
}

func (p *Parser) newLabelCaseIfF() *lexer.TuplaLex {
	return p.newLabel("caseIfF")
}

func (p *Parser) newLabelEndCase() *lexer.TuplaLex {
	return p.newLabel("endCase")
}
