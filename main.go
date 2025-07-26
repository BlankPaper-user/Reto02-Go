package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type ParseRequest struct {
	JSON string `json:"json"`
}

type ParseResponse struct {
	Success bool        `json:"success"`
	Result  interface{} `json:"result,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func main() {
	// Servir archivos estáticos (HTML, CSS, JS)
	http.Handle("/", http.FileServer(http.Dir("./static/")))

	// API endpoint para parsear JSON
	http.HandleFunc("/api/parse", parseJSONHandler)

	// Endpoint para obtener ejemplos
	http.HandleFunc("/api/examples", examplesHandler)

	fmt.Println("🚀 Servidor del Parser JSON iniciado")
	fmt.Println("📁 Sirviendo archivos desde: ./static/")
	fmt.Println("🌐 Accede a: http://localhost:8080")
	fmt.Println("🔧 API endpoints:")
	fmt.Println("   POST /api/parse    - Parsear JSON")
	fmt.Println("   GET  /api/examples - Obtener ejemplos")
	fmt.Println("⏹️  Presiona Ctrl+C para detener")
	fmt.Println(strings.Repeat("-", 50))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func parseJSONHandler(w http.ResponseWriter, r *http.Request) {
	// Configurar CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req ParseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := ParseResponse{
			Success: false,
			Error:   "Error al decodificar la solicitud: " + err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Limpiar espacios en blanco al inicio y final
	req.JSON = strings.TrimSpace(req.JSON)

	if req.JSON == "" {
		response := ParseResponse{
			Success: false,
			Error:   "El JSON no puede estar vacío",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Usar nuestro parser personalizado
	parser := &Parser{}
	result, err := parser.ParseJSON(req.JSON)

	if err != nil {
		response := ParseResponse{
			Success: false,
			Error:   err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := ParseResponse{
		Success: true,
		Result:  result,
	}
	json.NewEncoder(w).Encode(response)
}

func examplesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		return
	}

	examples := map[string]interface{}{
		"ejemplos": []map[string]interface{}{
			{
				"nombre": "Objeto simple",
				"json":   `{"name": "Juan", "age": 30, "active": true}`,
			},
			{
				"nombre": "Array de strings",
				"json":   `["manzana", "plátano", "cereza"]`,
			},
			{
				"nombre": "Objeto anidado",
				"json":   `{"usuario": {"nombre": "Ana", "preferencias": {"tema": "oscuro", "idioma": "es"}}}`,
			},
			{
				"nombre": "Array con objetos",
				"json":   `[{"id": 1, "nombre": "Producto A"}, {"id": 2, "nombre": "Producto B"}]`,
			},
			{
				"nombre": "Tipos mixtos",
				"json":   `{"string": "texto", "number": 42.5, "boolean": true, "null": null, "array": [1, 2, 3]}`,
			},
			{
				"nombre": "JSON complejo",
				"json":   `{"empresa": "TechCorp", "empleados": [{"nombre": "Carlos", "departamento": "IT", "salario": 75000}, {"nombre": "María", "departamento": "Marketing", "salario": 65000}], "fundada": 2010, "activa": true}`,
			},
			{
				"nombre": "Con caracteres escape",
				"json":   `{"mensaje": "Hola\nMundo", "ruta": "C:\\Users\\Juan", "comillas": "Dice \"Hola\""}`,
			},
			{
				"nombre": "Números especiales",
				"json":   `{"entero": 42, "decimal": 3.14159, "negativo": -123, "cientifico": 1.5e10}`,
			},
		},
		"ejemplos_invalidos": []map[string]interface{}{
			{
				"nombre": "Coma extra en objeto",
				"json":   `{"a": 1,}`,
			},
			{
				"nombre": "Comillas faltantes",
				"json":   `{name: "Juan"}`,
			},
			{
				"nombre": "Llave no cerrada",
				"json":   `{"a": 1`,
			},
			{
				"nombre": "Coma extra en array",
				"json":   `[1, 2,]`,
			},
			{
				"nombre": "String no cerrado",
				"json":   `{"mensaje": "hola mundo`,
			},
			{
				"nombre": "Caracteres extra",
				"json":   `{"a": 1} {"b": 2}`,
			},
		},
	}

	// Log para debugging
	fmt.Printf("Enviando ejemplos: %d válidos, %d inválidos\n",
		len(examples["ejemplos"].([]map[string]interface{})),
		len(examples["ejemplos_invalidos"].([]map[string]interface{})))

	if err := json.NewEncoder(w).Encode(examples); err != nil {
		http.Error(w, "Error al codificar ejemplos", http.StatusInternalServerError)
		return
	}
}
