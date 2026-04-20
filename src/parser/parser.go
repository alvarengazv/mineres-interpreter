package parser

import (
	"fmt"
	"mineres-interpreter/src/lexer"
	"mineres-interpreter/src/utils"
)

// Conjuntos de tokens para lookups O(1), substituindo cadeias de if/||
var typeTokens = map[lexer.TabelaPalavras]bool{
	lexer.Type_int:    true,
	lexer.Type_float:  true,
	lexer.Type_string: true,
	lexer.Type_bool:   true,
	lexer.Type_char:   true,
}

var relOpTokens = map[lexer.TabelaPalavras]bool{
	lexer.Op_eq:  true,
	lexer.Op_neq: true,
	lexer.Op_lt:  true,
	lexer.Op_gt:  true,
	lexer.Op_lte: true,
	lexer.Op_gte: true,
}

var mulOpTokens = map[lexer.TabelaPalavras]bool{
	lexer.Op_mul:     true,
	lexer.Op_div:     true,
	lexer.Op_mod:     true,
	lexer.Op_int_div: true,
}

var fatorTokens = map[lexer.TabelaPalavras]bool{
	lexer.Literal_string: true,
	lexer.Identifier:     true,
	lexer.Literal_int:    true,
	lexer.Literal_float:  true,
	lexer.Literal_char:   true,
	lexer.Literal_true:   true,
	lexer.Literal_false:  true,
	lexer.Literal_hex:    true,
	lexer.Literal_oct:    true,
}

var exprStartTokens = map[lexer.TabelaPalavras]bool{
	lexer.Identifier:     true,
	lexer.Literal_int:    true,
	lexer.Literal_float:  true,
	lexer.Literal_string: true,
	lexer.Literal_char:   true,
	lexer.Literal_true:   true,
	lexer.Literal_false:  true,
	lexer.Literal_hex:    true,
	lexer.Literal_oct:    true,
	lexer.Open_paren:     true,
	lexer.Op_add:         true,
	lexer.Op_sub:         true,
	lexer.Op_not:         true,
}

type Parser struct {
	tokens []lexer.Tupla
	pos    int
}

func NewParser(tokens []lexer.Tupla) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
	}
}

// Funções auxiliares para o manuseio do ponteiro do buffer

func (p *Parser) current() lexer.Tupla {

	if p.pos >= len(p.tokens) {
		linha, coluna := 1, 1
		if len(p.tokens) > 0 {
			ultimo := p.tokens[len(p.tokens)-1]
			linha = ultimo.Linha
			coluna = ultimo.Coluna + len(ultimo.Lexema)
		}
		return lexer.Tupla{Token: -1, Linha: linha, Coluna: coluna, Lexema: "EOF"}
	}

	return p.tokens[p.pos]
}

func (p *Parser) advance() {

	if p.pos < len(p.tokens) {
		p.pos++
	}

}

func (p *Parser) tokenToString(t lexer.TabelaPalavras) string {
	if t == -1 {
		return "EOF"
	}
	tokenF, _ := lexer.TabelaPalavrasFromInt(int(t))
	return tokenF.String()
}

func (p *Parser) consume(expected lexer.TabelaPalavras) {
	curr := p.current()

	if curr.Token == expected {
		p.advance()
	} else {
		expectedStr := p.tokenToString(expected)
		currentStr := p.tokenToString(curr.Token)

		utils.ThrowParserException(fmt.Sprintf("expected token '%v', got '%v'", expectedStr, currentStr), curr.Linha, curr.Coluna)
	}

}

// Funções de parsing

func (p *Parser) ParserFunction() {

	p.consume(lexer.Func_decl)
	p.consume(lexer.Main_function)
	p.consume(lexer.Open_paren)
	p.consume(lexer.Close_paren)
	p.parseBloco()

	if p.current().Token != -1 {
		utils.ThrowParserException("Unexpected token after main function: '"+ p.tokenToString(p.current().Token) + "'", p.current().Linha, p.current().Coluna)
	}

	fmt.Println("\nSyntactic analysis completed!")
}

func (p *Parser) parseBloco() {

	p.consume(lexer.Block_open)  // simbora
	p.parseStmtList()            // <stmt> <smtList> | &
	p.consume(lexer.Block_close) // cabo

}

