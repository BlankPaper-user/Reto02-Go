package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Parser estructura principal usando expresiones regulares optimizadas
type Parser struct {
	// Regex precompiladas para máximo rendimiento
	objectRegex     *regexp.Regexp
	arrayRegex      *regexp.Regexp
	stringRegex     *regexp.Regexp
	numberRegex     *regexp.Regexp
	booleanRegex    *regexp.Regexp
	nullRegex       *regexp.Regexp
	keyValueRegex   *regexp.Regexp
	escapeRegex     *regexp.Regexp
	whitespaceRegex *regexp.Regexp
	validationRegex *regexp.Regexp
	structureRegex  *regexp.Regexp
}

// NewParser crea un nuevo parser con todas las regex precompiladas
func NewParser() *Parser {
	return &Parser{
		// Regex para objetos completos
		objectRegex: regexp.MustCompile(`^\s*\{(.*)\}\s*$`),

		// Regex para arrays completos
		arrayRegex: regexp.MustCompile(`^\s*\[(.*)\]\s*$`),

		// Regex para strings con escape completo
		stringRegex: regexp.MustCompile(`^\s*"((?:[^"\\]|\\["\\\/bfnrt]|\\u[0-9a-fA-F]{4})*)"\s*$`),

		// Regex para números JSON válidos
		numberRegex: regexp.MustCompile(`^\s*(-?(?:0|[1-9]\d*)(?:\.\d+)?(?:[eE][+-]?\d+)?)\s*$`),

		// Regex para booleanos
		booleanRegex: regexp.MustCompile(`^\s*(true|false)\s*$`),

		// Regex para null
		nullRegex: regexp.MustCompile(`^\s*null\s*$`),

		// Regex para pares clave-valor
		keyValueRegex: regexp.MustCompile(`^\s*"((?:[^"\\]|\\.)*)"\s*:\s*`),

		// Regex para secuencias de escape
		escapeRegex: regexp.MustCompile(`\\(["\\\/bfnrt]|u[0-9a-fA-F]{4})`),

		// Regex para espacios en blanco
		whitespaceRegex: regexp.MustCompile(`\s+`),

		// Regex para validación general
		validationRegex: regexp.MustCompile(`^\s*(?:\{.*\}|\[.*\]|".*"|-?(?:0|[1-9]\d*)(?:\.\d+)?(?:[eE][+-]?\d+)?|true|false|null)\s*$`),

		// Regex para estructuras balanceadas
		structureRegex: regexp.MustCompile(`[\{\}\[\]]`),
	}
}

// ParseJSON función principal de parsing
func (p *Parser) ParseJSON(input string) (interface{}, error) {
	cleaned := strings.TrimSpace(input)
	if cleaned == "" {
		return nil, fmt.Errorf("entrada JSON vacía")
	}

	// Validar estructura general
	if !p.validationRegex.MatchString(cleaned) {
		return nil, fmt.Errorf("formato JSON inválido")
	}

	// Validar balance de estructuras
	if err := p.validateStructureBalance(cleaned); err != nil {
		return nil, err
	}

	// Parsear el valor
	return p.parseValue(cleaned)
}

// parseValue determina el tipo y parsea el valor
func (p *Parser) parseValue(input string) (interface{}, error) {
	input = strings.TrimSpace(input)

	// Detectar objetos
	if matches := p.objectRegex.FindStringSubmatch(input); matches != nil {
		return p.parseObject(matches[1])
	}

	// Detectar arrays
	if matches := p.arrayRegex.FindStringSubmatch(input); matches != nil {
		return p.parseArray(matches[1])
	}

	// Detectar strings
	if matches := p.stringRegex.FindStringSubmatch(input); matches != nil {
		return p.unescapeString(matches[1]), nil
	}

	// Detectar números
	if matches := p.numberRegex.FindStringSubmatch(input); matches != nil {
		return p.parseNumber(matches[1])
	}

	// Detectar booleanos
	if matches := p.booleanRegex.FindStringSubmatch(input); matches != nil {
		return matches[1] == "true", nil
	}

	// Detectar null
	if p.nullRegex.MatchString(input) {
		return nil, nil
	}

	return nil, fmt.Errorf("valor JSON no reconocido: %s", input)
}

