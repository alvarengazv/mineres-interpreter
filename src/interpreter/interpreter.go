package interpreter

import (
	"mineres-interpreter/src/parser"
)

type Interpreter struct {
	code []parser.TuplaMicrocode
	ip   int
	labels map[string]int
	memory    map[string]any
}

func NewInterpreter(code []parser.TuplaMicrocode) *Interpreter {
	interpreter := &Interpreter{
		code:         code,
		ip: 		  0,
		labels:       make(map[string]int),
		memory:    make(map[string]any),
	}

	interpreter.buildLabels()

	return interpreter
}

func (interpreter *Interpreter) Run() {
	for interpreter.ip < len(interpreter.code) {
		instrucao := interpreter.code[interpreter.ip]
		interpreter.execute(instrucao)
		interpreter.ip++
	}
}

