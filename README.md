<a id="readme-topo"></a>

<h1 align='center'>
  Minerês Interpreter
</h1>

<div align='center'>

[![SO][Ubuntu-badge]][Ubuntu-url]
[![IDE][vscode-badge]][vscode-url]
[![Go][Golang-badge]][Golang-url]

<b>
  Guilherme Alvarenga de Azevedo<br>
  Joaquim Cezar Santana da Cruz<br>
  Luan Gonçalves dos Santos<br>
</b>
  
<br>
Compiladores <br>
Eduardo Gabriel Reis Miranda<br>
Engenharia de Computação <br>
CEFET-MG Campus V <br>
2026/1

</div>

<details>
  <summary>
  <b style='font-size: 15px'>
    📑 Sumário
  </b>
  </summary>
  <ol>
    <li><a href="#-o-projeto">📚 O Projeto</a></li>
    <li><a href="#estrutura-e-explicacao-do-codigo">🏗️ Estrutura e Explicação do Código</a></li>
    <li><a href="#instalando">🔨 Instalando</a></li>
    <li><a href="#-ambiente-de-compilação-e-execução">🧪 Ambiente de Compilação e Execução</a></li>
    <li><a href="#-contato">📨 Contato</a></li>
    <li><a href="#referencias">📚 Referências</a></li>
  </ol>
</details>

## 📚 O Projeto

Neste repositório você encontrará o código fonte do projeto. O projeto foi desenvolvido em Go. Este trabalho também tem a produção de um TXT para relatar cada alteração nas etapas do projeto, que está disponível em [`release_notes.txt`](release_notes.txt).

<a id="estrutura-e-explicacao-do-codigo"></a>

## 🏗️ Estrutura e Explicação do Código

A arquitetura do projeto é dividida em módulos que separam a análise léxica, a análise sintática e as ferramentas utilitárias.

### ESTRUTURA

De uma forma compacta e organizada, os arquivos e diretórios estão dispostos da seguinte forma:

```.
mineres-interpreter/
├── data/
│   ├── lexerValidation*.uai
│   ├── parserValidation*.uai
│   ├── mineres.gmr
│   └── main.uai
├── src/
│   ├── lexer/
│   │   ├── lexer.go
│   │   ├── tabela.go
│   │   └── tokens.go
│   ├── parser/
│   │   └── parser.go
│   └── utils/
│       ├── error.go
│       └── file.go
├── .gitignore
├── go.mod
├── main.go
├── README.md
└── release_notes.txt
```

### 🔍 LEXER (Análise Léxica)

Responsável por transformar o código-fonte (texto bruto) em uma lista de unidades lógicas chamadas **Tokens**.

**`tokens.go`**: Define o "vocabulário" da linguagem Minerês por meio da enumeração `TabelaPalavras`. Mapeia termos regionais como `trem_di_numeru` para tipos de dados e `uai` para o fim de instruções.

**`tabela.go`**: Estrutura a saída do analisador através da struct `Tupla`, que armazena o lexema, o token e sua localização exata (linha e coluna) no código.

**`lexer.go`**: Contém o "motor" do analisador que utiliza expressões regulares para identificar números hexadecimais, octais, inteiros e floats. Também gera estados para strings e comentários de bloco (`causo ... fim_do_causo`).

### 📐 PARSER (Análise Sintática)

Responsável por verificar se a sequência de tokens segue as regras gramaticais da linguagem.

**`parser.go`**: Implementa um **Analisador Descendente Recursivo** que consome os tokens e valida estruturas. Ele processa laços `roda_esse_trem` (for), condicionais `uai_se` (if), declarações e a precedência de operadores.

### 🛠️ UTILS (Utilitários)

Funções de suporte para o funcionamento do sistema e tratamento de erros.

**`error.go`**: Centraliza o tratamento de exceções, fornecendo funções como `ThrowLexerException` e `ThrowParserException`. Exibe erros detalhados indicando a linha e a coluna exata do problema.

**`file.go`**: Gerencia a manipulação de arquivos, incluindo funções para verificar a existência e realizar a leitura completa de scripts `.uai`.

### 🚀 RAÍZ (Execução e Configuração)

**`main.go`**: Ponto de entrada do interpretador que coordena a leitura do arquivo, a geração de tokens pelo Lexer e a validação final pelo Parser.

**`go.mod`**: Define o nome do módulo do projeto e especifica a versão do Go (1.26.1) necessária para a compilação.

## Instalando

Para instalar o projeto, siga os passos abaixo:

<div align="justify">
  Com o ambiente preparado, os seguintes passos são para a instalação, compilação e execução do programa localmente:

1. Clone o repositório no diretório desejado:

```console
git clone https://github.com/alvarengazv/mineres-interpreter.git
cd mineres-interpreter
```

