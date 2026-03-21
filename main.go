package main

import (
	"fmt"
	"mineres-interpreter/src/lexer"
	"mineres-interpreter/src/utils"
)

func main() {
	fmt.Println("Interpretador de Mineirês em GO!")
	conteudo:= utils.LerArquivo("data/main.uai")
	listTupla:= lexer.AnalisarArquivo(conteudo);

	lexer.ListTuplaToString(listTupla);
}