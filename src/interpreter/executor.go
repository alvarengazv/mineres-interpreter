package interpreter

import (
	"bufio"
	"fmt"
	"math"
	"mineres-interpreter/src/lexer"
	"mineres-interpreter/src/parser"
	"mineres-interpreter/src/utils"
	"os"
	"strconv"
	"strings"
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
	case lexer.Literal_string:
		interpreter.memory[nome] = t.Lexema
	case lexer.Literal_char:
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
				utils.ThrowInterpreterException(
					"division by zero is not allowed",
					t.Linha,
					t.Coluna,
				)
			}
		case float64:
			if v == 0.0 {
				utils.ThrowInterpreterException(
					"division by zero is not allowed",
					t.Linha,
					t.Coluna,
				)
			}
		}
		return
	} else if t.Lexema == "0" || t.Lexema == "0.0" {
		utils.ThrowInterpreterException(
			"division by zero is not allowed",
			t.Linha,
			t.Coluna,
		)
	}
}

func (i *Interpreter) operationAdd(t1 *lexer.TuplaLex, t2 *lexer.TuplaLex) any {
	switch t1.Token {
	case lexer.Literal_int: // Op1 é int
		op1, _ := strconv.Atoi(t1.Lexema)
		switch t2.Token {
		case lexer.Literal_int: // Op1 é int e Op2 é int
			op2, _ := strconv.Atoi(t2.Lexema)
			return op1 + op2
		case lexer.Literal_float: // Op1 é int e Op2 é float
			op2, _ := strconv.ParseFloat(t2.Lexema, 64)
			return float64(op1) + op2
		case lexer.Identifier: // Op1 é int e Op2 está em memória
			op2, _ := i.memory[t2.Lexema]
			switch v := op2.(type) {
			case int: // Op1 é int e Op2 é int
				return op1 + v
			case float64: // Op1 é int e Op2 é float
				return float64(op1) + v
			}
		}
	case lexer.Literal_float: // Op1 é float
		op1, _ := strconv.ParseFloat(t1.Lexema, 64)
		switch t2.Token {
		case lexer.Literal_int: // Op1 é float e Op2 é int
			op2, _ := strconv.Atoi(t2.Lexema)
			return op1 + float64(op2)
		case lexer.Literal_float: // Op1 é float e Op2 é float
			op2, _ := strconv.ParseFloat(t2.Lexema, 64)
			return op1 + op2
		case lexer.Identifier: // Op1 é float e Op2 está em memória
			op2, _ := i.memory[t2.Lexema]
			switch v := op2.(type) {
			case int: // Op1 é float e Op2 é int
				return op1 + float64(v)
			case float64: // Op1 é float e Op2 é float
				return op1 + v
			}
		}
	case lexer.Literal_string: // Op1 é string
		op1 := t1.Lexema
		switch t2.Token {
		case lexer.Literal_string: // Op1 é string e Op2 é string
			op2 := t2.Lexema
			return op1 + op2
		case lexer.Literal_char: // Op1 é string e Op2 é char
			op2 := string(t2.Lexema)
			return op1 + op2
		case lexer.Identifier: // Op1 é string e Op2 está em memória
			op2, _ := i.memory[t2.Lexema]
			switch v := op2.(type) {
			case string: // Op1 é string e Op2 é string
				// Op1 é string e Op2 é char
				return op1 + v
			}
		}
	case lexer.Literal_char: // Op1 é char
		op1 := string(t1.Lexema)
		switch t2.Token {
		case lexer.Literal_string: // Op1 é char e Op2 é string
			op2 := t2.Lexema
			return op1 + op2
		case lexer.Literal_char: // Op1 é char e Op2 é char
			op2 := string(t2.Lexema)
			return op1 + op2
		case lexer.Identifier: // Op1 é char e Op2 está em memória
			op2, _ := i.memory[t2.Lexema]
			switch v := op2.(type) {
			case string: // Op1 é char e Op2 é string
				// Op1 é char e Op2 é char
				return op1 + v
			}
		}
	case lexer.Identifier: // Op1 esta em memoria
		op1, _ := i.memory[t1.Lexema]

		switch v := op1.(type) {
		case int: // Op1 é intc
			switch t2.Token {
			case lexer.Literal_int: // Op1 é int e Op2 é int
				op2, _ := strconv.Atoi(t2.Lexema)
				return v + op2
			case lexer.Literal_float: // Op1 é int e Op2 é float
				op2, _ := strconv.ParseFloat(t2.Lexema, 64)
				return float64(v) + op2
			case lexer.Identifier: // Op1 é int e Op2 está em memória
				op2, _ := i.memory[t2.Lexema]
				switch v2 := op2.(type) {
				case int: // Op1 é int e Op2 é int
					return v + v2
				case float64: // Op1 é int e Op2 é float
					return float64(v) + v2
				}
			}
		case float64: // Op1 é float
			switch t2.Token {
			case lexer.Literal_int: // Op1 é float e Op2 é int
				op2, _ := strconv.Atoi(t2.Lexema)
				return v + float64(op2)
			case lexer.Literal_float: // Op1 é float e Op2 é float
				op2, _ := strconv.ParseFloat(t2.Lexema, 64)
				return v + op2
			case lexer.Identifier: // Op1 é float e Op2 está em memória
				op2, _ := i.memory[t2.Lexema]
				switch v2 := op2.(type) {
				case int: // Op1 é float e Op2 é int
					return v + float64(v2)
				case float64: // Op1 é float e Op2 é float
					return v + v2
				}
			}
		case string:
			switch t2.Token {
			case lexer.Literal_char: // Op1 é string e Op2 é char
				op2 := string(t2.Lexema)
				return v + op2
			case lexer.Literal_string: // Op1 é string e Op2 é string
				op2 := t2.Lexema
				return v + op2
			case lexer.Identifier: // Op1 é string e Op2 está em memória
				op2, _ := i.memory[t2.Lexema]
				switch v2 := op2.(type) {
				case string: // Op1 é string e Op2 é string
					// Op1 é string e Op2 é char
					return v + v2
				}
			}
		}
	default:
		utils.ThrowException("executor.go", "operationAdd", "invalid token type for operationAdd")
	}
	return false
}

func (i *Interpreter) operationSub(t1 *lexer.TuplaLex, t2 *lexer.TuplaLex) any {
	switch t1.Token {
	case lexer.Literal_int: // Op1 é int
		op1, _ := strconv.Atoi(t1.Lexema)
		switch t2.Token {
		case lexer.Literal_int: // Op1 é int e Op2 é int
			op2, _ := strconv.Atoi(t2.Lexema)
			return op1 - op2
		case lexer.Literal_float: // Op1 é int e Op2 é float
			op2, _ := strconv.ParseFloat(t2.Lexema, 64)
			return float64(op1) - op2
		case lexer.Identifier: // Op1 é int e Op2 está em memória
			op2, _ := i.memory[t2.Lexema]
			switch v := op2.(type) {
			case int: // Op1 é int e Op2 é int
				return op1 - v
			case float64: // Op1 é int e Op2 é float
				return float64(op1) - v
			}
		}
	case lexer.Literal_float: // Op1 é float
		op1, _ := strconv.ParseFloat(t1.Lexema, 64)
		switch t2.Token {
		case lexer.Literal_int: // Op1 é float e Op2 é int
			op2, _ := strconv.Atoi(t2.Lexema)
			return op1 - float64(op2)
		case lexer.Literal_float: // Op1 é float e Op2 é float
			op2, _ := strconv.ParseFloat(t2.Lexema, 64)
			return op1 - op2
		case lexer.Identifier: // Op1 é float e Op2 está em memória
			op2, _ := i.memory[t2.Lexema]
			switch v := op2.(type) {
			case int: // Op1 é float e Op2 é int
				return op1 - float64(v)
			case float64: // Op1 é float e Op2 é float
				return op1 - v
			}
		}
	case lexer.Identifier: // Op1 esta em memoria
		op1, _ := i.memory[t1.Lexema]
		switch v := op1.(type) {
		case int: // Op1 é int
			switch t2.Token {
			case lexer.Literal_int: // Op1 é int e Op2 é int
				op2, _ := strconv.Atoi(t2.Lexema)
				return v - op2
			case lexer.Literal_float: // Op1 é int e Op2 é float
				op2, _ := strconv.ParseFloat(t2.Lexema, 64)
				return float64(v) - op2
			case lexer.Identifier: // Op1 é int e Op2 está em memória
				op2, _ := i.memory[t2.Lexema]
				switch v2 := op2.(type) {
				case int: // Op1 é int e Op2 é int
					return v - v2
				case float64: // Op1 é int e Op2 é float
					return float64(v) - v2
				}
			}
		case float64: // Op1 é float
			switch t2.Token {
			case lexer.Literal_int: // Op1 é float e Op2 é int
				op2, _ := strconv.Atoi(t2.Lexema)
				return v - float64(op2)
			case lexer.Literal_float: // Op1 é float e Op2 é float
				op2, _ := strconv.ParseFloat(t2.Lexema, 64)
				return v - op2
			case lexer.Identifier: // Op1 é float e Op2 está em memória
				op2, _ := i.memory[t2.Lexema]
				switch v2 := op2.(type) {
				case int: // Op1 é float e Op2 é int
					return v - float64(v2)
				case float64: // Op1 é float e Op2 é float
					return v - v2
				}
			}
		}
	default:
		utils.ThrowException("executor.go", "operationSub", "invalid type for operationSub")
	}
	return false
}

