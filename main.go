package main

import (
	"mineres-interpreter/src/lexer"
	"mineres-interpreter/src/parser"
	"mineres-interpreter/src/utils"
)

func main() {
	conteudo := utils.LerArquivo("data/parseValidation_for_while_if.uai")
	listTupla := lexer.AnalisarArquivo(conteudo)
	//lexer.ListTuplaToString(listTupla)

	parser := parser.NewParser(listTupla)
	parser.ParserFunction()

}

