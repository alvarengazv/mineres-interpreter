package lexer

import (
	"mineres-interpreter/src/utils"
	"testing"
)

// expectPanic executa fn e verifica se houve um panic do tipo utils.CompilerError.
func expectPanic(t *testing.T, fn func()) string {
	t.Helper()
	old := utils.ExitOnError
	utils.ExitOnError = false
	defer func() { utils.ExitOnError = old }()

	var msg string
	func() {
		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("expected panic but none occurred")
			}
			ce, ok := r.(utils.CompilerError)
			if !ok {
				t.Fatalf("expected CompilerError, got %T: %v", r, r)
			}
			msg = ce.Message
		}()
		fn()
	}()
	return msg
}

// ---------- Palavras Reservadas ----------

func TestLexer_PalavrasReservadas(t *testing.T) {
	tests := []struct {
		input string
		token TabelaPalavras
	}{
		{"trem_di_numeru", Type_int},
		{"trem_cum_virgula", Type_float},
		{"trem_discrita", Type_string},
		{"trem_discolhe", Type_bool},
		{"trosso", Type_char},
		{"uai_se", Conditional_if},
		{"uai_senao", Conditional_else},
		{"roda_esse_trem", Loop_for},
		{"enquanto_tiver_trem", Loop_while},
		{"dependenu", Conditional_switch},
		{"du_casu", Conditional_case},
		{"uai_so", Conditional_default},
		{"bora_cumpade", Func_decl},
		{"main", Main_function},
		{"para_o_trem", Loop_break},
		{"toca_o_trem", Loop_continue},
		{"eh", Literal_true},
		{"num_eh", Literal_false},
		{"simbora", Block_open},
		{"cabo", Block_close},
		{"uai", Stmt_end},
		{"fica_assim_entao", Op_assign},
		{"neh_nada", Op_neq},
		{"mema_coisa", Op_eq},
		{"quarque_um", Op_or},
		{"vam_marca", Op_not},
		{"tamem", Op_and},
		{"um_o_oto", Op_xor},
		{"veiz", Op_mul},
		{"sob", Op_div},
		{"xove", Io_scan},
		{"oia_proce_ve", Io_print},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			// Cada token isolado + espaço para finalizar
			tokens := AnalisarArquivo(tt.input + " ")
			if len(tokens) == 0 {
				t.Fatalf("no tokens produced for input %q", tt.input)
			}
			if tokens[0].Token != tt.token {
				t.Errorf("token = %d, want %d for input %q", tokens[0].Token, tt.token, tt.input)
			}
			if tokens[0].Lexema != tt.input {
				t.Errorf("lexema = %q, want %q", tokens[0].Lexema, tt.input)
			}
		})
	}
}

// ---------- Identificadores ----------

func TestLexer_Identificadores(t *testing.T) {
	tests := []string{"a", "var1", "minha_var", "xYz", "abc123"}
	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			tokens := AnalisarArquivo(input + " ")
			if len(tokens) == 0 {
				t.Fatal("no tokens")
			}
			if tokens[0].Token != Identifier {
				t.Errorf("expected Identifier, got %d for %q", tokens[0].Token, input)
			}
		})
	}
}

// ---------- Literais Numéricos ----------

func TestLexer_LiteralInt(t *testing.T) {
	tests := []string{"0", "42", "100", "999"}
	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			tokens := AnalisarArquivo(input + " ")
			if len(tokens) == 0 {
				t.Fatal("no tokens")
			}
			if tokens[0].Token != Literal_int {
				t.Errorf("expected Literal_int, got %d for %q", tokens[0].Token, input)
			}
		})
	}
}

func TestLexer_LiteralOctal(t *testing.T) {
	tests := []string{"07", "0123", "0777"}
	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			tokens := AnalisarArquivo(input + " ")
			if len(tokens) == 0 {
				t.Fatal("no tokens")
			}
			if tokens[0].Token != Literal_oct {
				t.Errorf("expected Literal_oct, got %d for %q", tokens[0].Token, input)
			}
		})
	}
}

