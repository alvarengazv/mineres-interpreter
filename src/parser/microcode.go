package parser

import "mineres-interpreter/src/lexer"

type TabelaMicrocodes int

const (
	// Math Operations
	add  TabelaMicrocodes = iota // 0  - (add, res, op1, op2)
	sub                          // 1  - (sub, res, op1, op2)
	mul                          // 2  - (mul, res, op1, op2)
	div                          // 3  - (div, res, op1, op2)
	mod                          // 4  - (mod, res, op1, op2)
	divI                         // 5  - (divI, res, op1, op2)

	// Conditional operations
	eq  // 6  - (eq, res, op1, op2)
	neq // 7  - (neq, res, op1, op2)
	lt  // 8  - (lt, res, op1, op2)
	gt  // 9  - (gt, res, op1, op2)
	lte // 10 - (lte, res, op1, op2)
	gte // 11 - (gte, res, op1, op2)

	// Logical operations
	and // 12 - (and, res, op1, op2)
	or  // 13 - (or, res, op1, op2)
	not // 14 - (not, res, op1, null)
	xor // 15 - (xor, res, op1, op2)

	// Others operations
	call  // 16 - (call, func, var/null, value/null)
	jump  // 17 - (jump, label, null, null)
	label // 18 - (label, null, null, null)
	if_eq // 19 - (if_eq, op, labelTrue, labelFalse)
	uno   // 20 - (uno, res, null, null)
	att   // 21 - (att, op1, op2, null)
)

var PalavrasReservadas = map[string]TabelaMicrocodes{
	"add":   add,
	"sub":   sub,
	"mul":   mul,
	"div":   div,
	"mod":   mod,
	"divI":  divI,
	"eq":    eq,
	"neq":   neq,
	"lt":    lt,
	"gt":    gt,
	"lte":   lte,
	"gte":   gte,
	"and":   and,
	"or":    or,
	"not":   not,
	"xor":   xor,
	"call":  call,
	"jump":  jump,
	"label": label,
	"if_eq": if_eq,
	"uno":   uno,
	"att":   att,
}

var PalavrasReservadasReverso = map[TabelaMicrocodes]string{
	add:   "add",
	sub:   "sub",
	mul:   "mul",
	div:   "div",
	mod:   "mod",
	divI:  "divI",
	eq:    "eq",
	neq:   "neq",
	lt:    "lt",
	gt:    "gt",
	lte:   "lte",
	gte:   "gte",
	and:   "and",
	or:    "or",
	not:   "not",
	xor:   "xor",
	call:  "call",
	jump:  "jump",
	label: "label",
	if_eq: "if_eq",
	uno:   "uno",
	att:   "att",
}

func getDefaultValue(typeToken lexer.TabelaPalavras) *lexer.TuplaLex {

	switch typeToken {
	case lexer.Type_int:
		return &lexer.TuplaLex{Token: lexer.Literal_int, Lexema: "0", Linha: 0, Coluna: 0}
	case lexer.Type_float:
		return &lexer.TuplaLex{Token: lexer.Literal_float, Lexema: "0.0", Linha: 0, Coluna: 0}
	case lexer.Type_string:
		return &lexer.TuplaLex{Token: lexer.Literal_string, Lexema: "", Linha: 0, Coluna: 0}
	case lexer.Type_bool:
		return &lexer.TuplaLex{Token: lexer.Literal_false, Lexema: "false", Linha: 0, Coluna: 0}
	case lexer.Type_char:
		return &lexer.TuplaLex{Token: lexer.Literal_char, Lexema: "", Linha: 0, Coluna: 0}
	default:
		return nil
	}

}
