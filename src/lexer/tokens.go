package lexer

import "fmt"

type TabelaPalavras int

const (
	// Data types
	Type_int    TabelaPalavras = iota // 0  - trem_di_numeru   → int
	Type_float                        // 1  - trem_cum_virgula → float
	Type_string                       // 2  - trem_discrita    → string
	Type_bool                         // 3  - trem_discolhe    → boolean
	Type_char                         // 4  - trosso           → char

	// Conditionals
	Conditional_if      // 5  - uai_se    → if
	Conditional_else    // 6  - uai_senao → else
	Loop_for            // 7  - roda_esse_trem      → for
	Loop_while          // 8  - enquanto_tiver_trem → while
	Conditional_switch  // 9  - dependenu → switch
	Conditional_case    // 10 - du_casu   → case
	Conditional_default // 11 - uai_so   → default

	// Functions & return
	Func_return   // 12 - ta_bao       → return
	Loop_break    // 13 - para_o_trem  → break
	Loop_continue // 14 - toca_o_trem  → continue
	Func_decl     // 15 - bora_cumpade → function / main

	// Symbols
	Open_paren  // 16 - abre_parentese  → (
	Close_paren // 17 - fecha_parentese → )

	// Boolean literals
	Literal_true  // 18 - eh     → true
	Literal_false // 19 - num_eh → false

	// Block delimiters
	Block_open  // 20 - simbora → começo de bloco  { }
	Block_close // 21 - cabo    → fim de bloco    { }
	Open_brace  // 22 - abre_chave  → {
	Close_brace // 23 - fecha_chave → }

	// Punctuation
	Comma    // 24 - virgula → ,
	Stmt_end // 25 - uai     → ; (fim da instrução)
	Colon    // 26 - dois pontos → :
	// Relational operators
	Op_lt  // 27 - <  → menor que
	Op_gt  // 28 - >  → maior que
	Op_lte // 29 - <= → menor ou igual
	Op_gte // 30 - >= → maior ou igual

	// Assignment & equality
	Op_assign // 31 - fica_assim_entao → = (atribuição)
	Op_neq    // 32 - neh_nada         → != (diferente de)
	Op_eq     // 33 - mema_coisa       → == (igual a)

	// Logical operators
	Op_or  // 34 - quarque_um → or
	Op_not // 35 - vam_marca  → not
	Op_and // 36 - tamem      → and
	Op_xor // 37 - um_o_oto   → xor

	// Arithmetic operators
	Op_add     // 38 - +    → adição
	Op_sub     // 39 - -    → subtração
	Op_mul     // 40 - veiz → multiplicação  (*)
	Op_div     // 41 - sob  → divisão        (/)
	Op_mod     // 42 - %    → módulo
	Op_int_div // 43 - /    → divisão inteira (//)

	// I/O
	Io_scan  // 44 - xove        → scan  / input
	Io_print // 45 - oia_proce_ve → print / output

	// Literals & tokens
	Literal_string // 46 - conteúdo string
	Literal_char   // 47 - conteúdo char
	Comment_line   // 48 - // comentário de linha

	Comment_block_open  // 49 - causo       → /* comentário de bloco
	Comment_block_close // 50 - fim_do_causo → */ comentário de bloco

	Literal_int   // 51 - conteúdo inteiro
	Literal_hex   // 52 - conteúdo hexadecimal (0x...)
	Literal_oct   // 53 - conteúdo octal (0...)
	Literal_float // 54 - conteúdo float

	Identifier    // 55 - variável / identificador
	Lexical_error // 56 - token inválido

	Stmt_end_for  // 57 - ; → fim de statement no for
	Main_function // 58 - main → função principal
)

