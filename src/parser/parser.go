package parser

import (
	"fmt"
	"mineres-interpreter/src/lexer"
	"mineres-interpreter/src/utils"
)

// Conjuntos de tokens para lookups O(1), substituindo cadeias de if
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

type TypeTable int

const (
	Type_int TypeTable = iota
	Type_float
	Type_string
	Type_bool
	Type_char
)

// conversao para mensagem de erro
func (t TypeTable) String() string {
	switch t {
	case Type_int:
		return "trem_di_numeru"
	case Type_float:
		return "trem_cum_virgula"
	case Type_string:
		return "trem_discrita"
	case Type_bool:
		return "trem_discolhe"
	case Type_char:
		return "trosso"
	default:
		return "unknown"
	}
}

type Parser struct {
	tokens []lexer.TuplaLex
	pos    int

	microcodes    []TuplaMicrocode
	tempCount     int
	labelTrue     int
	labelFalse    int
	labelEndIf    int
	labelLoopInit int
	labelForInc   int
	labelLoopEnd  int
	labelCaseIfT  int
	labelCaseIfF  int
	labelEndCase  int
	symbolTable   map[string]TypeTable

	breakStack    Stack
	continueStack Stack
}

func verifyTypeCompatibility(type1 TypeTable, type2 TypeTable, line int, column int) {

	if type1 != type2 && !(type1 == Type_float && type2 == Type_int) && !(type1 == Type_int && type2 == Type_float) {
		utils.ThrowParserException(fmt.Sprintf("Type mismatch: cannot operate between %v and %v", type1, type2), line, column)
	}

}

func (p *Parser) toType(t *lexer.TuplaLex) TypeTable {

	switch t.Token {
	case lexer.Type_int, lexer.Literal_int, lexer.Literal_hex, lexer.Literal_oct:
		return Type_int
	case lexer.Type_float, lexer.Literal_float:
		return Type_float
	case lexer.Type_string, lexer.Literal_string:
		return Type_string
	case lexer.Type_bool, lexer.Literal_true, lexer.Literal_false:
		return Type_bool
	case lexer.Type_char, lexer.Literal_char:
		return Type_char

	case lexer.Identifier:
		value, ok := p.symbolTable[t.Lexema]
		if !ok {
			utils.ThrowParserException(fmt.Sprintf("Variable '%s' not declared", t.Lexema), t.Linha, t.Coluna)
		}
		return value
	}
	return -1
}

func (p *Parser) inferType(op lexer.TabelaPalavras, t1, t2 TypeTable, linha, coluna int) TypeTable {
	switch op {
	case lexer.Op_add:
		if (t1 == Type_char && t2 == Type_char) ||
			(t1 == Type_char && t2 == Type_string) ||
			(t1 == Type_string && t2 == Type_char) ||
			(t1 == Type_string && t2 == Type_string) {
			return Type_string
		}

		if t1 == Type_int && t2 == Type_int {
			return Type_int
		}

		if (t1 == Type_float || t1 == Type_int) && (t2 == Type_float || t2 == Type_int) {
			return Type_float
		}

	case lexer.Op_sub, lexer.Op_mul, lexer.Op_div, lexer.Op_int_div, lexer.Op_mod:
		if t1 == Type_string || t2 == Type_string || t1 == Type_char || t2 == Type_char {
			utils.ThrowParserException("Arithmetic operation not allowed on strings or characters", linha, coluna)
		}

		if op == lexer.Op_mod || op == lexer.Op_int_div {
			if t1 != Type_int || t2 != Type_int {
				utils.ThrowParserException("Operation expected integer values", linha, coluna)
			}
			return Type_int
		}

		if t1 == Type_int && t2 == Type_int {
			return Type_int
		}
		return Type_float

	case lexer.Op_eq, lexer.Op_neq:
		if (t1 == Type_int || t1 == Type_float) && (t2 == Type_int || t2 == Type_float) {
			return Type_bool
		}
		if t1 == t2 {
			return Type_bool
		}
		utils.ThrowParserException(fmt.Sprintf("Equality operation not allowed between %v and %v", t1, t2), linha, coluna)

	case lexer.Op_gt, lexer.Op_lt, lexer.Op_gte, lexer.Op_lte:
		if (t1 == Type_int || t1 == Type_float) && (t2 == Type_int || t2 == Type_float) {
			return Type_bool
		}
		if t1 == Type_string && t2 == Type_string {
			return Type_bool
		}
		utils.ThrowParserException(fmt.Sprintf("Relational operation not allowed between %v and %v", t1, t2), linha, coluna)

	case lexer.Op_and, lexer.Op_or, lexer.Op_xor:
		if t1 == Type_bool && t2 == Type_bool {
			return Type_bool
		}
		utils.ThrowParserException(fmt.Sprintf("Logical operation not allowed between %v and %v", t1, t2), linha, coluna)
	}

	utils.ThrowParserException("Operação com tipos incompatíveis", linha, coluna)
	return -1
}