func (i *Interpreter) operationMul(t1 *lexer.TuplaLex, t2 *lexer.TuplaLex) any {
	switch t1.Token {
	case lexer.Literal_int: // Op1 é int
		op1, _ := strconv.Atoi(t1.Lexema)
		switch t2.Token {
		case lexer.Literal_int: // Op1 é int e Op2 é int
			op2, _ := strconv.Atoi(t2.Lexema)
			return op1 * op2
		case lexer.Literal_float: // Op1 é int e Op2 é float
			op2, _ := strconv.ParseFloat(t2.Lexema, 64)
			return float64(op1) * op2
		case lexer.Identifier: // Op1 é int e Op2 está em memória
			op2, _ := i.memory[t2.Lexema]
			switch v := op2.(type) {
			case int: // Op1 é int e Op2 é int
				return op1 * v
			case float64: // Op1 é int e Op2 é float
				return float64(op1) * v
			}
		}
	case lexer.Literal_float: // Op1 é float
		op1, _ := strconv.ParseFloat(t1.Lexema, 64)
		switch t2.Token {
		case lexer.Literal_int: // Op1 é float e Op2 é int
			op2, _ := strconv.Atoi(t2.Lexema)
			return op1 * float64(op2)
		case lexer.Literal_float: // Op1 é float e Op2 é float
			op2, _ := strconv.ParseFloat(t2.Lexema, 64)
			return op1 * op2
		case lexer.Identifier: // Op1 é float e Op2 está em memória
			op2, _ := i.memory[t2.Lexema]
			switch v := op2.(type) {
			case int: // Op1 é float e Op2 é int
				return op1 * float64(v)
			case float64: // Op1 é float e Op2 é float
				return op1 * v
			}
		}
	case lexer.Identifier: // Op1 esta em memoria
		op1, _ := i.memory[t1.Lexema]
		switch v := op1.(type) {
		case int: // Op1 é int
			switch t2.Token {
			case lexer.Literal_int: // Op1 é int e Op2 é int
				op2, _ := strconv.Atoi(t2.Lexema)
				return v * op2
			case lexer.Literal_float: // Op1 é int e Op2 é float
				op2, _ := strconv.ParseFloat(t2.Lexema, 64)
				return float64(v) * op2
			case lexer.Identifier: // Op1 é int e Op2 está em memória
				op2, _ := i.memory[t2.Lexema]
				switch v2 := op2.(type) {
				case int: // Op1 é int e Op2 é int
					return v * v2
				case float64: // Op1 é int e Op2 é float
					return float64(v) * v2
				}
			}
		case float64: // Op1 é float
			switch t2.Token {
			case lexer.Literal_int: // Op1 é float e Op2 é int
				op2, _ := strconv.Atoi(t2.Lexema)
				return v * float64(op2)
			case lexer.Literal_float: // Op1 é float e Op2 é float
				op2, _ := strconv.ParseFloat(t2.Lexema, 64)
				return v * op2
			case lexer.Identifier: // Op1 é float e Op2 está em memória
				op2, _ := i.memory[t2.Lexema]
				switch v2 := op2.(type) {
				case int: // Op1 é float e Op2 é int
					return v * float64(v2)
				case float64: // Op1 é float e Op2 é float
					return v * v2
				}
			}
		}
	default:
		utils.ThrowException("executor.go", "operationMul", "invalid type for operationMul")
	}
	return false
}

func (i *Interpreter) operationDiv(t1 *lexer.TuplaLex, t2 *lexer.TuplaLex) any {
	i.checkDivisionByZero(t2)
	switch t1.Token {
	case lexer.Literal_int: // Op1 é int
		op1, _ := strconv.Atoi(t1.Lexema)
		switch t2.Token {
		case lexer.Literal_int: // Op1 é int e Op2 é int
			op2, _ := strconv.Atoi(t2.Lexema)
			return float64(op1) / float64(op2)
		case lexer.Literal_float: // Op1 é int e Op2 é float
			op2, _ := strconv.ParseFloat(t2.Lexema, 64)
			return float64(op1) / op2
		case lexer.Identifier: // Op1 é int e Op2 está em memória
			op2, _ := i.memory[t2.Lexema]
			switch v := op2.(type) {
			case int: // Op1 é int e Op2 é int
				return float64(op1) / float64(v)
			case float64: // Op1 é int e Op2 é float
				return float64(op1) / v
			}
		}
	case lexer.Literal_float: // Op1 é float
		op1, _ := strconv.ParseFloat(t1.Lexema, 64)
		switch t2.Token {
		case lexer.Literal_int: // Op1 é float e Op2 é int
			op2, _ := strconv.Atoi(t2.Lexema)
			return op1 / float64(op2)
		case lexer.Literal_float: // Op1 é float e Op2 é float
			op2, _ := strconv.ParseFloat(t2.Lexema, 64)
			return op1 / op2
		case lexer.Identifier: // Op1 é float e Op2 está em memória
			op2, _ := i.memory[t2.Lexema]
			switch v := op2.(type) {
			case int: // Op1 é float e Op2 é int
				return op1 / float64(v)
			case float64: // Op1 é float e Op2 é float
				return op1 / v
			}
		}
	case lexer.Identifier: // Op1 esta em memoria
		op1, _ := i.memory[t1.Lexema]
		switch v := op1.(type) {
		case int: // Op1 é int
			switch t2.Token {
			case lexer.Literal_int: // Op1 é int e Op2 é int
				op2, _ := strconv.Atoi(t2.Lexema)
				return float64(v) / float64(op2)
			case lexer.Literal_float: // Op1 é int e Op2 é float
				op2, _ := strconv.ParseFloat(t2.Lexema, 64)
				return float64(v) / op2
			case lexer.Identifier: // Op1 é int e Op2 está em memória
				op2, _ := i.memory[t2.Lexema]
				switch v2 := op2.(type) {
				case int: // Op1 é int e Op2 é int
					return float64(v) / float64(v2)
				case float64: // Op1 é int e Op2 é float
					return float64(v) / v2
				}
			}
		case float64: // Op1 é float
			switch t2.Token {
			case lexer.Literal_int: // Op1 é float e Op2 é int
				op2, _ := strconv.Atoi(t2.Lexema)
				return v / float64(op2)
			case lexer.Literal_float: // Op1 é float e Op2 é float
				op2, _ := strconv.ParseFloat(t2.Lexema, 64)
				return v / op2
			case lexer.Identifier: // Op1 é float e Op2 está em memória
				op2, _ := i.memory[t2.Lexema]
				switch v2 := op2.(type) {
				case int: // Op1 é float e Op2 é int
					return v / float64(v2)
				case float64: // Op1 é float e Op2 é float
					return v / v2
				}
			}
		}
	default:
		utils.ThrowException("executor.go", "operationDiv", "invalid type for operationDiv")
	}
	return false
}