// parseObject parsea objetos JSON
func (p *Parser) parseObject(content string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	content = strings.TrimSpace(content)
	if content == "" {
		return result, nil
	}

	// VALIDAR COMAS EXTRA ANTES DE PARSEAR
	if err := p.validateNoTrailingCommas(content, '}'); err != nil {
		return nil, err
	}

	// Separar pares clave-valor respetando estructuras anidadas
	pairs, err := p.splitKeyValuePairs(content)
	if err != nil {
		return nil, err
	}

	for _, pair := range pairs {
		key, value, err := p.parseKeyValue(pair)
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

// parseArray parsea arrays JSON
func (p *Parser) parseArray(content string) ([]interface{}, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return []interface{}{}, nil
	}

	// VALIDAR COMAS EXTRA ANTES DE PARSEAR
	if err := p.validateNoTrailingCommas(content, ']'); err != nil {
		return nil, err
	}

	// Separar elementos respetando estructuras anidadas
	elements, err := p.splitArrayElements(content)
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, len(elements))
	for i, element := range elements {
		value, err := p.parseValue(strings.TrimSpace(element))
		if err != nil {
			return nil, fmt.Errorf("error parseando elemento del array '%s': %v", element, err)
		}
		result[i] = value
	}

	return result, nil
}

// splitKeyValuePairs separa pares clave-valor respetando anidación
func (p *Parser) splitKeyValuePairs(content string) ([]string, error) {
	var pairs []string
	var current strings.Builder
	var depth int
	var inString bool
	var lastWasEscape bool

	for i, char := range content {
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
					pairStr := strings.TrimSpace(current.String())
					if pairStr == "" {
						return nil, fmt.Errorf("par clave-valor vacío en posición %d", i)
					}
					pairs = append(pairs, pairStr)
					current.Reset()
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
		pairStr := strings.TrimSpace(current.String())
		if pairStr != "" {
			pairs = append(pairs, pairStr)
		}
	}

	return pairs, nil
}

// splitArrayElements separa elementos de array respetando anidación
func (p *Parser) splitArrayElements(content string) ([]string, error) {
	var elements []string
	var current strings.Builder
	var depth int
	var inString bool
	var lastWasEscape bool

	for i, char := range content {
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
					elemStr := strings.TrimSpace(current.String())
					if elemStr == "" {
						return nil, fmt.Errorf("elemento vacío en posición %d", i)
					}
					elements = append(elements, elemStr)
					current.Reset()
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
		elemStr := strings.TrimSpace(current.String())
		if elemStr != "" {
			elements = append(elements, elemStr)
		}
	}

	return elements, nil
}

// parseKeyValue parsea un par clave-valor
func (p *Parser) parseKeyValue(pair string) (string, interface{}, error) {
	pair = strings.TrimSpace(pair)

	// Encontrar el separador ':'
	matches := p.keyValueRegex.FindStringSubmatch(pair)
	if matches == nil {
		return "", nil, fmt.Errorf("formato de par clave-valor inválido: %s", pair)
	}

	key := p.unescapeString(matches[1])

	// Extraer el valor después del ':'
	keyEndIndex := p.keyValueRegex.FindStringIndex(pair)
	if keyEndIndex == nil {
		return "", nil, fmt.Errorf("no se pudo encontrar separador ':' en: %s", pair)
	}

	valueStr := strings.TrimSpace(pair[keyEndIndex[1]:])
	if valueStr == "" {
		return "", nil, fmt.Errorf("valor faltante para la clave '%s'", key)
	}

	value, err := p.parseValue(valueStr)
	if err != nil {
		return "", nil, err
	}

	return key, value, nil
}

