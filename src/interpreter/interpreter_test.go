package interpreter

import (
	"bytes"
	"mineres-interpreter/src/lexer"
	"mineres-interpreter/src/parser"
	"mineres-interpreter/src/utils"
	"os"
	"strings"
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

// runProgram executa um programa Minerês completo e retorna o interpreter para inspecionar memória.
func runProgram(input string) *Interpreter {
	tokens := lexer.AnalisarArquivo(input + " ")
	p := parser.NewParser(tokens)
	code := p.ParserFunction()
	interp := NewInterpreter(code)
	interp.Run()
	return interp
}

// captureOutput executa um programa e captura stdout.
func captureOutput(input string) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	runProgram(input)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

// ---------- Atribuição ----------

func TestInterpreter_AtribuicaoInt(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_di_numeru x uai
		x fica_assim_entao 42 uai
	cabo`)
	if interp.memory["x"] != 42 {
		t.Errorf("x = %v, want 42", interp.memory["x"])
	}
}

func TestInterpreter_AtribuicaoFloat(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_cum_virgula x uai
		x fica_assim_entao 3.14 uai
	cabo`)
	if interp.memory["x"] != 3.14 {
		t.Errorf("x = %v, want 3.14", interp.memory["x"])
	}
}

func TestInterpreter_AtribuicaoString(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_discrita x uai
		x fica_assim_entao "hello" uai
	cabo`)
	if interp.memory["x"] != "hello" {
		t.Errorf("x = %v, want 'hello'", interp.memory["x"])
	}
}

func TestInterpreter_AtribuicaoBool(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_discolhe x uai
		x fica_assim_entao eh uai
	cabo`)
	if interp.memory["x"] != true {
		t.Errorf("x = %v, want true", interp.memory["x"])
	}
}

func TestInterpreter_AtribuicaoChar(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trosso x uai
		x fica_assim_entao 'z' uai
	cabo`)
	if interp.memory["x"] != "z" {
		t.Errorf("x = %v, want 'z'", interp.memory["x"])
	}
}

// ---------- Adição ----------

func TestInterpreter_AddIntInt(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_di_numeru a uai
		a fica_assim_entao 5 + 3 uai
	cabo`)
	if interp.memory["a"] != 8 {
		t.Errorf("a = %v, want 8", interp.memory["a"])
	}
}

func TestInterpreter_AddIntFloat(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_cum_virgula a uai
		a fica_assim_entao 5 + 3.2 uai
	cabo`)
	if interp.memory["a"] != 8.2 {
		t.Errorf("a = %v, want 8.2", interp.memory["a"])
	}
}

func TestInterpreter_AddStringString(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_discrita a uai
		a fica_assim_entao "ab" + "cd" uai
	cabo`)
	if interp.memory["a"] != "abcd" {
		t.Errorf("a = %v, want 'abcd'", interp.memory["a"])
	}
}

func TestInterpreter_AddStringChar(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_discrita a uai
		a fica_assim_entao "ab" + 'c' uai
	cabo`)
	if interp.memory["a"] != "abc" {
		t.Errorf("a = %v, want 'abc'", interp.memory["a"])
	}
}

func TestInterpreter_AddVarVar(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_di_numeru a, b, c uai
		a fica_assim_entao 10 uai
		b fica_assim_entao 20 uai
		c fica_assim_entao a + b uai
	cabo`)
	if interp.memory["c"] != 30 {
		t.Errorf("c = %v, want 30", interp.memory["c"])
	}
}

// ---------- Subtração ----------

func TestInterpreter_SubIntInt(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_di_numeru a uai
		a fica_assim_entao 10 - 3 uai
	cabo`)
	if interp.memory["a"] != 7 {
		t.Errorf("a = %v, want 7", interp.memory["a"])
	}
}

// ---------- Multiplicação ----------

func TestInterpreter_MulIntInt(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_di_numeru a uai
		a fica_assim_entao 4 veiz 5 uai
	cabo`)
	if interp.memory["a"] != 20 {
		t.Errorf("a = %v, want 20", interp.memory["a"])
	}
}

// ---------- Divisão ----------

func TestInterpreter_DivIntInt(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_cum_virgula a uai
		a fica_assim_entao 10 sob 4 uai
	cabo`)
	if interp.memory["a"] != 2.5 {
		t.Errorf("a = %v, want 2.5", interp.memory["a"])
	}
}

func TestInterpreter_DivInteira(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_di_numeru a uai
		a fica_assim_entao 10 / 3 uai
	cabo`)
	if interp.memory["a"] != 3 {
		t.Errorf("a = %v, want 3", interp.memory["a"])
	}
}

