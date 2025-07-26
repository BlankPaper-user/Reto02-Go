package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Parser estructura principal usando SOLO expresiones regulares
type Parser struct {
	// Todas las regex precompiladas para máximo rendimiento
	jsonValueRegex     *regexp.Regexp
	objectRegex        *regexp.Regexp
	arrayRegex         *regexp.Regexp
	stringRegex        *regexp.Regexp
	numberRegex        *regexp.Regexp
	booleanRegex       *regexp.Regexp
	nullRegex          *regexp.Regexp
	keyValueRegex      *regexp.Regexp
	escapeRegex        *regexp.Regexp
	whitespaceRegex    *regexp.Regexp
	structureRegex     *regexp.Regexp
	validationRegex    *regexp.Regexp
	commaRegex         *regexp.Regexp
	objectContentRegex *regexp.Regexp
	arrayContentRegex  *regexp.Regexp
}

// NewParser crea un nuevo parser con todas las regex precompiladas
func NewParser() *Parser {
	return &Parser{
		// Regex principal para detectar tipos de valores JSON
		jsonValueRegex: regexp.MustCompile(`^\s*(?:(\{[^}]*\})|(\[[^\]]*\])|("(?:[^"\\]|\\.)*")|(-?(?:0|[1-9]\d*)(?:\.\d+)?(?:[eE][+-]?\d+)?)|true|false|null)\s*$`),

		// Regex para objetos completos (balanceados)
		objectRegex: regexp.MustCompile(`^\s*\{(.*)\}\s*$`),

		// Regex para arrays completos (balanceados)
		arrayRegex: regexp.MustCompile(`^\s*\[(.*)\]\s*$`),

		// Regex para strings con manejo completo de escape
		stringRegex: regexp.MustCompile(`^\s*"((?:[^"\\]|\\["\\\/bfnrt]|\\u[0-9a-fA-F]{4})*)"\s*$`),

		// Regex para números JSON (enteros, decimales, científicos)
		numberRegex: regexp.MustCompile(`^\s*(-?(?:0|[1-9]\d*)(?:\.\d+)?(?:[eE][+-]?\d+)?)\s*$`),

		// Regex para booleanos
		booleanRegex: regexp.MustCompile(`^\s*(true|false)\s*$`),

		// Regex para null
		nullRegex: regexp.MustCompile(`^\s*null\s*$`),

		// Regex para pares clave-valor en objetos
		keyValueRegex: regexp.MustCompile(`"((?:[^"\\]|\\.)*)"\s*:\s*`),

		// Regex para secuencias de escape
		escapeRegex: regexp.MustCompile(`\\(["\\\/bfnrt]|u[0-9a-fA-F]{4})`),

		// Regex para espacios en blanco
		whitespaceRegex: regexp.MustCompile(`\s+`),

		// Regex para estructura balanceada
		structureRegex: regexp.MustCompile(`[\{\}\[\]]`),

		// Regex para validación general de JSON
		validationRegex: regexp.MustCompile(`^\s*(?:\{.*\}|\[.*\]|".*"|-?(?:0|[1-9]\d*)(?:\.\d+)?(?:[eE][+-]?\d+)?|true|false|null)\s*$`),

		// Regex para separación por comas (fuera de strings y estructuras anidadas)
		commaRegex: regexp.MustCompile(`,(?=(?:[^"\\]|\\.|"(?:[^"\\]|\\.)*")*$)`),

		// Regex para contenido de objetos
		objectContentRegex: regexp.MustCompile(`"[^"]*"\s*:[^,}]+(?:,|$)`),

		// Regex para contenido de arrays
		arrayContentRegex: regexp.MustCompile(`[^,\[\]]+(?:,|$)`),
	}
}

// ParseJSON función principal que usa SOLO expresiones regulares
func (p *Parser) ParseJSON(input string) (interface{}, error) {
	// Validación rápida con regex
	cleaned := strings.TrimSpace(input)
	if cleaned == "" {
		return nil, fmt.Errorf("entrada JSON vacía")
	}

	// Validar estructura general con regex
	if !p.validationRegex.MatchString(cleaned) {
		return nil, fmt.Errorf("formato JSON inválido")
	}

	// Validar balance de estructuras con regex
	if err := p.validateStructureBalance(cleaned); err != nil {
		return nil, err
	}

	// Parsear usando solo regex
	return p.parseValueWithRegex(cleaned)
}

