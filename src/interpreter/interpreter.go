package interpreter

import "mineres-interpreter/src/parser"

type Interpreter struct {
	code []parser.TuplaMicrocode
	ip   int

	labels map[string]int
	memory map[string]interface{}
}

func NewInterpreter(code []parser.TuplaMicrocode) *Interpreter {
	i := &Interpreter{
		code:   code,
		ip:     0,
		labels: make(map[string]int),
		memory: make(map[string]interface{}),
	}

	//i.buildLabels()

	return i
}

func (i *Interpreter) Run() {
	for i.ip < len(i.code) {
		//instr := i.code[i.ip]
		//i.execute(instr)
		i.ip++
	}
}