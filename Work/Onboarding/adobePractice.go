package main

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

func main() {
	// Define the registry key path
	keyPath := `Software\Adobe\Adobe Acrobat\DC\Privileged`

	// Define the registry value to be modified
	valueName := "bProtectedMode"

	// Define the new value for the registry key
	newValue := uint32(0)

	// Open the registry key for editing
	key, err := registry.OpenKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
	if err != nil {
		fmt.Println("Error opening registry key:", err)
		return
	}
	defer key.Close()

	// Set the new value for the registry key
	err = key.SetDWordValue(valueName, newValue)
	if err != nil {
		fmt.Println("Error setting registry value:", err)
		return
	}

	fmt.Println("Registry value updated successfully!")
}