func (i *Interpreter) operationMod(t1 *lexer.TuplaLex, t2 *lexer.TuplaLex) any {
	i.checkDivisionByZero(t2)
	switch t1.Token {
	case lexer.Literal_int: // Op1 é int
		op1, _ := strconv.Atoi(t1.Lexema)
		switch t2.Token {
		case lexer.Literal_int: // Op1 é int e Op2 é int
			op2, _ := strconv.Atoi(t2.Lexema)
			return op1 % op2
		// case lexer.Literal_float: // Op1 é int e Op2 é float
		// 	op2, _ := strconv.ParseFloat(t2.Lexema, 64)
		// 	return op1 % int(op2)
		case lexer.Identifier: // Op1 é int e Op2 está em memória
			op2, _ := i.memory[t2.Lexema]
			switch v := op2.(type) {
			case int: // Op1 é int e Op2 é int
				return op1 % v
				// case float64: // Op1 é int e Op2 é float
				// 	return op1 % int(v)
			}
		}
	// case lexer.Literal_float: // Op1 é float
	// 	op1, _ := strconv.ParseFloat(t1.Lexema, 64)
	// 	switch t2.Token {
	// 	case lexer.Literal_int: // Op1 é float e Op2 é int
	// 		op2, _ := strconv.Atoi(t2.Lexema)
	// 		return int(op1) % op2
	// 	// case lexer.Literal_float: // Op1 é float e Op2 é float
	// 	// 	op2, _ := strconv.ParseFloat(t2.Lexema, 64)
	// 	// 	return int(op1) % int(op2)
	// 	case lexer.Identifier: // Op1 é float e Op2 está em memória
	// 		op2, _ := i.memory[t2.Lexema]
	// 		switch v := op2.(type) {
	// 		case int: // Op1 é float e Op2 é int
	// 			return int(op1) % v
	// 		case float64: // Op1 é float e Op2 é float
	// 			return int(op1) % int(v)
	// 		}
	// 	}
	case lexer.Identifier: // Op1 esta em memoria
		op1, _ := i.memory[t1.Lexema]
		switch v := op1.(type) {
		case int: // Op1 é int
			switch t2.Token {
			case lexer.Literal_int: // Op1 é int e Op2 é int
				op2, _ := strconv.Atoi(t2.Lexema)
				return v % op2
			// case lexer.Literal_float: // Op1 é int e Op2 é float
			// 	op2, _ := strconv.ParseFloat(t2.Lexema, 64)
			// 	return v % int(op2)
			case lexer.Identifier: // Op1 é int e Op2 está em memória
				op2, _ := i.memory[t2.Lexema]
				switch v2 := op2.(type) {
				case int: // Op1 é int e Op2 é int
					return v % v2
					// case float64: // Op1 é int e Op2 é float
					// 	return v % int(v2)
				}
			}
			// case float64: // Op1 é float
			// 	switch t2.Token {
			// 	case lexer.Literal_int: // Op1 é float e Op2 é int
			// 		op2, _ := strconv.Atoi(t2.Lexema)
			// 		return int(v) % op2
			// 	case lexer.Literal_float: // Op1 é float e Op2 é float
			// 		op2, _ := strconv.ParseFloat(t2.Lexema, 64)
			// 		return int(v) % int(op2)
			// 	case lexer.Identifier: // Op1 é float e Op2 está em memória
			// 		op2, _ := i.memory[t2.Lexema]
			// 		switch v2 := op2.(type) {
			// 		case int: // Op1 é float e Op2 é int
			// 			return int(v) % v2
			// 		case float64: // Op1 é float e Op2 é float
			// 			return int(v) % int(v2)
			// 		}
			// 	}
		}
	default:
		utils.ThrowException("executor.go", "operationMod", "invalid type for operationMod")
	}
	return false
}

func (i *Interpreter) operationDivI(t1 *lexer.TuplaLex, t2 *lexer.TuplaLex) any {
	i.checkDivisionByZero(t2)
	switch t1.Token {
	case lexer.Literal_int: // Op1 é int
		op1, _ := strconv.Atoi(t1.Lexema)
		switch t2.Token {
		case lexer.Literal_int: // Op1 é int e Op2 é int
			op2, _ := strconv.Atoi(t2.Lexema)
			return op1 / op2
		// case lexer.Literal_float: // Op1 é int e Op2 é float
		// 	op2, _ := strconv.ParseFloat(t2.Lexema, 64)
		// 	return op1 / int(op2)
		case lexer.Identifier: // Op1 é int e Op2 está em memória
			op2, _ := i.memory[t2.Lexema]
			switch v := op2.(type) {
			case int: // Op1 é int e Op2 é int
				return op1 / v
				// case float64: // Op1 é int e Op2 é float
				// 	return op1 / int(v)
			}
		}
	// case lexer.Literal_float: // Op1 é float
	// 	op1, _ := strconv.ParseFloat(t1.Lexema, 64)
	// 	switch t2.Token {
	// 	case lexer.Literal_int: // Op1 é float e Op2 é int
	// 		op2, _ := strconv.Atoi(t2.Lexema)
	// 		return int(op1) / op2
	// 	case lexer.Literal_float: // Op1 é float e Op2 é float
	// 		op2, _ := strconv.ParseFloat(t2.Lexema, 64)
	// 		return int(op1) / int(op2)
	// 	case lexer.Identifier: // Op1 é float e Op2 está em memória
	// 		op2, _ := i.memory[t2.Lexema]
	// 		switch v := op2.(type) {
	// 		case int: // Op1 é float e Op2 é int
	// 			return int(op1) / v
	// 		case float64: // Op1 é float e Op2 é float
	// 			return int(op1) / int(v)
	// 		}
	// 	}
	case lexer.Identifier: // Op1 esta em memoria
		op1, _ := i.memory[t1.Lexema]
		switch v := op1.(type) {
		case int: // Op1 é int
			switch t2.Token {
			case lexer.Literal_int: // Op1 é int e Op2 é int
				op2, _ := strconv.Atoi(t2.Lexema)
				return v / op2
			// case lexer.Literal_float: // Op1 é int e Op2 é float
			// 	op2, _ := strconv.ParseFloat(t2.Lexema, 64)
			// 	return v / int(op2)
			case lexer.Identifier: // Op1 é int e Op2 está em memória
				op2, _ := i.memory[t2.Lexema]
				switch v2 := op2.(type) {
				case int: // Op1 é int e Op2 é int
					return v / v2
					// case float64: // Op1 é int e Op2 é float
					// 	return v / int(v2)
				}
			}
			// case float64: // Op1 é float
			// 	switch t2.Token {
			// 	case lexer.Literal_int: // Op1 é float e Op2 é int
			// 		op2, _ := strconv.Atoi(t2.Lexema)
			// 		return int(v) / op2
			// 	case lexer.Literal_float: // Op1 é float e Op2 é float
			// 		op2, _ := strconv.ParseFloat(t2.Lexema, 64)
			// 		return int(v) / int(op2)
			// 	case lexer.Identifier: // Op1 é float e Op2 está em memória
			// 		op2, _ := i.memory[t2.Lexema]
			// 		switch v2 := op2.(type) {
			// 		case int: // Op1 é float e Op2 é int
			// 			return int(v) / v2
			// 		case float64: // Op1 é float e Op2 é float
			// 			return int(v) / int(v2)
			// 		}
			// 	}
		}
	default:
		utils.ThrowException("executor.go", "operationDivI", "invalid type for operationDivI")
	}
	return false
}

