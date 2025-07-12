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
}

// ParseJSON es la función principal que iniciará el proceso de deserialización.
func (p *Parser) ParseJSON(input string) (interface{}, error) {
	p.input = input
	p.index = 0

	// Inicia el parseo llamando a parseValue, que manejará la lógica de los diferentes tipos JSON.
	result, err := p.parseValue()
	if err != nil {
		return nil, err
	}

	// Después de parsear, asegúrate de que no haya caracteres extra
	p.skipWhitespace()
	if p.index < len(p.input) {
		return nil, fmt.Errorf("caracteres inesperados al final de la entrada en la posición %d", p.index)
	}

	return result, nil
}

// parseValue maneja diferentes tipos de valores en JSON: objetos, arrays, strings, números, etc.
func (p *Parser) parseValue() (interface{}, error) {
	// Ignorar espacios en blanco
	p.skipWhitespace()

	if p.index >= len(p.input) {
		return nil, fmt.Errorf("fin inesperado de la entrada JSON en la posición %d", p.index)
	}

	// Leer el primer carácter
	ch := rune(p.input[p.index])

	// Según el primer carácter, decidimos cómo procesar el valor.
	switch ch {
	case '{':
		p.index++
		return p.parseObject()
	case '[':
		p.index++
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
		return nil, fmt.Errorf("carácter inesperado: %c en la posición %d", ch, p.index)
	}
}

// skipWhitespace avanza el índice ignorando los espacios en blanco
func (p *Parser) skipWhitespace() {
	for p.index < len(p.input) && unicode.IsSpace(rune(p.input[p.index])) {
		p.index++
	}
}

// parseObject maneja los objetos JSON, que son de la forma { "clave": valor, ... }
func (p *Parser) parseObject() (map[string]interface{}, error) {
	obj := make(map[string]interface{})

	p.skipWhitespace()
	if p.index < len(p.input) && p.input[p.index] == '}' {
		p.index++
		return obj, nil
	}

	for {
		p.skipWhitespace()
		if p.index >= len(p.input) || p.input[p.index] != '"' {
			return nil, fmt.Errorf("se esperaba '\"' para la clave del objeto en la posición %d", p.index)
		}

		key, err := p.parseString()
		if err != nil {
			return nil, err
		}

		// Ignoramos los espacios en blanco y el siguiente carácter debe ser ":"
		p.skipWhitespace()
		if p.index >= len(p.input) || p.input[p.index] != ':' {
			return nil, fmt.Errorf("se esperaba ':' después de la clave en la posición %d", p.index)
		}
		p.index++ // Consume colon

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
			return nil, fmt.Errorf("objeto no terminado en la posición %d", p.index)
		}

		ch := p.input[p.index]
		if ch == '}' {
			p.index++
			break
		}
		if ch != ',' {
			return nil, fmt.Errorf("se esperaba ',' o '}' pero se obtuvo: %c en la posición %d", ch, p.index)
		}
		p.index++ // Consume comma
	}

	return obj, nil
}

// parseArray maneja los arrays JSON, que son de la forma [ valor, valor, ... ]
func (p *Parser) parseArray() ([]interface{}, error) {
	arr := []interface{}{}

	p.skipWhitespace()
	if p.index < len(p.input) && p.input[p.index] == ']' {
		p.index++
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
			return nil, fmt.Errorf("array no terminado en la posición %d", p.index)
		}
		ch := p.input[p.index]
		if ch == ']' {
			p.index++
			break
		}
		if ch != ',' {
			return nil, fmt.Errorf("se esperaba ',' o ']' pero se obtuvo: %c en la posición %d", ch, p.index)
		}
		p.index++ // Consume comma
	}

	return arr, nil
}

// parseString maneja las cadenas JSON (delimitadas por comillas)
func (p *Parser) parseString() (string, error) {
	if p.index >= len(p.input) || p.input[p.index] != '"' {
		return "", fmt.Errorf("se esperaba '\"' al inicio de la cadena en la posición %d", p.index)
	}
	p.index++ // consume opening quote

	var result string
	for {
		if p.index >= len(p.input) {
			return "", fmt.Errorf("cadena de texto no terminada en la posición %d", p.index)
		}
		ch := rune(p.input[p.index])
		p.index++

		// Si encontramos una comilla de cierre, hemos terminado de leer la cadena
		if ch == '"' {
			break
		}

		// Si encontramos un carácter de escape (\\), lo manejamos
		if ch == '\\' {
			if p.index >= len(p.input) {
				return "", fmt.Errorf("cadena de texto no terminada en la posición %d", p.index)
			}
			escapedChar := rune(p.input[p.index])
			p.index++

			switch escapedChar {
			case '"', '\\', '/', 'b', 'f', 'n', 'r', 't':
				result += string(escapedChar) // Agregar el carácter escapado
			case 'u':
				return "", fmt.Errorf("secuencias de escape unicode no implementadas en la posición %d", p.index-1)
			default:
				return "", fmt.Errorf("secuencia de escape inválida: \\%c en la posición %d", escapedChar, p.index-1)
			}
		} else {
			// Si no es un carácter escapado, agregamos el carácter a la cadena
			result += string(ch)
		}
	}

	return result, nil
}

// parseNumber maneja los números JSON (enteros o flotantes)
func (p *Parser) parseNumber() (interface{}, error) {
	start := p.index

	// Lee hasta que el carácter no sea un número o punto decimal
	for p.index < len(p.input) {
		ch := rune(p.input[p.index])
		if !unicode.IsDigit(ch) && ch != '.' && ch != '-' && ch != 'e' && ch != 'E' && ch != '+' {
			break
		}
		p.index++
	}

	// Subcadena del número encontrado
	numberStr := p.input[start:p.index]
	if numberStr == "" {
		return nil, fmt.Errorf("número inválido en la posición %d", start)
	}
	// Intentamos convertir el número a un flotante (en JSON siempre son float64)
	number, err := strconv.ParseFloat(numberStr, 64)
	if err != nil {
		return nil, fmt.Errorf("número inválido: %s en la posición %d", numberStr, start)
	}

	return number, nil
}

// parseBoolean maneja los valores booleanos JSON (true, false)
func (p *Parser) parseBoolean() (bool, error) {
	start := p.index

	// Lee los siguientes 4 caracteres para "true" o 5 caracteres para "false"
	if len(p.input) >= start+4 && p.input[start:start+4] == "true" {
		p.index += 4
		return true, nil
	}
	if len(p.input) >= start+5 && p.input[start:start+5] == "false" {
			p.index += 5
		return false, nil
	}

	return false, fmt.Errorf("valor booleano inválido en la posición %d", start)
}

// parseNull maneja el valor null JSON, que se convierte en nil en Go
func (p *Parser) parseNull() (interface{}, error) {
	start := p.index

	// Lee los 4 caracteres para "null"
	if len(p.input) >= start+4 && p.input[start:start+4] == "null" {
		p.index += 4
		return nil, nil
	}

	return nil, fmt.Errorf("valor nulo inválido en la posición %d", start)
}