func NewParser(tokens []lexer.TuplaLex) *Parser {
	return &Parser{
		tokens:      tokens,
		pos:         0,
		symbolTable: make(map[string]TypeTable),
	}
}

func (p *Parser) newTemp(typeVal TypeTable) *lexer.TuplaLex {

	p.tempCount++
	nomeTemp := fmt.Sprintf("$t%d", p.tempCount)
	p.symbolTable[nomeTemp] = typeVal

	return &lexer.TuplaLex{
		Token:  lexer.Identifier,
		Lexema: fmt.Sprintf("$t%d", p.tempCount),
		Linha:  0,
		Coluna: 0,
	}
}

// Funções auxiliares para o manuseio do ponteiro do buffer

func (p *Parser) current() lexer.TuplaLex {

	if p.pos >= len(p.tokens) {
		linha, coluna := 1, 1
		if len(p.tokens) > 0 {
			ultimo := p.tokens[len(p.tokens)-1]
			linha = ultimo.Linha
			coluna = ultimo.Coluna + len(ultimo.Lexema)
		}
		return lexer.TuplaLex{Token: -1, Linha: linha, Coluna: coluna, Lexema: "EOF"}
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

func (p *Parser) consume(expected lexer.TabelaPalavras) *lexer.TuplaLex {
	curr := p.current()

	if curr.Token == expected {
		p.advance()
		return &curr
	} else {
		expectedStr := p.tokenToString(expected)
		currentStr := p.tokenToString(curr.Token)

		utils.ThrowParserException(fmt.Sprintf("expected token '%v', got '%v'", expectedStr, currentStr), curr.Linha, curr.Coluna)

		return nil
	}

}

// Funções de parsing

func (p *Parser) ParserFunction() []TuplaMicrocode {

	p.consume(lexer.Func_decl)
	p.consume(lexer.Main_function)
	p.consume(lexer.Open_paren)
	p.consume(lexer.Close_paren)
	p.parseBloco()

	if p.current().Token != -1 {
		utils.ThrowParserException("Unexpected token after main function: '"+p.tokenToString(p.current().Token)+"'", p.current().Linha, p.current().Coluna)
	}

	// *
	// ListTuplaMicrocodeToString(p.microcodes)

	return p.microcodes
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
		p.parseBreak()

	case lexer.Loop_continue:
		p.parseContinue()

	case lexer.Type_int, lexer.Type_float, lexer.Type_string, lexer.Type_bool, lexer.Type_char:
		p.parseDeclaration()

	case lexer.Stmt_end:
		p.advance() // uai

	// case lexer.Func_return:
	// 	p.parseReturnStmt()

	default:
		if exprStartTokens[token] {
			_, code := p.parseAtrib()
			// ListTuplaMicrocodeToString(code)
			p.microcodes = append(p.microcodes, code...)
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
	resExpr, codeExpr := p.parseExpr()
	if p.toType(resExpr) != Type_bool {
		utils.ThrowParserException("If condition must be a boolean expression", resExpr.Linha, resExpr.Coluna)
	}
	p.microcodes = append(p.microcodes, codeExpr...)
	p.consume(lexer.Close_paren)

	labelTrue := p.newLabelTrue()
	labelFalse := p.newLabelFalse()
	labelEndIf := p.newLabelEndIf()

	codeIf := TuplaMicrocode{
		Operation: If_eq,
		Res:       resExpr,
		Op1:       labelTrue,
		Op2:       labelFalse,
	}

	p.microcodes = append(p.microcodes, codeIf)

	p.microcodes = append(p.microcodes, TuplaMicrocode{
		Operation: Label,
		Res:       labelTrue,
		Op1:       nil,
		Op2:       nil,
	})

	p.parseStmt()

	p.microcodes = append(p.microcodes, TuplaMicrocode{
		Operation: Jump,
		Res:       labelEndIf,
		Op1:       nil,
		Op2:       nil,
	})

	p.microcodes = append(p.microcodes, TuplaMicrocode{
		Operation: Label,
		Res:       labelFalse,
		Op1:       nil,
		Op2:       nil,
	})

	// se for seguido de else, consome o else e o bloco do else
	if p.current().Token == lexer.Conditional_else {
		p.consume(lexer.Conditional_else)
		p.parseStmt()
	}

	p.microcodes = append(p.microcodes, TuplaMicrocode{
		Operation: Label,
		Res:       labelEndIf,
		Op1:       nil,
		Op2:       nil,
	})

	// ListTuplaMicrocodeToString(p.microcodes)
}

// while
func (p *Parser) parseWhileStmt() {

	p.consume(lexer.Loop_while)
	p.consume(lexer.Open_paren)
	resExpr, codeExpr := p.parseExpr()

	labelLoopInit := p.newLabelLoopInit()
	labelTrue := p.newLabelTrue()
	labelEndLoop := p.newLabelLoopEnd()

	codeLabelInit := TuplaMicrocode{
		Operation: Label,
		Res:       labelLoopInit,
		Op1:       nil,
		Op2:       nil,
	}

	p.microcodes = append(p.microcodes, codeLabelInit)
	// ListTuplaMicrocodeToString(code)
	p.microcodes = append(p.microcodes, codeExpr...)

	p.breakStack.Push(labelEndLoop.Lexema)
	p.continueStack.Push(labelLoopInit.Lexema)

	codeIf := TuplaMicrocode{
		Operation: If_eq,
		Res:       resExpr,
		Op1:       labelTrue,
		Op2:       labelEndLoop,
	}

	p.microcodes = append(p.microcodes, codeIf)

	p.microcodes = append(p.microcodes, TuplaMicrocode{
		Operation: Label,
		Res:       labelTrue,
		Op1:       nil,
		Op2:       nil,
	})

	p.consume(lexer.Close_paren)
	p.parseStmt()

	p.microcodes = append(p.microcodes, TuplaMicrocode{
		Operation: Jump,
		Res:       labelLoopInit,
		Op1:       nil,
		Op2:       nil,
	})

	p.microcodes = append(p.microcodes, TuplaMicrocode{
		Operation: Label,
		Res:       labelEndLoop,
		Op1:       nil,
		Op2:       nil,
	})

	p.breakStack.Pop()
	p.continueStack.Pop()
}

// for
func (p *Parser) parseForStmt() {

	p.consume(lexer.Loop_for)
	p.consume(lexer.Open_paren)

	_, code1 := p.parseOptExpr()
	if code1 != nil {
		// ListTuplaMicrocodeToString(code1)
		p.microcodes = append(p.microcodes, code1...)
	}

	p.consume(lexer.Stmt_end_for)

	labelLoopInit := p.newLabelLoopInit()
	labelTrue := p.newLabelTrue()
	labelEndLoop := p.newLabelLoopEnd()
	labelForInc := p.newLabelForInc()

	p.breakStack.Push(labelEndLoop.Lexema)
	p.continueStack.Push(labelForInc.Lexema)

	codeLabelInit := TuplaMicrocode{
		Operation: Label,
		Res:       labelLoopInit,
		Op1:       nil,
		Op2:       nil,
	}

	p.microcodes = append(p.microcodes, codeLabelInit)

	resExpr2, codeExpr2 := p.parseOptExpr()
	if codeExpr2 != nil {
		// ListTuplaMicrocodeToString(code2)
		p.microcodes = append(p.microcodes, codeExpr2...)
	}

	p.consume(lexer.Stmt_end_for)
	_, codeExpr3 := p.parseOptExpr()

	codeIf := TuplaMicrocode{
		Operation: If_eq,
		Res:       resExpr2,
		Op1:       labelTrue,
		Op2:       labelEndLoop,
	}

	p.microcodes = append(p.microcodes, codeIf)

	p.microcodes = append(p.microcodes, TuplaMicrocode{
		Operation: Label,
		Res:       labelTrue,
		Op1:       nil,
		Op2:       nil,
	})

	p.consume(lexer.Close_paren)
	p.parseStmt()

	p.microcodes = append(p.microcodes, TuplaMicrocode{
		Operation: Label,
		Res:       labelForInc,
		Op1:       nil,
		Op2:       nil,
	})

	if codeExpr3 != nil {
		// ListTuplaMicrocodeToString(code3)
		p.microcodes = append(p.microcodes, codeExpr3...)
	}

	p.microcodes = append(p.microcodes, TuplaMicrocode{
		Operation: Jump,
		Res:       labelLoopInit,
		Op1:       nil,
		Op2:       nil,
	})

	p.microcodes = append(p.microcodes, TuplaMicrocode{
		Operation: Label,
		Res:       labelEndLoop,
		Op1:       nil,
		Op2:       nil,
	})

	p.breakStack.Pop()
	p.continueStack.Pop()
}

func (p *Parser) parseDeclaration() []TuplaMicrocode {

	typeToken := p.parseType()

	commandList := p.parseDeclarationList(typeToken)
	p.consume(lexer.Stmt_end)

	// ListTuplaMicrocodeToString(commandList)
	p.microcodes = append(p.microcodes, commandList...)

	return commandList
}

func (p *Parser) parseDeclarationList(typeToken lexer.TabelaPalavras) []TuplaMicrocode {
	commandList := []TuplaMicrocode{}

	currentType := p.toType(&lexer.TuplaLex{Token: typeToken})
	currentToken := p.consume(lexer.Identifier)
	value, ok := p.symbolTable[currentToken.Lexema]

	if ok {
		utils.ThrowParserException(fmt.Sprintf("Variable '%s' already declared with type %v", currentToken.Lexema, value), currentToken.Linha, currentToken.Coluna)
	}

	p.symbolTable[currentToken.Lexema] = currentType

	commandList = append(commandList, TuplaMicrocode{
		Operation: Att,
		Res:       currentToken,
		Op1:       getDefaultValue(typeToken),
		Op2:       nil,
	})

	for p.current().Token == lexer.Comma {
		p.advance()
		currentToken := p.consume(lexer.Identifier)
		value, ok := p.symbolTable[currentToken.Lexema]

		if ok {
			utils.ThrowParserException(fmt.Sprintf("Variable '%s' already declared with type %v", currentToken.Lexema, value), currentToken.Linha, currentToken.Coluna)
		}

		p.symbolTable[currentToken.Lexema] = currentType

		commandList = append(commandList, TuplaMicrocode{
			Operation: Att,
			Res:       currentToken,
			Op1:       getDefaultValue(typeToken),
			Op2:       nil,
		})
	}

	return commandList
}

// IO
func (p *Parser) parseIoStmt() []TuplaMicrocode {
	commandList := []TuplaMicrocode{}

	if p.current().Token == lexer.Io_scan {
		p.advance()
		p.consume(lexer.Open_paren)
		typeToken := p.parseType()
		p.consume(lexer.Comma)
		ident := p.consume(lexer.Identifier)
		p.consume(lexer.Close_paren)

		commandList = append(commandList, TuplaMicrocode{
			Operation: Call,
			Res: &lexer.TuplaLex{
				Lexema: "read",
			},
			Op1: ident,
			Op2: &lexer.TuplaLex{
				Token:  typeToken,
				Lexema: typeToken.String(),
			},
		})
	} else {
		p.consume(lexer.Io_print)
		p.consume(lexer.Open_paren)

		// Agora pega tudo
		args := p.parseOutputList()

		p.consume(lexer.Close_paren)

		// Depois de pega tudo faz um for para criar tupla de print para cada argumento
		for _, arg := range args {
			microcode := TuplaMicrocode{
				Operation: Call,
				Res: &lexer.TuplaLex{
					Lexema: "print",
				},
			}

			if arg.Token == lexer.Identifier {
				microcode.Op1 = arg
				microcode.Op2 = nil
			} else {
				microcode.Op1 = nil
				microcode.Op2 = arg
			}

			commandList = append(commandList, microcode)
		}
	}

	p.consume(lexer.Stmt_end)
	// ListTuplaMicrocodeToString(commandList)
	p.microcodes = append(p.microcodes, commandList...)
	return commandList
}

func (p *Parser) parseOutputList() []*lexer.TuplaLex {
	args := []*lexer.TuplaLex{}

	tupla, _ := p.parseFatorZin()
	if tupla != nil {
		args = append(args, tupla)
	}
	for p.current().Token == lexer.Comma {
		p.advance()
		tupla, _ := p.parseFatorZin()
		args = append(args, tupla)
	}

	return args
}

func (p *Parser) parseCaseStmt() {

	p.consume(lexer.Conditional_switch)
	p.consume(lexer.Open_paren)

	ident := p.consume(lexer.Identifier)

	p.consume(lexer.Close_paren)
	p.consume(lexer.Block_open)
	labelEndCase := p.newLabelEndCase()
	p.breakStack.Push(labelEndCase.Lexema)

	for p.current().Token == lexer.Conditional_case {

		p.parseDoCaso(ident, labelEndCase)
	}
	if p.current().Token == lexer.Conditional_default {
		p.consume(lexer.Conditional_default)
		p.consume(lexer.Colon)
		p.parseStmt()
	}
	p.consume(lexer.Block_close)
	p.breakStack.Pop()

	p.microcodes = append(p.microcodes, TuplaMicrocode{
		Operation: Label,
		Res:       labelEndCase,
		Op1:       nil,
		Op2:       nil,
	})

}

func (p *Parser) parseDoCaso(ident *lexer.TuplaLex, labelEndCase *lexer.TuplaLex) {

	p.consume(lexer.Conditional_case)

	labelTrue := p.newLabelCaseIfT()
	labelFalse := p.newLabelCaseIfF()

	valueCase, _ := p.parseFatorZin()

	temp := p.newTemp(Type_bool)

	p.microcodes = append(p.microcodes, TuplaMicrocode{
		Operation: Eq,
		Res:       temp,
		Op1:       ident,
		Op2:       valueCase,
	})

	p.microcodes = append(p.microcodes, TuplaMicrocode{
		Operation: If_eq,
		Res:       temp,
		Op1:       labelTrue,
		Op2:       labelFalse,
	})

	p.microcodes = append(p.microcodes, TuplaMicrocode{
		Operation: Label,
		Res:       labelTrue,
		Op1:       nil,
		Op2:       nil,
	})

	p.consume(lexer.Colon)
	p.parseStmt()

	p.microcodes = append(p.microcodes, TuplaMicrocode{
		Operation: Jump,
		Res:       labelEndCase,
		Op1:       nil,
		Op2:       nil,
	})

	p.microcodes = append(p.microcodes, TuplaMicrocode{
		Operation: Label,
		Res:       labelFalse,
		Op1:       nil,
		Op2:       nil,
	})

}

// func (p *Parser) parseReturnStmt() {

// 	p.consume(lexer.Func_return)
// 	if p.current().Token != lexer.Stmt_end {
// 		p.parseExpr()
// 	}
// 	p.consume(lexer.Stmt_end)

// }

func (p *Parser) parseType() lexer.TabelaPalavras {
	token := p.current().Token
	if typeTokens[token] {
		p.advance()
		return token
	} else {
		stringToken := p.tokenToString(token)
		utils.ThrowParserException(fmt.Sprintf("expected type, got '%v'", stringToken), p.current().Linha, p.current().Coluna)
		return -1
	}
}

// Precedencia de operadores

func (p *Parser) parseExpr() (*lexer.TuplaLex, []TuplaMicrocode) {
	return p.parseAtrib()
}

func (p *Parser) parseOptExpr() (*lexer.TuplaLex, []TuplaMicrocode) {
	if exprStartTokens[p.current().Token] {
		return p.parseAtrib()
	}
	return nil, nil
}

func (p *Parser) parseAtrib() (*lexer.TuplaLex, []TuplaMicrocode) {
	left, code := p.parseOR()
	if p.current().Token == lexer.Op_assign {
		if left.Token != lexer.Identifier {
			utils.ThrowParserException(
				"left side of assignment must be a variable",
				left.Linha,
				left.Coluna,
			)
		}
		p.consume(lexer.Op_assign)
		right, rightCode := p.parseAtrib()
		code = append(code, rightCode...)

		typeleft := p.toType(left)
		typeright := p.toType(right)
		verifyTypeCompatibility(typeleft, typeright, left.Linha, left.Coluna)

		code = append(code, TuplaMicrocode{
			Operation: Att,
			Res:       left,
			Op1:       right,
			Op2:       nil,
		})
		return left, code
	}
	return left, code
}

// <or> -> <xor> { 'quarque_um' <xor> }
func (p *Parser) parseOR() (*lexer.TuplaLex, []TuplaMicrocode) {
	left, commandList := p.parseXor()
	for p.current().Token == lexer.Op_or {
		p.advance()
		right, rightCommands := p.parseXor()
		if (p.toType(left) != Type_bool) || (p.toType(right) != Type_bool) {
			utils.ThrowParserException("Can not perform 'quarque_um' operation on non-boolean values", left.Linha, left.Coluna)
		}
		commandList = append(commandList, rightCommands...)
		typeLeft := p.toType(left)
		typeRight := p.toType(right)
		resType := p.inferType(lexer.Op_or, typeLeft, typeRight, left.Linha, left.Coluna)
		temp := p.newTemp(resType)
		commandList = append(commandList, TuplaMicrocode{
			Operation: Or,
			Res:       temp,
			Op1:       left,
			Op2:       right,
		})
		left = temp
	}
	return left, commandList
}

// <xor> -> <and> { 'um_o_oto' <and> }
func (p *Parser) parseXor() (*lexer.TuplaLex, []TuplaMicrocode) {
	left, commandList := p.parseAnd()
	for p.current().Token == lexer.Op_xor {
		p.advance()
		right, rightCommands := p.parseAnd()
		if (p.toType(left) != Type_bool) || (p.toType(right) != Type_bool) {
			utils.ThrowParserException("Can not perform 'um_o_oto' operation on non-boolean values", left.Linha, left.Coluna)
		}
		commandList = append(commandList, rightCommands...)
		typeLeft := p.toType(left)
		typeRight := p.toType(right)
		resType := p.inferType(lexer.Op_xor, typeLeft, typeRight, left.Linha, left.Coluna)
		temp := p.newTemp(resType)
		commandList = append(commandList, TuplaMicrocode{
			Operation: Xor,
			Res:       temp,
			Op1:       left,
			Op2:       right,
		})
		left = temp
	}
	return left, commandList
}

// <and> -> <not> { 'tamem' <not> }
func (p *Parser) parseAnd() (*lexer.TuplaLex, []TuplaMicrocode) {
	left, commandList := p.parseNot()
	for p.current().Token == lexer.Op_and {
		p.advance()
		right, rightCommands := p.parseNot()
		if (p.toType(left) != Type_bool) || (p.toType(right) != Type_bool) {
			utils.ThrowParserException("Can not perform 'tamem' operation on non-boolean values", left.Linha, left.Coluna)
		}
		commandList = append(commandList, rightCommands...)
		typeLeft := p.toType(left)
		typeRight := p.toType(right)
		resType := p.inferType(lexer.Op_and, typeLeft, typeRight, left.Linha, left.Coluna)
		temp := p.newTemp(resType)

		commandList = append(commandList, TuplaMicrocode{
			Operation: And,
			Res:       temp,
			Op1:       left,
			Op2:       right,
		})
		left = temp
	}
	return left, commandList
}

// <not> -> 'vam_marca' <not> | <rel>
func (p *Parser) parseNot() (*lexer.TuplaLex, []TuplaMicrocode) {
	if p.current().Token == lexer.Op_not {
		p.advance()
		value, commandList := p.parseNot()
		if p.toType(value) != Type_bool {
			utils.ThrowParserException("Can not perform 'vam_marca' operation on non-boolean values", value.Linha, value.Coluna)
		}
		temp := p.newTemp(Type_bool)
		commandList = append(commandList, TuplaMicrocode{
			Operation: Not,
			Res:       temp,
			Op1:       value,
			Op2:       nil,
		})
		return temp, commandList
	}
	return p.parseRel()
}

// <rel> -> <add> { ('mema_coisa' | 'neh_nada') <add> }
func (p *Parser) parseRel() (*lexer.TuplaLex, []TuplaMicrocode) {

	left, commandList := p.parseAdd()

	for relOpTokens[p.current().Token] {

		operator := p.current().Token
		p.advance()

		right, rightCommands := p.parseAdd()

		commandList = append(commandList, rightCommands...)

		typeLeft := p.toType(left)
		typeRight := p.toType(right)
		resType := p.inferType(operator, typeLeft, typeRight, left.Linha, left.Coluna)
		temp := p.newTemp(resType)

		var op TabelaMicrocodes

		switch operator {
		case lexer.Op_eq:
			op = Eq
		case lexer.Op_neq:
			op = Neq
		case lexer.Op_lt:
			op = Lt
		case lexer.Op_gt:
			op = Gt
		case lexer.Op_lte:
			op = Lte
		case lexer.Op_gte:
			op = Gte
		}

		commandList = append(commandList, TuplaMicrocode{
			Operation: op,
			Res:       temp,
			Op1:       left,
			Op2:       right,
		})

		left = temp
	}

	return left, commandList
}

// <add> -> <mul> { ('veiz' | 'sob') <mul> }
func (p *Parser) parseAdd() (*lexer.TuplaLex, []TuplaMicrocode) {

	left, commandList := p.parseMul()

	for p.current().Token == lexer.Op_add ||
		p.current().Token == lexer.Op_sub {

		operator := p.current().Token

		p.advance()

		right, rightCommands := p.parseMul()

		commandList = append(commandList, rightCommands...)

		typeLeft := p.toType(left)
		typeRight := p.toType(right)
		resType := p.inferType(operator, typeLeft, typeRight, left.Linha, left.Coluna)
		temp := p.newTemp(resType)

		var op TabelaMicrocodes

		if operator == lexer.Op_add {
			op = Add
		} else {
			op = Sub
		}

		commandList = append(commandList, TuplaMicrocode{
			Operation: op,
			Res:       temp,
			Op1:       left,
			Op2:       right,
		})

		left = temp
	}

	return left, commandList
}

// <mul> -> <fator> { ('veiz' | 'sob') <fator> }
func (p *Parser) parseMul() (*lexer.TuplaLex, []TuplaMicrocode) {

	left, commandList := p.parseUno()

	for mulOpTokens[p.current().Token] {

		operator := p.current().Token

		p.advance()

		right, rightCommands := p.parseUno()

		commandList = append(commandList, rightCommands...)

		typeLeft := p.toType(left)
		typeRight := p.toType(right)
		resType := p.inferType(operator, typeLeft, typeRight, left.Linha, left.Coluna)
		temp := p.newTemp(resType)

		var op TabelaMicrocodes

		switch operator {
		case lexer.Op_mul:
			op = Mul
		case lexer.Op_div:
			op = Div
		case lexer.Op_int_div:
			op = DivI
		case lexer.Op_mod:
			op = Mod
		}

		commandList = append(commandList, TuplaMicrocode{
			Operation: op,
			Res:       temp,
			Op1:       left,
			Op2:       right,
		})

		left = temp
	}

	return left, commandList
}

// <uno> -> '+' <uno> | '-' <uno> | <fatorZao>
func (p *Parser) parseUno() (*lexer.TuplaLex, []TuplaMicrocode) {

	if p.current().Token == lexer.Op_add {

		p.advance()

		return p.parseUno()
	}

	if p.current().Token == lexer.Op_sub {

		p.advance()

		value, commandList := p.parseUno()
		if p.toType(value) != Type_int && p.toType(value) != Type_float {
			utils.ThrowParserException("Can not perform unary '-' operation on non-numeric values", value.Linha, value.Coluna)
		}

		typeValue := p.toType(value)
		temp := p.newTemp(typeValue)

		commandList = append(commandList, TuplaMicrocode{
			Operation: Uno,
			Res:       temp,
			Op1:       value,
			Op2:       nil,
		})

		return temp, commandList
	}

	return p.parseFatorZao()
}

// <fatorZao> -> <fatorzin> | '(' <atrib> ')'
func (p *Parser) parseFatorZao() (*lexer.TuplaLex, []TuplaMicrocode) {

	if p.current().Token == lexer.Open_paren {

		p.advance()

		result, commandList := p.parseAtrib()

		p.consume(lexer.Close_paren)

		return result, commandList
	}

	return p.parseFatorZin()
}

// <fatorZin> -> 'STR' | 'IDENT' | 'NUMint' |...
func (p *Parser) parseFatorZin() (*lexer.TuplaLex, []TuplaMicrocode) {

	t := p.current()

	if fatorTokens[t.Token] {

		if t.Token == lexer.Identifier {
			_, ok := p.symbolTable[t.Lexema]
			if !ok {
				utils.ThrowParserException(fmt.Sprintf("Variable '%s' not declared", t.Lexema), t.Linha, t.Coluna)
			}
		}

		p.advance()

		return &t, []TuplaMicrocode{}
	}

	stringToken := p.tokenToString(t.Token)

	utils.ThrowParserException(
		fmt.Sprintf(
			"expected fator, got '%v'",
			stringToken,
		),
		t.Linha,
		t.Coluna,
	)

	return nil, nil
}

func (p *Parser) parseBreak() {
	if p.breakStack.IsEmpty() {
		utils.ThrowParserException("break statement outside of loop or switch", p.current().Linha, p.current().Coluna)
	}

	labelEnd := p.breakStack.Top()

	labelTuple := &lexer.TuplaLex{
		Token:  lexer.Literal_string,
		Lexema: labelEnd,
	}

	command := TuplaMicrocode{
		Operation: Jump,
		Res:       labelTuple,
		Op1:       nil,
		Op2:       nil,
	}

	// ListTuplaMicrocodeToString([]TuplaMicrocode{command})
	p.microcodes = append(p.microcodes, command)

	p.advance()
	p.consume(lexer.Stmt_end)
}

func (p *Parser) parseContinue() {
	if p.continueStack.IsEmpty() {
		utils.ThrowParserException("continue statement outside of loop", p.current().Linha, p.current().Coluna)
	}

	labelEnd := p.continueStack.Top()

	labelTuple := &lexer.TuplaLex{
		Token:  lexer.Literal_string,
		Lexema: labelEnd,
	}

	command := TuplaMicrocode{
		Operation: Jump,
		Res:       labelTuple,
		Op1:       nil,
		Op2:       nil,
	}

	// ListTuplaMicrocodeToString([]TuplaMicrocode{command})
	p.microcodes = append(p.microcodes, command)

	p.advance()
	p.consume(lexer.Stmt_end)
}
