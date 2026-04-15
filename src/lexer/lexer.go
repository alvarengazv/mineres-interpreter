package lexer

import (
	"mineres-interpreter/src/utils"
	"regexp"
	"unicode"
)

var (
	// 0x seguido de um ou mais números hexadecimais
	regexHexa = regexp.MustCompile(`^0x[0-9A-F]+$`)
	// 0 seguido de um ou mais números octais, não podendo ter 0 como segundo caractere
	regexOctal = regexp.MustCompile(`^0[1-7][0-7]*$`)
	// Um ou mais números seguidos de ponto e um ou mais números, ou ponto seguido de um ou mais números
	regexFloat = regexp.MustCompile(`^[0-9]*\.[0-9]+$|^[0-9]+\.[0-9]+$`)
	// Um ou mais zeros seguidos, OU um número diferente de 0 seguido de qualquer número
	regexInteiro = regexp.MustCompile(`^[0]+$|^[1-9][0-9]*$`)
	// Começa com letra, seguido de letras, números ou underscore
	regexVariavel = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)
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
	lendoChar            bool
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
	switch {
	case regexHexa.MatchString(lexema):
		return Literal_hex
	case regexOctal.MatchString(lexema):
		return Literal_oct
	case regexFloat.MatchString(lexema):
		return Literal_float
	case regexInteiro.MatchString(lexema):
		return Literal_int
	case regexVariavel.MatchString(lexema):
		return Identifier
	default:
		if e.lendoChar {
			utils.ThrowLexerException("Unknown token: Unterminated char literal", e.linha_inicio, e.coluna_inicio)
		} else if e.lendoString {
			utils.ThrowLexerException("Unknown token: Unterminated string literal", e.linha_inicio, e.coluna_inicio)
		} else if lexema[0] >= '0' && lexema[0] <= '9' {
			utils.ThrowLexerException("Unknown token: Invalid number value", e.linha_inicio, e.coluna_inicio)
		}
		utils.ThrowLexerException("Unknown token: Invalid character", e.linha_inicio, e.coluna_inicio)

		e.erro_lexico = true
		return Lexical_error
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
		if token == Comment_block_open {
			e.lendoComentarioBloco = true
			e.buffer = []rune{}
			return
		}
		if token == Comment_block_close && !e.lendoComentarioBloco {
			utils.ThrowLexerException("Malformed comment block: 'causo' is required before 'fim_do_causo'", e.linha_inicio, e.coluna_inicio)
			e.erro_lexico = true
			return
		}
	} else {
		token = e.classificarLexema(lexema)
	}

	e.tabela_lexica = append(e.tabela_lexica, Tupla{
		Lexema: lexema,
		Token:  token,
		Linha:  e.linha_inicio,
		Coluna: e.coluna_inicio,
	})
	e.buffer = []rune{}
}

