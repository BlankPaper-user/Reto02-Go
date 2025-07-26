package main

import (
	"fmt"
	"strconv"
	"unicode"
)

// Parser estructura para mantener el estado del parsing
type Parser struct {
	input string
	index int
	line  int
	col   int
}

// ParseJSON es la función principal que iniciará el proceso de deserialización.
func (p *Parser) ParseJSON(input string) (interface{}, error) {
	p.input = input
	p.index = 0
	p.line = 1
	p.col = 1

	// Inicia el parseo llamando a parseValue, que manejará la lógica de los diferentes tipos JSON.
	result, err := p.parseValue()
	if err != nil {
		return nil, err
	}

	// Después de parsear, asegúrate de que no haya caracteres extra
	p.skipWhitespace()
	if p.index < len(p.input) {
		return nil, fmt.Errorf("caracteres inesperados al final de la entrada en línea %d, columna %d", p.line, p.col)
	}

	return result, nil
}

// advance avanza el índice y actualiza línea y columna
func (p *Parser) advance() {
	if p.index < len(p.input) {
		if p.input[p.index] == '\n' {
			p.line++
			p.col = 1
		} else {
			p.col++
		}
		p.index++
	}
}

// parseValue maneja diferentes tipos de valores en JSON: objetos, arrays, strings, números, etc.
func (p *Parser) parseValue() (interface{}, error) {
	// Ignorar espacios en blanco
	p.skipWhitespace()

	if p.index >= len(p.input) {
		return nil, fmt.Errorf("fin inesperado de la entrada JSON en línea %d, columna %d", p.line, p.col)
	}

	// Leer el primer carácter
	ch := rune(p.input[p.index])

	// Según el primer carácter, decidimos cómo procesar el valor.
	switch ch {
	case '{':
		p.advance()
		return p.parseObject()
	case '[':
		p.advance()
		return p.parseArray()
	case '"':
		return p.parseString()
	case 't', 'f':
		return p.parseBoolean()
	case 'n':
		return p.parseNull()
	default:
		if unicode.IsDigit(ch) || ch == '-' {
			return p.parseNumber()
		}
		return nil, fmt.Errorf("carácter inesperado: '%c' en línea %d, columna %d", ch, p.line, p.col)
	}
}

// skipWhitespace avanza el índice ignorando los espacios en blanco
func (p *Parser) skipWhitespace() {
	for p.index < len(p.input) && unicode.IsSpace(rune(p.input[p.index])) {
		p.advance()
	}
}

// parseObject maneja los objetos JSON, que son de la forma { "clave": valor, ... }
func (p *Parser) parseObject() (map[string]interface{}, error) {
	obj := make(map[string]interface{})

	p.skipWhitespace()
	if p.index < len(p.input) && p.input[p.index] == '}' {
		p.advance()
		return obj, nil
	}

	for {
		p.skipWhitespace()
		if p.index >= len(p.input) || p.input[p.index] != '"' {
			return nil, fmt.Errorf("se esperaba '\"' para la clave del objeto en línea %d, columna %d", p.line, p.col)
		}

		key, err := p.parseString()
		if err != nil {
			return nil, err
		}

		// Validar que la clave no esté duplicada
		if _, exists := obj[key]; exists {
			return nil, fmt.Errorf("clave duplicada '%s' en línea %d, columna %d", key, p.line, p.col)
		}

		// Ignoramos los espacios en blanco y el siguiente carácter debe ser ":"
		p.skipWhitespace()
		if p.index >= len(p.input) || p.input[p.index] != ':' {
			return nil, fmt.Errorf("se esperaba ':' después de la clave en línea %d, columna %d", p.line, p.col)
		}
		p.advance() // Consume colon

		// Parseamos el valor
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		// Almacenamos la clave y el valor en el mapa
		obj[key] = value

		// Ignoramos los espacios en blanco y luego comprobamos si es la coma o el final del objeto
		p.skipWhitespace()
		if p.index >= len(p.input) {
			return nil, fmt.Errorf("objeto no terminado en línea %d, columna %d", p.line, p.col)
		}

		ch := p.input[p.index]
		if ch == '}' {
			p.advance()
			break
		}
		if ch != ',' {
			return nil, fmt.Errorf("se esperaba ',' o '}' pero se obtuvo: '%c' en línea %d, columna %d", ch, p.line, p.col)
		}
		p.advance() // Consume comma

		// Verificar que después de la coma no venga inmediatamente '}'
		p.skipWhitespace()
		if p.index < len(p.input) && p.input[p.index] == '}' {
			return nil, fmt.Errorf("coma extra antes de '}' en línea %d, columna %d", p.line, p.col)
		}
	}

	return obj, nil
}

// parseArray maneja los arrays JSON, que son de la forma [ valor, valor, ... ]
func (p *Parser) parseArray() ([]interface{}, error) {
	arr := []interface{}{}

	p.skipWhitespace()
	if p.index < len(p.input) && p.input[p.index] == ']' {
		p.advance()
		return arr, nil
	}

	for {
		// Parseamos el valor dentro del array
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		// Agregamos el valor al array
		arr = append(arr, value)

		// Ignoramos los espacios y luego comprobamos si es la coma o el final del array
		p.skipWhitespace()
		if p.index >= len(p.input) {
			return nil, fmt.Errorf("array no terminado en línea %d, columna %d", p.line, p.col)
		}
		ch := p.input[p.index]
		if ch == ']' {
			p.advance()
			break
		}
		if ch != ',' {
			return nil, fmt.Errorf("se esperaba ',' o ']' pero se obtuvo: '%c' en línea %d, columna %d", ch, p.line, p.col)
		}
		p.advance() // Consume comma

		// Verificar que después de la coma no venga inmediatamente ']'
		p.skipWhitespace()
		if p.index < len(p.input) && p.input[p.index] == ']' {
			return nil, fmt.Errorf("coma extra antes de ']' en línea %d, columna %d", p.line, p.col)
		}
	}

	return arr, nil
}

