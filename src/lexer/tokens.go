package lexer

type TabelaPalavras int

const (
	// Data types
	type_int    TabelaPalavras = iota // 0  - trem_di_numeru   → int
	type_float                        // 1  - trem_cum_virgula → float
	type_string                       // 2  - trem_discrita    → string
	type_bool                         // 3  - trem_discolhe    → boolean
	type_char                         // 4  - trosso           → char

	// Conditionals
	conditional_if     // 5  - uai_se    → if
	conditional_else   // 6  - uai_senao → else
	loop_for           // 7  - roda_esse_trem      → for
	loop_while         // 8  - enquanto_tiver_trem → while
	conditional_switch // 9  - dependenu → switch
	conditional_case   // 10 - du_casu   → case

	// Functions & return
	func_return   // 11 - ta_bao       → return
	loop_break    // 12 - para_o_trem  → break
	loop_continue // 13 - toca_o_trem  → continue
	func_decl     // 14 - bora_cumpade → function / main

	// Symbols
	open_paren  // 15 - abre_parentese  → (
	close_paren // 16 - fecha_parentese → )

	// Boolean literals
	literal_true  // 17 - eh     → true
	literal_false // 18 - num_eh → false

	// Block delimiters
	block_open  // 19 - simbora → começo de bloco  { }
	block_close // 20 - cabo    → fim de bloco    { }
	open_brace  // 21 - abre_chave  → {
	close_brace // 22 - fecha_chave → }
	open_quote  // 23 - abre_aspas  → "
	close_quote // 24 - fecha_aspas → "

	// Punctuation
	comma    // 25 - virgula → ,
	stmt_end // 26 - uai     → ; (fim da instrução)

	// Relational operators
	op_lt  // 27 - <  → menor que
	op_gt  // 28 - >  → maior que
	op_lte // 29 - <= → menor ou igual
	op_gte // 30 - >= → maior ou igual

	// Assignment & equality
	op_assign // 31 - fica_assim_entao → = (atribuição)
	op_neq    // 32 - neh_nada         → != (diferente de)
	op_eq     // 33 - mema_coisa       → == (igual a)

	// Logical operators
	op_or  // 34 - quarque_um → or
	op_not // 35 - vam_marca  → not
	op_and // 36 - tamem      → and
	op_xor // 37 - um_o_oto   → xor

	// Arithmetic operators
	op_add     // 38 - +    → adição
	op_sub     // 39 - -    → subtração
	op_mul     // 40 - veiz → multiplicação  (*)
	op_div     // 41 - sob  → divisão        (/)
	op_mod     // 42 - %    → módulo
	op_int_div // 43 - /    → divisão inteira (//)

	// I/O
	io_scan  // 44 - xove        → scan  / input
	io_print // 45 - oia_proce_ve → print / output

	// Literals & tokens
	literal_string // 46 - conteúdo string
	comment_line   // 47 - // comentário de linha

	comment_block_open  // 48 - causo       → /* comentário de bloco
	comment_block_close // 49 - fim_do_causo → */ comentário de bloco

	literal_int   // 50 - conteúdo inteiro
	literal_hex   // 51 - conteúdo hexadecimal (0x...)
	literal_oct   // 52 - conteúdo octal (0...)
	literal_float // 53 - conteúdo float

	identifier    // 54 - variável / identificador
	lexical_error // 55 - token inválido
)

var PalavrasReservadas = map[string]TabelaPalavras{
	// Data types
	"trem_di_numeru":   type_int,
	"trem_cum_virgula": type_float,
	"trem_discrita":    type_string,
	"trem_discolhe":    type_bool,
	"trosso":           type_char,

	// Conditionals
	"uai_se":    conditional_if,
	"uai_senao": conditional_else,
	"dependenu": conditional_switch,
	"du_casu":   conditional_case,

	// Loops
	"roda_esse_trem":      loop_for,
	"enquanto_tiver_trem": loop_while,
	"para_o_trem":         loop_break,
	"toca_o_trem":         loop_continue,

	// Functions & return
	"bora_cumpade": func_decl,
	"ta_bao":       func_return,

	// Symbols
	"abre_parentese":  open_paren,
	"fecha_parentese": close_paren,

	// Boolean literals
	"eh":     literal_true,
	"num_eh": literal_false,

	// Block delimiters
	"simbora":     block_open,
	"cabo":        block_close,
	"abre_chave":  open_brace,
	"fecha_chave": close_brace,
	"abre_aspas":  open_quote,
	"fecha_aspas": close_quote,

	// Punctuation
	"virgula": comma,
	"uai":     stmt_end,

	// Relational operators
	"<":  op_lt,
	">":  op_gt,
	"<=": op_lte,
	">=": op_gte,

	// Assignment & equality
	"fica_assim_entao": op_assign,
	"neh_nada":         op_neq,
	"mema_coisa":       op_eq,

	// Logical operators
	"quarque_um": op_or,
	"vam_marca":  op_not,
	"tamem":      op_and,
	"um_o_oto":   op_xor,

	// Arithmetic operators
	"+":    op_add,
	"-":    op_sub,
	"veiz": op_mul,
	"sob":  op_div,
	"%":    op_mod,
	"/":    op_int_div,

	// I/O
	"xove":         io_scan,
	"oia_proce_ve": io_print,

	// Comments
	"causo":        comment_block_open,
	"fim_do_causo": comment_block_close,
}
