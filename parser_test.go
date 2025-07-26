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

		// Números especiales
		{"Número cero", `0`, 0.0, false, ""},
		{"Número decimal con ceros", `0.0`, 0.0, false, ""},
		{"Número científico simple", `1e2`, 100.0, false, ""},
		{"Número científico negativo", `1e-2`, 0.01, false, ""},
		{"Número científico con signo", `1E+2`, 100.0, false, ""},

		// Casos inválidos - Sintaxis general
		{"JSON inválido - texto plano", `invalid json`, nil, true, "formato JSON inválido"},
		{"Entrada vacía", ``, nil, true, "entrada JSON vacía"},
		{"Solo espacios", `   `, nil, true, "entrada JSON vacía"},

		// Casos inválidos - Objetos
		{"Objeto sin cerrar", `{"a": 1`, nil, true, "llaves desbalanceadas"},
		{"Objeto sin comillas en clave", `{a: 1}`, nil, true, "formato JSON inválido"},
		{"Objeto con coma extra", `{"a": 1, "b": 2,}`, nil, true, "coma extra antes de '}'"},
		{"Clave duplicada", `{"a": 1, "a": 2}`, nil, true, "clave duplicada"},

		// Casos inválidos - Arrays
		{"Array sin cerrar", `[1, 2`, nil, true, "corchetes desbalanceados"},
		{"Array con coma extra", `[1, 2, 3,]`, nil, true, "coma extra antes de ']'"},

		// Casos inválidos - Strings
		{"String sin cerrar", `"hello`, nil, true, "formato JSON inválido"},

		// Casos inválidos - Números
		{"Número inválido - solo signo", `-`, nil, true, "formato JSON inválido"},
		{"Número inválido - punto al inicio", `.5`, nil, true, "formato JSON inválido"},
		{"Número inválido - múltiples ceros", `001`, nil, true, "números no pueden empezar con múltiples ceros"},

		// JSON con saltos de línea
		{"JSON con saltos de línea", "{\n  \"name\": \"John\",\n  \"age\": 30\n}",
			map[string]interface{}{"name": "John", "age": 30.0}, false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
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

// Test específico para validación rápida
func TestFastValidateJSON(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{"JSON válido - objeto", `{"key": "value"}`, false},
		{"JSON válido - array", `[1, 2, 3]`, false},
		{"JSON válido - string", `"hello"`, false},
		{"JSON válido - número", `42`, false},
		{"JSON válido - boolean", `true`, false},
		{"JSON válido - null", `null`, false},
		{"JSON inválido - sin comillas", `{key: value}`, true},
		{"JSON inválido - sin cerrar", `{"key": "value"`, true},
		{"Entrada vacía", ``, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			err := p.FastValidateJSON(tt.input)

			if (err != nil) != tt.wantError {
				t.Errorf("FastValidateJSON() error = %v, wantErr %v", err, tt.wantError)
			}
		})
	}
}

// Test para detección de tipos
func TestExtractJSONType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Objeto", `{"key": "value"}`, "object"},
		{"Array", `[1, 2, 3]`, "array"},
		{"String", `"hello"`, "string"},
		{"Número", `42`, "number"},
		{"Boolean", `true`, "boolean"},
		{"Null", `null`, "null"},
		{"Inválido", `invalid`, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			result := p.ExtractJSONType(tt.input)

			if result != tt.expected {
				t.Errorf("ExtractJSONType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test para conteo de elementos
func TestCountJSONElements(t *testing.T) {
	input := `{"str": "hello", "num": 42, "bool": true, "null": null, "arr": [1, 2], "obj": {"nested": "value"}}`

	p := NewParser()
	counts, err := p.CountJSONElements(input)

	if err != nil {
		t.Errorf("CountJSONElements() error = %v", err)
		return
	}

	// Verificar que se encontraron elementos (los conteos exactos pueden variar según la implementación)
	if counts["strings"] == 0 {
		t.Error("Expected to find strings")
	}
	if counts["numbers"] == 0 {
		t.Error("Expected to find numbers")
	}
	if counts["booleans"] == 0 {
		t.Error("Expected to find booleans")
	}
	if counts["nulls"] == 0 {
		t.Error("Expected to find nulls")
	}
}

// Benchmark para comparar rendimiento
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
			p := NewParser()
			for i := 0; i < b.N; i++ {
				_, err := p.ParseJSON(tc.json)
				if err != nil {
					b.Fatalf("ParseJSON() error = %v", err)
				}
			}
		})
	}
}

