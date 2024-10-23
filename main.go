package main

import (
	"fmt"
	"strconv"
	"strings"
)

func Calc(expression string) (float64, error) {
	expression = strings.ReplaceAll(expression, " ", "")
	rpn, err := toRPN(expression)
	if err != nil {
		return 0, err
	}
	return calculateRPN(rpn)
}

// Преобразуем строку выражения в стек записей обратной польской нотации
func toRPN(expression string) ([]string, error) {

	var output []string
	var operators []rune

	// Создаем мапу операций и их приоритетов в виде индексов
	precedence := map[rune]int{
		'*': 3,
		'/': 3,
		'+': 2,
		'-': 2,
		'(': 1,
	}

	// Разбираем в цикле символы из строки, заполняем слайс output значениями и операциями согласно обратной нотации(избавляемся от скобок)
	for i := 0; i < len(expression); i++ {

		char := rune(expression[i])

		switch {
		// Проверяем это число или точка(дробной части), разбираем следующие символы и складываем из них число, добавляем в стек output
		case char >= '0' && char <= '9' || char == '.':
			j := i
			for j < len(expression) && (expression[j] >= '0' && expression[j] <= '9' || expression[j] == '.') {
				j++
			}
			output = append(output, expression[i:j])
			i = j - 1

		// Если символ это знак операции то добавляем его в стек output и удаляем из буферного слайса operators
		case char == '+' || char == '-' || char == '*' || char == '/':
			for len(operators) > 0 && precedence[operators[len(operators)-1]] >= precedence[char] {
				output = append(output, string(operators[len(operators)-1]))
				operators = operators[:len(operators)-1]
			}
			operators = append(operators, char)

		// Если символ это открывающая скобка то добавляем его в слайс operators, для дальнейшего разбора
		case char == '(':
			operators = append(operators, char)

		// Если символ это закрывающая скобка то разбираем слайс operators
		case char == ')':
			for len(operators) > 0 && operators[len(operators)-1] != '(' {
				output = append(output, string(operators[len(operators)-1]))
				operators = operators[:len(operators)-1]
			}
			if len(operators) == 0 {
				return nil, fmt.Errorf("mismatched parentheses")
			}
			operators = operators[:len(operators)-1]

		default:
			return nil, fmt.Errorf("invalid character: %c", char)
		}
	}

	// Разбираем оставшиеся символы операций из слайса operators
	for len(operators) > 0 {
		if operators[len(operators)-1] == '(' {
			return nil, fmt.Errorf("mismatched parentheses")
		}
		output = append(output, string(operators[len(operators)-1]))
		operators = operators[:len(operators)-1]
	}

	return output, nil
}

// Выполняем вычисления по стеку в обратной польской нотации
func calculateRPN(rpn []string) (float64, error) {
	var stack []float64

	for _, token := range rpn {
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			stack = append(stack, num)
		} else {
			if len(stack) < 2 {
				return 0, fmt.Errorf("invalid RPN expression")
			}
			n2 := stack[len(stack)-1]
			n1 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			switch token {
			case "+":
				stack = append(stack, n1+n2)
			case "-":
				stack = append(stack, n1-n2)
			case "*":
				stack = append(stack, n1*n2)
			case "/":
				if n2 == 0 {
					return 0, fmt.Errorf("division by zero")
				}
				stack = append(stack, n1/n2)
			default:
				return 0, fmt.Errorf("invalid operator: %s", token)
			}
		}
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("invalid RPN expression")
	}

	return stack[0], nil
}

func main() {
	fmt.Println(Calc("20.5-1+3*(6*2+4.5)/2"))
}
