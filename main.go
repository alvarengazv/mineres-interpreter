package main

import (
	"mineres-interpreter/src/lexer"
	"mineres-interpreter/src/utils"
)

func main() {
	conteudo := utils.LerArquivo("data/lexerValidation_numero-invalido.uai")
	listTupla := lexer.AnalisarArquivo(conteudo)

	lexer.ListTuplaToString(listTupla)
}