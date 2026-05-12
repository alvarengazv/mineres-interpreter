package parser

import (
	"fmt"
	"mineres-interpreter/src/lexer"
)

type TuplaMicrocode struct {
	operation TabelaMicrocodes
	res       *lexer.TuplaLex
	op1       *lexer.TuplaLex
	op2       *lexer.TuplaLex
}

func TuplaMicrocodeToString(t TuplaMicrocode) {
	resStr := "null"
	if t.res != nil {
		resStr = t.res.Lexema
	}
	op1Str := "null"
	if t.op1 != nil {
		op1Str = t.op1.Lexema
	}
	op2Str := "null"
	if t.op2 != nil {
		op2Str = t.op2.Lexema
	}

	fmt.Printf("Operation: %s | Res: %s | Op1: %s | Op2: %s\n",
		PalavrasReservadasReverso[t.operation],
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
