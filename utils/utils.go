package utils

import (
	"fmt"
	"os"
	"os/exec"
)

func EnsureDirectoryExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

func OpenFileWithEditor(editor string, filename string, date string) error {
	// Check if the file already exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// File does not exist, create and write default content
    defaultContent := fmt.Sprintf("> %s.plan\n\n`accomplished:`\n+ \n\n`todo`\n+ \n\n`ideas:`\n+ ", date)
		err := os.WriteFile(filename, []byte(defaultContent), 0644)
		if err != nil {
			return fmt.Errorf("error creating file: %w", err)
		}
		fmt.Printf("Created new file: %s with default content.\n", filename)
	}

	// Open the file with the specified editor
	cmd := exec.Command(editor, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