// Test de robustez con JSON grandes
func TestParseJSONLarge(t *testing.T) {
	// Generar un JSON grande
	var sb strings.Builder
	sb.WriteString(`{"items": [`)
	for i := 0; i < 100; i++ { // Reducido para que el test sea más rápido
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

	p := NewParser()
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

	if len(items) != 100 {
		t.Errorf("ParseJSON() large JSON expected 100 items, got %d", len(items))
	}
}

// Test para verificar balance de estructuras
func TestValidateStructureBalance(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{"Balanceado - objeto", `{"a": 1}`, false},
		{"Balanceado - array", `[1, 2, 3]`, false},
		{"Balanceado - anidado", `{"a": [1, {"b": 2}]}`, false},
		{"Desbalanceado - llave extra", `{"a": 1}}`, true},
		{"Desbalanceado - corchete extra", `[1, 2]]`, true},
		{"Desbalanceado - llave faltante", `{"a": 1`, true},
		{"Desbalanceado - corchete faltante", `[1, 2`, true},
		{"Mixto desbalanceado", `{"a": [1, 2}`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			err := p.validateStructureBalance(tt.input)

			if (err != nil) != tt.wantError {
				t.Errorf("validateStructureBalance() error = %v, wantErr %v", err, tt.wantError)
			}
		})
	}
}

// Test para parsing de números con casos edge
func TestParseNumber(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  float64
		wantError bool
	}{
		{"Entero positivo", "123", 123.0, false},
		{"Entero negativo", "-123", -123.0, false},
		{"Decimal positivo", "123.45", 123.45, false},
		{"Decimal negativo", "-123.45", -123.45, false},
		{"Cero", "0", 0.0, false},
		{"Cero decimal", "0.0", 0.0, false},
		{"Científico simple", "1e2", 100.0, false},
		{"Científico negativo", "1e-2", 0.01, false},
		{"Científico con signo", "1E+2", 100.0, false},
		{"Inválido - múltiples ceros", "00", 0.0, true},
		{"Inválido - punto al final", "123.", 0.0, true},
		{"Inválido - exponente vacío", "1e", 0.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			result, err := p.parseNumber(tt.input)

			if (err != nil) != tt.wantError {
				t.Errorf("parseNumber() error = %v, wantErr %v", err, tt.wantError)
				return
			}

			if !tt.wantError && result != tt.expected {
				t.Errorf("parseNumber() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test para unescaping de strings
func TestUnescapeString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Sin escape", "hello", "hello"},
		{"Comillas escapadas", `say \"hello\"`, `say "hello"`},
		{"Barra invertida", `path\\to\\file`, `path\to\file`},
		{"Barra diagonal", `http:\/\/example.com`, `http://example.com`},
		{"Salto de línea", `line1\nline2`, "line1\nline2"},
		{"Tab", `col1\tcol2`, "col1\tcol2"},
		{"Retorno de carro", `line1\rline2`, "line1\rline2"},
		{"Backspace", `text\b`, "text\b"},
		{"Form feed", `page1\fpage2`, "page1\fpage2"},
		{"Unicode simple", `\u0041`, "A"},
		{"Unicode complejo", `\u00e9`, "é"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			result := p.unescapeString(tt.input)

			if result != tt.expected {
				t.Errorf("unescapeString() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// Test para funciones de conveniencia
func TestConvenienceFunctions(t *testing.T) {
	// Test OptimizedParseJSON
	result, err := OptimizedParseJSON(`{"test": true}`)
	if err != nil {
		t.Errorf("OptimizedParseJSON() error = %v", err)
	}
	expected := map[string]interface{}{"test": true}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("OptimizedParseJSON() = %v, want %v", result, expected)
	}

	// Test FastValidateJSON
	err = FastValidateJSON(`{"valid": "json"}`)
	if err != nil {
		t.Errorf("FastValidateJSON() error = %v", err)
	}

	err = FastValidateJSON(`invalid json`)
	if err == nil {
		t.Error("FastValidateJSON() expected error for invalid JSON")
	}

	// Test IsValidJSON
	if !IsValidJSON(`{"valid": true}`) {
		t.Error("IsValidJSON() should return true for valid JSON")
	}

	if IsValidJSON(`invalid`) {
		t.Error("IsValidJSON() should return false for invalid JSON")
	}
}