// tratarComentarioBloco trata caracteres enquanto estiver dentro de um
// comentário de bloco (causo ... fim_do_causo).
// Retorna o novo índice i (possivelmente avançado).
func (e *estadoLexer) tratarComentarioBloco(i int) int {
	char := e.runes[i]

	if char == 'f' && i+e.tamanhoFimCauso <= len(e.runes) && string(e.runes[i:i+e.tamanhoFimCauso]) == e.FimCauso {
		/*
		 * O -1 é para não pular o último caractere do "fim_do_causo", que é o 'o'.
		 * Então, se houver algum comando logo após o "fim_do_causo", ele será lido corretamente.
		 * Como, por exemplo, fim_do_causouai, o 'u' será lido corretamente.
		 */
		avanco := e.tamanhoFimCauso - 1
		e.coluna += e.tamanhoFimCauso
		e.lendoComentarioBloco = false
		return i + avanco
	}

	// Somar a coluna e a linha caso não seja o fim do causo
	if char == '\n' {
		e.linha++
		e.coluna = 1
	} else {
		e.coluna++
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

// tratarSequenciaEscape verifica se o caractere atual e o próximo formam uma
// sequência de escape válida (ex: \n, \t, \', \", \\). Se sim, trata de acordo com o contexto:
// - Dentro de string/char: converte para o caractere real e acumula no buffer.
// - Fora de string/char: retorna false para que a barra seja tratada como caractere normal.
func (e *estadoLexer) tratarSequenciaEscape(i int) (bool, int) {
	if i < len(e.runes) && e.runes[i] == '\\' && i+1 < len(e.runes) {
		// Só processa sequências de escape se estivermos dentro de uma string ou char
		if !e.lendoString && !e.lendoChar {
			return false, i
		}

		prox := e.runes[i+1]

		// Verifica se é uma sequência de escape conhecida e qual o seu caractere real
		var runeReal rune
		var ehEscapeValido bool

		switch prox {
		case 'n':
			runeReal = '\n'
			ehEscapeValido = true
		case 't':
			runeReal = '\t'
			ehEscapeValido = true
		case '\'', '"', '\\':
			runeReal = prox
			ehEscapeValido = true
		default:
			ehEscapeValido = false
		}

		if ehEscapeValido {
			// Acumula o caractere real no buffer
			if len(e.buffer) == 0 {
				e.linha_inicio = e.linha
				e.coluna_inicio = e.coluna
			}
			e.buffer = append(e.buffer, runeReal)

			i++ // Pula o próximo caractere (n, t, ', ", \)
			e.coluna += 2
			return true, i
		}
	}
	return false, i
}

// tratarString trata caracteres enquanto estiver lendo o conteúdo de uma string.
func (e *estadoLexer) tratarString(char rune, i int) int {
	if char == '\n' || char == '\r' {
		utils.ThrowLexerException("Unterminated string literal", e.linha_inicio, e.coluna_inicio)
	} else if char == '"' {
		e.tabela_lexica = append(e.tabela_lexica, Tupla{
			Lexema: string(e.buffer),
			Token:  Literal_string,
			Linha:  e.linha_inicio,
			Coluna: e.coluna_inicio,
		})
		e.buffer = []rune{}
		e.lendoString = false
		e.coluna++
	} else if detectado, novoI := e.tratarSequenciaEscape(i); detectado {
		return novoI
	} else {
		if len(e.buffer) == 0 {
			e.linha_inicio = e.linha
			e.coluna_inicio = e.coluna
		}
		e.buffer = append(e.buffer, char)
		e.coluna++
	}
	return i
}

// tratarChar trata caracteres enquanto estiver lendo o conteúdo de um char literal.
func (e *estadoLexer) tratarChar(char rune, i int) int {
	if char == '\n' || char == '\r' {
		utils.ThrowLexerException("Unterminated char literal", e.linha_inicio, e.coluna_inicio)
	} else if char == '\'' {
		e.tabela_lexica = append(e.tabela_lexica, Tupla{
			Lexema: string(e.buffer),
			Token:  Literal_char,
			Linha:  e.linha_inicio,
			Coluna: e.coluna_inicio,
		})
		e.buffer = []rune{}
		e.lendoChar = false
		e.coluna++
	} else if detectado, novoI := e.tratarSequenciaEscape(i); detectado {
		return novoI
	} else if len(e.buffer) >= 1 {
		utils.ThrowLexerException("Unknown token: char literal has more than one character", e.linha_inicio, e.coluna_inicio)
		e.erro_lexico = true
		return i
	} else {
		if len(e.buffer) == 0 {
			e.linha_inicio = e.linha
			e.coluna_inicio = e.coluna
		}
		e.buffer = append(e.buffer, char)
		e.coluna++
	}
	return i
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
	e.lendoString = true
	e.linha_inicio = e.linha
	e.coluna_inicio = e.coluna + 1
	e.coluna++
}

// detectarInicioChar trata o início de um char (aspas simples).
func (e *estadoLexer) detectarInicioChar() {
	e.processarBuffer()
	e.lendoChar = true
	e.linha_inicio = e.linha
	e.coluna_inicio = e.coluna + 1
	e.coluna++
}

// tratarEspacosDelimitadores trata espaços em branco e quebras de linha.
func (e *estadoLexer) tratarEspacosDelimitadores(i int) int {
	e.processarBuffer()

	// Enquanto o próximo também for espaço, a gente só vai incrementando
	for i < len(e.runes) && unicode.IsSpace(e.runes[i]) {
		if e.runes[i] == '\n' {
			e.linha++
			e.coluna = 1
		} else {
			e.coluna++
		}
		if i+1 < len(e.runes) && unicode.IsSpace(e.runes[i+1]) {
			i++ // Pula pro próximo espaço
		} else {
			break // Sai se o próximo não for mais espaço
		}
	}
	return i
}

// tratarSimbolosEspeciais trata símbolos especiais como parênteses, vírgulas, chaves, etc.
// Para < e >, verifica se o próximo caractere é '=' para formar <= ou >=.
// Retorna o novo índice i (possivelmente avançado).
func (e *estadoLexer) tratarSimbolosEspeciais(char rune, i int) int {
	e.processarBuffer()
	var token TabelaPalavras
	lexema := string(char)

	switch char {
	case '(':
		token = Open_paren
	case ')':
		token = Close_paren
	case ',':
		token = Comma
	case ':':
		token = Colon
	case '{':
		token = Open_brace
	case '}':
		token = Close_brace
	case '+':
		token = Op_add
	case '-':
		token = Op_sub
	case '<', '>':
		if i+1 < len(e.runes) && e.runes[i+1] == '=' {
			lexema = string(char) + "="
			if char == '<' {
				token = Op_lte
			} else {
				token = Op_gte
			}
			i++
		} else if char == '<' {
			token = Op_lt
		} else {
			token = Op_gt
		}
	case '%':
		token = Op_mod
	case '/':
		token = Op_int_div
	case ';':
		token = Stmt_end_for
	}

	e.tabela_lexica = append(e.tabela_lexica, Tupla{Lexema: lexema, Token: token, Linha: e.linha, Coluna: e.coluna})
	e.coluna += len([]rune(lexema))
	return i
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

func AnalisarArquivo(conteudo string) []Tupla {
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
			i = e.tratarString(char, i)
		} else if e.lendoChar {
			// Tratamento de Caracteres
			i = e.tratarChar(char, i)
		} else if char == '/' && i+1 < len(e.runes) && e.runes[i+1] == '/' {
			// Detectar início de comentário de linha //
			i = e.detectarInicioComentarioLinha(i)
		} else if char == '"' {
			// Detectar início de string "
			e.detectarInicioString()
		} else if char == '\'' {
			// Detectar início de caractere '
			e.detectarInicioChar()
		} else if detectado, novoI := e.tratarSequenciaEscape(i); detectado {
			// Tratamento de sequências de escape fora de strings/chars
			i = novoI
		} else if unicode.IsSpace(char) {
			// Delimitadores e Espaços
			i = e.tratarEspacosDelimitadores(i)
		} else if char == '(' || char == ')' || char == ',' || char == '{' || char == '}' || char == '+' || char == '-' || char == '%' || char == '/' || char == '<' || char == '>' || char == ';' || char == ':' {
			// Símbolos especiais (inclui operadores compostos <= e >=)
			i = e.tratarSimbolosEspeciais(char, i)
		} else {
			// Acumular caractere no buffer
			e.acumularBuffer(char)
		}
	}

	// Processa o que sobrou no buffer
	e.processarBuffer()

	if e.erro_lexico {
		utils.ThrowException("lexor.go", "AnalisarArquivo", "Erro desconhecido encontrado")
	}

	/* Verifica se ficou algum comentário de bloco aberto
	   Deve ser ao final pois anteriormente poderia não considerar um fim de bloco ao final do arquivo
	*/
	if e.lendoComentarioBloco {
		utils.ThrowLexerException("Unterminated comment block: 'fim_do_causo' is required after using 'causo'", e.linha_inicio, e.coluna_inicio)
	}

	return e.tabela_lexica
}