2. Instale as dependências do [Go](https://go.dev/dl/) para o seu SO e certifique-se de que a versão é compatível com o projeto (1.26.1 ou superior):

```console
go version
```

3. Execute o programa:

```console
  go run .
```

</div>

<p align="right">(<a href="#readme-topo">voltar ao topo</a>)</p>

## 🧪 Ambiente de Compilação e Execução

<div align="justify">

O trabalho foi desenvolvido e testado em várias configurações de hardware. Podemos destacar algumas configurações de Sistema Operacional e Compilador, pois as demais configurações não influenciam diretamente no desempenho do programa.

</div>

<div align='center'>

[![SO][Ubuntu-badge]][Ubuntu-url]
[![IDE][vscode-badge]][vscode-url]
[![Go][Golang-badge]][Golang-url]

|        _SO_        |     _Compilador_     |
| :----------------: | :------------------: |
| Ubuntu 24.04.4 LTS | go1.26.1 linux/amd64 |

</div>

> [!IMPORTANT]
> Para que os testes tenham validade, considere as especificações
> do ambiente de compilação e execução do programa.

<p align="right">(<a href="#readme-topo">voltar ao topo</a>)</p>

## 📨 Contato

<div align="center">
  <br><br>
     <i>Guilherme Alvarenga de Azevedo - Graduando - 7º Período de Engenharia de Computação @ CEFET-MG</i>
  <br><br>
  
  [![Gmail][gmail-badge]][gmail-autor1]
  [![Linkedin][linkedin-badge]][linkedin-autor1]
  [![Telegram][telegram-badge]][telegram-autor1]
  
  
  <br><br>
     <i>Joaquim Cezar Santana da Cruz - Graduando - 7º Período de Engenharia de Computação @ CEFET-MG</i>
  <br><br>
  
  [![Gmail][gmail-badge]][gmail-autor2]
  [![Linkedin][linkedin-badge]][linkedin-autor2]
  [![Telegram][telegram-badge]][telegram-autor2]

<br><br>
<i>Luan Gonçalves dos Santos - Graduando - Engenharia de Computação @ CEFET-MG</i>
<br><br>

[![Gmail][gmail-badge]][gmail-autor3]
[![Linkedin][linkedin-badge]][linkedin-autor3]
[![Telegram][telegram-badge]][telegram-autor3]

</div>

<p align="right">(<a href="#readme-topo">voltar ao topo</a>)</p>

<a id="referencias">📚 Referências</a>

1. AZEVEDO, Guilherme A. CRUZ, Joaquim C. S. SANTOS, Luan G. . **Compiladores**: Lexer - Minerês Interpreter. 2026. Disponível em: [https://github.com/alvarengazv/mineres-interpreter](https://github.com/alvarengazv/mineres-interpreter) Acesso em: 05 abr. 2026.

[vscode-badge]: https://img.shields.io/badge/Visual%20Studio%20Code-0078d7.svg?style=for-the-badge&logo=visual-studio-code&logoColor=white
[vscode-url]: https://code.visualstudio.com/docs/?dv=linux64_deb
[make-badge]: https://img.shields.io/badge/_-MAKEFILE-427819.svg?style=for-the-badge
[make-url]: https://www.gnu.org/software/make/manual/make.html
[cpp-badge]: https://img.shields.io/badge/c++-%2300599C.svg?style=for-the-badge&logo=c%2B%2B&logoColor=white
[cpp-url]: https://en.cppreference.com/w/cpp
[trabalho-url]: https://drive.google.com/file/d/1-IHbGaA1BIC6_CMBydOC-NbV2bCERc8r/view?usp=sharing
[github-prof]: https://github.com/mpiress
[main-ref]: src/main.cpp
[branchAMM-url]: https://github.com/alvarengazv/trabalhosAEDS1/tree/AlgoritmosMinMax
[makefile]: ./makefile
[bash-url]: https://www.hostgator.com.br/blog/o-que-e-bash/
[lenovo-badge]: https://img.shields.io/badge/lenovo%20laptop-E2231A?style=for-the-badge&logo=lenovo&logoColor=white
[ubuntu-badge]: https://img.shields.io/badge/Ubuntu-E95420?style=for-the-badge&logo=ubuntu&logoColor=white
[Ubuntu-url]: https://ubuntu.com/
[ryzen5500-badge]: https://img.shields.io/badge/AMD%20Ryzen_5_5500U-ED1C24?style=for-the-badge&logo=amd&logoColor=white
[ryzen3500-badge]: https://img.shields.io/badge/AMD%20Ryzen_5_3500X-ED1C24?style=for-the-badge&logo=amd&logoColor=white
[windows-badge]: https://img.shields.io/badge/Windows-0078D6?style=for-the-badge&logo=windows&logoColor=white
[gcc-badge]: https://img.shields.io/badge/GCC-5C6EB8?style=for-the-badge&logo=gnu&logoColor=white
[Golang-badge]: https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white
[Golang-url]: https://go.dev/
[linkedin-autor1]: https://www.linkedin.com/in/guilherme-alvarenga-de-azevedo-959474201/
[telegram-autor1]: https://t.me/alvarengazv
[gmail-autor1]: mailto:gui.alvarengas234@gmail.com
[linkedin-autor2]: https://www.linkedin.com/in/joaquim-cruz-b760bb350/
[telegram-autor2]: https://t.me/joaquim1333
[gmail-autor2]: mailto:joaquimcezar930@gmail.com
[linkedin-autor3]: https://www.linkedin.com/in/luan-santos-9bb01920b/
[telegram-autor3]: https://t.me/LuanLuL_SO7
[gmail-autor3]: mailto:luanzinholulus@gmail.com
[linkedin-badge]: https://img.shields.io/badge/-LinkedIn-0077B5?style=for-the-badge&logo=Linkedin&logoColor=white
[telegram-badge]: https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white
[gmail-badge]: https://img.shields.io/badge/-Gmail-D14836?style=for-the-badge&logo=Gmail&logoColor=white