// ---------- Módulo ----------

func TestInterpreter_ModIntInt(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_di_numeru a uai
		a fica_assim_entao 10 % 3 uai
	cabo`)
	if interp.memory["a"] != 1 {
		t.Errorf("a = %v, want 1", interp.memory["a"])
	}
}

// ---------- Divisão por Zero ----------

func TestInterpreter_DivisaoPorZeroLiteral(t *testing.T) {
	msg := expectPanic(t, func() {
		runProgram(`bora_cumpade main() simbora
			trem_di_numeru a uai
			a fica_assim_entao 10 / 0 uai
		cabo`)
	})
	if !strings.Contains(msg, "division by zero") {
		t.Errorf("expected 'division by zero' in message, got: %s", msg)
	}
}

func TestInterpreter_DivisaoPorZeroVariavel(t *testing.T) {
	msg := expectPanic(t, func() {
		runProgram(`bora_cumpade main() simbora
			trem_di_numeru a, b uai
			b fica_assim_entao 0 uai
			a fica_assim_entao 10 / b uai
		cabo`)
	})
	if !strings.Contains(msg, "division by zero") {
		t.Errorf("expected 'division by zero' in message, got: %s", msg)
	}
}

// ---------- Comparações ----------

func TestInterpreter_EqTrue(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_discolhe a uai
		a fica_assim_entao 5 mema_coisa 5 uai
	cabo`)
	if interp.memory["a"] != true {
		t.Errorf("a = %v, want true", interp.memory["a"])
	}
}

func TestInterpreter_EqFalse(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_discolhe a uai
		a fica_assim_entao 5 mema_coisa 3 uai
	cabo`)
	if interp.memory["a"] != false {
		t.Errorf("a = %v, want false", interp.memory["a"])
	}
}

func TestInterpreter_Neq(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_discolhe a uai
		a fica_assim_entao 5 neh_nada 3 uai
	cabo`)
	if interp.memory["a"] != true {
		t.Errorf("a = %v, want true", interp.memory["a"])
	}
}

func TestInterpreter_Lt(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_discolhe a uai
		a fica_assim_entao 3 < 5 uai
	cabo`)
	if interp.memory["a"] != true {
		t.Errorf("a = %v, want true", interp.memory["a"])
	}
}

func TestInterpreter_Gt(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_discolhe a uai
		a fica_assim_entao 5 > 3 uai
	cabo`)
	if interp.memory["a"] != true {
		t.Errorf("a = %v, want true", interp.memory["a"])
	}
}

func TestInterpreter_Lte(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_discolhe a, b uai
		a fica_assim_entao 5 <= 5 uai
		b fica_assim_entao 3 <= 5 uai
	cabo`)
	if interp.memory["a"] != true {
		t.Errorf("a = %v, want true", interp.memory["a"])
	}
	if interp.memory["b"] != true {
		t.Errorf("b = %v, want true", interp.memory["b"])
	}
}

func TestInterpreter_Gte(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_discolhe a, b uai
		a fica_assim_entao 5 >= 5 uai
		b fica_assim_entao 5 >= 3 uai
	cabo`)
	if interp.memory["a"] != true {
		t.Errorf("a = %v, want true", interp.memory["a"])
	}
	if interp.memory["b"] != true {
		t.Errorf("b = %v, want true", interp.memory["b"])
	}
}

// ---------- Lógicos ----------

func TestInterpreter_And(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_discolhe a, b uai
		a fica_assim_entao eh tamem eh uai
		b fica_assim_entao eh tamem num_eh uai
	cabo`)
	if interp.memory["a"] != true {
		t.Errorf("a = %v, want true", interp.memory["a"])
	}
	if interp.memory["b"] != false {
		t.Errorf("b = %v, want false", interp.memory["b"])
	}
}

func TestInterpreter_Or(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_discolhe a, b uai
		a fica_assim_entao num_eh quarque_um eh uai
		b fica_assim_entao num_eh quarque_um num_eh uai
	cabo`)
	if interp.memory["a"] != true {
		t.Errorf("a = %v, want true", interp.memory["a"])
	}
	if interp.memory["b"] != false {
		t.Errorf("b = %v, want false", interp.memory["b"])
	}
}