func (i *Interpreter) operationEq(t1 *lexer.TuplaLex, t2 *lexer.TuplaLex) bool {
	switch t1.Token {
	case lexer.Literal_int: // Op1 é int
		op1, _ := strconv.Atoi(t1.Lexema)
		switch t2.Token {
		case lexer.Literal_int: // Op1 é int e Op2 é int
			op2, _ := strconv.Atoi(t2.Lexema)
			return op1 == op2
		case lexer.Literal_float: // Op1 é int e Op2 é float
			op2, _ := strconv.ParseFloat(t2.Lexema, 64)
			return float64(op1) == op2
		case lexer.Identifier: // Op1 é int e Op2 está em memória
			op2, _ := i.memory[t2.Lexema]
			switch v := op2.(type) {
			case int: // Op1 é int e Op2 é int
				return op1 == v
			case float64: // Op1 é int e Op2 é float
				return float64(op1) == v
			}
		}
	case lexer.Literal_float: // Op1 é float
		op1, _ := strconv.ParseFloat(t1.Lexema, 64)
		switch t2.Token {
		case lexer.Literal_int: // Op1 é float e Op2 é int
			op2, _ := strconv.Atoi(t2.Lexema)
			return op1 == float64(op2)
		case lexer.Literal_float: // Op1 é float e Op2 é float
			op2, _ := strconv.ParseFloat(t2.Lexema, 64)
			return op1 == op2
		case lexer.Identifier: // Op1 é float e Op2 está em memória
			op2, _ := i.memory[t2.Lexema]
			switch v := op2.(type) {
			case int: // Op1 é float e Op2 é int
				return op1 == float64(v)
			case float64: // Op1 é float e Op2 é float
				return op1 == v
			}
		}
	case lexer.Literal_false: // Op1 é false
		switch t2.Token {
		case lexer.Literal_true:
			return false // Op1 é false e Op2 é true
		case lexer.Literal_false:
			return true // Op1 é false e Op2 é false
		case lexer.Identifier: // Opt 1 é false e Op2 está em memória
			op2, _ := i.memory[t2.Lexema]
			switch op2 {
			case true:
				return false //Op1 é false e Op2 é true
			case false:
				return true //Op1 é false e Op2 é false
			}

		}
	case lexer.Literal_true: // Op1 é true
		switch t2.Token {
		case lexer.Literal_true:
			return true // Op1 é true e Op2 é true
		case lexer.Literal_false:
			return false // Op1 é true e Op2 é false
		case lexer.Identifier: // Opt 1 é true e Op2 está em memória
			op2, _ := i.memory[t2.Lexema]
			switch op2 {
			case true:
				return true //Op1 é true e Op2 é true
			case false:
				return false //Op1 é true e Op2 é false
			}

		}
	case lexer.Literal_string: // Op1 é string
		op1 := t1.Lexema
		switch t2.Token {
		case lexer.Literal_string: // Op1 é string e Op2 é string
			op2 := t2.Lexema
			return op1 == op2
		case lexer.Identifier: // Op1 é string e Op2 está em memória
			op2, _ := i.memory[t2.Lexema]
			switch v := op2.(type) { // Op1 é string e Op2 é string
			case string:
				return op1 == v
			}
		}
	case lexer.Literal_char: // Op1 é char
		op1 := t1.Lexema
		switch t2.Token {
		case lexer.Literal_char: // Op1 é char e Op2 é char
			op2 := t2.Lexema
			return op1 == op2
		case lexer.Identifier: //Op1 é char e Op2 está em memória
			op2, _ := i.memory[t2.Lexema]
			switch v := op2.(type) {
			case string:
				return op1 == v
			}
		}
	case lexer.Identifier: // Op1 esta em memoria
		op1, _ := i.memory[t1.Lexema]
		switch v := op1.(type) {
		case int: // Op1 é int
			switch t2.Token {
			case lexer.Literal_int: // Op1 é int e Op2 é int
				op2, _ := strconv.Atoi(t2.Lexema)
				return v == op2
			case lexer.Literal_float: // Op1 é int e Op2 é float
				op2, _ := strconv.ParseFloat(t2.Lexema, 64)
				return float64(v) == op2
			case lexer.Identifier: // Op1 é int e Op2 está em memória
				op2, _ := i.memory[t2.Lexema]
				switch v2 := op2.(type) {
				case int: // Op1 é int e Op2 é int
					return v == v2
				case float64: // Op1 é int e Op2 é float
					return float64(v) == v2
				}
			}
		case float64: // Op1 é float
			switch t2.Token {
			case lexer.Literal_int: // Op1 é float e Op2 é int
				op2, _ := strconv.Atoi(t2.Lexema)
				return v == float64(op2)
			case lexer.Literal_float: // Op1 é float e Op2 é float
				op2, _ := strconv.ParseFloat(t2.Lexema, 64)
				return v == op2
			case lexer.Identifier: // Op1 é float e Op2 está em memória
				op2, _ := i.memory[t2.Lexema]
				switch v2 := op2.(type) {
				case int: // Op1 é float e Op2 é int
					return v == float64(v2)
				case float64: // Op1 é float e Op2 é float
					return v == v2
				}
			}
		case bool: //Op1 é bool
			switch op1 {
			case false: // Op1 é false
				switch t2.Token {
				case lexer.Literal_true:
					return false // Op1 é false e Op2 é true
				case lexer.Literal_false:
					return true // Op1 é false e Op2 é false
				case lexer.Identifier: // Opt 1 é false e Op2 está em memória
					op2, _ := i.memory[t2.Lexema]
					switch op2 {
					case true:
						return false //Op1 é false e Op2 é true
					case false:
						return true //Op1 é false e Op2 é false
					}
				}
			case true: // Op1 é true
				switch t2.Token {
				case lexer.Literal_true:
					return true // Op1 é true e Op2 é true
				case lexer.Literal_false:
					return false // Op1 é true e Op2 é false
				case lexer.Identifier: // Opt 1 é true e Op2 está em memória
					op2, _ := i.memory[t2.Lexema]
					switch op2 {
					case true:
						return true //Op1 é true e Op2 é true
					case false:
						return false //Op1 é true e Op2 é false
					}

				}
			}
		case string: // Op1 é string
			switch t2.Token {
			case lexer.Literal_string: // Op1 é string e Op2 é string
				op2 := t2.Lexema
				return op1 == op2
			case lexer.Identifier: // Op1 é string e Op2 está em memória
				op2, _ := i.memory[t2.Lexema]
				switch v := op2.(type) { // Op1 é string e Op2 é string
				case string:
					return op1 == v
				}
			}
		}
	default:
		utils.ThrowException("executor.go", "operationEq", "invalid token type for operationEq")
	}
	return false
}

