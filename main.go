package main

import (
	"mineres-interpreter/src/interpreter"
	"mineres-interpreter/src/lexer"
	"mineres-interpreter/src/parser"
	"mineres-interpreter/src/utils"
)

func main() {
    conteudo := utils.LerArquivo("data/tests/operationAndTest.uai")
    listTupla := lexer.AnalisarArquivo(conteudo)

    parser := parser.NewParser(listTupla)
    codigoIntermediario := parser.ParserFunction()

	interpreter := interpreter.NewInterpreter(codigoIntermediario)
	interpreter.Run()
}