// parseNumber parsea números con validaciones JSON estrictas
func (p *Parser) parseNumber(numberStr string) (float64, error) {
	// Validaciones JSON estrictas
	if strings.HasPrefix(numberStr, "00") || strings.HasPrefix(numberStr, "-00") {
		return 0, fmt.Errorf("números no pueden empezar con múltiples ceros: %s", numberStr)
	}

	if strings.HasSuffix(numberStr, ".") {
		return 0, fmt.Errorf("número decimal mal formado: %s", numberStr)
	}

	// Verificar notación científica válida
	if strings.Contains(numberStr, "e") || strings.Contains(numberStr, "E") {
		parts := regexp.MustCompile(`[eE]`).Split(numberStr, 2)
		if len(parts) == 2 && (parts[1] == "" || parts[1] == "+" || parts[1] == "-") {
			return 0, fmt.Errorf("exponente inválido en notación científica: %s", numberStr)
		}
	}

	number, err := strconv.ParseFloat(numberStr, 64)
	if err != nil {
		return 0, fmt.Errorf("número inválido: %s", numberStr)
	}

	return number, nil
}

// unescapeString procesa secuencias de escape
func (p *Parser) unescapeString(s string) string {
	return p.escapeRegex.ReplaceAllStringFunc(s, func(match string) string {
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
			return match
		}
	})
}

