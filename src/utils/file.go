package utils

import (
	"fmt"
	"os"
)

/**
 * ArquivoExiste(resumo = "Verificar a existencia do arquivo .uai",
 *	 	Parâmetros = {
 *	 		caminho(
 *					description = "Caminho em que se encontra o arquivo .uai",
 *					example = "data/helloWolrd.uai"
 *				)
 *		},
 *		Retorno = {
 *			Se o arquivo existir = bool(true),
 *			Se o arquivo não existir ou ocorrer erro ao acessá-lo = bool(false)
 *		}
 * )
 */
func ArquivoExiste(caminho string) bool {
	_, err := os.Stat(caminho)
	return err == nil
}

/**
  * LerArquivo(Resumo = "Realizar a leitura do arquivo .uai",
  *		Parâmetros = {
  * 		caminho(
  *				description = "Caminho em que se encontra o arquivo .uai",
  *				example = "data/helloWolrd.uai"
  *			)
  *		},
  *		Retorno = {
  *			Conteúdo completo do arquivo = string
  *		}
  * )
  */
func LerArquivo(caminho string) (string) {
	if !ArquivoExiste(caminho) {
		ThrowException("file.go", "LerArquivo",
			fmt.Sprintf("arquivo %s não encontrado", caminho))
	}

	bytes, err := os.ReadFile(caminho)
	if err != nil {
		ThrowException("file.go", "LerArquivo",
			fmt.Sprintf("erro ao ler arquivo %s", caminho))
	}
	
	return string(bytes)
}