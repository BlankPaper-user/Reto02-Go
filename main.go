package main

import (
	"fmt"
	"log"
)

func main() {

	parser := &Parser{}

	// Prueba con JSON simple
	// json := `{"name": "John", "age": 30, "active": true, "interests": ["football", "music"]}`
	// json := `{"nombre":"Branz Saul", "password": "1234"}`
	//json :=`{"a": 1,}`
	json :=`{"a": 1,adadadd}`
	result, err := parser.ParseJSON(json)
	if err != nil {
		log.Fatal(err) // Si hay un error, lo mostramos
	}

	// Imprimir el resultado
	fmt.Printf("%#v\n", result)
}
