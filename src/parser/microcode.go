package parser

import "mineres-interpreter/src/lexer"

type TabelaMicrocodes int

const (
	// Math Operations
	Add  TabelaMicrocodes = iota // 0  - (add, res, op1, op2)
	Sub                          // 1  - (sub, res, op1, op2)
	Mul                          // 2  - (mul, res, op1, op2)
	Div                          // 3  - (div, res, op1, op2)
	Mod                          // 4  - (mod, res, op1, op2)
	DivI                         // 5  - (divI, res, op1, op2)

	// Conditional operations
	Eq  // 6  - (eq, res, op1, op2)
	Neq // 7  - (neq, res, op1, op2)
	Lt  // 8  - (lt, res, op1, op2)
	Gt  // 9  - (gt, res, op1, op2)
	Lte // 10 - (lte, res, op1, op2)
	Gte // 11 - (gte, res, op1, op2)

	// Logical operations
	And // 12 - (and, res, op1, op2)
	Or  // 13 - (or, res, op1, op2)
	Not // 14 - (not, res, op1, null)
	Xor // 15 - (xor, res, op1, op2)

	// Others operations
	Call  // 16 - (call, func, var/null, value/null)
	Jump  // 17 - (jump, label, null, null)
	Label // 18 - (label, null, null, null)
	If_eq // 19 - (if_eq, op, labelTrue, labelFalse)
	Uno   // 20 - (uno, res, null, null)
	Att   // 21 - (att, op1, op2, null)
)

var PalavrasReservadas = map[string]TabelaMicrocodes{
	"add":   Add,
	"sub":   Sub,
	"mul":   Mul,
	"div":   Div,
	"mod":   Mod,
	"divI":  DivI,
	"eq":    Eq,
	"neq":   Neq,
	"lt":    Lt,
	"gt":    Gt,
	"lte":   Lte,
	"gte":   Gte,
	"and":   And,
	"or":    Or,
	"not":   Not,
	"xor":   Xor,
	"call":  Call,
	"jump":  Jump,
	"label": Label,
	"if_eq": If_eq,
	"uno":   Uno,
	"att":   Att,
}

var PalavrasReservadasReverso = map[TabelaMicrocodes]string{
	Add:   "add",
	Sub:   "sub",
	Mul:   "mul",
	Div:   "div",
	Mod:   "mod",
	DivI:  "divI",
	Eq:    "eq",
	Neq:   "neq",
	Lt:    "lt",
	Gt:    "gt",
	Lte:   "lte",
	Gte:   "gte",
	And:   "and",
	Or:    "or",
	Not:   "not",
	Xor:   "xor",
	Call:  "call",
	Jump:  "jump",
	Label: "label",
	If_eq: "if_eq",
	Uno:   "uno",
	Att:   "att",
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