// Pro Lexer
var PalavrasReservadas = map[string]TabelaPalavras{
	// Data types
	"trem_di_numeru":   Type_int,
	"trem_cum_virgula": Type_float,
	"trem_discrita":    Type_string,
	"trem_discolhe":    Type_bool,
	"trosso":           Type_char,

	// Conditionals
	"uai_se":    Conditional_if,
	"uai_senao": Conditional_else,
	"dependenu": Conditional_switch,
	"du_casu":   Conditional_case,

	// Loops
	"roda_esse_trem":      Loop_for,
	"enquanto_tiver_trem": Loop_while,
	"para_o_trem":         Loop_break,
	"toca_o_trem":         Loop_continue,

	// Functions & return
	"bora_cumpade": Func_decl,
	// "ta_bao":       Func_return,
	"main":         Main_function,

	// Symbols — COMENTADOS: tratados diretamente em tratarSimbolosEspeciais, nunca chegam ao buffer.
	// Se descomentadas, nomes de variáveis como "abre_parentese" seriam confundidos com tokens.
	// "abre_parentese":  Open_paren,
	// "fecha_parentese": Close_paren,

	// Boolean literals
	"eh":     Literal_true,
	"num_eh": Literal_false,

	// Block delimiters
	"simbora": Block_open,
	"cabo":    Block_close,
	"uai_so":  Conditional_default,
	// COMENTADOS: tratados diretamente pelo lexer (chars '{', '}', '"', etc.)
	// Se descomentadas, variáveis com esses nomes seriam mascaradas.
	"{": Open_brace,
	"}": Close_brace,
	"(": Open_paren,
	")": Close_paren,
	// "abre_aspas":          Open_quote,
	// "fecha_aspas":         Close_quote,
	// "abre_aspas_simples":  Open_squote,
	// "fecha_aspas_simples": Close_squote,

	// Punctuation
	// COMENTADO: ',' é tratado em tratarSimbolosEspeciais. Uma variável "virgula" seria mascarada.
	",": Comma,
	"uai": Stmt_end,
	";":   Stmt_end_for,
	":":   Colon,

	// Assignment & equality
	"fica_assim_entao": Op_assign,
	"neh_nada":         Op_neq,
	"mema_coisa":       Op_eq,

	// Logical operators
	"quarque_um": Op_or,
	"vam_marca":  Op_not,
	"tamem":      Op_and,
	"um_o_oto":   Op_xor,

	// Arithmetic operators (only words, not symbols)
	// COMENTADOS: tratados em tratarSimbolosEspeciais, nunca chegam ao buffer.
	"+":  Op_add,
	"-":  Op_sub,
	"veiz": Op_mul,
	"sob":  Op_div,
	"%":  Op_mod,
	"/":  Op_int_div,

	// COMENTADOS: '<' e '>' são tratados em tratarSimbolosEspeciais (inclui <= e >=).
	"<":  Op_lt,
	">":  Op_gt,
	"<=": Op_lte,
	">=": Op_gte,

	// I/O
	"xove":         Io_scan,
	"oia_proce_ve": Io_print,

	// Comments
	"causo":        Comment_block_open,
	"fim_do_causo": Comment_block_close,
}

