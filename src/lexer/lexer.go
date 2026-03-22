package lexer

import (
	"regexp"
	"unicode"
)

// estadoLexer mantém todo o estado mutável do analisador léxico,
// permitindo que as funções auxiliares compartilhem e modifiquem o estado.
type estadoLexer struct {
	runes         []rune
	tabela_lexica []Tupla
	linha         int
	coluna        int
	linha_inicio  int
	coluna_inicio int
	buffer        []rune
	erro_lexico   bool

	// Estados de modo do Lexer
	lendoString          bool
	lendoComentarioLinha bool
	lendoComentarioBloco bool

	// Constantes para comentário de bloco
	FimCauso        string
	tamanhoFimCauso int
}

func novoEstadoLexer(conteudo string) *estadoLexer {
	return &estadoLexer{
		runes:           []rune(conteudo),
		linha:           1,
		coluna:          1,
		FimCauso:        "fim_do_causo",
		tamanhoFimCauso: len("fim_do_causo"),
	}
}

// classificarLexema classifica um lexema que não é palavra reservada,
// verificando se é um literal numérico (hexa, octal, float, inteiro) ou variável.
func (e *estadoLexer) classificarLexema(lexema string) TabelaPalavras {
	// 0x seguido de um ou mais números hexadecimais
	regexHexa := regexp.MustCompile(`^0x[0-9A-F]+$`)
	// 0 seguido de um ou mais números octais, não podendo ter 0 como segundo caractere
	regexOctal := regexp.MustCompile(`^0[1-7][0-7]+$`)
	// Um ou mais números seguidos de ponto e um ou mais números, ou ponto seguido de um ou mais números
	regexFloat := regexp.MustCompile(`^[0-9]*\.[0-9]+$|^[0-9]+\.[0-9]+$`)
	// Um ou mais zeros seguidos, OU um número diferente de 0 seguido de qualquer número
	regexInteiro := regexp.MustCompile(`^[0]+$|^[1-9][0-9]*$`)
	// Começa com letra, seguido de letras, números ou underscore
	regexVariavel := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)
	switch {
	case regexHexa.MatchString(lexema):
		return conteudo_hexa
	case regexOctal.MatchString(lexema):
		return conteudo_octal
	case regexFloat.MatchString(lexema):
		return conteudo_float
	case regexInteiro.MatchString(lexema):
		return conteudo_inteiro
	case regexVariavel.MatchString(lexema):
		return variavel
	default:
		e.erro_lexico = true
		return erro_lexico
	}
}

// processarBuffer processa o conteúdo acumulado no buffer,
// identificando palavras reservadas ou variáveis.
func (e *estadoLexer) processarBuffer() {
	if len(e.buffer) == 0 {
		return
	}

	lexema := string(e.buffer)
	var token TabelaPalavras

	if t, existe := PalavrasReservadas[lexema]; existe {
		token = t
		if token == causo {
			e.lendoComentarioBloco = true
		}
	} else {
		token = e.classificarLexema(lexema)
	}

	e.tabela_lexica = append(e.tabela_lexica, Tupla{
		lexema: lexema,
		token:  token,
		linha:  e.linha_inicio,
		coluna: e.coluna_inicio,
	})
	e.buffer = []rune{}
}

// tratarComentarioBloco trata caracteres enquanto estiver dentro de um
// comentário de bloco (causo ... fim_do_causo).
// Retorna o novo índice i (possivelmente avançado).
func (e *estadoLexer) tratarComentarioBloco(i int) int {
	char := e.runes[i]

	if char == '\n' {
		e.linha++
		e.coluna = 1
	} else {
		e.coluna++
	}

	if char == 'f' && i+e.tamanhoFimCauso < len(e.runes) && string(e.runes[i:i+e.tamanhoFimCauso]) == e.FimCauso {
		e.tabela_lexica = append(e.tabela_lexica, Tupla{
			lexema: "fim_do_causo",
			token:  fim_do_causo,
			linha:  e.linha,
			coluna: e.coluna,
		})
		/*
		 * Esse i + 11 aqui é para ele saltar o restante do "fim_do_causo" depois de reconhecer ele,
		 * para não ficar lendo ele como parte do comentário de bloco.
		 * O número 11 é o tamanho da string "fim_do_causo".
		 */
		i += 11
		e.coluna += 11
		e.lendoComentarioBloco = false
	}

	return i
}

