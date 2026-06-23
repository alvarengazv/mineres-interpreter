package parser

import (
	"mineres-interpreter/src/lexer"
	"mineres-interpreter/src/utils"
	"testing"
)

// expectPanic executa fn e verifica se houve um panic do tipo utils.CompilerError.
func expectPanic(t *testing.T, fn func()) string {
	t.Helper()
	old := utils.ExitOnError
	utils.ExitOnError = false
	defer func() { utils.ExitOnError = old }()

	var msg string
	func() {
		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("expected panic but none occurred")
			}
			ce, ok := r.(utils.CompilerError)
			if !ok {
				t.Fatalf("expected CompilerError, got %T: %v", r, r)
			}
			msg = ce.Message
		}()
		fn()
	}()
	return msg
}

// parseProgram é um helper que tokeniza e parseia um programa Minerês completo.
func parseProgram(input string) []TuplaMicrocode {
	tokens := lexer.AnalisarArquivo(input + " ")
	p := NewParser(tokens)
	return p.ParserFunction()
}

// countOp conta quantas instruções de um dado tipo existem no microcode.
func countOp(code []TuplaMicrocode, op TabelaMicrocodes) int {
	count := 0
	for _, c := range code {
		if c.Operation == op {
			count++
		}
	}
	return count
}

// ---------- Programa Mínimo ----------

func TestParser_ProgramaMinimo(t *testing.T) {
	code := parseProgram("bora_cumpade main() simbora cabo")
	// Programa vazio deve gerar 0 microcodes
	if len(code) != 0 {
		t.Errorf("expected 0 microcodes for empty program, got %d", len(code))
	}
}

// ---------- Declaração de Variável ----------

func TestParser_DeclaracaoSimples(t *testing.T) {
	code := parseProgram("bora_cumpade main() simbora trem_di_numeru x uai cabo")
	attCount := countOp(code, Att)
	if attCount != 1 {
		t.Errorf("expected 1 Att instruction, got %d", attCount)
	}
}

func TestParser_DeclaracaoMultipla(t *testing.T) {
	code := parseProgram("bora_cumpade main() simbora trem_di_numeru a, b, c uai cabo")
	attCount := countOp(code, Att)
	if attCount != 3 {
		t.Errorf("expected 3 Att instructions, got %d", attCount)
	}
}

func TestParser_DeclaracaoTodosOsTipos(t *testing.T) {
	prog := `bora_cumpade main() simbora
		trem_di_numeru a uai
		trem_cum_virgula b uai
		trem_discrita c uai
		trem_discolhe d uai
		trosso e uai
	cabo`
	code := parseProgram(prog)
	attCount := countOp(code, Att)
	if attCount != 5 {
		t.Errorf("expected 5 Att instructions (one per type), got %d", attCount)
	}
}

// ---------- Atribuição ----------

func TestParser_AtribuicaoSimples(t *testing.T) {
	prog := `bora_cumpade main() simbora
		trem_di_numeru x uai
		x fica_assim_entao 5 uai
	cabo`
	code := parseProgram(prog)
	// 1 Att para declaração + 1 Att para atribuição
	attCount := countOp(code, Att)
	if attCount != 2 {
		t.Errorf("expected 2 Att instructions, got %d", attCount)
	}
}

// ---------- Expressões Aritméticas ----------

func TestParser_ExpressaoAritmetica(t *testing.T) {
	prog := `bora_cumpade main() simbora
		trem_di_numeru a, b, c uai
		a fica_assim_entao 10 uai
		b fica_assim_entao 5 uai
		c fica_assim_entao a + b uai
	cabo`
	code := parseProgram(prog)
	addCount := countOp(code, Add)
	if addCount != 1 {
		t.Errorf("expected 1 Add instruction, got %d", addCount)
	}
}

func TestParser_PrecedenciaOperadores(t *testing.T) {
	// a + b veiz c deve gerar Mul primeiro, depois Add (mul tem maior precedência)
	prog := `bora_cumpade main() simbora
		trem_di_numeru a, b, c, d uai
		a fica_assim_entao 2 uai
		b fica_assim_entao 3 uai
		c fica_assim_entao 4 uai
		d fica_assim_entao a + b veiz c uai
	cabo`
	code := parseProgram(prog)
	mulCount := countOp(code, Mul)
	addCount := countOp(code, Add)
	if mulCount != 1 {
		t.Errorf("expected 1 Mul instruction, got %d", mulCount)
	}
	if addCount != 1 {
		t.Errorf("expected 1 Add instruction, got %d", addCount)
	}
}

func TestParser_Parenteses(t *testing.T) {
	// (a + b) veiz c → Add primeiro, Mul depois
	prog := `bora_cumpade main() simbora
		trem_di_numeru a, b, c, d uai
		a fica_assim_entao 2 uai
		b fica_assim_entao 3 uai
		c fica_assim_entao 4 uai
		d fica_assim_entao (a + b) veiz c uai
	cabo`
	code := parseProgram(prog)
	if countOp(code, Add) != 1 || countOp(code, Mul) != 1 {
		t.Error("expected 1 Add + 1 Mul instruction")
	}
}

