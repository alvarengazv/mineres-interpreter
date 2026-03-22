package lexer

import (
	"unicode"
)

func AnalisarArquivo(conteudo string) []Tupla {
	var tabela_lexica []Tupla
	var linha int = 1
	var coluna int = 1
	var linha_inicio int
	var coluna_inicio int
	var buffer []rune

	// Estados do Lexer
	var lendoString bool
	var lendoComentarioLinha bool
	var lendoComentarioBloco bool
	var FimCauso string = "fim_do_causo"
	var tamanhoFimCauso int = len(FimCauso)


	processarBuffer := func() {
		if len(buffer) == 0 {
			return
		}

		lexema := string(buffer)
		var token TabelaPalavras

		if t, existe := PalavrasReservadas[lexema]; existe {
			token = t
			if token == causo {
				lendoComentarioBloco = true
			}
		} else {
			token = variavel
		}

		tabela_lexica = append(tabela_lexica, Tupla{
			lexema: lexema,
			token:  token,
			linha:  linha_inicio,
			coluna: coluna_inicio,
		})
		buffer = []rune{}
	}

	runes := []rune(conteudo)
	for i := 0; i < len(runes); i++ {
		char := runes[i]

		// Tratamento de Comentário de Bloco (causo ... fim_do_causo)
		if lendoComentarioBloco {
			if char == '\n' {
				linha++
				coluna = 1
			} else {
				coluna++
			}
			
			if char == 'f' && i+tamanhoFimCauso < len(runes) && string(runes[i:i+tamanhoFimCauso]) == FimCauso {
				tabela_lexica = append(tabela_lexica, Tupla{
					lexema: "fim_do_causo",
					token:  fim_do_causo,
					linha:  linha,
					coluna: coluna,
				})
				/*
				* Esse i + 11 aqui é para ele saltar o restante do "fim_do_causo" depois de reconhecer ele, 
				* para não ficar lendo ele como parte do comentário de bloco. 
				* O número 11 é o tamanho da string "fim_do_causo".
				*/
				i += 11 
				coluna += 11
				lendoComentarioBloco = false
			}
			continue
		}

		// Tratamento de Comentário de Linha (//)
		if lendoComentarioLinha {
			if char == '\n' {
				lendoComentarioLinha = false
				linha++
				coluna = 1
			} else {
				coluna++
			}
			continue
		}

		//Tratamento de Strings 
		if lendoString {
			if char == '"' {
				
				tabela_lexica = append(tabela_lexica, Tupla{
					lexema: string(buffer),
					token:  conteudo_string,
					linha:  linha_inicio,
					coluna: coluna_inicio,
				})
				buffer = []rune{}
				lendoString = false
				tabela_lexica = append(tabela_lexica, Tupla{lexema: "\"", token: fecha_aspas, linha: linha, coluna: coluna})
			} else {
				buffer = append(buffer, char)
			}
			coluna++
			continue
		}

		// Detectar início de comentário de linha //
		if char == '/' && i+1 < len(runes) && runes[i+1] == '/' {
			processarBuffer()
			lendoComentarioLinha = true
			i++ // Pula o segundo /
			coluna += 2
			continue
		}

		// Detectar início de string "
		if char == '"' {
			processarBuffer()
			tabela_lexica = append(tabela_lexica, Tupla{lexema: "\"", token: abre_aspas, linha: linha, coluna: coluna})
			lendoString = true
			linha_inicio = linha
			coluna_inicio = coluna + 1
			coluna++
			continue
		}

		// Delimitadores e Espaços
		if unicode.IsSpace(char) {
			processarBuffer()
			if char == '\n' {
				linha++
				coluna = 1
			} else {
				coluna++
			}
			continue
		}

		// Símbolos especiais
		if char == '(' || char == ')' || char == ',' || char == '{' || char == '}' || char == '+' || char == '<' {
			processarBuffer()
			var token TabelaPalavras
			switch char {
			case '(': token = abre_parentese
			case ')': token = fecha_parentese
			case ',': token = virgula
			case '{': token = abre_chave
			case '}': token = fecha_chave
			case '+': token = soma
			case '<': token = menor_que
			}
			tabela_lexica = append(tabela_lexica, Tupla{lexema: string(char), token: token, linha: linha, coluna: coluna})
			coluna++
			continue
		}

		if len(buffer) == 0 {
			linha_inicio = linha
			coluna_inicio = coluna
		}
		buffer = append(buffer, char)
		coluna++
	}

	processarBuffer()
	return tabela_lexica
}