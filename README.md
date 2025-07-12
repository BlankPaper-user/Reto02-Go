# Reto #2: Parser JSON de Bajo Nivel en Go

Este proyecto es una implementación de un parser de JSON de bajo nivel en Go, desarrollado como parte del Reto #2 del documento "Retos Finales - TLP Go.pdf".

## Enfoque General del Parseo

El parser se ha construido desde cero utilizando únicamente paquetes estándar de Go, sin depender de la librería `encoding/json`. El enfoque principal es un **parser recursivo descendente**. La lógica principal se encuentra en la `struct Parser`, que mantiene el estado del proceso de análisis (la cadena de entrada y la posición actual).

El proceso de parseo comienza con una llamada a la función `ParseJSON`, que a su vez llama a `parseValue`. Esta última función actúa como un despachador que, basándose en el carácter actual, decide qué tipo de dato JSON se va a procesar (objeto, array, string, etc.).

## Manejo de Tipos de Datos JSON

El parser es capaz de manejar los siguientes tipos de datos JSON:

*   **Objetos (`{}`):** Se parsean en un `map[string]interface{}`. La lógica itera sobre pares clave-valor, esperando comas como separadores y una llave de cierre para terminar.
*   **Arrays (`[]`):** Se parsean en un `[]interface{}`. La lógica itera sobre los valores, esperando comas como separadores y un corchete de cierre para terminar.
*   **Strings (`""`):** Se parsean en un `string` de Go. Se manejan los caracteres de escape básicos.
*   **Números:** Se parsean en un `float64` de Go. Se admiten números enteros, decimales y negativos.
*   **Booleanos:** Se parsean en un `bool` de Go (`true` y `false`).
*   **Null:** Se parsea como `nil` en Go.

## Estrategia para el Manejo de Errores

Se ha puesto un énfasis especial en un manejo de errores robusto y detallado. Cuando el parser encuentra un error de sintaxis, devuelve un error que incluye:

1.  **Una descripción clara del problema** (ej. `unexpected character`, `unterminated object`).
2.  **La posición exacta (índice) en la cadena de entrada donde se detectó el error.**

Esto facilita enormemente la depuración de JSON inválido. Además, el parser está diseñado para fallar de forma segura y predecible ante cualquier desviación del formato JSON estándar, como comas extra, llaves o corchetes sin cerrar, o caracteres inesperados.

## Estructura del Proyecto

El proyecto está organizado de la siguiente manera:

*   `main.go`: Contiene la función `main` que sirve como punto de entrada para probar el parser con un ejemplo.
*   `parser.go`: Contiene la lógica principal del parser de JSON.
*   `parser_test.go`: Contiene un conjunto de pruebas unitarias exhaustivas que cubren tanto casos de JSON válido como inválido.
*   `README.md`: Este archivo, que documenta el proyecto.

## Cómo Ejecutar el Proyecto

### Ejecutar las Pruebas

Para verificar que toda la lógica del parser funciona correctamente, ejecuta el siguiente comando:

```bash
go test
```

### Ejecutar el Programa Principal

Para ver el parser en acción con el ejemplo definido en `main.go`, ejecuta:

```bash
go run .
```