func (i *Interpreter) operationNeq(t1 *lexer.TuplaLex, t2 *lexer.TuplaLex) bool {
	switch t1.Token {
	case lexer.Literal_int: // Op1 é int
		op1, _ := strconv.Atoi(t1.Lexema)
		switch t2.Token {
		case lexer.Literal_int: // Op1 é int e Op2 é int
			op2, _ := strconv.Atoi(t2.Lexema)
			return op1 != op2
		case lexer.Literal_float: // Op1 é int e Op2 é float
			op2, _ := strconv.ParseFloat(t2.Lexema, 64)
			return float64(op1) != op2
		case lexer.Identifier: // Op1 é int e Op2 está em memória
			op2, _ := i.memory[t2.Lexema]
			switch v := op2.(type) {
			case int: // Op1 é int e Op2 é int
				return op1 != v
			case float64: // Op1 é int e Op2 é float
				return float64(op1) != v
			}
		}
	case lexer.Literal_float: // Op1 é float
		op1, _ := strconv.ParseFloat(t1.Lexema, 64)
		switch t2.Token {
		case lexer.Literal_int: // Op1 é float e Op2 é int
			op2, _ := strconv.Atoi(t2.Lexema)
			return op1 != float64(op2)
		case lexer.Literal_float: // Op1 é float e Op2 é float
			op2, _ := strconv.ParseFloat(t2.Lexema, 64)
			return op1 != op2
		case lexer.Identifier: // Op1 é float e Op2 está em memória
			op2, _ := i.memory[t2.Lexema]
			switch v := op2.(type) {
			case int: // Op1 é float e Op2 é int
				return op1 != float64(v)
			case float64: // Op1 é float e Op2 é float
				return op1 != v
			}
		}
	case lexer.Literal_false: // Op1 é false
		switch t2.Token {
		case lexer.Literal_true:
			return true // Op1 é false e Op2 é true
		case lexer.Literal_false:
			return false // Op1 é false e Op2 é false
		case lexer.Identifier: // Opt 1 é false e Op2 está em memória
			op2, _ := i.memory[t2.Lexema]
			switch op2 {
			case true:
				return true //Op1 é false e Op2 é true
			case false:
				return false //Op1 é false e Op2 é false
			}

		}
	case lexer.Literal_true: // Op1 é true
		switch t2.Token {
		case lexer.Literal_true:
			return false // Op1 é true e Op2 é true
		case lexer.Literal_false:
			return true // Op1 é true e Op2 é false
		case lexer.Identifier: // Opt 1 é true e Op2 está em memória
			op2, _ := i.memory[t2.Lexema]
			switch op2 {
			case true:
				return false //Op1 é true e Op2 é true
			case false:
				return true //Op1 é true e Op2 é false
			}

		}
	case lexer.Literal_string: // Op1 é string
		op1 := t1.Lexema
		switch t2.Token {
		case lexer.Literal_string: // Op1 é string e Op2 é string
			op2 := t2.Lexema
			return op1 != op2
		case lexer.Identifier: // Op1 é string e Op2 está em memória
			op2, _ := i.memory[t2.Lexema]
			switch v := op2.(type) { // Op1 é string e Op2 é string
			case string:
				return op1 != v
			}
		}
	case lexer.Literal_char: // Op1 é char
		op1 := t1.Lexema
		switch t2.Token {
		case lexer.Literal_char: // Op1 é char e Op2 é char
			op2 := t2.Lexema
			return op1 != op2
		case lexer.Identifier: //Op1 é char e Op2 está em memória
			op2, _ := i.memory[t2.Lexema]
			switch v := op2.(type) {
			case string:
				return op1 != v
			}
		}
	case lexer.Identifier: // Op1 esta em memoria
		op1, _ := i.memory[t1.Lexema]
		switch v := op1.(type) {
		case int: // Op1 é int
			switch t2.Token {
			case lexer.Literal_int: // Op1 é int e Op2 é int
				op2, _ := strconv.Atoi(t2.Lexema)
				return v != op2
			case lexer.Literal_float: // Op1 é int e Op2 é float
				op2, _ := strconv.ParseFloat(t2.Lexema, 64)
				return float64(v) != op2
			case lexer.Identifier: // Op1 é int e Op2 está em memória
				op2, _ := i.memory[t2.Lexema]
				switch v2 := op2.(type) {
				case int: // Op1 é int e Op2 é int
					return v != v2
				case float64: // Op1 é int e Op2 é float
					return float64(v) != v2
				}
			}
		case float64: // Op1 é float
			switch t2.Token {
			case lexer.Literal_int: // Op1 é float e Op2 é int
				op2, _ := strconv.Atoi(t2.Lexema)
				return v != float64(op2)
			case lexer.Literal_float: // Op1 é float e Op2 é float
				op2, _ := strconv.ParseFloat(t2.Lexema, 64)
				return v != op2
			case lexer.Identifier: // Op1 é float e Op2 está em memória
				op2, _ := i.memory[t2.Lexema]
				switch v2 := op2.(type) {
				case int: // Op1 é float e Op2 é int
					return v != float64(v2)
				case float64: // Op1 é float e Op2 é float
					return v != v2
				}
			}
		case bool: //Op1 é bool
			switch op1 {
			case false: // Op1 é false
				switch t2.Token {
				case lexer.Literal_true:
					return true // Op1 é false e Op2 é true
				case lexer.Literal_false:
					return false // Op1 é false e Op2 é false
				case lexer.Identifier: // Opt 1 é false e Op2 está em memória
					op2, _ := i.memory[t2.Lexema]
					switch op2 {
					case true:
						return true //Op1 é false e Op2 é true
					case false:
						return false //Op1 é false e Op2 é false
					}

				}
			case true: // Op1 é true
				switch t2.Token {
				case lexer.Literal_true:
					return false // Op1 é true e Op2 é true
				case lexer.Literal_false:
					return true // Op1 é true e Op2 é false
				case lexer.Identifier: // Opt 1 é true e Op2 está em memória
					op2, _ := i.memory[t2.Lexema]
					switch op2 {
					case true:
						return false //Op1 é true e Op2 é true
					case false:
						return true //Op1 é true e Op2 é false
					}

				}
			}
		case string: // Op1 é string
			switch t2.Token {
			case lexer.Literal_string: // Op1 é string e Op2 é string
				op2 := t2.Lexema
				return op1 != op2
			case lexer.Identifier: // Op1 é string e Op2 está em memória
				op2, _ := i.memory[t2.Lexema]
				switch v := op2.(type) { // Op1 é string e Op2 é string
				case string:
					return op1 != v
				}
			}
		}
	default:
		utils.ThrowException("executor.go", "operationNeq", "invalid token type for operationNeq")
	}
	return false
}

func (i *Interpreter) operationLt(t1 *lexer.TuplaLex, t2 *lexer.TuplaLex) bool {
	switch t1.Token {
	case lexer.Literal_int: // Op1 é int
		op1, _ := strconv.Atoi(t1.Lexema)
		switch t2.Token {
		case lexer.Literal_int: // Op1 é int e Op2 é int
			op2, _ := strconv.Atoi(t2.Lexema)
			return op1 < op2
		case lexer.Literal_float: // Op1 é int e Op2 é float
			op2, _ := strconv.ParseFloat(t2.Lexema, 64)
			return float64(op1) < op2
		case lexer.Identifier: // Op1 é int e Op2 está em memória
			op2, _ := i.memory[t2.Lexema]
			switch v := op2.(type) {
			case int: // Op1 é int e Op2 é int
				return op1 < v
			case float64: // Op1 é int e Op2 é float
				return float64(op1) < v
			}
		}
	case lexer.Literal_float: // Op1 é float
		op1, _ := strconv.ParseFloat(t1.Lexema, 64)
		switch t2.Token {
		case lexer.Literal_int: // Op1 é float e Op2 é int
			op2, _ := strconv.Atoi(t2.Lexema)
			return op1 < float64(op2)
		case lexer.Literal_float: // Op1 é float e Op2 é float
			op2, _ := strconv.ParseFloat(t2.Lexema, 64)
			return op1 < op2
		case lexer.Identifier: // Op1 é float e Op2 está em memória
			op2, _ := i.memory[t2.Lexema]
			switch v := op2.(type) {
			case int: // Op1 é float e Op2 é int
				return op1 < float64(v)
			case float64: // Op1 é float e Op2 é float
				return op1 < v
			}
		}
	case lexer.Identifier: // Op1 esta em memoria
		op1, _ := i.memory[t1.Lexema]
		switch v := op1.(type) {
		case int: // Op1 é int
			switch t2.Token {
			case lexer.Literal_int: // Op1 é int e Op2 é int
				op2, _ := strconv.Atoi(t2.Lexema)
				return v < op2
			case lexer.Literal_float: // Op1 é int e Op2 é float
				op2, _ := strconv.ParseFloat(t2.Lexema, 64)
				return float64(v) < op2
			case lexer.Identifier: // Op1 é int e Op2 está em memória
				op2, _ := i.memory[t2.Lexema]
				switch v2 := op2.(type) {
				case int: // Op1 é int e Op2 é int
					return v < v2
				case float64: // Op1 é int e Op2 é float
					return float64(v) < v2
				}
			}
		case float64: // Op1 é float
			switch t2.Token {
			case lexer.Literal_int: // Op1 é float e Op2 é int
				op2, _ := strconv.Atoi(t2.Lexema)
				return v < float64(op2)
			case lexer.Literal_float: // Op1 é float e Op2 é float
				op2, _ := strconv.ParseFloat(t2.Lexema, 64)
				return v < op2
			case lexer.Identifier: // Op1 é float e Op2 está em memória
				op2, _ := i.memory[t2.Lexema]
				switch v2 := op2.(type) {
				case int: // Op1 é float e Op2 é int
					return v < float64(v2)
				case float64: // Op1 é float e Op2 é float
					return v < v2
				}
			}
		}
	default:
		utils.ThrowException("executor.go", "operationLt", "invalid token type for operationLt")
	}
	return false
}

