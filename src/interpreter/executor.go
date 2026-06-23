package interpreter

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"mineres-interpreter/src/lexer"
	"mineres-interpreter/src/parser"
	"mineres-interpreter/src/utils"
)

func RoundHalfDown(x float64) float64 {
	frac := x - math.Floor(x)
	if frac <= 0.5 {
		return math.Floor(x)
	}
	return math.Ceil(x)
}

func (interpreter *Interpreter) setMemory(nome string, t *lexer.TuplaLex) {
	switch t.Token {
	case lexer.Literal_oct:
		v, _ := strconv.ParseInt(t.Lexema[1:], 8, 64)
		interpreter.memory[nome] = int(v)
	case lexer.Literal_hex:
		v, _ := strconv.ParseInt(t.Lexema[2:], 16, 64)
		interpreter.memory[nome] = int(v)
	case lexer.Literal_int:
		v, _ := strconv.Atoi(t.Lexema)
		interpreter.memory[nome] = v
	case lexer.Literal_float:
		v, _ := strconv.ParseFloat(t.Lexema, 64)
		interpreter.memory[nome] = v
	case lexer.Literal_true:
		interpreter.memory[nome] = true
	case lexer.Literal_false:
		interpreter.memory[nome] = false
	case lexer.Literal_string, lexer.Literal_char:
		interpreter.memory[nome] = t.Lexema
	case lexer.Identifier:
		variavel, _ := interpreter.memory[nome]
		switch variavel.(type) {
		case int:
			switch v := interpreter.memory[t.Lexema].(type) {
			case float64:
				interpreter.memory[nome] = RoundHalfDown(v)
				return
			}
		}
		interpreter.memory[nome] = interpreter.memory[t.Lexema]
	default:
		utils.ThrowInterpreterException(
			"invalid token type for value resolution",
			t.Linha,
			t.Coluna,
		)
	}
}

func (i *Interpreter) checkDivisionByZero(t *lexer.TuplaLex) {
	if t.Token == lexer.Identifier {
		switch v := i.memory[t.Lexema].(type) {
		case int:
			if v == 0 {
				utils.ThrowInterpreterException("division by zero is not allowed", t.Linha, t.Coluna)
			}
		case float64:
			if v == 0.0 {
				utils.ThrowInterpreterException("division by zero is not allowed", t.Linha, t.Coluna)
			}
		}
		return
	} else if t.Lexema == "0" || t.Lexema == "0.0" {
		utils.ThrowInterpreterException("division by zero is not allowed", t.Linha, t.Coluna)
	}
}

func (i *Interpreter) resolveValue(t *lexer.TuplaLex) any {
	if t == nil {
		return nil
	}
	switch t.Token {
	case lexer.Literal_int:
		v, _ := strconv.Atoi(t.Lexema)
		return v
	case lexer.Literal_float:
		v, _ := strconv.ParseFloat(t.Lexema, 64)
		return v
	case lexer.Literal_string, lexer.Literal_char:
		return t.Lexema
	case lexer.Literal_true:
		return true
	case lexer.Literal_false:
		return false
	case lexer.Identifier:
		v, _ := i.memory[t.Lexema]
		return v
	}
	return nil
}

func (i *Interpreter) resolveNumericPair(t1, t2 *lexer.TuplaLex) (float64, float64, bool) {
	v1 := i.resolveValue(t1)
	v2 := i.resolveValue(t2)

	isFloat := false
	var f1, f2 float64

	switch val := v1.(type) {
	case int:
		f1 = float64(val)
	case float64:
		f1 = val
		isFloat = true
	}

	switch val := v2.(type) {
	case int:
		f2 = float64(val)
	case float64:
		f2 = val
		isFloat = true
	}

	return f1, f2, isFloat
}

func (i *Interpreter) resolveStringPair(t1, t2 *lexer.TuplaLex) (string, string, bool) {
	v1 := i.resolveValue(t1)
	v2 := i.resolveValue(t2)

	var s1, s2 string
	isString := false

	if s, ok := v1.(string); ok {
		s1 = s
		isString = true
	}
	if s, ok := v2.(string); ok {
		s2 = s
		isString = true
	}

	return s1, s2, isString
}

