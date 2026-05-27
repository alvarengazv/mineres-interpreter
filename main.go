package main

import (
	"mineres-interpreter/src/interpreter"
	"mineres-interpreter/src/lexer"
	"mineres-interpreter/src/parser"
	"mineres-interpreter/src/utils"
)

func main() {
    conteudo := utils.LerArquivo("data/gramatica.uai")
    listTupla := lexer.AnalisarArquivo(conteudo)

    parser := parser.NewParser(listTupla)
    ir := parser.ParserFunction()

	interpreter := interpreter.NewInterpreter(ir)
	interpreter.Run()
}