func TestInterpreter_Not(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_discolhe a, b uai
		a fica_assim_entao vam_marca eh uai
		b fica_assim_entao vam_marca num_eh uai
	cabo`)
	if interp.memory["a"] != false {
		t.Errorf("a = %v, want false", interp.memory["a"])
	}
	if interp.memory["b"] != true {
		t.Errorf("b = %v, want true", interp.memory["b"])
	}
}

func TestInterpreter_Xor(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_discolhe a, b, c uai
		a fica_assim_entao eh um_o_oto num_eh uai
		b fica_assim_entao eh um_o_oto eh uai
		c fica_assim_entao num_eh um_o_oto num_eh uai
	cabo`)
	if interp.memory["a"] != true {
		t.Errorf("a = %v, want true", interp.memory["a"])
	}
	if interp.memory["b"] != false {
		t.Errorf("b = %v, want false", interp.memory["b"])
	}
	if interp.memory["c"] != false {
		t.Errorf("c = %v, want false", interp.memory["c"])
	}
}

// ---------- Unário ----------

func TestInterpreter_UnarioNegativo(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_di_numeru a uai
		a fica_assim_entao -5 uai
	cabo`)
	if interp.memory["a"] != -5 {
		t.Errorf("a = %v, want -5", interp.memory["a"])
	}
}

func TestInterpreter_UnarioPositivo(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_di_numeru a uai
		a fica_assim_entao +5 uai
	cabo`)
	if interp.memory["a"] != 5 {
		t.Errorf("a = %v, want 5", interp.memory["a"])
	}
}

// ---------- If-Else ----------

func TestInterpreter_IfTrue(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_di_numeru x uai
		x fica_assim_entao 0 uai
		uai_se(eh) simbora
			x fica_assim_entao 1 uai
		cabo
	cabo`)
	if interp.memory["x"] != 1 {
		t.Errorf("x = %v, want 1 (if-true branch)", interp.memory["x"])
	}
}

func TestInterpreter_IfFalse(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_di_numeru x uai
		x fica_assim_entao 0 uai
		uai_se(num_eh) simbora
			x fica_assim_entao 1 uai
		cabo
	cabo`)
	if interp.memory["x"] != 0 {
		t.Errorf("x = %v, want 0 (if-false, no else)", interp.memory["x"])
	}
}

func TestInterpreter_IfElse(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_di_numeru x uai
		x fica_assim_entao 0 uai
		uai_se(num_eh) simbora
			x fica_assim_entao 1 uai
		cabo
		uai_senao simbora
			x fica_assim_entao 2 uai
		cabo
	cabo`)
	if interp.memory["x"] != 2 {
		t.Errorf("x = %v, want 2 (else branch)", interp.memory["x"])
	}
}

// ---------- While ----------

func TestInterpreter_While(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_di_numeru i uai
		i fica_assim_entao 0 uai
		enquanto_tiver_trem(i < 5) simbora
			i fica_assim_entao i + 1 uai
		cabo
	cabo`)
	if interp.memory["i"] != 5 {
		t.Errorf("i = %v, want 5", interp.memory["i"])
	}
}

// ---------- For ----------

func TestInterpreter_For(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_di_numeru sum, i uai
		sum fica_assim_entao 0 uai
		roda_esse_trem(i fica_assim_entao 1; i <= 5; i fica_assim_entao i + 1) simbora
			sum fica_assim_entao sum + i uai
		cabo
	cabo`)
	if interp.memory["sum"] != 15 {
		t.Errorf("sum = %v, want 15 (1+2+3+4+5)", interp.memory["sum"])
	}
}

// ---------- Break ----------

func TestInterpreter_Break(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_di_numeru i uai
		i fica_assim_entao 0 uai
		enquanto_tiver_trem(i < 100) simbora
			uai_se(i mema_coisa 3) simbora
				para_o_trem uai
			cabo
			i fica_assim_entao i + 1 uai
		cabo
	cabo`)
	if interp.memory["i"] != 3 {
		t.Errorf("i = %v, want 3 (broken at i==3)", interp.memory["i"])
	}
}

// ---------- Continue ----------

func TestInterpreter_Continue(t *testing.T) {
	// Soma apenas os ímpares de 0 a 4
	interp := runProgram(`bora_cumpade main() simbora
		trem_di_numeru sum, i uai
		sum fica_assim_entao 0 uai
		roda_esse_trem(i fica_assim_entao 0; i < 5; i fica_assim_entao i + 1) simbora
			uai_se(i % 2 mema_coisa 0) simbora
				toca_o_trem uai
			cabo
			sum fica_assim_entao sum + i uai
		cabo
	cabo`)
	// sum = 1 + 3 = 4
	if interp.memory["sum"] != 4 {
		t.Errorf("sum = %v, want 4 (1+3, skipping evens)", interp.memory["sum"])
	}
}

// ---------- Print ----------