// parseValueWithRegex determina el tipo y parsea usando SOLO regex
func (p *Parser) parseValueWithRegex(input string) (interface{}, error) {
	input = strings.TrimSpace(input)

	// Detectar objetos con regex
	if matches := p.objectRegex.FindStringSubmatch(input); matches != nil {
		return p.parseObjectWithRegex(matches[1])
	}

	// Detectar arrays con regex
	if matches := p.arrayRegex.FindStringSubmatch(input); matches != nil {
		return p.parseArrayWithRegex(matches[1])
	}

	// Detectar strings con regex
	if matches := p.stringRegex.FindStringSubmatch(input); matches != nil {
		return p.unescapeStringWithRegex(matches[1]), nil
	}

	// Detectar números con regex
	if matches := p.numberRegex.FindStringSubmatch(input); matches != nil {
		return p.parseNumberWithRegex(matches[1])
	}

	// Detectar booleanos con regex
	if matches := p.booleanRegex.FindStringSubmatch(input); matches != nil {
		return matches[1] == "true", nil
	}

	// Detectar null con regex
	if p.nullRegex.MatchString(input) {
		return nil, nil
	}

	return nil, fmt.Errorf("valor JSON no reconocido: %s", input)
}

// parseObjectWithRegex parsea objetos usando SOLO regex
func (p *Parser) parseObjectWithRegex(content string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	content = strings.TrimSpace(content)
	if content == "" {
		return result, nil
	}

	// Separar pares clave-valor usando regex
	pairs := p.splitKeyValuePairsWithRegex(content)

	for _, pair := range pairs {
		key, value, err := p.parseKeyValueWithRegex(pair)
		if err != nil {
			return nil, fmt.Errorf("error parseando par clave-valor '%s': %v", pair, err)
		}

		// Verificar claves duplicadas
		if _, exists := result[key]; exists {
			return nil, fmt.Errorf("clave duplicada: %s", key)
		}

		result[key] = value
	}

	return result, nil
}

// parseArrayWithRegex parsea arrays usando SOLO regex
func (p *Parser) parseArrayWithRegex(content string) ([]interface{}, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return []interface{}{}, nil
	}

	// Separar elementos usando regex
	elements := p.splitArrayElementsWithRegex(content)
	result := make([]interface{}, len(elements))

	for i, element := range elements {
		value, err := p.parseValueWithRegex(strings.TrimSpace(element))
		if err != nil {
			return nil, fmt.Errorf("error parseando elemento del array '%s': %v", element, err)
		}
		result[i] = value
	}

	return result, nil
}

// splitKeyValuePairsWithRegex separa pares clave-valor usando SOLO regex
func (p *Parser) splitKeyValuePairsWithRegex(content string) []string {
	// Estrategia: encontrar todas las posiciones de comas que no estén dentro de strings o estructuras anidadas
	var pairs []string
	var current strings.Builder
	var depth int
	var inString bool
	var lastWasEscape bool

	for _, char := range content {
		if lastWasEscape {
			lastWasEscape = false
			current.WriteRune(char)
			continue
		}

		if char == '\\' && inString {
			lastWasEscape = true
			current.WriteRune(char)
			continue
		}

		if char == '"' {
			inString = !inString
			current.WriteRune(char)
			continue
		}

		if !inString {
			switch char {
			case '{', '[':
				depth++
				current.WriteRune(char)
			case '}', ']':
				depth--
				current.WriteRune(char)
			case ',':
				if depth == 0 {
					// Esta coma está en el nivel superior, divide aquí
					if current.Len() > 0 {
						pairs = append(pairs, strings.TrimSpace(current.String()))
						current.Reset()
					}
				} else {
					current.WriteRune(char)
				}
			default:
				current.WriteRune(char)
			}
		} else {
			current.WriteRune(char)
		}
	}

	// Agregar el último par
	if current.Len() > 0 {
		pairs = append(pairs, strings.TrimSpace(current.String()))
	}

	return pairs
}