func TestLexer_LiteralHex(t *testing.T) {
	tests := []string{"0xFF", "0x1A", "0x0"}
	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			tokens := AnalisarArquivo(input + " ")
			if len(tokens) == 0 {
				t.Fatal("no tokens")
			}
			if tokens[0].Token != Literal_hex {
				t.Errorf("expected Literal_hex, got %d for %q", tokens[0].Token, input)
			}
		})
	}
}

func TestLexer_LiteralFloat(t *testing.T) {
	tests := []string{"3.14", "0.0", "10.5"}
	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			tokens := AnalisarArquivo(input + " ")
			if len(tokens) == 0 {
				t.Fatal("no tokens")
			}
			if tokens[0].Token != Literal_float {
				t.Errorf("expected Literal_float, got %d for %q", tokens[0].Token, input)
			}
		})
	}
}

// ---------- Literais String e Char ----------

func TestLexer_LiteralString(t *testing.T) {
	tokens := AnalisarArquivo(`"hello world" `)
	if len(tokens) == 0 {
		t.Fatal("no tokens")
	}
	if tokens[0].Token != Literal_string {
		t.Errorf("expected Literal_string, got %d", tokens[0].Token)
	}
	if tokens[0].Lexema != "hello world" {
		t.Errorf("lexema = %q, want %q", tokens[0].Lexema, "hello world")
	}
}

func TestLexer_LiteralStringVazia(t *testing.T) {
	tokens := AnalisarArquivo(`"" `)
	if len(tokens) == 0 {
		t.Fatal("no tokens")
	}
	if tokens[0].Token != Literal_string {
		t.Errorf("expected Literal_string, got %d", tokens[0].Token)
	}
	if tokens[0].Lexema != "" {
		t.Errorf("lexema = %q, want empty string", tokens[0].Lexema)
	}
}

func TestLexer_LiteralChar(t *testing.T) {
	tokens := AnalisarArquivo("'a' ")
	if len(tokens) == 0 {
		t.Fatal("no tokens")
	}
	if tokens[0].Token != Literal_char {
		t.Errorf("expected Literal_char, got %d", tokens[0].Token)
	}
	if tokens[0].Lexema != "a" {
		t.Errorf("lexema = %q, want %q", tokens[0].Lexema, "a")
	}
}

// ---------- Sequências de Escape ----------

func TestLexer_SequenciasEscape(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"newline", `"\n" `, "\n"},
		{"tab", `"\t" `, "\t"},
		{"backslash", `"\\" `, "\\"},
		{"escaped_quote", `"\"" `, "\""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := AnalisarArquivo(tt.input)
			if len(tokens) == 0 {
				t.Fatal("no tokens")
			}
			if tokens[0].Lexema != tt.expected {
				t.Errorf("lexema = %q, want %q", tokens[0].Lexema, tt.expected)
			}
		})
	}
}

// ---------- Comentários ----------

func TestLexer_ComentarioLinha(t *testing.T) {
	tokens := AnalisarArquivo("uai // isso é um comentario\nuai ")
	if len(tokens) != 2 {
		t.Fatalf("expected 2 tokens, got %d", len(tokens))
	}
	if tokens[0].Token != Stmt_end || tokens[1].Token != Stmt_end {
		t.Error("comments should be ignored")
	}
}

func TestLexer_ComentarioBloco(t *testing.T) {
	tokens := AnalisarArquivo("uai causo isso tudo é ignorado fim_do_causo uai ")
	if len(tokens) != 2 {
		t.Fatalf("expected 2 tokens, got %d", len(tokens))
	}
}

// ---------- Operadores Compostos ----------

func TestLexer_OperadoresCompostos(t *testing.T) {
	tests := []struct {
		input string
		token TabelaPalavras
	}{
		{"<= ", Op_lte},
		{">= ", Op_gte},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tokens := AnalisarArquivo(tt.input)
			if len(tokens) == 0 {
				t.Fatal("no tokens")
			}
			if tokens[0].Token != tt.token {
				t.Errorf("token = %d, want %d for %q", tokens[0].Token, tt.token, tt.input)
			}
		})
	}
}

// ---------- Operadores Simples ----------