func (i *Interpreter) operationGt(t1 *lexer.TuplaLex, t2 *lexer.TuplaLex) bool {
	switch t1.Token {
	case lexer.Literal_int: // Op1 é int
		op1, _ := strconv.Atoi(t1.Lexema)
		switch t2.Token {
		case lexer.Literal_int: // Op1 é int e Op2 é int
			op2, _ := strconv.Atoi(t2.Lexema)
			return op1 > op2
		case lexer.Literal_float: // Op1 é int e Op2 é float
			op2, _ := strconv.ParseFloat(t2.Lexema, 64)
			return float64(op1) > op2
		case lexer.Identifier: // Op1 é int e Op2 está em memória
			op2, _ := i.memory[t2.Lexema]
			switch v := op2.(type) {
			case int: // Op1 é int e Op2 é int
				return op1 > v
			case float64: // Op1 é int e Op2 é float
				return float64(op1) > v
			}
		}
	case lexer.Literal_float: // Op1 é float
		op1, _ := strconv.ParseFloat(t1.Lexema, 64)
		switch t2.Token {
		case lexer.Literal_int: // Op1 é float e Op2 é int
			op2, _ := strconv.Atoi(t2.Lexema)
			return op1 > float64(op2)
		case lexer.Literal_float: // Op1 é float e Op2 é float
			op2, _ := strconv.ParseFloat(t2.Lexema, 64)
			return op1 > op2
		case lexer.Identifier: // Op1 é float e Op2 está em memória
			op2, _ := i.memory[t2.Lexema]
			switch v := op2.(type) {
			case int: // Op1 é float e Op2 é int
				return op1 > float64(v)
			case float64: // Op1 é float e Op2 é float
				return op1 > v
			}
		}
	case lexer.Identifier: // Op1 esta em memoria
		op1, _ := i.memory[t1.Lexema]
		switch v := op1.(type) {
		case int: // Op1 é int
			switch t2.Token {
			case lexer.Literal_int: // Op1 é int e Op2 é int
				op2, _ := strconv.Atoi(t2.Lexema)
				return v > op2
			case lexer.Literal_float: // Op1 é int e Op2 é float
				op2, _ := strconv.ParseFloat(t2.Lexema, 64)
				return float64(v) > op2
			case lexer.Identifier: // Op1 é int e Op2 está em memória
				op2, _ := i.memory[t2.Lexema]
				switch v2 := op2.(type) {
				case int: // Op1 é int e Op2 é int
					return v > v2
				case float64: // Op1 é int e Op2 é float
					return float64(v) > v2
				}
			}
		case float64: // Op1 é float
			switch t2.Token {
			case lexer.Literal_int: // Op1 é float e Op2 é int
				op2, _ := strconv.Atoi(t2.Lexema)
				return v > float64(op2)
			case lexer.Literal_float: // Op1 é float e Op2 é float
				op2, _ := strconv.ParseFloat(t2.Lexema, 64)
				return v > op2
			case lexer.Identifier: // Op1 é float e Op2 está em memória
				op2, _ := i.memory[t2.Lexema]
				switch v2 := op2.(type) {
				case int: // Op1 é float e Op2 é int
					return v > float64(v2)
				case float64: // Op1 é float e Op2 é float
					return v > v2
				}
			}
		}
	default:
		utils.ThrowException("executor.go", "operationGt", "invalid token type for operationGt")
	}
	return false
}

func (i *Interpreter) operationLte(t1 *lexer.TuplaLex, t2 *lexer.TuplaLex) bool {
	switch t1.Token {
	case lexer.Literal_int: // Op1 é int
		op1, _ := strconv.Atoi(t1.Lexema)
		switch t2.Token {
		case lexer.Literal_int: // Op1 é int e Op2 é int
			op2, _ := strconv.Atoi(t2.Lexema)
			return op1 <= op2
		case lexer.Literal_float: // Op1 é int e Op2 é float
			op2, _ := strconv.ParseFloat(t2.Lexema, 64)
			return float64(op1) <= op2
		case lexer.Identifier: // Op1 é int e Op2 está em memória
			op2, _ := i.memory[t2.Lexema]
			switch v := op2.(type) {
			case int: // Op1 é int e Op2 é int
				return op1 <= v
			case float64: // Op1 é int e Op2 é float
				return float64(op1) <= v
			}
		}
	case lexer.Literal_float: // Op1 é float
		op1, _ := strconv.ParseFloat(t1.Lexema, 64)
		switch t2.Token {
		case lexer.Literal_int: // Op1 é float e Op2 é int
			op2, _ := strconv.Atoi(t2.Lexema)
			return op1 <= float64(op2)
		case lexer.Literal_float: // Op1 é float e Op2 é float
			op2, _ := strconv.ParseFloat(t2.Lexema, 64)
			return op1 <= op2
		case lexer.Identifier: // Op1 é float e Op2 está em memória
			op2, _ := i.memory[t2.Lexema]
			switch v := op2.(type) {
			case int: // Op1 é float e Op2 é int
				return op1 <= float64(v)
			case float64: // Op1 é float e Op2 é float
				return op1 <= v
			}
		}
	case lexer.Identifier: // Op1 esta em memoria
		op1, _ := i.memory[t1.Lexema]
		switch v := op1.(type) {
		case int: // Op1 é int
			switch t2.Token {
			case lexer.Literal_int: // Op1 é int e Op2 é int
				op2, _ := strconv.Atoi(t2.Lexema)
				return v <= op2
			case lexer.Literal_float: // Op1 é int e Op2 é float
				op2, _ := strconv.ParseFloat(t2.Lexema, 64)
				return float64(v) <= op2
			case lexer.Identifier: // Op1 é int e Op2 está em memória
				op2, _ := i.memory[t2.Lexema]
				switch v2 := op2.(type) {
				case int: // Op1 é int e Op2 é int
					return v <= v2
				case float64: // Op1 é int e Op2 é float
					return float64(v) <= v2
				}
			}
		case float64: // Op1 é float
			switch t2.Token {
			case lexer.Literal_int: // Op1 é float e Op2 é int
				op2, _ := strconv.Atoi(t2.Lexema)
				return v <= float64(op2)
			case lexer.Literal_float: // Op1 é float e Op2 é float
				op2, _ := strconv.ParseFloat(t2.Lexema, 64)
				return v <= op2
			case lexer.Identifier: // Op1 é float e Op2 está em memória
				op2, _ := i.memory[t2.Lexema]
				switch v2 := op2.(type) {
				case int: // Op1 é float e Op2 é int
					return v <= float64(v2)
				case float64: // Op1 é float e Op2 é float
					return v <= v2
				}
			}
		}
	default:
		utils.ThrowException("executor.go", "operationLte", "invalid token type for operationLte")
	}
	return false
}

func (i *Interpreter) operationGte(t1 *lexer.TuplaLex, t2 *lexer.TuplaLex) bool {
	switch t1.Token {
	case lexer.Literal_int: // Op1 é int
		op1, _ := strconv.Atoi(t1.Lexema)
		switch t2.Token {
		case lexer.Literal_int: // Op1 é int e Op2 é int
			op2, _ := strconv.Atoi(t2.Lexema)
			return op1 >= op2
		case lexer.Literal_float: // Op1 é int e Op2 é float
			op2, _ := strconv.ParseFloat(t2.Lexema, 64)
			return float64(op1) >= op2
		case lexer.Identifier: // Op1 é int e Op2 está em memória
			op2, _ := i.memory[t2.Lexema]
			switch v := op2.(type) {
			case int: // Op1 é int e Op2 é int
				return op1 >= v
			case float64: // Op1 é int e Op2 é float
				return float64(op1) >= v
			}
		}
	case lexer.Literal_float: // Op1 é float
		op1, _ := strconv.ParseFloat(t1.Lexema, 64)
		switch t2.Token {
		case lexer.Literal_int: // Op1 é float e Op2 é int
			op2, _ := strconv.Atoi(t2.Lexema)
			return op1 >= float64(op2)
		case lexer.Literal_float: // Op1 é float e Op2 é float
			op2, _ := strconv.ParseFloat(t2.Lexema, 64)
			return op1 >= op2
		case lexer.Identifier: // Op1 é float e Op2 está em memória
			op2, _ := i.memory[t2.Lexema]
			switch v := op2.(type) {
			case int: // Op1 é float e Op2 é int
				return op1 >= float64(v)
			case float64: // Op1 é float e Op2 é float
				return op1 >= v
			}
		}
	case lexer.Identifier: // Op1 esta em memoria
		op1, _ := i.memory[t1.Lexema]
		switch v := op1.(type) {
		case int: // Op1 é int
			switch t2.Token {
			case lexer.Literal_int: // Op1 é int e Op2 é int
				op2, _ := strconv.Atoi(t2.Lexema)
				return v >= op2
			case lexer.Literal_float: // Op1 é int e Op2 é float
				op2, _ := strconv.ParseFloat(t2.Lexema, 64)
				return float64(v) >= op2
			case lexer.Identifier: // Op1 é int e Op2 está em memória
				op2, _ := i.memory[t2.Lexema]
				switch v2 := op2.(type) {
				case int: // Op1 é int e Op2 é int
					return v >= v2
				case float64: // Op1 é int e Op2 é float
					return float64(v) >= v2
				}
			}
		case float64: // Op1 é float
			switch t2.Token {
			case lexer.Literal_int: // Op1 é float e Op2 é int
				op2, _ := strconv.Atoi(t2.Lexema)
				return v >= float64(op2)
			case lexer.Literal_float: // Op1 é float e Op2 é float
				op2, _ := strconv.ParseFloat(t2.Lexema, 64)
				return v >= op2
			case lexer.Identifier: // Op1 é float e Op2 está em memória
				op2, _ := i.memory[t2.Lexema]
				switch v2 := op2.(type) {
				case int: // Op1 é float e Op2 é int
					return v >= float64(v2)
				case float64: // Op1 é float e Op2 é float
					return v >= v2
				}
			}
		}
	default:
		utils.ThrowException("executor.go", "operationGte", "invalid token type for operationGte")
	}
	return false
}