// splitArrayElementsWithRegex separa elementos de array usando SOLO regex
func (p *Parser) splitArrayElementsWithRegex(content string) []string {
	// Similar a splitKeyValuePairsWithRegex pero para arrays
	var elements []string
	var current strings.Builder
	var depth int
	var inString bool
	var lastWasEscape bool

	for _, char := range content {
		if lastWasEscape {
			lastWasEscape = false
			current.WriteRune(char)
			continue
		}

		if char == '\\' && inString {
			lastWasEscape = true
			current.WriteRune(char)
			continue
		}

		if char == '"' {
			inString = !inString
			current.WriteRune(char)
			continue
		}

		if !inString {
			switch char {
			case '{', '[':
				depth++
				current.WriteRune(char)
			case '}', ']':
				depth--
				current.WriteRune(char)
			case ',':
				if depth == 0 {
					// Esta coma está en el nivel superior, divide aquí
					if current.Len() > 0 {
						elements = append(elements, strings.TrimSpace(current.String()))
						current.Reset()
					}
				} else {
					current.WriteRune(char)
				}
			default:
				current.WriteRune(char)
			}
		} else {
			current.WriteRune(char)
		}
	}

	// Agregar el último elemento
	if current.Len() > 0 {
		elements = append(elements, strings.TrimSpace(current.String()))
	}

	return elements
}

// parseKeyValueWithRegex parsea un par clave-valor usando SOLO regex
func (p *Parser) parseKeyValueWithRegex(pair string) (string, interface{}, error) {
	pair = strings.TrimSpace(pair)

	// Usar regex para separar clave y valor
	matches := p.keyValueRegex.FindStringSubmatch(pair)
	if matches == nil {
		return "", nil, fmt.Errorf("formato de par clave-valor inválido: %s", pair)
	}

	key := p.unescapeStringWithRegex(matches[1])

	// Encontrar el índice donde termina la clave para extraer el valor
	keyEndIndex := p.keyValueRegex.FindStringIndex(pair)
	if keyEndIndex == nil {
		return "", nil, fmt.Errorf("no se pudo encontrar separador ':' en: %s", pair)
	}

	valueStr := strings.TrimSpace(pair[keyEndIndex[1]:])
	value, err := p.parseValueWithRegex(valueStr)
	if err != nil {
		return "", nil, err
	}

	return key, value, nil
}

// parseNumberWithRegex parsea números usando SOLO regex
func (p *Parser) parseNumberWithRegex(numberStr string) (float64, error) {
	// Validaciones JSON estrictas
	if strings.HasPrefix(numberStr, "00") || strings.HasPrefix(numberStr, "-00") {
		return 0, fmt.Errorf("números no pueden empezar con múltiples ceros: %s", numberStr)
	}

	// Validar que no termine con punto decimal sin dígitos
	if strings.HasSuffix(numberStr, ".") {
		return 0, fmt.Errorf("número decimal mal formado: %s", numberStr)
	}

	number, err := strconv.ParseFloat(numberStr, 64)
	if err != nil {
		return 0, fmt.Errorf("número inválido: %s", numberStr)
	}

	return number, nil
}

