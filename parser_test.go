package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestParseJSON(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expected      interface{}
		hasError      bool
		errorContains string
	}{
		// Casos válidos básicos
		{"Objeto válido simple", `{"name": "John", "age": 30}`,
			map[string]interface{}{"name": "John", "age": 30.0}, false, ""},
		{"Array válido", `["apple", "banana", "cherry"]`,
			[]interface{}{"apple", "banana", "cherry"}, false, ""},
		{"String válido", `"hello"`, "hello", false, ""},
		{"Número entero", `123`, 123.0, false, ""},
		{"Número decimal", `123.45`, 123.45, false, ""},
		{"Número negativo", `-42`, -42.0, false, ""},
		{"Número decimal negativo", `-123.45`, -123.45, false, ""},
		{"Booleano true", `true`, true, false, ""},
		{"Booleano false", `false`, false, false, ""},
		{"Null válido", `null`, nil, false, ""},

		// Casos válidos con espacios
		{"Objeto con espacios", `{ "name" : "John" , "age" : 30 }`,
			map[string]interface{}{"name": "John", "age": 30.0}, false, ""},
		{"Array con espacios", `[ "a" , "b" , "c" ]`,
			[]interface{}{"a", "b", "c"}, false, ""},

		// Casos válidos vacíos
		{"Objeto vacío", `{}`, map[string]interface{}{}, false, ""},
		{"Array vacío", `[]`, []interface{}{}, false, ""},

		// Casos válidos anidados
		{"Objeto anidado", `{"a": {"b": 1}}`,
			map[string]interface{}{"a": map[string]interface{}{"b": 1.0}}, false, ""},
		{"Array anidado", `[1, [2, 3], 4]`,
			[]interface{}{1.0, []interface{}{2.0, 3.0}, 4.0}, false, ""},
		{"Estructura compleja", `{"users": [{"name": "Ana", "active": true}, {"name": "Carlos", "active": false}], "count": 2}`,
			map[string]interface{}{
				"users": []interface{}{
					map[string]interface{}{"name": "Ana", "active": true},
					map[string]interface{}{"name": "Carlos", "active": false},
				},
				"count": 2.0,
			}, false, ""},

		// Casos con caracteres de escape
		{"String con escape", `"Hello\nWorld"`, "Hello\nWorld", false, ""},
		{"String con comillas escapadas", `"Say \"Hello\""`, "Say \"Hello\"", false, ""},
		{"String con barra invertida", `"C:\\Users"`, "C:\\Users", false, ""},

		// Casos válidos con múltiples tipos
		{"Tipos mixtos", `{"string": "text", "number": 42.5, "boolean": true, "null": null, "array": [1, 2]}`,
			map[string]interface{}{
				"string":  "text",
				"number":  42.5,
				"boolean": true,
				"null":    nil,
				"array":   []interface{}{1.0, 2.0},
			}, false, ""},

		// Casos inválidos - Sintaxis general
		{"JSON inválido - texto plano", `invalid json`, nil, true, "carácter inesperado"},
		{"Entrada vacía", ``, nil, true, "fin inesperado"},
		{"Solo espacios", `   `, nil, true, "fin inesperado"},

		// Casos inválidos - Objetos
		{"Objeto sin cerrar", `{"a": 1`, nil, true, "objeto no terminado"},
		{"Objeto sin dos puntos", `{"a" 1}`, nil, true, "se esperaba ':'"},
		{"Objeto sin comillas en clave", `{a: 1}`, nil, true, "se esperaba '\"'"},
		{"Objeto con coma extra", `{"a": 1,}`, nil, true, "coma extra antes de '}'"},
		{"Objeto sin coma", `{"a": 1 "b": 2}`, nil, true, "se esperaba ',' o '}'"},
		{"Clave duplicada", `{"a": 1, "a": 2}`, nil, true, "clave duplicada"},

		// Casos inválidos - Arrays
		{"Array sin cerrar", `[1, 2`, nil, true, "array no terminado"},
		{"Array con coma extra", `[1, 2,]`, nil, true, "coma extra antes de ']'"},
		{"Array sin coma", `[1 2]`, nil, true, "se esperaba ',' o ']'"},

		// Casos inválidos - Strings
		{"String sin cerrar", `"hello`, nil, true, "cadena de texto no terminada"},
		{"String con escape inválido", `"hello\x"`, nil, true, "secuencia de escape inválida"},
		{"String con escape incompleto", `"hello\`, nil, true, "cadena de texto no terminada"},

		// Casos inválidos - Números
		{"Número inválido - solo signo", `-`, nil, true, "número inválido"},
		{"Número inválido - punto al inicio", `.5`, nil, true, "carácter inesperado"},
		{"Número inválido - doble punto", `12.34.56`, nil, true, "caracteres inesperados al final"},

		// Casos inválidos - Booleanos y null
		{"Boolean inválido", `tru`, nil, true, "valor booleano inválido"},
		{"Null inválido", `nul`, nil, true, "valor nulo inválido"},

		// Casos inválidos - Caracteres extra
		{"Caracteres extra al final", `{"a": 1} extra`, nil, true, "caracteres inesperados al final"},
		{"Múltiples valores", `1 2`, nil, true, "caracteres inesperados al final"},

		// Casos con espacios en blanco variados
		{"JSON con saltos de línea", "{\n  \"name\": \"John\",\n  \"age\": 30\n}",
			map[string]interface{}{"name": "John", "age": 30.0}, false, ""},
		{"JSON con tabs", "\t{\t\"a\"\t:\t1\t}\t",
			map[string]interface{}{"a": 1.0}, false, ""},

		// Números especiales
		{"Número cero", `0`, 0.0, false, ""},
		{"Número decimal con ceros", `0.0`, 0.0, false, ""},
		{"Número científico simple", `1e2`, 100.0, false, ""},
		{"Número científico negativo", `1e-2`, 0.01, false, ""},
		{"Número científico con signo", `1E+2`, 100.0, false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{}
			result, err := p.ParseJSON(tt.input)

			// Verificar si se esperaba un error
			if (err != nil) != tt.hasError {
				t.Errorf("ParseJSON() error = %v, wantErr %v", err, tt.hasError)
				return
			}

			// Si se esperaba un error, verificar el mensaje
			if tt.hasError && err != nil {
				if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("ParseJSON() error = %q, want to contain %q", err.Error(), tt.errorContains)
				}
				return
			}

			// Si no se esperaba error, verificar el resultado
			if !tt.hasError && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ParseJSON() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test específico para verificar el manejo de líneas y columnas en errores
func TestParserErrorPositions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantLine int
		wantCol  int
	}{
		{"Error en línea 1", `{"a":}`, 1, 6},
		{"Error en línea 2", "{\n  \"a\":\n}", 3, 1},
		{"Error en línea 3", "{\n  \"a\": 1,\n  \"b\":\n}", 4, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{}
			_, err := p.ParseJSON(tt.input)

			if err == nil {
				t.Errorf("ParseJSON() expected error but got none")
				return
			}

			errorMsg := err.Error()
			expectedLine := strings.Contains(errorMsg, fmt.Sprintf("línea %d", tt.wantLine))

			if !expectedLine {
				t.Errorf("ParseJSON() error = %q, want to contain line %d", errorMsg, tt.wantLine)
			}
		})
	}
}

// Benchmark para medir el rendimiento del parser
func BenchmarkParseJSON(b *testing.B) {
	testCases := []struct {
		name string
		json string
	}{
		{"Simple object", `{"name": "John", "age": 30}`},
		{"Array", `[1, 2, 3, 4, 5]`},
		{"Complex nested", `{"users": [{"name": "Ana", "data": {"score": 95.5, "active": true}}, {"name": "Carlos", "data": {"score": 87.2, "active": false}}]}`},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				p := &Parser{}
				_, err := p.ParseJSON(tc.json)
				if err != nil {
					b.Fatalf("ParseJSON() error = %v", err)
				}
			}
		})
	}
}

// Test de robustez con JSON muy grandes
func TestParseJSONLarge(t *testing.T) {
	// Generar un JSON grande
	var sb strings.Builder
	sb.WriteString(`{"items": [`)
	for i := 0; i < 1000; i++ {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(`{"id": `)
		sb.WriteString(fmt.Sprintf("%d", i))
		sb.WriteString(`, "name": "item`)
		sb.WriteString(fmt.Sprintf("%d", i))
		sb.WriteString(`"}`)
	}
	sb.WriteString(`]}`)

	p := &Parser{}
	result, err := p.ParseJSON(sb.String())

	if err != nil {
		t.Errorf("ParseJSON() large JSON error = %v", err)
		return
	}

	// Verificar que el resultado es el esperado
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Errorf("ParseJSON() large JSON result is not a map")
		return
	}

	items, ok := resultMap["items"].([]interface{})
	if !ok {
		t.Errorf("ParseJSON() large JSON items is not an array")
		return
	}

	if len(items) != 1000 {
		t.Errorf("ParseJSON() large JSON expected 1000 items, got %d", len(items))
	}
}
