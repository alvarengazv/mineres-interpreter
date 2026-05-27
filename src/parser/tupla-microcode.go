package parser

import (
	"fmt"
	"mineres-interpreter/src/lexer"
)

type TuplaMicrocode struct {
	Operation TabelaMicrocodes
	Res       *lexer.TuplaLex
	Op1       *lexer.TuplaLex
	Op2       *lexer.TuplaLex
}

func TuplaMicrocodeToString(t TuplaMicrocode) {
	resStr := "null"
	if t.Res != nil {
		resStr = t.Res.Lexema
	}
	op1Str := "null"
	if t.Op1 != nil {
		op1Str = t.Op1.Lexema
	}
	op2Str := "null"
	if t.Op2 != nil {
		op2Str = t.Op2.Lexema
	}

	fmt.Printf("Operation: %s | Res: %s | Op1: %s | Op2: %s\n",
		PalavrasReservadasReverso[t.Operation],
		resStr,
		op1Str,
		op2Str,
	)
}

func ListTuplaMicrocodeToString(lista []TuplaMicrocode) {
	for _, item := range lista {
		TuplaMicrocodeToString(item)
	}
}