func (i *Interpreter) operationAnd(t1 *lexer.TuplaLex, t2 *lexer.TuplaLex) bool {
	switch t1.Token {
	case lexer.Literal_false: // Op1 é false
		return false
	case lexer.Literal_true: // Op1 é true
		switch t2.Token {
		case lexer.Literal_false: // Op1 é true e Op2 é false
			return false
		case lexer.Literal_true: // Op1 é true e Op2 é true
			return true
		case lexer.Identifier: // Op1 é true e Op2 está em memória
			op2, _ := i.memory[t2.Lexema]
			switch op2 {
			case false: // Op1 é true e Op2 é false
				return false
			case true: // Op1 é true e Op2 é true
				return true
			}
		}
	case lexer.Identifier: // Op1 esta em memoria
		op1, _ := i.memory[t1.Lexema]
		switch op1 {
		case false: // Op1 é false
			return false
		case true: // Op1 é true
			switch t2.Token {
			case lexer.Literal_false: // Op1 é true e Op2 é false
				return false
			case lexer.Literal_true: // Op1 é true e Op2 é true
				return true
			case lexer.Identifier: // Op1 é true e Op2 está em memória
				op2, _ := i.memory[t2.Lexema]
				switch op2 {
				case false: // Op1 é true e Op2 é false
					return false
				case true: // Op1 é true e Op2 é true
					return true
				}
			}
		}
	default:
		utils.ThrowException("executor.go", "operationAnd", "invalid token type for operationAnd")
	}
	return false
}

func (i *Interpreter) operationOr(t1 *lexer.TuplaLex, t2 *lexer.TuplaLex) bool {
	switch t1.Token {
	case lexer.Literal_false: // Op1 é false
		switch t2.Token {
		case lexer.Literal_false: // Op1 é false e Op2 é false
			return false
		case lexer.Literal_true: // Op1 é false e Op2 é true
			return true
		case lexer.Identifier: // Op1 é false e Op2 está em memória
			op2, _ := i.memory[t2.Lexema]
			switch op2 {
			case false: // Op1 é false e Op2 é false
				return false
			case true: // Op1 é false e Op2 é true
				return true
			}
		}
	case lexer.Literal_true: // Op1 é true
		return true
	case lexer.Identifier: // Op1 esta em memoria
		op1, _ := i.memory[t1.Lexema]
		switch op1 {
		case false: // Op1 é false
			switch t2.Token {
			case lexer.Literal_false: // Op1 é false e Op2 é false
				return false
			case lexer.Literal_true: // Op1 é false e Op2 é true
				return true
			case lexer.Identifier: // Op1 é false e Op2 está em memória
				op2, _ := i.memory[t2.Lexema]
				switch op2 {
				case false: // Op1 é false e Op2 é false
					return false
				case true: // Op1 é false e Op2 é true
					return true
				}
			}
		case true: // Op1 é true
			return true
		}
	default:
		utils.ThrowException("executor.go", "operationOr", "invalid token type for operationOr")
	}
	return false
}

func (i *Interpreter) operationNot(t1 *lexer.TuplaLex) bool {
	switch t1.Token {
	case lexer.Literal_false: // Op1 é false
		return true
	case lexer.Literal_true: // Op1 é true
		return false
	case lexer.Identifier: // Op1 esta em memoria
		op1, _ := i.memory[t1.Lexema]
		switch op1 {
		case false: // Op1 é false
			return true
		case true: // Op1 é true
			return false
		}
	default:
		utils.ThrowException("executor.go", "operationNot", "invalid token type for operationNot")
	}
	return false
}

func (i *Interpreter) operationXor(t1 *lexer.TuplaLex, t2 *lexer.TuplaLex) bool {
	switch t1.Token {
	case lexer.Literal_false: // Op1 é false
		switch t2.Token {
		case lexer.Literal_false: // Op1 é false e Op2 é false
			return false
		case lexer.Literal_true: // Op1 é false e Op2 é true
			return true
		case lexer.Identifier: // Op1 é false e Op2 está em memória
			op2, _ := i.memory[t2.Lexema]
			switch op2 {
			case false: // Op1 é false e Op2 é false
				return false
			case true: // Op1 é false e Op2 é true
				return true
			}
		}
	case lexer.Literal_true: // Op1 é true
		switch t2.Token {
		case lexer.Literal_false: // Op1 é true e Op2 é false
			return true
		case lexer.Literal_true: // Op1 é true e Op2 é true
			return false
		case lexer.Identifier: // Op1 é true e Op2 está em memória
			op2, _ := i.memory[t2.Lexema]
			switch op2 {
			case false: // Op1 é true e Op2 é false
				return true
			case true: // Op1 é true e Op2 é true
				return false
			}
		}
	case lexer.Identifier: // Op1 esta em memoria
		op1, _ := i.memory[t1.Lexema]
		switch op1 {
		case false: // Op1 é false
			switch t2.Token {
			case lexer.Literal_false: // Op1 é false e Op2 é false
				return false
			case lexer.Literal_true: // Op1 é false e Op2 é true
				return true
			case lexer.Identifier: // Op1 é false e Op2 está em memória
				op2, _ := i.memory[t2.Lexema]
				switch op2 {
				case false: // Op1 é false e Op2 é false
					return false
				case true: // Op1 é false e Op2 é true
					return true
				}
			}
		case true: // Op1 é true
			switch t2.Token {
			case lexer.Literal_false: // Op1 é true e Op2 é false
				return true
			case lexer.Literal_true: // Op1 é true e Op2 é true
				return false
			case lexer.Identifier: // Op1 é true e Op2 está em memória
				op2, _ := i.memory[t2.Lexema]
				switch op2 {
				case false: // Op1 é true e Op2 é false
					return true
				case true: // Op1 é true e Op2 é true
					return false
				}
			}
		}
	default:
		utils.ThrowException("executor.go", "operationXor", "invalid token type for operationXor")
	}
	return false
}

