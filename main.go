package main

import "fmt"

type TabelaPalavras int

const (
	trem_di_numeru TabelaPalavras = iota
	trem_cum_virgula
	trem_discrita
	trem_discolhe
	trosso
)

func main() {
	fmt.Println("Interpretador de Mineirês em GO!")

}
