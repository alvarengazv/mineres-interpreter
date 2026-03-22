package main

import (
	"fmt"
	"mineres-interpreter/src/lexer"
	"mineres-interpreter/src/utils"
)

func main() {
	fmt.Println("Interpretador de Mineirês em GO!")
	conteudo := utils.LerArquivo("data/main.uai")
	listTupla, erro := lexer.AnalisarArquivo(conteudo)

	if erro {
		lexer.ListTuplaToString(listTupla)
		fmt.Println("Erro léxico encontrado!")
		return
	}

}