func (i *Interpreter) resolveBoolPair(t1, t2 *lexer.TuplaLex) (bool, bool) {
	v1 := i.resolveValue(t1)
	v2 := i.resolveValue(t2)
	b1, _ := v1.(bool)
	b2, _ := v2.(bool)
	return b1, b2
}

func (i *Interpreter) operationAdd(t1 *lexer.TuplaLex, t2 *lexer.TuplaLex) any {
	s1, s2, isString := i.resolveStringPair(t1, t2)
	if isString {
		return s1 + s2
	}
	f1, f2, isFloat := i.resolveNumericPair(t1, t2)
	if isFloat {
		return f1 + f2
	}
	return int(f1) + int(f2)
}

func (i *Interpreter) operationSub(t1 *lexer.TuplaLex, t2 *lexer.TuplaLex) any {
	f1, f2, isFloat := i.resolveNumericPair(t1, t2)
	if isFloat {
		return f1 - f2
	}
	return int(f1) - int(f2)
}

func (i *Interpreter) operationMul(t1 *lexer.TuplaLex, t2 *lexer.TuplaLex) any {
	f1, f2, isFloat := i.resolveNumericPair(t1, t2)
	if isFloat {
		return f1 * f2
	}
	return int(f1) * int(f2)
}

func (i *Interpreter) operationDiv(t1 *lexer.TuplaLex, t2 *lexer.TuplaLex) any {
	i.checkDivisionByZero(t2)
	f1, f2, isFloat := i.resolveNumericPair(t1, t2)
	if isFloat {
		return f1 / f2
	}
	if f1 == 0 && f2 == 0 { return 0.0 }
	return f1 / f2
}

func (i *Interpreter) operationMod(t1 *lexer.TuplaLex, t2 *lexer.TuplaLex) any {
	i.checkDivisionByZero(t2)
	f1, f2, _ := i.resolveNumericPair(t1, t2)
	return int(f1) % int(f2)
}

func (i *Interpreter) operationDivI(t1 *lexer.TuplaLex, t2 *lexer.TuplaLex) any {
	i.checkDivisionByZero(t2)
	f1, f2, _ := i.resolveNumericPair(t1, t2)
	return int(f1) / int(f2)
}

func (i *Interpreter) operationEq(t1 *lexer.TuplaLex, t2 *lexer.TuplaLex) bool {
	v1 := i.resolveValue(t1)
	v2 := i.resolveValue(t2)

	switch v1.(type) {
	case string, bool:
		return v1 == v2
	}

	f1, f2, _ := i.resolveNumericPair(t1, t2)
	return f1 == f2
}

func (i *Interpreter) operationNeq(t1 *lexer.TuplaLex, t2 *lexer.TuplaLex) bool {
	return !i.operationEq(t1, t2)
}

func (i *Interpreter) operationLt(t1 *lexer.TuplaLex, t2 *lexer.TuplaLex) bool {
	f1, f2, _ := i.resolveNumericPair(t1, t2)
	return f1 < f2
}

func (i *Interpreter) operationGt(t1 *lexer.TuplaLex, t2 *lexer.TuplaLex) bool {
	f1, f2, _ := i.resolveNumericPair(t1, t2)
	return f1 > f2
}

func (i *Interpreter) operationLte(t1 *lexer.TuplaLex, t2 *lexer.TuplaLex) bool {
	f1, f2, _ := i.resolveNumericPair(t1, t2)
	return f1 <= f2
}

func (i *Interpreter) operationGte(t1 *lexer.TuplaLex, t2 *lexer.TuplaLex) bool {
	f1, f2, _ := i.resolveNumericPair(t1, t2)
	return f1 >= f2
}

func (i *Interpreter) operationAnd(t1 *lexer.TuplaLex, t2 *lexer.TuplaLex) bool {
	b1, b2 := i.resolveBoolPair(t1, t2)
	return b1 && b2
}

func (i *Interpreter) operationOr(t1 *lexer.TuplaLex, t2 *lexer.TuplaLex) bool {
	b1, b2 := i.resolveBoolPair(t1, t2)
	return b1 || b2
}

func (i *Interpreter) operationNot(t1 *lexer.TuplaLex) bool {
	b1, _ := i.resolveBoolPair(t1, nil)
	return !b1
}