func (p *Parser) parseStmtList() {

	for p.current().Token != lexer.Block_close && p.current().Token != -1 {
		p.parseStmt() // função principal que vai chamar as outras funções de stmt
	}

}

func (p *Parser) parseStmt() {

	token := p.current().Token

	switch token {

	case lexer.Loop_for:
		p.parseForStmt()

	case lexer.Conditional_if:
		p.parseIfStmt()

	case lexer.Conditional_switch:
		p.parseCaseStmt()

	case lexer.Loop_while:
		p.parseWhileStmt()

	case lexer.Block_open:
		p.parseBloco()

	case lexer.Io_scan, lexer.Io_print:
		p.parseIoStmt()

	case lexer.Loop_break:
		p.advance()
		p.consume(lexer.Stmt_end)

	case lexer.Loop_continue:
		p.advance()
		p.consume(lexer.Stmt_end)

	case lexer.Type_int, lexer.Type_float, lexer.Type_string, lexer.Type_bool, lexer.Type_char:
		p.parseDeclaration()

	case lexer.Stmt_end:
		p.advance() // uai

	// case lexer.Func_return:
	// 	p.parseReturnStmt()

	default:
		if exprStartTokens[token] {
			p.parseAtrib()
			p.consume(lexer.Stmt_end)
		} else {
			stringToken := p.tokenToString(token)
			utils.ThrowParserException(fmt.Sprintf("Expected statement, got '%v'", stringToken), p.current().Linha, p.current().Coluna)
		}
	}
}

// if
func (p *Parser) parseIfStmt() {

	p.consume(lexer.Conditional_if)
	p.consume(lexer.Open_paren)
	p.parseExpr()
	p.consume(lexer.Close_paren)
	p.parseStmt()

	// se for seguido de else, consome o else e o bloco do else
	if p.current().Token == lexer.Conditional_else {
		p.consume(lexer.Conditional_else)
		p.parseStmt()
	}
}

// while
func (p *Parser) parseWhileStmt() {

	p.consume(lexer.Loop_while)
	p.consume(lexer.Open_paren)
	p.parseExpr()
	p.consume(lexer.Close_paren)
	p.parseStmt()
}

// for
func (p *Parser) parseForStmt() {

	p.consume(lexer.Loop_for)
	p.consume(lexer.Open_paren)
	p.parseOptExpr()
	p.consume(lexer.Stmt_end_for)
	p.parseOptExpr()
	p.consume(lexer.Stmt_end_for)
	p.parseOptExpr()
	p.consume(lexer.Close_paren)
	p.parseStmt()

}

func (p *Parser) parseDeclaration() {

	p.parseType()
	p.parseDeclarationList()
	p.consume(lexer.Stmt_end)

}

func (p *Parser) parseDeclarationList() {

	p.consume(lexer.Identifier)
	for p.current().Token == lexer.Comma {
		p.advance()
		p.consume(lexer.Identifier)
	}

}

// IO
func (p *Parser) parseIoStmt() {

	if p.current().Token == lexer.Io_scan {
		p.advance()
		p.consume(lexer.Open_paren)
		p.parseType()
		p.consume(lexer.Comma)
		p.consume(lexer.Identifier)
		p.consume(lexer.Close_paren)
	} else {
		p.consume(lexer.Io_print)
		p.consume(lexer.Open_paren)
		p.parseOutputList()
		p.consume(lexer.Close_paren)
	}

	p.consume(lexer.Stmt_end)

}

func (p *Parser) parseOutputList() {

	p.parseFatorZin()
	for p.current().Token == lexer.Comma {
		p.advance()
		p.parseFatorZin()
	}

}

func (p *Parser) parseCaseStmt() {

	p.consume(lexer.Conditional_switch)
	p.consume(lexer.Open_paren)
	p.consume(lexer.Identifier)
	p.consume(lexer.Close_paren)
	p.consume(lexer.Block_open)
	for p.current().Token == lexer.Conditional_case {
		p.parseDoCaso()
	}
	if p.current().Token == lexer.Conditional_default {
		p.consume(lexer.Conditional_default)
		p.consume(lexer.Colon)
		p.parseStmt()
	}
	p.consume(lexer.Block_close)

}

