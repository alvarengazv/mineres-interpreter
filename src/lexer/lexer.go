package lexer

import (
	"fmt"
)

/**
 * AnalisarArquivo(resumo = "Analisa e separa o arquivo .uai",
 * 		Parâmetros = {
 * 			conteudo(
 *				description = "texto a ser verificado",
 *				example = "causo\r\n  Programa de teste pro interpretador em Minerês.\r\n"
 *			)
 *		},
 *		Retorno = {
 *			Conteúdo separado em tuplas = Array(Tupla)
 *		}
 * )
 */
func AnalisarArquivo(conteudo string) ([]Tupla) {
	var tabela []Tupla
	var tabela_lexica []Tupla
	var linha int = 0
	var coluna int = 0
	var linha_inicio int = 0
	var coluna_inicio int = 0
	var buffer []rune

	// Percorrendo o arquivo
	for _, char := range string(conteudo) {
		// Se for espaço em branco
		if char == ' ' && len(buffer) > 0 {
			//fmt.Println(string(buffer))
			if string(buffer) == "causo" {
				fmt.Println("Comentário")
				// Ler até encontrar fim_do_causo
				tabela_lexica = append(tabela_lexica, Tupla{
					lexema: string(buffer),
					token:  causo,
					linha:  linha_inicio,
					coluna: coluna_inicio,
				})
				fmt.Println(tabela_lexica)
			}

			buffer = []rune{}
		} else if char == ' ' {
			// Enquanto houver espaço em branco, incrementar a coluna, parar quando não houver mais espaço em branco
			coluna++
		} else if char == '\n' {
			linha++
		} else {
			if len(buffer) == 0 {
				linha_inicio = linha
				coluna_inicio = coluna
			}
			buffer = append(buffer, char)
		}

	}

	return tabela
}