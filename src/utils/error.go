package utils

import (
	"fmt"
)

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
 *			erro(
 *				description = "Erro no qual o sistema deve apresentar",
 *				example = "Arquivo não encontrado"
 *			),
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
		fmt.Println(err.Error())
		return
	}
}