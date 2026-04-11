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
	Conditional_if     // 5  - uai_se    → if
	Conditional_else   // 6  - uai_senao → else
	Loop_for           // 7  - roda_esse_trem      → for
	Loop_while         // 8  - enquanto_tiver_trem → while
	Conditional_switch // 9  - dependenu → switch
	Conditional_case   // 10 - du_casu   → case

	// Functions & return
	Func_return   // 11 - ta_bao       → return
	Loop_break    // 12 - para_o_trem  → break
	Loop_continue // 13 - toca_o_trem  → continue
	Func_decl     // 14 - bora_cumpade → function / main

	// Symbols
	Open_paren  // 15 - abre_parentese  → (
	Close_paren // 16 - fecha_parentese → )

	// Boolean literals
	Literal_true  // 17 - eh     → true
	Literal_false // 18 - num_eh → false

	// Block delimiters
	Block_open  // 19 - simbora → começo de bloco  { }
	Block_close // 20 - cabo    → fim de bloco    { }
	Open_brace  // 21 - abre_chave  → {
	Close_brace // 22 - fecha_chave → }

	// Punctuation
	Comma    // 23 - virgula → ,
	Stmt_end // 24 - uai     → ; (fim da instrução)
	Colon    // 25 - dois pontos → :
	// Relational operators
	Op_lt  // 26 - <  → menor que
	Op_gt  // 27 - >  → maior que
	Op_lte // 28 - <= → menor ou igual
	Op_gte // 29 - >= → maior ou igual

	// Assignment & equality
	Op_assign // 30 - fica_assim_entao → = (atribuição)
	Op_neq    // 31 - neh_nada         → != (diferente de)
	Op_eq     // 32 - mema_coisa       → == (igual a)

	// Logical operators
	Op_or  // 33 - quarque_um → or
	Op_not // 34 - vam_marca  → not
	Op_and // 35 - tamem      → and
	Op_xor // 36 - um_o_oto   → xor

	// Arithmetic operators
	Op_add     // 37 - +    → adição
	Op_sub     // 38 - -    → subtração
	Op_mul     // 39 - veiz → multiplicação  (*)
	Op_div     // 40 - sob  → divisão        (/)
	Op_mod     // 41 - %    → módulo
	Op_int_div // 42 - /    → divisão inteira (//)

	// I/O
	Io_scan  // 43 - xove        → scan  / input
	Io_print // 44 - oia_proce_ve → print / output

	// Literals & tokens
	Literal_string // 45 - conteúdo string
	Literal_char   // 46 - conteúdo char
	Comment_line   // 47 - // comentário de linha

	Comment_block_open  // 48 - causo       → /* comentário de bloco
	Comment_block_close // 49 - fim_do_causo → */ comentário de bloco

	Literal_int   // 50 - conteúdo inteiro
	Literal_hex   // 51 - conteúdo hexadecimal (0x...)
	Literal_oct   // 52 - conteúdo octal (0...)
	Literal_float // 53 - conteúdo float

	Identifier    // 54 - variável / identificador
	Lexical_error // 55 - token inválido

	Stmt_end_for  // 56 - ; → fim de statement no for
	Main_function // 57 - main → função principal
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
	"ta_bao":       Func_return,
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
	// COMENTADOS: tratados diretamente pelo lexer (chars '{', '}', '"', etc.)
	// Se descomentadas, variáveis com esses nomes seriam mascaradas.
	// "abre_chave":          Open_brace,
	// "fecha_chave":         Close_brace,
	// "abre_aspas":          Open_quote,
	// "fecha_aspas":         Close_quote,
	// "abre_aspas_simples":  Open_squote,
	// "fecha_aspas_simples": Close_squote,

	// Punctuation
	// COMENTADO: ',' é tratado em tratarSimbolosEspeciais. Uma variável "virgula" seria mascarada.
	// "virgula": Comma,
	"uai": Stmt_end,
	";":   Stmt_end_for,

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
	// "+":  Op_add,
	// "-":  Op_sub,
	"veiz": Op_mul,
	"sob":  Op_div,
	// "%":  Op_mod,
	// "/":  Op_int_div,

	// COMENTADOS: '<' e '>' são tratados em tratarSimbolosEspeciais (inclui <= e >=).
	// "<":  Op_lt,
	// ">":  Op_gt,
	// "<=": Op_lte,
	// ">=": Op_gte,

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
	Type_int:    "trem_di_numeru",  	// 0  - trem_di_numeru   → int
	Type_float:  "trem_cum_virgula",	// 1  - trem_cum_virgula → float
	Type_string: "trem_discrita",		// 2  - trem_discrita    → string
	Type_bool:   "trem_discolhe",		// 3  - trem_discolhe    → bool
	Type_char:   "trosso",				// 4  - trosso           → char

	// Conditionals
	Conditional_if:     "uai_se",		// 5  - uai_se    → if
	Conditional_else:   "uai_senao",    // 6  - uai_senao → else
	Conditional_switch: "dependenu",	// 9  - dependenu → switch
	Conditional_case:   "du_casu",		// 10 - du_casu   → case

	// Loops
	Loop_for:           "roda_esse_trem",		// 7  - roda_esse_trem      → for
	Loop_while:         "enquanto_tiver_trem",	// 8  - enquanto_tiver_trem → while

	// Functions & return
	Func_decl:     "bora_cumpade",		// 14 - bora_cumpade → function declaration
	Func_return:   "ta_bao",			// 11 - ta_bao       → return
	Main_function: "main",				// 57 - main         → main function
	Loop_break:    "para_o_trem",		// 12 - para_o_trem  → break
	Loop_continue: "toca_o_trem",		// 13 - toca_o_trem  → continue


	// Boolean literals
	Literal_true:  "eh",				// 17 - eh     → true
	Literal_false: "num_eh",			// 18 - num_eh → false

	// Block delimiters
	Block_open:  "simbora",			// 15 - simbora → {
	Block_close: "cabo",				// 16 - cabo → }
	Open_brace:  "{",		// 21 - abre_chave → {
	Close_brace: "}",		// 22 - fecha_chave → }
	
	// Parentheses
	Open_paren:  "(",		// 19 - abre_parentese  → (
	Close_paren: ")",		// 20 - fecha_parentese → )

	// Punctuation
	Stmt_end:     "uai",	// 24 - uai     → ; (fim da instrução)
	Stmt_end_for: ";",		// 56 - ; → fim de statement no for
	Comma:        ",",		// 23 - virgula → ,
	Colon:        ":",		// 25 - dois pontos → :

	// Assignment & equality
	Op_assign: "fica_assim_entao", 		// 30 - fica_assim_entao → = (atribuição)
	Op_neq:    "neh_nada",				// 31 - neh_nada         → != (diferente de)
	Op_eq:     "mema_coisa",			// 32 - mema_coisa       → == (igual a)

	// Logical operators
	Op_or:  "quarque_um",				// 33 - quarque_um → or
	Op_not: "vam_marca",				// 34 - vam_marca  → not
	Op_and: "tamem",					// 35 - tamem      → and
	Op_xor: "um_o_oto",					// 36 - um_o_oto   → xor

	// Arithmetic operators (only words, not symbols)
	Op_add: "+",						// 37 - + → +
	Op_sub: "-",						// 38 - - → -
	Op_mul: "veiz",					// 39 - veiz → *
	Op_div: "sob",						// 40 - sob → /
	Op_mod: "%",						// 41 - % → %
	Op_int_div: "/",					// 42 - / → /
	
	// Relational operators
	Op_lt:  "<",						// 26 - <  → menor que
	Op_gt:  ">",						// 27 - >  → maior que
	Op_lte: "<=",						// 28 - <= → menor ou igual
	Op_gte: ">=",						// 29 - >= → maior ou igual

	// I/O
	Io_scan:  "xove",					// 43 - xove        → scan  / input
	Io_print: "oia_proce_ve",			// 44 - oia_proce_ve → print / output

	// Comments
	Comment_block_open:  "causo",		// 48 - causo       → /* comentário de bloco
	Comment_block_close: "fim_do_causo",// 49 - fim_do_causo → */ comentário de bloco

	// Literals & tokens
	Literal_string: "conteúdo string",	// 45 - conteúdo string
	Literal_char:   "conteúdo char",	// 46 - conteúdo char
	Comment_line:   "// comentário de linha", // 47 - // comentário de linha

	Literal_int:   "conteúdo inteiro",	// 50 - conteúdo inteiro
	Literal_hex:   "conteúdo hexadecimal (0x...)", // 51 - conteúdo hexadecimal (0x...)
	Literal_oct:   "conteúdo octal (0...)", // 52 - conteúdo octal (0...)
	Literal_float: "conteúdo float",	// 53 - conteúdo float

	Identifier:    "variável / identificador", // 54 - variável / identificador
	Lexical_error: "token inválido",	// 55 - token inválido


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