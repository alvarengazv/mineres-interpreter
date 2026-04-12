package parser

import (
	"fmt"
	"mineres-interpreter/src/lexer"
	"mineres-interpreter/src/utils"
)

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
		utils.ThrowException("parser.go", "current", "unexpected end of file")
	}

	return p.tokens[p.pos]
}

func (p *Parser) advance() {

	if p.pos < len(p.tokens) {
		p.pos++
	}

}

func (p *Parser) consume(expected lexer.TabelaPalavras) {

	if p.current().Token == expected {
		p.advance()
	} else {
		// Get the expected and current tokens as strings for better error messages
		token, _ := lexer.TabelaPalavrasFromInt(int(expected))
		expectedStr := token.String()

		currentToken, _ := lexer.TabelaPalavrasFromInt(int(p.current().Token))
		currentStr := currentToken.String()
		utils.ThrowParserException(fmt.Sprintf("expected token '%v', got '%v'", expectedStr, currentStr), p.current().Linha, p.current().Coluna)
	}

}

// Funções de parsing

func (p *Parser) ParserFunction() {

	p.consume(lexer.Func_decl)
	p.consume(lexer.Main_function)
	p.consume(lexer.Open_paren)
	p.consume(lexer.Close_paren)
	p.parseBloco()

	fmt.Println("Syntactic analysis completed!")
}

func (p *Parser) parseBloco() {

	p.consume(lexer.Block_open)  // simbora
	p.parseStmtList()            // <stmt> <smtList> | &
	p.consume(lexer.Block_close) // cabo

}

func (p *Parser) parseStmtList() {

	for p.current().Token != lexer.Block_close {
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
		if p.isStartOfExpr(token) {
			p.parseAtrib()
			p.consume(lexer.Stmt_end)
		} else {
			tokenF, _ := lexer.TabelaPalavrasFromInt(int(token))
			stringToken := tokenF.String()
			utils.ThrowParserException(fmt.Sprintf("unexpected token '%v'", stringToken), p.current().Linha, p.current().Coluna)
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
	p.parseAtrib()
	p.consume(lexer.Stmt_end_for)
	p.parseExpr()
	p.consume(lexer.Stmt_end_for)
	p.parseAtrib()
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

	if token == lexer.Type_int || token == lexer.Type_float || token == lexer.Type_string || token == lexer.Type_bool || token == lexer.Type_char {
		p.advance()
	} else {
		tokenF, _ := lexer.TabelaPalavrasFromInt(int(token))
		stringToken := tokenF.String()
		utils.ThrowParserException(fmt.Sprintf("expected type, got '%v'", stringToken), p.current().Linha, p.current().Coluna)
	}

}

// Precedencia de operadores

func (p *Parser) parseExpr() {

	p.parseAtrib()

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
	if p.current().Token == lexer.Op_eq || p.current().Token == lexer.Op_lt || p.current().Token == lexer.Op_gt || p.current().Token == lexer.Op_lte || p.current().Token == lexer.Op_gte || p.current().Token == lexer.Op_neq {
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
	t := p.current().Token
	for t == lexer.Op_mul || t == lexer.Op_div || t == lexer.Op_mod || t == lexer.Op_int_div {
		p.advance()
		p.parseUno()
		t = p.current().Token
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
	if t == lexer.Literal_string || t == lexer.Identifier || 
	   t == lexer.Literal_int || t == lexer.Literal_float || 
	   t == lexer.Literal_char || t == lexer.Literal_true || 
	   t == lexer.Literal_false || t == lexer.Literal_hex ||
	   t == lexer.Literal_oct {
		p.advance()
	} else {
		tokenF, _ := lexer.TabelaPalavrasFromInt(int(t))
		stringToken := tokenF.String()
		utils.ThrowParserException(fmt.Sprintf("expected 'STR' or 'IDENT' or 'NUMint' or 'NUMfloat' or 'NUMhex' or 'NUMoct' or 'valorBooleano' or 'valorChar', got '%v'", stringToken), p.current().Linha, p.current().Coluna)
	}
}

func (p *Parser) isStartOfExpr(t lexer.TabelaPalavras) bool {

	return t == lexer.Identifier || t == lexer.Literal_int || t == lexer.Literal_float || t == lexer.Literal_string ||
		t == lexer.Literal_char || t == lexer.Literal_true || t == lexer.Literal_false || t == lexer.Open_paren ||
		t == lexer.Op_add || t == lexer.Op_sub || t == lexer.Op_not || t == lexer.Literal_hex || t == lexer.Literal_oct

}
