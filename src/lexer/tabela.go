package lexer

import "fmt"

type Tupla struct {
	Lexema string
	Token  TabelaPalavras
	Linha  int
	Coluna int
}

/**
* TuplaToString(resumo = "Printar uma tupla",
*	 	Parâmetros = {
*	 		tupla(
*					description = "struct do tipo Tupla",
*					example = {
 *						lexama: "trem_cum_virgula",
*						token: 1,
*						linha: 2,
*						coluna: 14
*					}
*				)
*		},
*		Retorno = {}
* )
*/
func TuplaToString(t Tupla) {
	fmt.Printf("Lexema: %s | Token: %d | Linha: %d | Coluna: %d\n",
		t.Lexema,
		t.Token,
		t.Linha,
		t.Coluna,
	)
}

/**
 * ListTuplaToString(resumo = "Printar uma lista de tuplas",
 *	 	Parâmetros = {
 *	 		lista(
 *				description = "List de structs do tipo Tupla",
 *				example = [
 *					{
 *						lexama: "trem_cum_virgula",
 *						token: 1,
 *						linha: 2,
 *						coluna: 14
 *					},
 *					{
 *						lexama: "ta_bao",
 *						token: 11,
 *						linha: 4,
 *						coluna: 6
 *					}
 *				]
 *			)
 *		},
 *		Retorno = {}
 * )
 */
func ListTuplaToString(lista []Tupla) {
	for _, item := range lista {
		TuplaToString(item)
	}
}