func (p *Parser) parseDoCaso() {

	p.consume(lexer.Conditional_case)
	p.parseFatorZin()
	p.consume(lexer.Colon)
	p.parseStmt()

}

// func (p *Parser) parseReturnStmt() {

// 	p.consume(lexer.Func_return)
// 	if p.current().Token != lexer.Stmt_end {
// 		p.parseExpr()
// 	}
// 	p.consume(lexer.Stmt_end)

// }

func (p *Parser) parseType() {

	token := p.current().Token

	if typeTokens[token] {
		p.advance()
	} else {
		stringToken := p.tokenToString(token)
		utils.ThrowParserException(fmt.Sprintf("expected type, got '%v'", stringToken), p.current().Linha, p.current().Coluna)
	}

}

// Precedencia de operadores

func (p *Parser) parseExpr() {

	p.parseAtrib()

}

func (p *Parser) parseOptExpr() {
	if exprStartTokens[p.current().Token] {
		p.parseAtrib()
	}
}

func (p *Parser) parseAtrib() {

	p.parseOR()
	if p.current().Token == lexer.Op_assign {
		p.advance()
		p.parseAtrib() // recursão a direita para permitir varias atribuições
	}
}

// <or> -> <xor> { 'quarque_um' <xor> }
func (p *Parser) parseOR() {

	p.parseXor()
	for p.current().Token == lexer.Op_or {
		p.advance()
		p.parseXor()
	}

}

// <xor> -> <and> { 'um_o_oto' <and> }
func (p *Parser) parseXor() {

	p.parseAnd()
	for p.current().Token == lexer.Op_xor {
		p.advance()
		p.parseAnd()
	}
}

// <and> -> <not> { 'tamem' <not> }
func (p *Parser) parseAnd() {

	p.parseNot()
	for p.current().Token == lexer.Op_and {
		p.advance()
		p.parseNot()
	}
}

// <not> -> 'vam_marca' <not> | <rel>
func (p *Parser) parseNot() {

	if p.current().Token == lexer.Op_not {
		p.advance()
		p.parseNot()
	} else {
		p.parseRel()
	}

}

// <rel> -> <add> { ('mema_coisa' | 'neh_nada') <add> }
func (p *Parser) parseRel() {

	p.parseAdd()
	if relOpTokens[p.current().Token] {
		p.advance()
		p.parseAdd()
	}

}

// <add> -> <mul> { ('veiz' | 'sob') <mul> }
func (p *Parser) parseAdd() {

	p.parseMul()
	for p.current().Token == lexer.Op_add || p.current().Token == lexer.Op_sub {
		p.advance()
		p.parseMul()
	}

}

// <mul> -> <fator> { ('veiz' | 'sob') <fator> }
func (p *Parser) parseMul() {
	p.parseUno()
	for mulOpTokens[p.current().Token] {
		p.advance()
		p.parseUno()
	}

}

// <uno> -> '+' <uno> | '-' <uno> | <fatorZao>
func (p *Parser) parseUno() {

	if p.current().Token == lexer.Op_add || p.current().Token == lexer.Op_sub {
		p.advance()
		p.parseUno()
	} else {
		p.parseFatorZao()
	}
}

// <fatorZao> -> <fatorzin> | '(' <atrib> ')'
func (p *Parser) parseFatorZao() {

	if p.current().Token == lexer.Open_paren {
		p.advance()
		p.parseAtrib()
		p.consume(lexer.Close_paren)
	} else {
		p.parseFatorZin()
	}

}

// <fatorZin> -> 'STR' | 'IDENT' | 'NUMint' |...
func (p *Parser) parseFatorZin() {

	t := p.current().Token
	if fatorTokens[t] {
		p.advance()
	} else {
		stringToken := p.tokenToString(t)
		utils.ThrowParserException(fmt.Sprintf("expected 'STR' or 'IDENT' or 'NUMint' or 'NUMfloat' or 'NUMhex' or 'NUMoct' or 'valorBooleano' or 'valorChar', got '%v'", stringToken), p.current().Linha, p.current().Coluna)
	}
}