func TestInterpreter_PrintLiteral(t *testing.T) {
	output := captureOutput(`bora_cumpade main() simbora
		oia_proce_ve("hello") uai
	cabo`)
	if !strings.Contains(output, "hello") {
		t.Errorf("output = %q, want to contain 'hello'", output)
	}
}

func TestInterpreter_PrintVariavel(t *testing.T) {
	output := captureOutput(`bora_cumpade main() simbora
		trem_di_numeru x uai
		x fica_assim_entao 42 uai
		oia_proce_ve(x) uai
	cabo`)
	if !strings.Contains(output, "42") {
		t.Errorf("output = %q, want to contain '42'", output)
	}
}

func TestInterpreter_PrintMultiplo(t *testing.T) {
	output := captureOutput(`bora_cumpade main() simbora
		oia_proce_ve("a", "b", "c") uai
	cabo`)
	if !strings.Contains(output, "a") || !strings.Contains(output, "b") || !strings.Contains(output, "c") {
		t.Errorf("output = %q, want to contain a, b, c", output)
	}
}

// ---------- Hex e Octal ----------

func TestInterpreter_AtribuicaoHex(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_di_numeru x uai
		x fica_assim_entao 0xFF uai
	cabo`)
	if interp.memory["x"] != 255 {
		t.Errorf("x = %v, want 255 (0xFF)", interp.memory["x"])
	}
}

func TestInterpreter_AtribuicaoOctal(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_di_numeru x uai
		x fica_assim_entao 077 uai
	cabo`)
	if interp.memory["x"] != 63 {
		t.Errorf("x = %v, want 63 (077 octal)", interp.memory["x"])
	}
}

// ---------- Jump e Labels ----------

func TestInterpreter_JumpLabel(t *testing.T) {
	// Testa que o jump funciona corretamente via if-else
	interp := runProgram(`bora_cumpade main() simbora
		trem_di_numeru x uai
		uai_se(eh) simbora
			x fica_assim_entao 10 uai
		cabo
		uai_senao simbora
			x fica_assim_entao 20 uai
		cabo
	cabo`)
	if interp.memory["x"] != 10 {
		t.Errorf("x = %v, want 10", interp.memory["x"])
	}
}

// ---------- Expressões Complexas ----------

func TestInterpreter_ExpressaoComplexa(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_di_numeru a uai
		a fica_assim_entao (2 + 3) veiz (4 - 1) uai
	cabo`)
	// (2+3) * (4-1) = 5 * 3 = 15
	if interp.memory["a"] != 15 {
		t.Errorf("a = %v, want 15", interp.memory["a"])
	}
}

func TestInterpreter_NeqStringLiteral(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_discolhe a, b uai
		a fica_assim_entao "abc" neh_nada "def" uai
		b fica_assim_entao "abc" neh_nada "abc" uai
	cabo`)
	if interp.memory["a"] != true {
		t.Errorf("a = %v, want true ('abc' != 'def')", interp.memory["a"])
	}
	if interp.memory["b"] != false {
		t.Errorf("b = %v, want false ('abc' != 'abc')", interp.memory["b"])
	}
}

// ---------- Switch Case ----------

func TestInterpreter_SwitchCase(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_di_numeru x, result uai
		x fica_assim_entao 2 uai
		dependenu(x) simbora
			du_casu 1: result fica_assim_entao 10 uai
			du_casu 2: result fica_assim_entao 20 uai
			du_casu 3: result fica_assim_entao 30 uai
		cabo
	cabo`)
	if interp.memory["result"] != 20 {
		t.Errorf("result = %v, want 20 (case 2)", interp.memory["result"])
	}
}

// ---------- Float para Int (RoundHalfDown) ----------

func TestInterpreter_FloatToInt(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_cum_virgula f uai
		trem_di_numeru i uai
		f fica_assim_entao 3.7 uai
		i fica_assim_entao f uai
	cabo`)
	v, ok := interp.memory["i"].(float64)
	if !ok {
		t.Fatalf("i is not float64, is %T: %v", interp.memory["i"], interp.memory["i"])
	}
	// RoundHalfDown(3.7) = 4.0
	if v != 4.0 {
		t.Errorf("i = %v, want 4 (RoundHalfDown(3.7))", v)
	}
}

// ---------- Comparação de strings com variáveis ----------

func TestInterpreter_EqStringVar(t *testing.T) {
	interp := runProgram(`bora_cumpade main() simbora
		trem_discrita s uai
		trem_discolhe r uai
		s fica_assim_entao "hello" uai
		r fica_assim_entao s mema_coisa "hello" uai
	cabo`)
	if interp.memory["r"] != true {
		t.Errorf("r = %v, want true", interp.memory["r"])
	}
}