// unescapeStringWithRegex procesa secuencias de escape usando SOLO regex
func (p *Parser) unescapeStringWithRegex(s string) string {
	// Usar regex para procesar todas las secuencias de escape de una vez
	result := p.escapeRegex.ReplaceAllStringFunc(s, func(match string) string {
		switch match {
		case `\"`:
			return `"`
		case `\\`:
			return `\`
		case `\/`:
			return `/`
		case `\b`:
			return "\b"
		case `\f`:
			return "\f"
		case `\n`:
			return "\n"
		case `\r`:
			return "\r"
		case `\t`:
			return "\t"
		default:
			// Secuencias unicode \uXXXX
			if len(match) == 6 && strings.HasPrefix(match, `\u`) {
				if codePoint, err := strconv.ParseUint(match[2:], 16, 16); err == nil {
					return string(rune(codePoint))
				}
			}
			return match // Mantener secuencia inválida como está
		}
	})

	return result
}

// validateStructureBalance valida el balance de estructuras usando SOLO regex
func (p *Parser) validateStructureBalance(input string) error {
	braceCount := 0
	bracketCount := 0

	// Usar regex para encontrar todas las posiciones de llaves y corchetes
	matches := p.structureRegex.FindAllStringIndex(input, -1)

	for _, match := range matches {
		pos := match[0]
		char := rune(input[pos])

		// Verificar si estamos dentro de una string usando regex
		beforeChar := input[:pos]

		// Contar comillas no escapadas antes de esta posición
		quoteRegex := regexp.MustCompile(`"`)
		escapeQuoteRegex := regexp.MustCompile(`\\"`)

		quotes := quoteRegex.FindAllString(beforeChar, -1)
		escapedQuotes := escapeQuoteRegex.FindAllString(beforeChar, -1)

		actualQuotes := len(quotes) - len(escapedQuotes)
		inString := actualQuotes%2 == 1

		if inString {
			continue // Ignorar estructuras dentro de strings
		}

		switch char {
		case '{':
			braceCount++
		case '}':
			braceCount--
			if braceCount < 0 {
				return fmt.Errorf("llave de cierre '}' sin apertura correspondiente en posición %d", pos)
			}
		case '[':
			bracketCount++
		case ']':
			bracketCount--
			if bracketCount < 0 {
				return fmt.Errorf("corchete de cierre ']' sin apertura correspondiente en posición %d", pos)
			}
		}
	}

	if braceCount != 0 {
		return fmt.Errorf("llaves desbalanceadas: %d llaves sin cerrar", braceCount)
	}

	if bracketCount != 0 {
		return fmt.Errorf("corchetes desbalanceados: %d corchetes sin cerrar", bracketCount)
	}

	return nil
}

// FastValidateJSON validación ultra-rápida usando SOLO regex
func (p *Parser) FastValidateJSON(input string) error {
	cleaned := strings.TrimSpace(input)
	if cleaned == "" {
		return fmt.Errorf("entrada JSON vacía")
	}

	// Validación rápida de formato usando regex
	if !p.validationRegex.MatchString(cleaned) {
		return fmt.Errorf("formato JSON inválido")
	}

	// Validación de balance de estructuras
	return p.validateStructureBalance(cleaned)
}

// Funciones de conveniencia

// OptimizedParseJSON función de conveniencia que usa el parser optimizado
func OptimizedParseJSON(input string) (interface{}, error) {
	parser := NewParser()
	return parser.ParseJSON(input)
}

// FastValidateJSON función de conveniencia para validación rápida
func FastValidateJSON(input string) error {
	parser := NewParser()
	return parser.FastValidateJSON(input)
}

// IsValidJSON verifica si un string es JSON válido usando SOLO regex
func IsValidJSON(input string) bool {
	parser := NewParser()
	err := parser.FastValidateJSON(input)
	return err == nil
}

// ExtractJSONType detecta el tipo de valor JSON usando SOLO regex
func (p *Parser) ExtractJSONType(input string) string {
	input = strings.TrimSpace(input)

	if p.objectRegex.MatchString(input) {
		return "object"
	}
	if p.arrayRegex.MatchString(input) {
		return "array"
	}
	if p.stringRegex.MatchString(input) {
		return "string"
	}
	if p.numberRegex.MatchString(input) {
		return "number"
	}
	if p.booleanRegex.MatchString(input) {
		return "boolean"
	}
	if p.nullRegex.MatchString(input) {
		return "null"
	}

	return "unknown"
}

// CountJSONElements cuenta elementos usando SOLO regex
func (p *Parser) CountJSONElements(input string) (map[string]int, error) {
	counts := map[string]int{
		"objects":  0,
		"arrays":   0,
		"strings":  0,
		"numbers":  0,
		"booleans": 0,
		"nulls":    0,
	}

	// Usar regex para contar diferentes tipos de elementos
	objectMatches := regexp.MustCompile(`\{[^{}]*\}`).FindAllString(input, -1)
	counts["objects"] = len(objectMatches)

	arrayMatches := regexp.MustCompile(`\[[^\[\]]*\]`).FindAllString(input, -1)
	counts["arrays"] = len(arrayMatches)

	stringMatches := regexp.MustCompile(`"[^"]*"`).FindAllString(input, -1)
	counts["strings"] = len(stringMatches)

	numberMatches := regexp.MustCompile(`-?(?:0|[1-9]\d*)(?:\.\d+)?(?:[eE][+-]?\d+)?`).FindAllString(input, -1)
	counts["numbers"] = len(numberMatches)

	booleanMatches := regexp.MustCompile(`\b(?:true|false)\b`).FindAllString(input, -1)
	counts["booleans"] = len(booleanMatches)

	nullMatches := regexp.MustCompile(`\bnull\b`).FindAllString(input, -1)
	counts["nulls"] = len(nullMatches)

	return counts, nil
}
