# 🚀 Parser JSON de Bajo Nivel

Un parser JSON completo implementado desde cero en Go, sin dependencias externas. Este proyecto forma parte del **Reto #2** y demuestra la implementación de un analizador sintáctico robusto con interfaz web interactiva.

## ✨ Características

### 🔧 Parser JSON Personalizado
- **Implementación desde cero** - Sin usar `encoding/json` de Go
- **Tipos soportados**:
  - ✅ Objetos JSON → `map[string]interface{}`
  - ✅ Arrays JSON → `[]interface{}`
  - ✅ Strings con caracteres de escape
  - ✅ Números (enteros, decimales, notación científica)
  - ✅ Booleanos (`true`/`false`)
  - ✅ Valores `null`
  - ✅ Estructuras anidadas complejas

### 🎯 Características Avanzadas
- **Manejo detallado de errores** con línea y columna exacta
- **Validación estricta** de sintaxis JSON
- **Detección de errores comunes** (comas extra, llaves sin cerrar, etc.)
- **Soporte completo para caracteres de escape** (`\n`, `\"`, `\\`, etc.)
- **Validación de claves duplicadas** en objetos
- **Números con notación científica** (1e10, 1.5E-3)

### 🌐 Interfaz Web Interactiva
- **Editor JSON** con estadísticas en tiempo real
- **Ejemplos predefinidos** válidos e inválidos
- **Formateo automático** de JSON
- **Respuesta en tiempo real** con detalles del error
- **Interfaz responsiva** con Bootstrap 5

## 🏗️ Arquitectura del Proyecto

```
📁 Reto02-Go/
├── 📄 main.go          # Servidor HTTP y endpoints API
├── 📄 parser.go        # Implementación del parser JSON
├── 📄 parser_test.go   # Suite completa de tests
├── 📄 go.mod           # Dependencias del módulo
├── 📁 static/
│   └── 📄 index.html   # Interfaz web interactiva
└── 📄 README.md        # Este archivo
```

## 🚀 Instalación y Uso

### Prerrequisitos
- Go 1.24.4 o superior

### Instalación
```bash
# Clonar el repositorio (cualquiera de las dos versiones)
git clone https://github.com/BrSilvinha/Reto02-Go.git
# o
git clone https://github.com/BlankPaper-user/Reto02-Go.git

cd Reto02-Go

# Verificar dependencias
go mod tidy

# Ejecutar tests (opcional)
go test -v

# Iniciar el servidor
go run .
```

### Acceder a la aplicación
```
🌐 Interfaz Web: http://localhost:8080
🔧 API Endpoint: http://localhost:8080/api/parse
📚 Ejemplos: http://localhost:8080/api/examples
```

## 🔌 API Reference

### POST `/api/parse`
Parsea una cadena JSON y devuelve el resultado estructurado.

**Request:**
```json
{
  "json": "{\"name\": \"Juan\", \"age\": 30}"
}
```

**Response (éxito):**
```json
{
  "success": true,
  "result": {
    "name": "Juan",
    "age": 30
  }
}
```

**Response (error):**
```json
{
  "success": false,
  "error": "se esperaba ',' o '}' pero se obtuvo: 'x' en línea 1, columna 15"
}
```

### GET `/api/examples`
Obtiene ejemplos JSON válidos e inválidos para pruebas.

## 🧪 Testing

El proyecto incluye una suite completa de tests que cubren:

- ✅ **Casos válidos**: Objetos, arrays, strings, números, booleanos, null
- ✅ **Casos inválidos**: Sintaxis incorrecta, caracteres extra, estructuras malformadas
- ✅ **Casos edge**: JSON vacío, espacios en blanco, estructuras anidadas complejas
- ✅ **Posición de errores**: Verificación de línea y columna en errores
- ✅ **Benchmarks**: Pruebas de rendimiento
- ✅ **Robustez**: JSON de gran tamaño (1000+ elementos)

### Ejecutar tests
```bash
# Todos los tests
go test -v

# Tests con cobertura
go test -cover

# Benchmarks
go test -bench=.

# Test específico
go test -run TestParseJSON
```

## 📊 Ejemplos de Uso

### Programático (Go)
```go
parser := &Parser{}

// JSON válido
result, err := parser.ParseJSON(`{"name": "Ana", "age": 25}`)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

// result es map[string]interface{}
fmt.Printf("Resultado: %+v\n", result)
```

### API REST (curl)
```bash
# Parsear JSON válido
curl -X POST http://localhost:8080/api/parse \
  -H "Content-Type: application/json" \
  -d '{"json": "{\"hello\": \"world\"}"}'

# Obtener ejemplos
curl http://localhost:8080/api/examples
```

## 🐛 Manejo de Errores

El parser proporciona mensajes de error detallados con ubicación exacta:

```
❌ Error: se esperaba ',' o '}' pero se obtuvo: 'x' en línea 2, columna 15
❌ Error: coma extra antes de '}' en línea 1, columna 10  
❌ Error: cadena de texto no terminada en línea 3, columna 5
❌ Error: clave duplicada 'name' en línea 1, columna 25
```

## 🎯 Casos de Prueba Incluidos

### ✅ JSON Válidos
- Objetos simples y anidados
- Arrays con tipos mixtos
- Strings con caracteres de escape
- Números en todas las formas (enteros, decimales, científicos)
- Estructuras complejas multi-nivel

### ❌ JSON Inválidos
- Comas extra en objetos y arrays
- Llaves y corchetes sin cerrar
- Comillas faltantes en claves
- Caracteres extra al final
- Secuencias de escape inválidas

## 🏆 Características Técnicas

### Rendimiento
- **Parser recursivo descendente** optimizado
- **Manejo eficiente de memoria** sin copias innecesarias
- **Validación en una sola pasada** del input
- **Benchmarks incluidos** para medir rendimiento

### Robustez
- **Validación estricta** según especificación JSON
- **Manejo graceful de errores** sin panics
- **Soporte para archivos grandes** (probado con 1000+ elementos)
- **Detección de casos edge** (JSON vacío, solo espacios, etc.)

## 👥 Autores

Este proyecto fue desarrollado en colaboración por:

- **[@BrSilvinha](https://github.com/BrSilvinha)** - [Reto02-Go](https://github.com/BrSilvinha/Reto02-Go)
- **[@BlankPaper-user](https://github.com/BlankPaper-user)** - [Reto02-Go](https://github.com/BlankPaper-user/Reto02-Go)

## 🤝 Contribuir

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📝 Licencia

Este proyecto está bajo la Licencia MIT. Ver `LICENSE` para más detalles.

## 🎓 Aprendizajes

Este proyecto demuestra:
- **Implementación de parsers** desde cero
- **Manejo de estado** en parsing
- **Diseño de APIs REST** en Go
- **Testing exhaustivo** con casos edge
- **Interfaz web moderna** con JavaScript vanilla
- **Arquitectura modular** y mantenible

---

<div align="center">

**🚀 Desarrollado como parte del Reto #2 - TLP Go**

*Parser JSON implementado desde cero sin dependencias externas*

**👥 Creado por [@BrSilvinha](https://github.com/BrSilvinha) y [@BlankPaper-user](https://github.com/BlankPaper-user)**

</div>