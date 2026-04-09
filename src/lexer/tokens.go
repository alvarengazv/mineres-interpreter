package lexer

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

	// Relational operators
	Op_lt  // 25 - <  → menor que
	Op_gt  // 26 - >  → maior que
	Op_lte // 27 - <= → menor ou igual
	Op_gte // 28 - >= → maior ou igual

	// Assignment & equality
	Op_assign // 29 - fica_assim_entao → = (atribuição)
	Op_neq    // 30 - neh_nada         → != (diferente de)
	Op_eq     // 31 - mema_coisa       → == (igual a)

	// Logical operators
	Op_or  // 32 - quarque_um → or
	Op_not // 33 - vam_marca  → not
	Op_and // 34 - tamem      → and
	Op_xor // 35 - um_o_oto   → xor

	// Arithmetic operators
	Op_add     // 36 - +    → adição
	Op_sub     // 37 - -    → subtração
	Op_mul     // 38 - veiz → multiplicação  (*)
	Op_div     // 39 - sob  → divisão        (/)
	Op_mod     // 40 - %    → módulo
	Op_int_div // 41 - /    → divisão inteira (//)

	// I/O
	Io_scan  // 42 - xove        → scan  / input
	Io_print // 43 - oia_proce_ve → print / output

	// Literals & tokens
	Literal_string // 44 - conteúdo string
	Literal_char   // 45 - conteúdo char
	Comment_line   // 46 - // comentário de linha

	Comment_block_open  // 47 - causo       → /* comentário de bloco
	Comment_block_close // 48 - fim_do_causo → */ comentário de bloco

	Literal_int   // 49 - conteúdo inteiro
	Literal_hex   // 50 - conteúdo hexadecimal (0x...)
	Literal_oct   // 51 - conteúdo octal (0...)
	Literal_float // 52 - conteúdo float

	Identifier    // 53 - variável / identificador
	Lexical_error // 54 - token inválido

	Stmt_end_for  // 55 - ; → fim de statement no for
	Main_function // 56 - main → função principal
)

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