func (i *Interpreter) operationXor(t1 *lexer.TuplaLex, t2 *lexer.TuplaLex) bool {
	b1, b2 := i.resolveBoolPair(t1, t2)
	return b1 != b2
}

func (i *Interpreter) operationCall(res *lexer.TuplaLex, t1 *lexer.TuplaLex, t2 *lexer.TuplaLex) {
	switch res.Lexema {
	case "read":
		entrada, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		entrada = strings.TrimSpace(entrada)

		valAnterior, declarada := i.memory[t1.Lexema]
		if !declarada {
			utils.ThrowInterpreterException(
				fmt.Sprintf("Variable '%s' not declared", t1.Lexema),
				t1.Linha, t1.Coluna,
			)
		}

		switch valAnterior.(type) {
		case int:
			if t2.Token != lexer.Type_int {
				utils.ThrowInterpreterException(
					fmt.Sprintf("Type mismatch: cannot scan type '%s' into variable '%s' of type 'trem_di_numeru'", t2.Lexema, t1.Lexema),
					t1.Linha, t1.Coluna,
				)
			}
		case float64:
			if t2.Token != lexer.Type_float {
				utils.ThrowInterpreterException(
					fmt.Sprintf("Type mismatch: cannot scan type '%s' into variable '%s' of type 'trem_cum_virgula'", t2.Lexema, t1.Lexema),
					t1.Linha, t1.Coluna,
				)
			}
		case bool:
			if t2.Token != lexer.Type_bool {
				utils.ThrowInterpreterException(
					fmt.Sprintf("Type mismatch: cannot scan type '%s' into variable '%s' of type 'trem_discolhe'", t2.Lexema, t1.Lexema),
					t1.Linha, t1.Coluna,
				)
			}
		case string:
			if t2.Token != lexer.Type_string && t2.Token != lexer.Type_char {
				utils.ThrowInterpreterException(
					fmt.Sprintf("Type mismatch: cannot scan type '%s' into variable '%s' of type 'trem_discrita'/'trosso'", t2.Lexema, t1.Lexema),
					t1.Linha, t1.Coluna,
				)
			}
		}

		switch t2.Token {
		case lexer.Type_int:
			v, err := strconv.Atoi(entrada)
			if err != nil {
				utils.ThrowInterpreterException(
					fmt.Sprintf("Invalid input for type '%s': '%s'", t2.Lexema, entrada),
					t1.Linha, t1.Coluna,
				)
			}
			i.memory[t1.Lexema] = v
		case lexer.Type_float:
			v, err := strconv.ParseFloat(entrada, 64)
			if err != nil {
				utils.ThrowInterpreterException(
					fmt.Sprintf("Invalid input for type '%s': '%s'", t2.Lexema, entrada),
					t1.Linha, t1.Coluna,
				)
			}
			i.memory[t1.Lexema] = v
		case lexer.Type_bool:
			switch entrada {
			case "eh":
				i.memory[t1.Lexema] = true
			case "num_eh":
				i.memory[t1.Lexema] = false
			default:
				utils.ThrowInterpreterException(
					fmt.Sprintf("Invalid input for type '%s' (expected 'eh' or 'num_eh'): '%s'", t2.Lexema, entrada),
					t1.Linha, t1.Coluna,
				)
			}
		case lexer.Type_char:
			runes := []rune(entrada)
			if len(runes) != 1 {
				utils.ThrowInterpreterException(
					fmt.Sprintf("Invalid input for type '%s' (expected a single character): '%s'", t2.Lexema, entrada),
					t1.Linha, t1.Coluna,
				)
			}
			i.memory[t1.Lexema] = string(runes[0])
		case lexer.Type_string:
			i.memory[t1.Lexema] = entrada
		default:
			utils.ThrowInterpreterException(
				fmt.Sprintf("Unsupported scan type: '%s'", t2.Lexema),
				t1.Linha, t1.Coluna,
			)
		}
	default:
		if t1 != nil && t1.Token == lexer.Identifier {
			valor, _ := i.memory[t1.Lexema]
			fmt.Print(valor)
			return
		}
		if t2 != nil {
			switch t2.Token {
			case lexer.Literal_true:
				fmt.Print(true)
			case lexer.Literal_false:
				fmt.Print(false)
			default:
				fmt.Print(t2.Lexema)
			}
			return
		}
	}
}

