package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type ParseRequest struct {
	JSON string `json:"json"`
}

type ParseResponse struct {
	Success      bool           `json:"success"`
	Result       interface{}    `json:"result,omitempty"`
	Error        string         `json:"error,omitempty"`
	ParseTime    string         `json:"parse_time,omitempty"`
	Method       string         `json:"method"`
	Performance  string         `json:"performance,omitempty"`
	JSONType     string         `json:"json_type,omitempty"`
	ElementCount map[string]int `json:"element_count,omitempty"`
}

// Parser global para reutilizar regex compiladas (máximo rendimiento)
var globalParser = NewParser()

func main() {
	// Verificar que el parser está funcionando correctamente
	fmt.Println("🔥 Inicializando Parser JSON con Expresiones Regulares...")
	testResult, testErr := globalParser.ParseJSON(`{"test": "working"}`)
	if testErr != nil {
		log.Fatalf("❌ Error inicializando parser: %v", testErr)
	}
	fmt.Printf("✅ Parser inicializado correctamente: %+v\n", testResult)

	// Servir archivos estáticos
	http.Handle("/", http.FileServer(http.Dir("./static/")))

	// API endpoints optimizados
	http.HandleFunc("/api/parse", parseHandler)
	http.HandleFunc("/api/validate", validateHandler)
	http.HandleFunc("/api/analyze", analyzeJSONHandler)
	http.HandleFunc("/api/benchmark", benchmarkHandler)
	http.HandleFunc("/api/examples", examplesHandler)
	http.HandleFunc("/api/convert-to-go", convertToGoHandler)

	fmt.Println("🚀 PARSER JSON CON EXPRESIONES REGULARES")
	fmt.Println("📁 Sirviendo archivos desde: ./static/")
	fmt.Println("🌐 Accede a: http://localhost:8080")
	fmt.Println()
	fmt.Println("⚡ OPTIMIZACIONES IMPLEMENTADAS:")
	fmt.Println("   • 🔥 Expresiones Regulares Precompiladas")
	fmt.Println("   • ⚡ Detección Rápida de Tipos JSON")
	fmt.Println("   • 🧠 Parsing Estructural Inteligente")
	fmt.Println("   • 🛡️ Validación Completa con Regex")
	fmt.Println("   • 📊 Análisis de Elementos")
	fmt.Println("   • 🚀 Rendimiento Optimizado")
	fmt.Println("   • 📄 Conversión de TXT a Go")
	fmt.Println()
	fmt.Println("🔧 API ENDPOINTS DISPONIBLES:")
	fmt.Println("   POST /api/parse        - Parsing con regex")
	fmt.Println("   POST /api/validate     - Validación rápida")
	fmt.Println("   POST /api/analyze      - Análisis completo del JSON")
	fmt.Println("   POST /api/benchmark    - Comparación de rendimiento")
	fmt.Println("   POST /api/convert-to-go - Convierte TXT a código Go")
	fmt.Println("   GET  /api/examples     - Ejemplos de prueba")
	fmt.Println()
	fmt.Println("💡 TÉCNICAS DE OPTIMIZACIÓN:")
	fmt.Println("   • Regex patterns para extracción directa")
	fmt.Println("   • Eliminación de parsing manual caracter por caracter")
	fmt.Println("   • Procesamiento inteligente de estructuras anidadas")
	fmt.Println("   • Validación estructural con patterns")
	fmt.Println("   • Manejo eficiente de escape sequences")
	fmt.Println()
	fmt.Println("⏹️  Presiona Ctrl+C para detener el servidor")
	fmt.Println(strings.Repeat("=", 80))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func parseHandler(w http.ResponseWriter, r *http.Request) {
	setupCORS(w)
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req ParseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Error al decodificar la solicitud: "+err.Error(), "regex_parser")
		return
	}

	req.JSON = strings.TrimSpace(req.JSON)
	if req.JSON == "" {
		respondWithError(w, "El JSON no puede estar vacío", "regex_parser")
		return
	}

	// PARSING CON REGEX - MÁXIMO RENDIMIENTO
	startTime := time.Now()
	result, err := globalParser.ParseJSON(req.JSON)
	parseTime := time.Since(startTime)

	// Análisis adicional del JSON
	jsonType := globalParser.ExtractJSONType(req.JSON)
	elementCount, _ := globalParser.CountJSONElements(req.JSON)

	if err != nil {
		response := ParseResponse{
			Success:      false,
			Error:        err.Error(),
			ParseTime:    parseTime.String(),
			Method:       "regex_parser",
			Performance:  "error",
			JSONType:     jsonType,
			ElementCount: elementCount,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := ParseResponse{
		Success:      true,
		Result:       result,
		ParseTime:    parseTime.String(),
		Method:       "regex_parser",
		Performance:  determinePerformanceLevel(parseTime),
		JSONType:     jsonType,
		ElementCount: elementCount,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	setupCORS(w)
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req ParseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Error al decodificar la solicitud: "+err.Error(), "regex_validator")
		return
	}

	req.JSON = strings.TrimSpace(req.JSON)
	if req.JSON == "" {
		respondWithError(w, "El JSON no puede estar vacío", "regex_validator")
		return
	}

	// VALIDACIÓN ULTRA-RÁPIDA CON REGEX
	startTime := time.Now()
	err := globalParser.FastValidateJSON(req.JSON)
	validateTime := time.Since(startTime)

	jsonType := globalParser.ExtractJSONType(req.JSON)

	if err != nil {
		response := ParseResponse{
			Success:     false,
			Error:       err.Error(),
			ParseTime:   validateTime.String(),
			Method:      "regex_validator",
			Performance: "validation_error",
			JSONType:    jsonType,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := ParseResponse{
		Success:     true,
		Result:      fmt.Sprintf("JSON %s válido (validado con regex patterns)", jsonType),
		ParseTime:   validateTime.String(),
		Method:      "regex_validator",
		Performance: determinePerformanceLevel(validateTime),
		JSONType:    jsonType,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func analyzeJSONHandler(w http.ResponseWriter, r *http.Request) {
	setupCORS(w)
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req ParseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Error al decodificar la solicitud: "+err.Error(), "regex_analyzer")
		return
	}

	req.JSON = strings.TrimSpace(req.JSON)
	if req.JSON == "" {
		respondWithError(w, "El JSON no puede estar vacío", "regex_analyzer")
		return
	}

	// ANÁLISIS COMPLETO CON REGEX
	startTime := time.Now()

	// Validación
	validationErr := globalParser.FastValidateJSON(req.JSON)

	// Detección de tipo
	jsonType := globalParser.ExtractJSONType(req.JSON)

	// Conteo de elementos
	elementCount, _ := globalParser.CountJSONElements(req.JSON)

	// Parsing completo si es válido
	var parseResult interface{}
	var parseErr error
	if validationErr == nil {
		parseResult, parseErr = globalParser.ParseJSON(req.JSON)
	}

	analysisTime := time.Since(startTime)

	analysis := map[string]interface{}{
		"validation": map[string]interface{}{
			"is_valid": validationErr == nil,
			"error":    getErrorString(validationErr),
		},
		"structure": map[string]interface{}{
			"type":          jsonType,
			"element_count": elementCount,
			"size_bytes":    len(req.JSON),
			"size_chars":    len([]rune(req.JSON)),
		},
		"parsing": map[string]interface{}{
			"success": parseErr == nil,
			"error":   getErrorString(parseErr),
			"result":  parseResult,
		},
		"performance": map[string]interface{}{
			"analysis_time": analysisTime.String(),
			"speed_level":   determinePerformanceLevel(analysisTime),
		},
		"regex_optimizations": []string{
			"Detección de tipo con regex patterns",
			"Validación estructural sin parsing manual",
			"Conteo de elementos con regex",
			"Extracción directa de contenido",
			"Balance de estructuras optimizado",
		},
	}

	response := map[string]interface{}{
		"success": true,
		"method":  "regex_analyzer",
		"data":    analysis,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func benchmarkHandler(w http.ResponseWriter, r *http.Request) {
	setupCORS(w)
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req ParseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Error al decodificar la solicitud: "+err.Error(), "benchmark")
		return
	}

	req.JSON = strings.TrimSpace(req.JSON)
	if req.JSON == "" {
		respondWithError(w, "El JSON no puede estar vacío", "benchmark")
		return
	}

	// BENCHMARK COMPREHENSIVO
	results := make(map[string]interface{})

	// 1. Parser con Regex
	startTime := time.Now()
	regexResult, regexErr := globalParser.ParseJSON(req.JSON)
	regexTime := time.Since(startTime)

	// 2. Validación rápida
	startTime = time.Now()
	validationErr := globalParser.FastValidateJSON(req.JSON)
	validationTime := time.Since(startTime)

	// 3. Parser nativo de Go (para comparación)
	startTime = time.Now()
	var nativeResult interface{}
	nativeErr := json.Unmarshal([]byte(req.JSON), &nativeResult)
	nativeTime := time.Since(startTime)

	// 4. Análisis completo
	startTime = time.Now()
	jsonType := globalParser.ExtractJSONType(req.JSON)
	elementCount, _ := globalParser.CountJSONElements(req.JSON)
	analysisTime := time.Since(startTime)

	// Calcular mejoras de rendimiento
	var speedupVsNative float64
	if nativeTime.Nanoseconds() > 0 {
		speedupVsNative = float64(nativeTime.Nanoseconds()) / float64(regexTime.Nanoseconds())
	}

	var validationSpeedup float64
	if regexTime.Nanoseconds() > 0 {
		validationSpeedup = float64(regexTime.Nanoseconds()) / float64(validationTime.Nanoseconds())
	}

	results["benchmark_results"] = map[string]interface{}{
		"regex_parser": map[string]interface{}{
			"method":      "regex_parsing",
			"time":        regexTime.String(),
			"time_ns":     regexTime.Nanoseconds(),
			"success":     regexErr == nil,
			"error":       getErrorString(regexErr),
			"description": "Parsing con expresiones regulares optimizadas",
		},
		"validation_only": map[string]interface{}{
			"method":      "regex_validation",
			"time":        validationTime.String(),
			"time_ns":     validationTime.Nanoseconds(),
			"success":     validationErr == nil,
			"error":       getErrorString(validationErr),
			"description": "Solo validación de estructura con regex",
		},
		"analysis_complete": map[string]interface{}{
			"method":      "regex_analysis",
			"time":        analysisTime.String(),
			"time_ns":     analysisTime.Nanoseconds(),
			"success":     true,
			"description": "Análisis completo con detección de tipo y conteo",
		},
		"go_native_parser": map[string]interface{}{
			"method":      "go_standard_library",
			"time":        nativeTime.String(),
			"time_ns":     nativeTime.Nanoseconds(),
			"success":     nativeErr == nil,
			"error":       getErrorString(nativeErr),
			"description": "Parser nativo de Go (encoding/json)",
		},
	}

	results["performance_analysis"] = map[string]interface{}{
		"regex_vs_native": map[string]interface{}{
			"speedup_factor":    fmt.Sprintf("%.2fx", speedupVsNative),
			"time_improvement":  (nativeTime - regexTime).String(),
			"percentage_faster": fmt.Sprintf("%.1f%%", (speedupVsNative-1)*100),
		},
		"validation_efficiency": map[string]interface{}{
			"speedup_factor":     fmt.Sprintf("%.2fx", validationSpeedup),
			"validation_vs_full": fmt.Sprintf("%.1f%% del tiempo de parsing completo", float64(validationTime.Nanoseconds())/float64(regexTime.Nanoseconds())*100),
		},
		"fastest_operation": determineFastestOperation(regexTime, validationTime, nativeTime, analysisTime),
	}

	results["json_analysis"] = map[string]interface{}{
		"type":          jsonType,
		"element_count": elementCount,
		"complexity":    determineComplexity(elementCount),
		"size_metrics": map[string]interface{}{
			"bytes":      len(req.JSON),
			"characters": len([]rune(req.JSON)),
			"lines":      strings.Count(req.JSON, "\n") + 1,
		},
	}

	results["accuracy_check"] = map[string]interface{}{
		"regex_vs_native":   compareResults(regexResult, nativeResult),
		"results_identical": compareResults(regexResult, nativeResult),
	}

	response := map[string]interface{}{
		"success": true,
		"method":  "comprehensive_regex_benchmark",
		"data":    results,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func examplesHandler(w http.ResponseWriter, r *http.Request) {
	setupCORS(w)
	if r.Method == "OPTIONS" {
		return
	}

	examples := map[string]interface{}{
		"ejemplos": []map[string]interface{}{
			{
				"nombre": "Objeto simple",
				"json":   `{"name": "Juan", "age": 30, "active": true}`,
				"tipo":   "object",
			},
			{
				"nombre": "Array de números",
				"json":   `[1, 2, 3, 4, 5]`,
				"tipo":   "array",
			},
			{
				"nombre": "String con escape",
				"json":   `"Hola\nMundo con \"comillas\""`,
				"tipo":   "string",
			},
			{
				"nombre": "Número científico",
				"json":   `1.23e-4`,
				"tipo":   "number",
			},
			{
				"nombre": "Boolean verdadero",
				"json":   `true`,
				"tipo":   "boolean",
			},
			{
				"nombre": "Valor null",
				"json":   `null`,
				"tipo":   "null",
			},
		},
		"ejemplos_complejos": []map[string]interface{}{
			{
				"nombre": "Estructura anidada",
				"json":   `{"empresa": {"nombre": "TechCorp", "empleados": [{"id": 1, "nombre": "Ana"}, {"id": 2, "nombre": "Carlos"}]}}`,
				"tipo":   "object",
			},
			{
				"nombre": "Array multidimensional",
				"json":   `[[1, 2], ["a", "b"], [true, false]]`,
				"tipo":   "array",
			},
			{
				"nombre": "JSON con todos los tipos",
				"json":   `{"string": "texto", "number": 42.5, "boolean": true, "null": null, "array": [1, 2], "object": {"nested": "value"}}`,
				"tipo":   "object",
			},
		},
		"ejemplos_invalidos": []map[string]interface{}{
			{
				"nombre": "Coma extra en objeto",
				"json":   `{"a": 1, "b": 2,}`,
				"error":  "coma extra",
			},
			{
				"nombre": "Comillas faltantes en clave",
				"json":   `{name: "Juan"}`,
				"error":  "clave sin comillas",
			},
			{
				"nombre": "Llave no cerrada",
				"json":   `{"a": 1, "b": 2`,
				"error":  "estructura no cerrada",
			},
			{
				"nombre": "String no terminado",
				"json":   `{"mensaje": "hola mundo`,
				"error":  "string no cerrado",
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(examples)
}

// Funciones auxiliares

func setupCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
}

func respondWithError(w http.ResponseWriter, errorMsg, method string) {
	response := ParseResponse{
		Success: false,
		Error:   errorMsg,
		Method:  method,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getErrorString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func determinePerformanceLevel(duration time.Duration) string {
	switch {
	case duration < 10*time.Microsecond:
		return "ultra_lightning"
	case duration < 50*time.Microsecond:
		return "lightning_fast"
	case duration < 100*time.Microsecond:
		return "ultra_fast"
	case duration < 1*time.Millisecond:
		return "very_fast"
	case duration < 10*time.Millisecond:
		return "fast"
	default:
		return "moderate"
	}
}

func determineFastestOperation(regex, validation, native, analysis time.Duration) string {
	operations := map[string]time.Duration{
		"regex_parsing":    regex,
		"regex_validation": validation,
		"native_parsing":   native,
		"regex_analysis":   analysis,
	}

	fastest := "regex_parsing"
	fastestTime := regex

	for operation, duration := range operations {
		if duration < fastestTime {
			fastest = operation
			fastestTime = duration
		}
	}

	return fastest
}

func determineComplexity(elementCount map[string]int) string {
	total := 0
	for _, count := range elementCount {
		total += count
	}

	switch {
	case total <= 5:
		return "simple"
	case total <= 20:
		return "moderate"
	case total <= 100:
		return "complex"
	default:
		return "very_complex"
	}
}

func convertToGoHandler(w http.ResponseWriter, r *http.Request) {
	setupCORS(w)
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Parsear el formulario multipart
	err := r.ParseMultipartForm(10 << 20) // 10 MB máximo
	if err != nil {
		respondWithError(w, "Error al parsear el formulario: "+err.Error(), "txt_to_go_converter")
		return
	}

	// Obtener el archivo del formulario
	file, header, err := r.FormFile("txtFile")
	if err != nil {
		respondWithError(w, "Error al obtener el archivo: "+err.Error(), "txt_to_go_converter")
		return
	}
	defer file.Close()

	// Verificar que sea un archivo .txt
	if !strings.HasSuffix(strings.ToLower(header.Filename), ".txt") {
		respondWithError(w, "Solo se permiten archivos .txt", "txt_to_go_converter")
		return
	}

	// Leer el contenido del archivo
	content := make([]byte, header.Size)
	_, err = file.Read(content)
	if err != nil {
		respondWithError(w, "Error al leer el archivo: "+err.Error(), "txt_to_go_converter")
		return
	}

	// Obtener parámetros opcionales del formulario
	packageName := r.FormValue("packageName")
	if packageName == "" {
		packageName = "main"
	}

	variableName := r.FormValue("variableName")
	if variableName == "" {
		variableName = "textContent"
	}

	conversionType := r.FormValue("conversionType")
	if conversionType == "" {
		conversionType = "variable"
	}

	// Convertir a código Go
	startTime := time.Now()
	goCode := convertTextToGo(string(content), packageName, variableName, conversionType, header.Filename)
	conversionTime := time.Since(startTime)

	// Responder con el código Go generado
	response := map[string]interface{}{
		"success":         true,
		"method":          "txt_to_go_converter",
		"original_file":   header.Filename,
		"file_size":       header.Size,
		"conversion_time": conversionTime.String(),
		"go_code":         goCode,
		"parameters": map[string]interface{}{
			"package_name":    packageName,
			"variable_name":   variableName,
			"conversion_type": conversionType,
		},
		"download_filename": generateGoFilename(header.Filename, conversionType),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func convertTextToGo(content, packageName, variableName, conversionType, originalFilename string) string {
	var builder strings.Builder

	// Header del archivo Go
	builder.WriteString(fmt.Sprintf("package %s\n\n", packageName))
	builder.WriteString(fmt.Sprintf("// Archivo generado automáticamente desde: %s\n", originalFilename))
	builder.WriteString(fmt.Sprintf("// Generado el: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	switch conversionType {
	case "variable":
		builder.WriteString(fmt.Sprintf("// %s contiene el contenido del archivo de texto\n", variableName))
		builder.WriteString(fmt.Sprintf("var %s = `%s`\n", variableName, content))

	case "const":
		builder.WriteString(fmt.Sprintf("// %s contiene el contenido del archivo de texto como constante\n", variableName))
		builder.WriteString(fmt.Sprintf("const %s = `%s`\n", variableName, content))

	case "function":
		funcName := strings.Title(variableName)
		builder.WriteString(fmt.Sprintf("// Get%s retorna el contenido del archivo de texto\n", funcName))
		builder.WriteString(fmt.Sprintf("func Get%s() string {\n", funcName))
		builder.WriteString(fmt.Sprintf("\treturn `%s`\n", content))
		builder.WriteString("}\n")

	case "struct":
		structName := strings.Title(variableName)
		builder.WriteString(fmt.Sprintf("// %s contiene datos de texto estructurados\n", structName))
		builder.WriteString(fmt.Sprintf("type %s struct {\n", structName))
		builder.WriteString("\tContent  string\n")
		builder.WriteString("\tFilename string\n")
		builder.WriteString("\tSize     int\n")
		builder.WriteString("}\n\n")
		builder.WriteString(fmt.Sprintf("// New%s crea una nueva instancia con el contenido del archivo\n", structName))
		builder.WriteString(fmt.Sprintf("func New%s() *%s {\n", structName, structName))
		builder.WriteString(fmt.Sprintf("\treturn &%s{\n", structName))
		builder.WriteString(fmt.Sprintf("\t\tContent:  `%s`,\n", content))
		builder.WriteString(fmt.Sprintf("\t\tFilename: \"%s\",\n", originalFilename))
		builder.WriteString(fmt.Sprintf("\t\tSize:     %d,\n", len(content)))
		builder.WriteString("\t}\n")
		builder.WriteString("}\n")

	case "slice":
		lines := strings.Split(content, "\n")
		builder.WriteString(fmt.Sprintf("// %s contiene las líneas del archivo como slice\n", variableName))
		builder.WriteString(fmt.Sprintf("var %s = []string{\n", variableName))
		for _, line := range lines {
			// Escapar caracteres especiales
			escapedLine := strings.ReplaceAll(line, "`", "` + \"`\" + `")
			builder.WriteString(fmt.Sprintf("\t`%s`,\n", escapedLine))
		}
		builder.WriteString("}\n")

	case "map":
		lines := strings.Split(content, "\n")
		builder.WriteString(fmt.Sprintf("// %s contiene las líneas del archivo como map[int]string\n", variableName))
		builder.WriteString(fmt.Sprintf("var %s = map[int]string{\n", variableName))
		for i, line := range lines {
			escapedLine := strings.ReplaceAll(line, "`", "` + \"`\" + `")
			builder.WriteString(fmt.Sprintf("\t%d: `%s`,\n", i+1, escapedLine))
		}
		builder.WriteString("}\n")

	default:
		// Por defecto, usar variable
		builder.WriteString(fmt.Sprintf("var %s = `%s`\n", variableName, content))
	}

	return builder.String()
}

func generateGoFilename(originalFilename, conversionType string) string {
	baseName := strings.TrimSuffix(originalFilename, ".txt")
	baseName = strings.ReplaceAll(baseName, " ", "_")
	baseName = strings.ReplaceAll(baseName, "-", "_")

	// Convertir a snake_case válido para Go
	var result strings.Builder
	for i, char := range baseName {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '_' {
			if i == 0 && char >= '0' && char <= '9' {
				result.WriteString("file_")
			}
			result.WriteRune(char)
		} else {
			result.WriteString("_")
		}
	}

	return strings.ToLower(result.String()) + ".go"
}

func compareResults(result1, result2 interface{}) bool {
	json1, err1 := json.Marshal(result1)
	json2, err2 := json.Marshal(result2)

	if err1 != nil || err2 != nil {
		return false
	}

	return string(json1) == string(json2)
}