// validateStructureBalance valida el balance de estructuras
func (p *Parser) validateStructureBalance(input string) error {
	braceCount := 0
	bracketCount := 0
	inString := false
	lastWasEscape := false

	for i, char := range input {
		if lastWasEscape {
			lastWasEscape = false
			continue
		}

		if char == '\\' && inString {
			lastWasEscape = true
			continue
		}

		if char == '"' {
			inString = !inString
			continue
		}

		if inString {
			continue
		}

		switch char {
		case '{':
			braceCount++
		case '}':
			braceCount--
			if braceCount < 0 {
				return fmt.Errorf("llave de cierre '}' sin apertura correspondiente en posición %d", i)
			}
		case '[':
			bracketCount++
		case ']':
			bracketCount--
			if bracketCount < 0 {
				return fmt.Errorf("corchete de cierre ']' sin apertura correspondiente en posición %d", i)
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

// FastValidateJSON validación rápida usando regex
func (p *Parser) FastValidateJSON(input string) error {
	cleaned := strings.TrimSpace(input)
	if cleaned == "" {
		return fmt.Errorf("entrada JSON vacía")
	}

	// Validación rápida de formato
	if !p.validationRegex.MatchString(cleaned) {
		return fmt.Errorf("formato JSON inválido")
	}

	// Validación de balance de estructuras
	if err := p.validateStructureBalance(cleaned); err != nil {
		return err
	}

	// Validación específica de comas extra con regex
	if err := p.validateNoTrailingCommasRegex(cleaned); err != nil {
		return err
	}

	return nil
}

// ExtractJSONType detecta el tipo de valor JSON
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

// CountJSONElements cuenta elementos usando regex
func (p *Parser) CountJSONElements(input string) (map[string]int, error) {
	counts := map[string]int{
		"objects":  0,
		"arrays":   0,
		"strings":  0,
		"numbers":  0,
		"booleans": 0,
		"nulls":    0,
	}

	// Contar objetos (simplificado para estructura básica)
	objectMatches := regexp.MustCompile(`\{[^{}]*\}`).FindAllString(input, -1)
	counts["objects"] = len(objectMatches)

	// Contar arrays (simplificado para estructura básica)
	arrayMatches := regexp.MustCompile(`\[[^\[\]]*\]`).FindAllString(input, -1)
	counts["arrays"] = len(arrayMatches)

	// Contar strings
	stringMatches := regexp.MustCompile(`"[^"]*"`).FindAllString(input, -1)
	counts["strings"] = len(stringMatches)

	// Contar números
	numberMatches := regexp.MustCompile(`\b-?(?:0|[1-9]\d*)(?:\.\d+)?(?:[eE][+-]?\d+)?\b`).FindAllString(input, -1)
	counts["numbers"] = len(numberMatches)

	// Contar booleanos
	booleanMatches := regexp.MustCompile(`\b(?:true|false)\b`).FindAllString(input, -1)
	counts["booleans"] = len(booleanMatches)

	// Contar nulls
	nullMatches := regexp.MustCompile(`\bnull\b`).FindAllString(input, -1)
	counts["nulls"] = len(nullMatches)

	return counts, nil
}

// Funciones de conveniencia

// OptimizedParseJSON función de conveniencia
func OptimizedParseJSON(input string) (interface{}, error) {
	parser := NewParser()
	return parser.ParseJSON(input)
}

// FastValidateJSON función de conveniencia para validación
func FastValidateJSON(input string) error {
	parser := NewParser()
	return parser.FastValidateJSON(input)
}

// IsValidJSON verifica si un string es JSON válido
func IsValidJSON(input string) bool {
	parser := NewParser()
	err := parser.FastValidateJSON(input)
	return err == nil
}

// validateNoTrailingCommasRegex usa regex para detectar comas extra
func (p *Parser) validateNoTrailingCommasRegex(input string) error {
	// Regex para detectar comas seguidas solo de espacios y luego } o ]
	trailingCommaRegex := regexp.MustCompile(`,\s*[}\]]`)

	// Necesitamos verificar si la coma está fuera de strings
	matches := trailingCommaRegex.FindAllStringIndex(input, -1)

	for _, match := range matches {
		pos := match[0] // Posición de la coma

		// Verificar si esta coma está dentro de un string
		beforeComma := input[:pos]

		// Contar comillas no escapadas antes de la coma
		quoteCount := 0
		escaped := false

		for _, char := range beforeComma {
			if escaped {
				escaped = false
				continue
			}

			if char == '\\' {
				escaped = true
				continue
			}

			if char == '"' {
				quoteCount++
			}
		}

		// Si hay un número par de comillas, estamos fuera de un string
		// Si hay un número impar, estamos dentro de un string
		if quoteCount%2 == 0 {
			// Estamos fuera de string, esta coma es inválida
			closingChar := input[match[1]-1] // El carácter después de los espacios
			return fmt.Errorf("coma extra antes de '%c' en posición %d", closingChar, pos)
		}
	}

	return nil
}

// validateNoTrailingCommas valida que no haya comas extra antes de cerrar
func (p *Parser) validateNoTrailingCommas(content string, closingChar rune) error {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil
	}

	// Buscar comas que no estén seguidas de contenido válido
	var inString bool
	var lastWasEscape bool
	var depth int
	var lastNonWhitespace rune
	var lastNonWhitespacePos int

	for i, char := range content {
		if lastWasEscape {
			lastWasEscape = false
			continue
		}

		if char == '\\' && inString {
			lastWasEscape = true
			continue
		}

		if char == '"' {
			inString = !inString
			if !inString {
				lastNonWhitespace = char
				lastNonWhitespacePos = i
			}
			continue
		}

		if inString {
			continue
		}

		switch char {
		case '{', '[':
			depth++
			lastNonWhitespace = char
			lastNonWhitespacePos = i
		case '}', ']':
			depth--
			lastNonWhitespace = char
			lastNonWhitespacePos = i
		case ',':
			if depth == 0 {
				lastNonWhitespace = char
				lastNonWhitespacePos = i
			}
		case ' ', '\t', '\n', '\r':
			// Ignorar espacios en blanco
		default:
			lastNonWhitespace = char
			lastNonWhitespacePos = i
		}
	}

	// Si el último carácter no-espacio es una coma, es inválido
	if lastNonWhitespace == ',' {
		return fmt.Errorf("coma extra antes de '%c' en posición %d", closingChar, lastNonWhitespacePos)
	}

	return nil
}
