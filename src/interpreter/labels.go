package interpreter

import (
	"fmt"
	"mineres-interpreter/src/parser"
)

func (i *Interpreter) buildLabels() {
	for idx, instr := range i.code {
		if instr.Operation == parser.Label {
			if instr.Res != nil {
				i.labels[instr.Res.Lexema] = idx
			}
		}
	}
}


/**
 * PrintLabels(Resumo = "Imprime todos os labels do código intermediário",
 *		Parâmetros = {
 * 			*interprete(
 *				description = "ponteiro para o interpretador",
 *				example = "file.go"
 *			),
 *		},
 *		Retorno = {}
 * )
 */
func (i *Interpreter) PrintLabels() {
	fmt.Println("=== Labels ===")

	for label, index := range i.labels {
		fmt.Printf("%s -> %d\n", label, index)
	}
}