// ---------- Operadores Lógicos ----------

func TestParser_OperadoresLogicos(t *testing.T) {
	prog := `bora_cumpade main() simbora
		trem_discolhe a, b, c uai
		a fica_assim_entao eh uai
		b fica_assim_entao num_eh uai
		c fica_assim_entao a tamem b uai
	cabo`
	code := parseProgram(prog)
	if countOp(code, And) != 1 {
		t.Errorf("expected 1 And instruction")
	}
}

func TestParser_OperadorNot(t *testing.T) {
	prog := `bora_cumpade main() simbora
		trem_discolhe a, b uai
		a fica_assim_entao eh uai
		b fica_assim_entao vam_marca a uai
	cabo`
	code := parseProgram(prog)
	if countOp(code, Not) != 1 {
		t.Errorf("expected 1 Not instruction")
	}
}

// ---------- Operador Unário ----------

func TestParser_UnarioNegativo(t *testing.T) {
	prog := `bora_cumpade main() simbora
		trem_di_numeru a uai
		a fica_assim_entao -5 uai
	cabo`
	code := parseProgram(prog)
	if countOp(code, Uno) != 1 {
		t.Errorf("expected 1 Uno instruction")
	}
}

// ---------- If-Else ----------

func TestParser_IfSimples(t *testing.T) {
	prog := `bora_cumpade main() simbora
		uai_se(eh) simbora cabo
	cabo`
	code := parseProgram(prog)
	if countOp(code, If_eq) != 1 {
		t.Errorf("expected 1 If_eq instruction")
	}
	if countOp(code, Label) < 2 {
		t.Errorf("expected at least 2 Label instructions for if")
	}
}

func TestParser_IfElse(t *testing.T) {
	prog := `bora_cumpade main() simbora
		uai_se(eh) simbora cabo
		uai_senao simbora cabo
	cabo`
	code := parseProgram(prog)
	if countOp(code, If_eq) != 1 {
		t.Errorf("expected 1 If_eq instruction")
	}
	// if-else gera labelTrue, labelFalse e labelEndIf
	if countOp(code, Label) < 3 {
		t.Errorf("expected at least 3 Label instructions for if-else")
	}
}

// ---------- While ----------

func TestParser_While(t *testing.T) {
	prog := `bora_cumpade main() simbora
		enquanto_tiver_trem(eh) simbora cabo
	cabo`
	code := parseProgram(prog)
	if countOp(code, If_eq) != 1 {
		t.Errorf("expected 1 If_eq instruction for while condition")
	}
	if countOp(code, Jump) != 1 {
		t.Errorf("expected 1 Jump instruction for while loop-back")
	}
}

// ---------- For ----------

func TestParser_ForCompleto(t *testing.T) {
	prog := `bora_cumpade main() simbora
		trem_di_numeru i uai
		roda_esse_trem(i fica_assim_entao 0; i < 10; i fica_assim_entao i + 1) simbora cabo
	cabo`
	code := parseProgram(prog)
	if countOp(code, If_eq) != 1 {
		t.Errorf("expected 1 If_eq instruction for for condition")
	}
	if countOp(code, Jump) != 1 {
		t.Errorf("expected 1 Jump instruction for for loop-back")
	}
	if countOp(code, Lt) != 1 {
		t.Errorf("expected 1 Lt instruction for for condition")
	}
}

func TestParser_ForVazio(t *testing.T) {
	prog := `bora_cumpade main() simbora
		roda_esse_trem(;;) simbora cabo
	cabo`
	// Não deve crashar
	code := parseProgram(prog)
	if countOp(code, Jump) < 1 {
		t.Errorf("expected at least 1 Jump instruction for infinite for loop")
	}
}

// ---------- Switch-Case ----------

func TestParser_SwitchCase(t *testing.T) {
	prog := `bora_cumpade main() simbora
		trem_di_numeru x uai
		x fica_assim_entao 1 uai
		dependenu(x) simbora
			du_casu 1: oia_proce_ve("um") uai
			du_casu 2: oia_proce_ve("dois") uai
		cabo
	cabo`
	code := parseProgram(prog)
	eqCount := countOp(code, Eq)
	if eqCount != 2 {
		t.Errorf("expected 2 Eq instructions for 2 case comparisons, got %d", eqCount)
	}
}

func TestParser_SwitchCaseComDefault(t *testing.T) {
	prog := `bora_cumpade main() simbora
		trem_di_numeru x uai
		x fica_assim_entao 1 uai
		dependenu(x) simbora
			du_casu 1: oia_proce_ve("um") uai
			uai_so: oia_proce_ve("outro") uai
		cabo
	cabo`
	code := parseProgram(prog)
	if countOp(code, Eq) != 1 {
		t.Errorf("expected 1 Eq instruction for 1 case comparison")
	}
}

