package util

import (
	"encoding/json"
	"os"
)

func WriteJson(filename string, obj any) error {
	// Open a file for writing
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Encode the struct to JSON and write to the file
	encoder := json.NewEncoder(file)
	return encoder.Encode(obj)
}

func ReadJson[T any](filename string) (T, error) {
	var obj T

	// Open the file for reading
	file, err := os.Open(filename)
	if err != nil {
		return obj, err
	}
	defer file.Close()

	// Decode the JSON data into a map
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&obj)
	return obj, err
}