func (i *Interpreter) operationCall(res *lexer.TuplaLex, t1 *lexer.TuplaLex, t2 *lexer.TuplaLex) {
	switch res.Lexema {
	case "read":
		entrada, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		entrada = strings.TrimSpace(entrada)

		// Verificar se a variável foi devidamente declarada na memória
		valAnterior, declarada := i.memory[t1.Lexema]
		if !declarada {
			utils.ThrowInterpreterException(
				fmt.Sprintf("Variable '%s' not declared", t1.Lexema),
				t1.Linha,
				t1.Coluna,
			)
		}

		// Verificar compatibilidade do tipo no scan (t2.Token) com o tipo da variável declarada
		switch valAnterior.(type) {
		case int:
			if t2.Token != lexer.Type_int {
				utils.ThrowInterpreterException(
					fmt.Sprintf("Type mismatch: cannot scan type '%s' into variable '%s' of type 'trem_di_numeru'", t2.Lexema, t1.Lexema),
					t1.Linha,
					t1.Coluna,
				)
			}
		case float64:
			if t2.Token != lexer.Type_float {
				utils.ThrowInterpreterException(
					fmt.Sprintf("Type mismatch: cannot scan type '%s' into variable '%s' of type 'trem_cum_virgula'", t2.Lexema, t1.Lexema),
					t1.Linha,
					t1.Coluna,
				)
			}
		case bool:
			if t2.Token != lexer.Type_bool {
				utils.ThrowInterpreterException(
					fmt.Sprintf("Type mismatch: cannot scan type '%s' into variable '%s' of type 'trem_discolhe'", t2.Lexema, t1.Lexema),
					t1.Linha,
					t1.Coluna,
				)
			}
		case string:
			if t2.Token != lexer.Type_string && t2.Token != lexer.Type_char {
				utils.ThrowInterpreterException(
					fmt.Sprintf("Type mismatch: cannot scan type '%s' into variable '%s' of type 'trem_discrita'/'trosso'", t2.Lexema, t1.Lexema),
					t1.Linha,
					t1.Coluna,
				)
			}
		}

		// Fazer o parser do valor inserido pelo usuário de acordo com o tipo do scan
		switch t2.Token {
		case lexer.Type_int:
			v, err := strconv.Atoi(entrada)
			if err != nil {
				utils.ThrowInterpreterException(
					fmt.Sprintf("Invalid input for type '%s': '%s'", t2.Lexema, entrada),
					t1.Linha,
					t1.Coluna,
				)
			}
			i.memory[t1.Lexema] = v
		case lexer.Type_float:
			v, err := strconv.ParseFloat(entrada, 64)
			if err != nil {
				utils.ThrowInterpreterException(
					fmt.Sprintf("Invalid input for type '%s': '%s'", t2.Lexema, entrada),
					t1.Linha,
					t1.Coluna,
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
					t1.Linha,
					t1.Coluna,
				)
			}
		case lexer.Type_char:
			runes := []rune(entrada)
			if len(runes) != 1 {
				utils.ThrowInterpreterException(
					fmt.Sprintf("Invalid input for type '%s' (expected a single character): '%s'", t2.Lexema, entrada),
					t1.Linha,
					t1.Coluna,
				)
			}
			i.memory[t1.Lexema] = string(runes[0])
		case lexer.Type_string:
			i.memory[t1.Lexema] = entrada
		default:
			utils.ThrowInterpreterException(
				fmt.Sprintf("Unsupported scan type: '%s'", t2.Lexema),
				t1.Linha,
				t1.Coluna,
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

func (i *Interpreter) operationIf_eq(nome string, l1 string, l2 string) {
	comparation, _ := i.memory[nome]
	switch comparation {
	case true:
		i.ip = i.labels[l1] - 1
	case false:
		i.ip = i.labels[l2] - 1
	}
}

func (i *Interpreter) operationUno(t1 *lexer.TuplaLex) any {
	switch t1.Token {
	case lexer.Literal_int: // Op1 é int
		op1, _ := strconv.Atoi(t1.Lexema)
		return op1 * (-1)
	case lexer.Literal_float: // Op1 é float
		op1, _ := strconv.ParseFloat(t1.Lexema, 64)
		return op1 * (-1)
	case lexer.Identifier: // Op1 esta em memoria
		op1, _ := i.memory[t1.Lexema]
		switch v := op1.(type) {
		case int: // Op1 é int
			return v * (-1)
		case float64: // Op1 é float
			return v * (-1)
		}
	default:
		utils.ThrowException("executor.go", "operationUno", "invalid type for operationUno")
	}
	return false
}

func (i *Interpreter) operationAtt(nome string, operation *lexer.TuplaLex) {
	i.setMemory(nome, operation)
}

func (interpreter *Interpreter) execute(instrucao parser.TuplaMicrocode) {

	// fmt.Printf("\n\nPonteiro executando operação: %d => ", interpreter.ip)

	switch instrucao.Operation {
	case parser.Add:
		// fmt.Printf("Add\n")
		interpreter.memory[instrucao.Res.Lexema] = interpreter.operationAdd(instrucao.Op1, instrucao.Op2)
	case parser.Sub:
		// fmt.Printf("Sub\n")
		interpreter.memory[instrucao.Res.Lexema] = interpreter.operationSub(instrucao.Op1, instrucao.Op2)
	case parser.Mul:
		// fmt.Printf("Mul\n")
		interpreter.memory[instrucao.Res.Lexema] = interpreter.operationMul(instrucao.Op1, instrucao.Op2)
	case parser.Div:
		// fmt.Printf("Div")
		interpreter.memory[instrucao.Res.Lexema] = interpreter.operationDiv(instrucao.Op1, instrucao.Op2)
	case parser.Mod:
		// fmt.Printf("Mod\n")
		interpreter.memory[instrucao.Res.Lexema] = interpreter.operationMod(instrucao.Op1, instrucao.Op2)
	case parser.DivI:
		// fmt.Printf("DivI")
		interpreter.memory[instrucao.Res.Lexema] = interpreter.operationDivI(instrucao.Op1, instrucao.Op2)
	case parser.Eq:
		// fmt.Printf("Eq\n")
		interpreter.memory[instrucao.Res.Lexema] = interpreter.operationEq(instrucao.Op1, instrucao.Op2)
	case parser.Neq:
		// fmt.Printf("Neq\n")
		interpreter.memory[instrucao.Res.Lexema] = interpreter.operationNeq(instrucao.Op1, instrucao.Op2)
	case parser.Lt:
		// fmt.Printf("Lt\n")
		interpreter.memory[instrucao.Res.Lexema] = interpreter.operationLt(instrucao.Op1, instrucao.Op2)
	case parser.Gt:
		// fmt.Printf("Gt\n")
		interpreter.memory[instrucao.Res.Lexema] = interpreter.operationGt(instrucao.Op1, instrucao.Op2)
	case parser.Lte:
		// fmt.Printf("Lte\n")
		interpreter.memory[instrucao.Res.Lexema] = interpreter.operationLte(instrucao.Op1, instrucao.Op2)
	case parser.Gte:
		// fmt.Printf("Gte\n")
		interpreter.memory[instrucao.Res.Lexema] = interpreter.operationGte(instrucao.Op1, instrucao.Op2)
	case parser.And:
		// fmt.Printf("And\n")
		interpreter.memory[instrucao.Res.Lexema] = interpreter.operationAnd(instrucao.Op1, instrucao.Op2)
	case parser.Or:
		// fmt.Printf("Or\n")
		interpreter.memory[instrucao.Res.Lexema] = interpreter.operationOr(instrucao.Op1, instrucao.Op2)
	case parser.Not:
		// fmt.Printf("Not\n")
		interpreter.memory[instrucao.Res.Lexema] = interpreter.operationNot(instrucao.Op1)
	case parser.Xor:
		// fmt.Printf("Xor\n")
		interpreter.memory[instrucao.Res.Lexema] = interpreter.operationXor(instrucao.Op1, instrucao.Op2)
	case parser.Call:
		// fmt.Printf("Call\n")
		interpreter.operationCall(instrucao.Res, instrucao.Op1, instrucao.Op2)
	case parser.Jump:
		// fmt.Printf("Jump\n")
		interpreter.operationJump(instrucao.Res.Lexema)
	case parser.Label:
		// fmt.Printf("Label\n")
	case parser.If_eq:
		// fmt.Printf("If_eq")
		interpreter.operationIf_eq(instrucao.Res.Lexema, instrucao.Op1.Lexema, instrucao.Op2.Lexema)
	case parser.Uno:
		// fmt.Printf("Uno")
		interpreter.memory[instrucao.Res.Lexema] = interpreter.operationUno(instrucao.Op1)
	case parser.Att:
		// fmt.Printf("Att\n")
		interpreter.operationAtt(instrucao.Res.Lexema, instrucao.Op1)
	}

	// fmt.Printf("\n\n")
	// for nome, valor := range interpreter.memory {
	// 	fmt.Printf("%s -> %v, %T\n", nome, valor, valor)
	// }
}
