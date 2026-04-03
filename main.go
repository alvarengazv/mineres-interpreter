package main

import (
	"mineres-interpreter/src/lexer"
	"mineres-interpreter/src/utils"
)

func main() {
	conteudo := utils.LerArquivo("data/main.uai")
	listTupla := lexer.AnalisarArquivo(conteudo)

	lexer.ListTuplaToString(listTupla)
}
