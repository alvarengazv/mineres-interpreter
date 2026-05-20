package parser

import (
	"fmt"
	"mineres-interpreter/src/lexer"
)

func (p *Parser) newLabelTrue() *lexer.TuplaLex {

	label := &lexer.TuplaLex{
		Token:  lexer.Identifier,
		Lexema: fmt.Sprintf("V%d", p.labelTrue),
	}

	p.labelTrue++
	return label
}

func (p *Parser) newLabelFalse() *lexer.TuplaLex {

	label := &lexer.TuplaLex{
		Token:  lexer.Identifier,
		Lexema: fmt.Sprintf("F%d", p.labelFalse),
	}

	p.labelFalse++
	return label
}

func (p *Parser) newLabelEndIf() *lexer.TuplaLex {

	label := &lexer.TuplaLex{
		Token:  lexer.Identifier,
		Lexema: fmt.Sprintf("endIf%d", p.labelEndIf),
	}

	p.labelEndIf++
	return label
}

func (p *Parser) newLabelLoopInit() *lexer.TuplaLex {

	label := &lexer.TuplaLex{
		Token:  lexer.Identifier,
		Lexema: fmt.Sprintf("loopInit%d", p.labelLoopInit),
	}

	p.labelLoopInit++
	return label
}

func (p *Parser) newLabelForInc() *lexer.TuplaLex {

	label := &lexer.TuplaLex{
		Token:  lexer.Identifier,
		Lexema: fmt.Sprintf("forInc%d", p.labelForInc),
	}

	p.labelForInc++
	return label
}

func (p *Parser) newLabelLoopEnd() *lexer.TuplaLex {

	label := &lexer.TuplaLex{
		Token:  lexer.Identifier,
		Lexema: fmt.Sprintf("loopEnd%d", p.labelLoopEnd),
	}

	p.labelLoopEnd++
	return label
}

func (p *Parser) newLabelCaseIfT() *lexer.TuplaLex {

	label := &lexer.TuplaLex{
		Token:  lexer.Identifier,
		Lexema: fmt.Sprintf("caseIfT%d", p.labelCaseIfT),
	}

	p.labelCaseIfT++
	return label
}

func (p *Parser) newLabelCaseIfF() *lexer.TuplaLex {

	label := &lexer.TuplaLex{
		Token:  lexer.Identifier,
		Lexema: fmt.Sprintf("caseIfF%d", p.labelCaseIfF),
	}

	p.labelCaseIfF++
	return label
}

func (p *Parser) newLabelEndCase() *lexer.TuplaLex {

	label := &lexer.TuplaLex{
		Token:  lexer.Identifier,
		Lexema: fmt.Sprintf("endCase%d", p.labelEndCase),
	}

	p.labelEndCase++
	return label
}
