package parser

import "mineres-interpreter/src/lexer"

func Function(listTupla []lexer.Tupla) {
	for i := 0; i < len(listTupla); i++ {
		token := listTupla[i]
		switch i {
		case 0:
			if token.Token != lexer.Func_decl {
				// TO-DO: parser error
				return
			}
		case 1:
			if token.Token != lexer.Main_function {
				// TO-DO: parser error
				return
			}
		case 2:
			if token.Token != lexer.Open_paren {
				// TO-DO: parser error
				return
			}
		case 3:
			if token.Token != lexer.Close_paren {
				// TO-DO: parser error
				return
			}
		default:
			i = bloco_function(listTupla, i)
		}
	}
}

func bloco_function(listTupla []lexer.Tupla, i int) int {
	i1 := i
	for i := i; i < len(listTupla); i++ {
		token := listTupla[i]
		switch i {
		case i1:
			if token.Token != lexer.Block_open {
				// TO-DO: parser error
				return i
			}
		default:
			if token.Token != lexer.Block_close {
				i = stmtList_function(listTupla, i)
			} else {
				return i
			}
		}
	}
	return i
}

func stmtList_function(listTupla []lexer.Tupla, i int) int {
	for i := i; i < len(listTupla); i++ {
		token := listTupla[i]
		switch token.Token {
		case lexer.Loop_for:
			if token.Token != lexer.Block_open {
				// TO-DO: parser error
				return i
			}
		case lexer.Loop_while:
			if token.Token != lexer.Block_open {
				// TO-DO: parser error
				return i
			}
		case lexer.Conditional_if:
			if token.Token != lexer.Block_open {
				// TO-DO: parser error
				return i
			}
		case lexer.Conditional_case:
			if token.Token != lexer.Block_open {
				// TO-DO: parser error
				return i
			}
		case lexer.Block_open:
			if token.Token != lexer.Block_open {
				// TO-DO: parser error
				return i
			}
		case lexer.Loop_break:
			if token.Token != lexer.Block_open {
				// TO-DO: parser error
				return i
			}
		case lexer.Loop_continue:
			if token.Token != lexer.Block_open {
				// TO-DO: parser error
				return i
			}

		default:
			// TO-DO: <atrib> uai, <declaration>, uai
		}
	}
	return i
}

func stmt_function(listTupla []lexer.Tupla, i int) int {
	return i
}

func type_function(listTupla []lexer.Tupla, i int) int {
	return i
}
