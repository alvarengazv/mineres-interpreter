package utils

import (
	"fmt"
	"os"
)

/**
 * ThrowException(Resumo = "Criar e mostrar erros do lexer",
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
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}