func (i *Interpreter) operationJump(label string) {
	i.ip = i.labels[label]
}

func (i *Interpreter) operationIf_eq(res *lexer.TuplaLex, l1 string, l2 string) {
	comparation := i.resolveValue(res)
	if b, ok := comparation.(bool); ok {
		if b {
			i.ip = i.labels[l1] - 1
		} else {
			i.ip = i.labels[l2] - 1
		}
	}
}

func (i *Interpreter) operationUno(t1 *lexer.TuplaLex) any {
	f1, _, isFloat := i.resolveNumericPair(t1, nil)
	if isFloat {
		return f1 * -1
	}
	return int(f1) * -1
}

func (i *Interpreter) operationAtt(nome string, operation *lexer.TuplaLex) {
	i.setMemory(nome, operation)
}

func (interpreter *Interpreter) execute(instrucao parser.TuplaMicrocode) {
	switch instrucao.Operation {
	case parser.Add:
		interpreter.memory[instrucao.Res.Lexema] = interpreter.operationAdd(instrucao.Op1, instrucao.Op2)
	case parser.Sub:
		interpreter.memory[instrucao.Res.Lexema] = interpreter.operationSub(instrucao.Op1, instrucao.Op2)
	case parser.Mul:
		interpreter.memory[instrucao.Res.Lexema] = interpreter.operationMul(instrucao.Op1, instrucao.Op2)
	case parser.Div:
		interpreter.memory[instrucao.Res.Lexema] = interpreter.operationDiv(instrucao.Op1, instrucao.Op2)
	case parser.Mod:
		interpreter.memory[instrucao.Res.Lexema] = interpreter.operationMod(instrucao.Op1, instrucao.Op2)
	case parser.DivI:
		interpreter.memory[instrucao.Res.Lexema] = interpreter.operationDivI(instrucao.Op1, instrucao.Op2)
	case parser.Eq:
		interpreter.memory[instrucao.Res.Lexema] = interpreter.operationEq(instrucao.Op1, instrucao.Op2)
	case parser.Neq:
		interpreter.memory[instrucao.Res.Lexema] = interpreter.operationNeq(instrucao.Op1, instrucao.Op2)
	case parser.Lt:
		interpreter.memory[instrucao.Res.Lexema] = interpreter.operationLt(instrucao.Op1, instrucao.Op2)
	case parser.Gt:
		interpreter.memory[instrucao.Res.Lexema] = interpreter.operationGt(instrucao.Op1, instrucao.Op2)
	case parser.Lte:
		interpreter.memory[instrucao.Res.Lexema] = interpreter.operationLte(instrucao.Op1, instrucao.Op2)
	case parser.Gte:
		interpreter.memory[instrucao.Res.Lexema] = interpreter.operationGte(instrucao.Op1, instrucao.Op2)
	case parser.And:
		interpreter.memory[instrucao.Res.Lexema] = interpreter.operationAnd(instrucao.Op1, instrucao.Op2)
	case parser.Or:
		interpreter.memory[instrucao.Res.Lexema] = interpreter.operationOr(instrucao.Op1, instrucao.Op2)
	case parser.Not:
		interpreter.memory[instrucao.Res.Lexema] = interpreter.operationNot(instrucao.Op1)
	case parser.Xor:
		interpreter.memory[instrucao.Res.Lexema] = interpreter.operationXor(instrucao.Op1, instrucao.Op2)
	case parser.Call:
		interpreter.operationCall(instrucao.Res, instrucao.Op1, instrucao.Op2)
	case parser.Jump:
		interpreter.operationJump(instrucao.Res.Lexema)
	case parser.Label:
	case parser.If_eq:
		interpreter.operationIf_eq(instrucao.Res, instrucao.Op1.Lexema, instrucao.Op2.Lexema)
	case parser.Uno:
		interpreter.memory[instrucao.Res.Lexema] = interpreter.operationUno(instrucao.Op1)
	case parser.Att:
		interpreter.operationAtt(instrucao.Res.Lexema, instrucao.Op1)
	}
}