func TestLexer_OperadoresSimples(t *testing.T) {
	tests := []struct {
		input string
		token TabelaPalavras
	}{
		{"+ ", Op_add},
		{"- ", Op_sub},
		{"% ", Op_mod},
		{"< ", Op_lt},
		{"> ", Op_gt},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tokens := AnalisarArquivo(tt.input)
			if len(tokens) == 0 {
				t.Fatal("no tokens")
			}
			if tokens[0].Token != tt.token {
				t.Errorf("token = %d, want %d for %q", tokens[0].Token, tt.token, tt.input)
			}
		})
	}
}

// ---------- Posição (Linha/Coluna) ----------

func TestLexer_Posicao(t *testing.T) {
	// O token "uai" na segunda linha, coluna 5
	tokens := AnalisarArquivo("    \n    uai ")
	if len(tokens) == 0 {
		t.Fatal("no tokens")
	}
	if tokens[0].Linha != 2 {
		t.Errorf("Linha = %d, want 2", tokens[0].Linha)
	}
	if tokens[0].Coluna != 5 {
		t.Errorf("Coluna = %d, want 5", tokens[0].Coluna)
	}
}

// ---------- Erros Léxicos ----------

func TestLexer_Erro_StringNaoFechada(t *testing.T) {
	msg := expectPanic(t, func() {
		AnalisarArquivo("\"string sem fechar\n")
	})
	if msg == "" {
		t.Error("expected error message")
	}
}

func TestLexer_Erro_CharNaoFechado(t *testing.T) {
	msg := expectPanic(t, func() {
		AnalisarArquivo("'abc' ")
	})
	if msg == "" {
		t.Error("expected error message")
	}
}

func TestLexer_Erro_ComentarioBlocoNaoFechado(t *testing.T) {
	msg := expectPanic(t, func() {
		AnalisarArquivo("causo comentario sem fim")
	})
	if msg == "" {
		t.Error("expected error message")
	}
}

func TestLexer_Erro_NumeroInvalido(t *testing.T) {
	msg := expectPanic(t, func() {
		AnalisarArquivo("123abc ")
	})
	if msg == "" {
		t.Error("expected error message")
	}
}

// ---------- Programa Completo ----------

func TestLexer_ProgramaCompleto(t *testing.T) {
	input := "bora_cumpade main() simbora cabo "
	tokens := AnalisarArquivo(input)

	expectedTokens := []TabelaPalavras{
		Func_decl,     // bora_cumpade
		Main_function, // main
		Open_paren,    // (
		Close_paren,   // )
		Block_open,    // simbora
		Block_close,   // cabo
	}

	if len(tokens) != len(expectedTokens) {
		t.Fatalf("expected %d tokens, got %d", len(expectedTokens), len(tokens))
	}

	for i, expected := range expectedTokens {
		if tokens[i].Token != expected {
			t.Errorf("token[%d] = %d, want %d (lexema: %q)", i, tokens[i].Token, expected, tokens[i].Lexema)
		}
	}
}

func TestLexer_DeclaracaoComVirgula(t *testing.T) {
	input := "trem_di_numeru a, b uai "
	tokens := AnalisarArquivo(input)

	expectedTokens := []TabelaPalavras{
		Type_int,   // trem_di_numeru
		Identifier, // a
		Comma,      // ,
		Identifier, // b
		Stmt_end,   // uai
	}

	if len(tokens) != len(expectedTokens) {
		t.Fatalf("expected %d tokens, got %d", len(expectedTokens), len(tokens))
	}

	for i, expected := range expectedTokens {
		if tokens[i].Token != expected {
			t.Errorf("token[%d] = %d (%q), want %d", i, tokens[i].Token, tokens[i].Lexema, expected)
		}
	}
}

// ---------- Símbolos ----------

func TestLexer_Parenteses(t *testing.T) {
	tokens := AnalisarArquivo("() ")
	if len(tokens) != 2 {
		t.Fatalf("expected 2 tokens, got %d", len(tokens))
	}
	if tokens[0].Token != Open_paren {
		t.Errorf("token[0] = %d, want Open_paren", tokens[0].Token)
	}
	if tokens[1].Token != Close_paren {
		t.Errorf("token[1] = %d, want Close_paren", tokens[1].Token)
	}
}