// tratarComentarioLinha trata caracteres enquanto estiver dentro de um
// comentário de linha (// até o final da linha).
func (e *estadoLexer) tratarComentarioLinha(char rune) {
	if char == '\n' {
		e.lendoComentarioLinha = false
		e.linha++
		e.coluna = 1
	} else {
		e.coluna++
	}
}

// tratarString trata caracteres enquanto estiver lendo o conteúdo de uma string.
func (e *estadoLexer) tratarString(char rune) {
	if char == '"' {
		e.tabela_lexica = append(e.tabela_lexica, Tupla{
			lexema: string(e.buffer),
			token:  conteudo_string,
			linha:  e.linha_inicio,
			coluna: e.coluna_inicio,
		})
		e.buffer = []rune{}
		e.lendoString = false
		e.tabela_lexica = append(e.tabela_lexica, Tupla{lexema: "\"", token: fecha_aspas, linha: e.linha, coluna: e.coluna})
	} else {
		e.buffer = append(e.buffer, char)
	}
	e.coluna++
}

// detectarInicioComentarioLinha detecta o início de um comentário de linha (//).
// Retorna o novo índice i (avançado para pular o segundo /).
func (e *estadoLexer) detectarInicioComentarioLinha(i int) int {
	e.processarBuffer()
	e.lendoComentarioLinha = true
	i++ // Pula o segundo /
	e.coluna += 2
	return i
}

// detectarInicioString trata o início de uma string literal (aspas duplas).
func (e *estadoLexer) detectarInicioString() {
	e.processarBuffer()
	e.tabela_lexica = append(e.tabela_lexica, Tupla{lexema: "\"", token: abre_aspas, linha: e.linha, coluna: e.coluna})
	e.lendoString = true
	e.linha_inicio = e.linha
	e.coluna_inicio = e.coluna + 1
	e.coluna++
}

// tratarEspacosDelimitadores trata espaços em branco e quebras de linha.
func (e *estadoLexer) tratarEspacosDelimitadores(char rune) {
	e.processarBuffer()
	if char == '\n' {
		e.linha++
		e.coluna = 1
	} else {
		e.coluna++
	}
}

// tratarSimbolosEspeciais trata símbolos especiais como parênteses, vírgulas, chaves, etc.
func (e *estadoLexer) tratarSimbolosEspeciais(char rune) {
	e.processarBuffer()
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
	e.tabela_lexica = append(e.tabela_lexica, Tupla{lexema: string(char), token: token, linha: e.linha, coluna: e.coluna})
	e.coluna++
}

// acumularBuffer adiciona um caractere ao buffer de leitura,
// registrando a posição de início se o buffer estiver vazio.
func (e *estadoLexer) acumularBuffer(char rune) {
	if len(e.buffer) == 0 {
		e.linha_inicio = e.linha
		e.coluna_inicio = e.coluna
	}
	e.buffer = append(e.buffer, char)
	e.coluna++
}

func AnalisarArquivo(conteudo string) ([]Tupla, bool) {
	e := novoEstadoLexer(conteudo)

	for i := 0; i < len(e.runes); i++ {
		if e.erro_lexico {
			break
		}
		char := e.runes[i]

		if e.lendoComentarioBloco {
			// Tratamento de Comentário de Bloco (causo ... fim_do_causo)
			i = e.tratarComentarioBloco(i)
		} else if e.lendoComentarioLinha {
			// Tratamento de Comentário de Linha (//)
			e.tratarComentarioLinha(char)
		} else if e.lendoString {
			// Tratamento de Strings
			e.tratarString(char)
		} else if char == '/' && i+1 < len(e.runes) && e.runes[i+1] == '/' {
			// Detectar início de comentário de linha //
			i = e.detectarInicioComentarioLinha(i)
		} else if char == '"' {
			// Detectar início de string "
			e.detectarInicioString()
		} else if unicode.IsSpace(char) {
			// Delimitadores e Espaços
			e.tratarEspacosDelimitadores(char)
		} else if char == '(' || char == ')' || char == ',' || char == '{' || char == '}' || char == '+' || char == '<' {
			// Símbolos especiais
			e.tratarSimbolosEspeciais(char)
		} else {
			// Acumular caractere no buffer
			e.acumularBuffer(char)
		}
	}

	// Processa o que sobrou no buffer
	e.processarBuffer()
	return e.tabela_lexica, e.erro_lexico
}