// Pro parser, para encontrar o lexema esperado da linguagem minerês (palavras ou símbolos), a partir do número do token lido
var PalavrasReservadasReverso = map[TabelaPalavras]string{
	// Data types
	Type_int:    "trem_di_numeru",   // 0  - trem_di_numeru   → int
	Type_float:  "trem_cum_virgula", // 1  - trem_cum_virgula → float
	Type_string: "trem_discrita",    // 2  - trem_discrita    → string
	Type_bool:   "trem_discolhe",    // 3  - trem_discolhe    → bool
	Type_char:   "trosso",           // 4  - trosso           → char

	// Conditionals
	Conditional_if:      "uai_se",    // 5  - uai_se    → if
	Conditional_else:    "uai_senao", // 6  - uai_senao → else
	Conditional_switch:  "dependenu", // 9  - dependenu → switch
	Conditional_case:    "du_casu",   // 10 - du_casu   → case
	Conditional_default: "uai_so",    // 11 - uai_so    → default

	// Loops
	Loop_for:   "roda_esse_trem",      // 7  - roda_esse_trem      → for
	Loop_while: "enquanto_tiver_trem", // 8  - enquanto_tiver_trem → while

	// Functions & return
	Func_decl:     "bora_cumpade", // 15 - bora_cumpade → function declaration
	Func_return:   "ta_bao",       // 12 - ta_bao       → return
	Main_function: "main",         // 58 - main         → main function
	Loop_break:    "para_o_trem",  // 13 - para_o_trem  → break
	Loop_continue: "toca_o_trem",  // 14 - toca_o_trem  → continue

	// Boolean literals
	Literal_true:  "eh",     // 18 - eh     → true
	Literal_false: "num_eh", // 19 - num_eh → false

	// Block delimiters
	Block_open:  "simbora", // 20 - simbora → {
	Block_close: "cabo",    // 21 - cabo → }
	Open_brace:  "{",       // 22 - abre_chave → {
	Close_brace: "}",       // 23 - fecha_chave → }

	// Parentheses
	Open_paren:  "(", // 16 - abre_parentese  → (
	Close_paren: ")", // 17 - fecha_parentese → )

	// Punctuation
	Stmt_end:     "uai", // 25 - uai     → ; (fim da instrução)
	Stmt_end_for: ";",   // 57 - ; → fim de statement no for
	Comma:        ",",   // 24 - virgula → ,
	Colon:        ":",   // 26 - dois pontos → :

	// Assignment & equality
	Op_assign: "fica_assim_entao", // 31 - fica_assim_entao → = (atribuição)
	Op_neq:    "neh_nada",         // 32 - neh_nada         → != (diferente de)
	Op_eq:     "mema_coisa",       // 33 - mema_coisa       → == (igual a)

	// Logical operators
	Op_or:  "quarque_um", // 34 - quarque_um → or
	Op_not: "vam_marca",  // 35 - vam_marca  → not
	Op_and: "tamem",      // 36 - tamem      → and
	Op_xor: "um_o_oto",   // 37 - um_o_oto   → xor

	// Arithmetic operators (only words, not symbols)
	Op_add:     "+",    // 38 - + → +
	Op_sub:     "-",    // 39 - - → -
	Op_mul:     "veiz", // 40 - veiz → *
	Op_div:     "sob",  // 41 - sob → /
	Op_mod:     "%",    // 42 - % → %
	Op_int_div: "/",    // 43 - / → /

	// Relational operators
	Op_lt:  "<",  // 27 - <  → menor que
	Op_gt:  ">",  // 28 - >  → maior que
	Op_lte: "<=", // 29 - <= → menor ou igual
	Op_gte: ">=", // 30 - >= → maior ou igual

	// I/O
	Io_scan:  "xove",         // 44 - xove        → scan  / input
	Io_print: "oia_proce_ve", // 45 - oia_proce_ve → print / output

	// Comments
	Comment_block_open:  "causo",        // 49 - causo       → /* comentário de bloco
	Comment_block_close: "fim_do_causo", // 50 - fim_do_causo → */ comentário de bloco

	// Literals & tokens
	Literal_string: "trem_discrita content",        // 46 - conteúdo string
	Literal_char:   "trosso content",          // 47 - conteúdo char
	Comment_line:   "inline_comment", // 48 - // comentário de linha

	Literal_int:   "trem_di_numeru content",             // 51 - conteúdo inteiro
	Literal_hex:   "trem_di_numeru_hex content", // 52 - conteúdo hexadecimal (0x...)
	Literal_oct:   "trem_di_numeru_oct content",        // 53 - conteúdo octal (0...)
	Literal_float: "trem_cum_virgula content",               // 54 - conteúdo float

	Identifier:    "variable / identifier", // 55 - variável / identificador
	Lexical_error: "invalid token",           // 56 - token inválido

}

// TabelaPalavrasFromInt converte um valor numérico para o token tipado.
// Retorna false quando não existe equivalência no mapa reverso.
func TabelaPalavrasFromInt(value int) (TabelaPalavras, bool) {
	token := TabelaPalavras(value)
	_, ok := PalavrasReservadasReverso[token]
	return token, ok
}

// String retorna o lexema minerês associado ao token.
func (t TabelaPalavras) String() string {
	if lexema, ok := PalavrasReservadasReverso[t]; ok {
		return lexema
	}

	return fmt.Sprintf("TabelaPalavras(%d)", int(t))
}