// parseString maneja las cadenas JSON (delimitadas por comillas)
func (p *Parser) parseString() (string, error) {
	if p.index >= len(p.input) || p.input[p.index] != '"' {
		return "", fmt.Errorf("se esperaba '\"' al inicio de la cadena en línea %d, columna %d", p.line, p.col)
	}
	p.advance() // consume opening quote

	var result string
	for {
		if p.index >= len(p.input) {
			return "", fmt.Errorf("cadena de texto no terminada en línea %d, columna %d", p.line, p.col)
		}
		ch := rune(p.input[p.index])

		// Si encontramos una comilla de cierre, hemos terminado de leer la cadena
		if ch == '"' {
			p.advance()
			break
		}

		// Si encontramos un carácter de escape (\), lo manejamos
		if ch == '\\' {
			p.advance()
			if p.index >= len(p.input) {
				return "", fmt.Errorf("cadena de texto no terminada en línea %d, columna %d", p.line, p.col)
			}
			escapedChar := rune(p.input[p.index])
			p.advance()

			switch escapedChar {
			case '"':
				result += "\""
			case '\\':
				result += "\\"
			case '/':
				result += "/"
			case 'b':
				result += "\b"
			case 'f':
				result += "\f"
			case 'n':
				result += "\n"
			case 'r':
				result += "\r"
			case 't':
				result += "\t"
			case 'u':
				return "", fmt.Errorf("secuencias de escape unicode no implementadas en línea %d, columna %d", p.line, p.col)
			default:
				return "", fmt.Errorf("secuencia de escape inválida: \\%c en línea %d, columna %d", escapedChar, p.line, p.col)
			}
		} else {
			// Si no es un carácter escapado, agregamos el carácter a la cadena
			result += string(ch)
			p.advance()
		}
	}

	return result, nil
}

// parseNumber maneja los números JSON (enteros o flotantes)
func (p *Parser) parseNumber() (interface{}, error) {
	start := p.index
	startLine := p.line
	startCol := p.col

	// Manejar signo negativo
	if p.index < len(p.input) && p.input[p.index] == '-' {
		p.advance()
	}

	// Debe haber al menos un dígito después del signo
	if p.index >= len(p.input) || !unicode.IsDigit(rune(p.input[p.index])) {
		return nil, fmt.Errorf("número inválido en línea %d, columna %d", startLine, startCol)
	}

	// Lee la parte entera
	for p.index < len(p.input) && unicode.IsDigit(rune(p.input[p.index])) {
		p.advance()
	}

	// Manejar parte decimal
	if p.index < len(p.input) && p.input[p.index] == '.' {
		p.advance()
		if p.index >= len(p.input) || !unicode.IsDigit(rune(p.input[p.index])) {
			return nil, fmt.Errorf("número decimal inválido en línea %d, columna %d", startLine, startCol)
		}
		for p.index < len(p.input) && unicode.IsDigit(rune(p.input[p.index])) {
			p.advance()
		}
	}

	// Manejar notación científica
	if p.index < len(p.input) && (p.input[p.index] == 'e' || p.input[p.index] == 'E') {
		p.advance()
		if p.index < len(p.input) && (p.input[p.index] == '+' || p.input[p.index] == '-') {
			p.advance()
		}
		if p.index >= len(p.input) || !unicode.IsDigit(rune(p.input[p.index])) {
			return nil, fmt.Errorf("número en notación científica inválido en línea %d, columna %d", startLine, startCol)
		}
		for p.index < len(p.input) && unicode.IsDigit(rune(p.input[p.index])) {
			p.advance()
		}
	}

	// Subcadena del número encontrado
	numberStr := p.input[start:p.index]

	// Intentamos convertir el número a un flotante (en JSON siempre son float64)
	number, err := strconv.ParseFloat(numberStr, 64)
	if err != nil {
		return nil, fmt.Errorf("número inválido: %s en línea %d, columna %d", numberStr, startLine, startCol)
	}

	return number, nil
}

// parseBoolean maneja los valores booleanos JSON (true, false)
func (p *Parser) parseBoolean() (bool, error) {
	startLine := p.line
	startCol := p.col

	// Lee los siguientes 4 caracteres para "true" o 5 caracteres para "false"
	if p.index+4 <= len(p.input) && p.input[p.index:p.index+4] == "true" {
		for i := 0; i < 4; i++ {
			p.advance()
		}
		return true, nil
	}
	if p.index+5 <= len(p.input) && p.input[p.index:p.index+5] == "false" {
		for i := 0; i < 5; i++ {
			p.advance()
		}
		return false, nil
	}

	return false, fmt.Errorf("valor booleano inválido en línea %d, columna %d", startLine, startCol)
}

// parseNull maneja el valor null JSON, que se convierte en nil en Go
func (p *Parser) parseNull() (interface{}, error) {
	startLine := p.line
	startCol := p.col

	// Lee los 4 caracteres para "null"
	if p.index+4 <= len(p.input) && p.input[p.index:p.index+4] == "null" {
		for i := 0; i < 4; i++ {
			p.advance()
		}
		return nil, nil
	}

	return nil, fmt.Errorf("valor nulo inválido en línea %d, columna %d", startLine, startCol)
}
