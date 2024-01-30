package index

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const DEFAULT_SECRETS_PATH = "/etc/secrets/sapcp/"

// readK8SServices reads and returns the secrets from a directory
func ReadK8SServices() map[string]interface{} {
	log.Println("readK8SServices")
	// Define the default secrets path

	// Assign the secrets path to a variable
	secretsPath := DEFAULT_SECRETS_PATH

	// Declare a variable to store the result
	var result map[string]interface{}

	// Check if the secrets path exists
	if _, err := os.Stat(secretsPath); !os.IsNotExist(err) {
		// Read the secrets from the path
		result = ReadSecrets(secretsPath)
	}

	// Return the result
	return result
}

// readSecrets reads and returns the secrets from a directory
func ReadSecrets(secretsPath string) map[string]interface{} {
	// Check if the secrets path is a directory
	info, err := os.Stat(secretsPath)
	if err != nil || !info.IsDir() {
		// Handle the error
		log.Fatalf("secrets path must be a directory: %v", err)
	}

	// Create a map to store the result
	result := make(map[string]interface{})

	// Read the directory entries
	entries, err := os.ReadDir(secretsPath)
	if err != nil {
		// Handle the error
		log.Fatalf("error reading directory: %v", err)
	}

	// Loop over the entries
	for _, entry := range entries {
		// Get the service name
		serviceName := entry.Name()

		// Join the secrets path and the service name
		servicePath := filepath.Join(secretsPath, serviceName)

		// Check if the service path is a directory
		info, err := os.Stat(servicePath)
		if err != nil || !info.IsDir() {
			// Skip the entry
			continue
		}

		// Read the service instances and merge them with the result
		serviceInstances := readServiceInstances(serviceName, servicePath)
		for key, value := range serviceInstances {
			result[key] = value
		}
	}

	// Return the result
	return result
}

// readServiceInstances reads and returns the service instances from a directory
func readServiceInstances(serviceName, servicePath string) map[string]interface{} {
	// Create a map to store the result
	result := make(map[string]interface{})

	// Read the directory entries
	entries, err := os.ReadDir(servicePath)
	if err != nil {
		// Handle the error
		log.Fatalf("error reading directory: %v", err)
	}

	// Loop over the entries
	for _, entry := range entries {
		// Get the instance name
		instanceName := entry.Name()

		// Join the service path and the instance name
		instancePath := filepath.Join(servicePath, instanceName)

		// Check if the instance path is a directory
		info, err := entry.Info()
		if err != nil || !info.IsDir() {
			// Skip the entry
			continue
		}

		// Read the instance and assign it to the result
		result[instanceName] = readInstance(serviceName, instanceName, instancePath)
	}

	// Return the result
	return result
}

// readInstance reads and returns the instance from a directory
func readInstance(serviceName, instanceName, instancePath string) map[string]interface{} {
	// Read the files from the instance path
	credentials := readFiles(instancePath)

	// Create a map to store the instance
	instance := make(map[string]interface{})
	instance["credentials"] = credentials
	instance["name"] = instanceName
	instance["label"] = serviceName

	// // Convert the instance to a JSON string
	// data, err := json.Marshal(instance)
	// if err != nil {
	// 	// Handle the error
	// 	log.Fatalf("error converting instance to JSON: %v", err)
	// }

	// Return the JSON string
	//return string(data)
	return instance
}

// readFiles reads and returns the files from a directory
func readFiles(dirPath string) map[string]string {
	// Create a map to store the result
	result := make(map[string]string)

	// Read the directory entries
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		// Handle the error
		log.Fatalf("error reading directory: %v", err)
	}

	// Loop over the entries
	for _, entry := range entries {
		// Get the file name
		file := entry.Name()

		// Join the directory path and the file name
		filePath := filepath.Join(dirPath, file)

		// Check if the file path is a file
		info, err := entry.Info()
		if err != nil || info.IsDir() {
			// Skip the entry
			continue
		}

		// Read the file content and assign it to the result
		result[file] = readFileContent(filePath)
	}

	// Return the result
	return result
}

// isJsonObject checks if a string is a valid JSON string
func isJsonObject(str string) bool {
	// Convert the string to a byte slice
	data := []byte(str)

	// Use the json.Valid function to check the validity
	return json.Valid(data)
}

// readFileContent reads and parses the content of a file
func readFileContent(filePath string) string {
	// Read the file content into a byte slice
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		// Handle the error
		log.Printf("Error reading file %s: %v", filePath, err)
		return ""
	}

	// Parse the JSON content into a Go value
	value := string(content)

	// Return the parsed value
	return value
}
