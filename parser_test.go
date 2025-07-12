package main

import (
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
		errorContains string // Nuevo campo para verificar el mensaje de error
	}{
		{"Objeto válido", `{"name": "John", "age": 30, "active": true, "interests": ["football", "music"]}`,
			map[string]interface{}{"name": "John", "age": 30.0, "active": true, "interests": []interface{}{"football", "music"}}, false, ""},
		{"Array válido", `["apple", "banana", "cherry"]`, []interface{}{"apple", "banana", "cherry"}, false, ""},
		{"String válido", `"hello"`, "hello", false, ""},
		{"Número válido", `123`, 123.0, false, ""},
		{"Booleano true válido", `true`, true, false, ""},
		{"Booleano false válido", `false`, false, false, ""},
		{"Null válido", `null`, nil, false, ""},
		{"Objeto vacío", `{}`, map[string]interface{}{}, false, ""},
		{"Array vacío", `[]`, []interface{}{}, false, ""},
		{"Objeto anidado", `{"a": {"b": 1}}`, map[string]interface{}{"a": map[string]interface{}{"b": 1.0}}, false, ""},
		{"JSON inválid	o", `invalid json`, nil, true, "carácter inesperado"},
		{"Falta llave de cierre", `{"a": 1`, nil, true, "objeto no terminado"},
		{"Falta corchete de cierre", `[1, 2`, nil, true, "array no terminado"},
		{"Coma extra en objeto", `{"a": 1,}`, nil, true, "se esperaba '\"'"},
		{"Coma extra en array", `[1, 2,]`, nil, true, "carácter inesperado"},
		{"Caracteres extra al final", `[1, 2] extra`, nil, true, "caracteres inesperados al final"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{}
			result, err := p.ParseJSON(tt.input)

			if (err != nil) != tt.hasError {
				t.Errorf("ParseJSON() error = %v, wantErr %v", err, tt.hasError)
				return
			}

			if tt.hasError && err != nil {
				if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("ParseJSON() error = %q, want to contain %q", err.Error(), tt.errorContains)
				}
			}

			if !tt.hasError && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ParseJSON() = %v, want %v", result, tt.expected)
			}
		})
	}
}
