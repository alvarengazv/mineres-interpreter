package utils

import (
	"fmt"
	"os"
)

// ExitOnError controla o comportamento em caso de erro.
// Em produção (true): imprime o erro e chama os.Exit(1).
// Em testes (false): faz panic com CompilerError para que o teste possa usar recover().
var ExitOnError = true

// CompilerError é o tipo usado para panic em modo de teste,
// permitindo que os testes capturem e inspecionem erros do compilador.
type CompilerError struct {
	Message string
}

/**
 * ThrowInterpreterException(Resumo = "Criar e mostrar erros do interpreter",
 *		Parâmetros = {
 *			mensagem(
 *				description = "Erro no qual o sistema deve apresentar",
 *				example = "Token esperado: ';'"
 *			),
 *			line(
 *				description = "Linha na qual o erro aconteceu"
 *				example = 6
 *			),
 *			column(
 *				description = "Coluna na qual o erro aconteceu"
 *				example = 35
 *			)
 *		},
 *		Retorno = {}
 * )
 */
func ThrowInterpreterException(mensagem string, line, column int) {
	executeException(fmt.Errorf(`
Interpreter error on (%d::%d) => %s
`,
		line,
		column,
		mensagem,
	))
}


/**
 * ThrowLexerException(Resumo = "Criar e mostrar erros do lexer",
 *		Parâmetros = {
 *			mensagem(
 *				description = "Erro no qual o sistema deve apresentar",
 *				example = "Arquivo não encontrado"
 *			),
 *			line(
 *				description = "Linha na qual o erro aconteceu"
 *				example = 6
 *			),
 *			column(
 *				description = "Coluna na qual o erro aconteceu"
 *				example = 35
 *			)
 *		},
 *		Retorno = {}
 * )
 */
func ThrowLexerException(mensagem string, line, column int) {
	executeException(fmt.Errorf(`
Lexer parsed error on (%d::%d) => %s
`,
	line,
	column,
	mensagem,
))
}

/**
 * ThrowException(Resumo = "Criar e mostrar erros gerais",
 *		Parâmetros = {
 * 			arquivo(
 *				description = "arquivo que lançou o erro",
 *				example = "file.go"
 *			),
 *			function(
 *				description = "função que lançou o erro",
 *				example = "LerArquivo"
 *			),
 *			mensagem(
 *				description = "Erro no qual o sistema deve apresentar",
 *				example = "Arquivo não encontrado"
 *			)
 *		},
 *		Retorno = {}
 * )
 */
func ThrowException(arquivo, funcao, mensagem string) {
	executeException(fmt.Errorf(`
ThrowException::%s (%s) => %s
`,
	arquivo,
	funcao,
	mensagem,
))
}

/**
 * ThrowParserException(Resumo = "Criar e mostrar erros do parser",
 *		Parâmetros = {
 *			mensagem(
 *				description = "Erro no qual o sistema deve apresentar",
 *				example = "Token esperado: ';'"
 *			),
 *			line(
 *				description = "Linha na qual o erro aconteceu"
 *				example = 6
 *			),
 *			column(
 *				description = "Coluna na qual o erro aconteceu"
 *				example = 35
 *			)
 *		},
 *		Retorno = {}
 * )
 */
func ThrowParserException(mensagem string, line, column int) {
	executeException(fmt.Errorf(`
Parser error on (%d::%d) => %s
`,
		line,
		column,
		mensagem,
	))
}

/**
  * executeException(Resumo = "Lançar o erro no sistema",
  *		Parâmetros = {
  *			err(
  *				description = "Erro no qual o sistema deve apresentar",
  *				example = type(error)
  *			),
  *		},
  *		Retorno = {}
  * )
  */
func executeException(err error) {
	if err != nil {
		msg := err.Error()
		fmt.Fprintln(os.Stderr, msg)
		if ExitOnError {
			os.Exit(1)
		}
		panic(CompilerError{Message: msg})
	}
}