// ---------- Break e Continue ----------

func TestParser_BreakNoLoop(t *testing.T) {
	prog := `bora_cumpade main() simbora
		enquanto_tiver_trem(eh) simbora
			para_o_trem uai
		cabo
	cabo`
	code := parseProgram(prog)
	// Break deve gerar Jump para o labelEnd do loop
	jumpCount := countOp(code, Jump)
	if jumpCount < 2 {
		t.Errorf("expected at least 2 Jump instructions (loop-back + break), got %d", jumpCount)
	}
}

func TestParser_ContinueNoLoop(t *testing.T) {
	prog := `bora_cumpade main() simbora
		enquanto_tiver_trem(eh) simbora
			toca_o_trem uai
		cabo
	cabo`
	code := parseProgram(prog)
	jumpCount := countOp(code, Jump)
	if jumpCount < 2 {
		t.Errorf("expected at least 2 Jump instructions (loop-back + continue), got %d", jumpCount)
	}
}

// ---------- I/O ----------

func TestParser_Print(t *testing.T) {
	prog := `bora_cumpade main() simbora
		oia_proce_ve("oi") uai
	cabo`
	code := parseProgram(prog)
	if countOp(code, Call) != 1 {
		t.Errorf("expected 1 Call instruction for print")
	}
}

func TestParser_Scan(t *testing.T) {
	prog := `bora_cumpade main() simbora
		trem_di_numeru x uai
		xove(trem_di_numeru, x) uai
	cabo`
	code := parseProgram(prog)
	if countOp(code, Call) != 1 {
		t.Errorf("expected 1 Call instruction for scan")
	}
}

// ---------- Erros ----------

func TestParser_Erro_BreakForaDoLoop(t *testing.T) {
	msg := expectPanic(t, func() {
		parseProgram("bora_cumpade main() simbora para_o_trem uai cabo")
	})
	if msg == "" {
		t.Error("expected error for break outside loop")
	}
}

func TestParser_Erro_ContinueForaDoLoop(t *testing.T) {
	msg := expectPanic(t, func() {
		parseProgram("bora_cumpade main() simbora toca_o_trem uai cabo")
	})
	if msg == "" {
		t.Error("expected error for continue outside loop")
	}
}

func TestParser_Erro_VariavelNaoDeclarada(t *testing.T) {
	msg := expectPanic(t, func() {
		parseProgram("bora_cumpade main() simbora y fica_assim_entao 5 uai cabo")
	})
	if msg == "" {
		t.Error("expected error for undeclared variable")
	}
}

func TestParser_Erro_ReDeclaracao(t *testing.T) {
	msg := expectPanic(t, func() {
		parseProgram(`bora_cumpade main() simbora
			trem_di_numeru x uai
			trem_di_numeru x uai
		cabo`)
	})
	if msg == "" {
		t.Error("expected error for variable re-declaration")
	}
}

func TestParser_Erro_TipoIncompativel(t *testing.T) {
	msg := expectPanic(t, func() {
		parseProgram(`bora_cumpade main() simbora
			trem_di_numeru x uai
			x fica_assim_entao "string" uai
		cabo`)
	})
	if msg == "" {
		t.Error("expected error for type mismatch")
	}
}

func TestParser_Erro_TokenInesperadoAposMain(t *testing.T) {
	msg := expectPanic(t, func() {
		parseProgram("bora_cumpade main() simbora cabo uai")
	})
	if msg == "" {
		t.Error("expected error for unexpected token after main")
	}
}

func TestParser_Erro_IfCondicaoNaoBool(t *testing.T) {
	msg := expectPanic(t, func() {
		parseProgram(`bora_cumpade main() simbora
			trem_di_numeru x uai
			x fica_assim_entao 5 uai
			uai_se(x) simbora cabo
		cabo`)
	})
	if msg == "" {
		t.Error("expected error for non-boolean if condition")
	}
}

// ---------- Operadores Relacionais ----------

func TestParser_OperadoresRelacionais(t *testing.T) {
	tests := []struct {
		name string
		op   string
		code TabelaMicrocodes
	}{
		{"igual", "mema_coisa", Eq},
		{"diferente", "neh_nada", Neq},
		{"menor", "<", Lt},
		{"maior", ">", Gt},
		{"menor_igual", "<=", Lte},
		{"maior_igual", ">=", Gte},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prog := `bora_cumpade main() simbora
				trem_di_numeru a, b uai
				trem_discolhe c uai
				a fica_assim_entao 5 uai
				b fica_assim_entao 3 uai
				c fica_assim_entao a ` + tt.op + ` b uai
			cabo`
			code := parseProgram(prog)
			if countOp(code, tt.code) != 1 {
				t.Errorf("expected 1 %v instruction", tt.code)
			}
		})
	}
}
