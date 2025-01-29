package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
  homedir "github.com/mitchellh/go-homedir"
)

func EnsureDirectoryExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

func getLastPlanFile(dir string) (string, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.md"))
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "", fmt.Errorf("no plan files found")
	}

	// Sort files by name (which are in ddmmyyyy.md format)
	sort.Strings(files)
	// The last file in the sorted list is the latest one
	return files[len(files)-1], nil
}

func extractTodoAndIdeas(content string) (string, string) {
	todoSection := ""
	ideasSection := ""

	lines := strings.Split(content, "\n")
	for i := 0; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "`todo`") {
			for j := i + 1; j < len(lines); j++ {
				if strings.HasPrefix(lines[j], "`ideas:`") {
					break
				}
				todoSection += lines[j] + "\n"
			}
		}
		if strings.HasPrefix(lines[i], "`ideas:`") {
			for j := i + 1; j < len(lines); j++ {
				if strings.HasPrefix(lines[j], "`") {
					break
				}
				ideasSection += lines[j] + "\n"
			}
		}
	}

	return todoSection, ideasSection
}

func OpenFileWithEditor(editor string, filename string, date string, carryOver bool) error {
	// Check if the file already exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// File does not exist, create and write default content
		defaultContent := fmt.Sprintf("> %s.plan\n\n`accomplished:`\n+ \n\n`todo`\n+ \n\n`ideas:`\n+ ", date)

		// Get the last plan file
    if carryOver {
      lastPlanFile, err := getLastPlanFile(filepath.Dir(filename))
      if err == nil {
        content, err := os.ReadFile(lastPlanFile)
        if err == nil {
          todo, ideas := extractTodoAndIdeas(string(content))
          defaultContent = fmt.Sprintf("> %s.plan\n\n`accomplished:`\n+ \n\n`todo`\n%s\n`ideas:`\n%s", date, todo, ideas)
        }
      }
    }

		err = os.WriteFile(filename, []byte(defaultContent), 0644)
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

func DisplayFileContents(filename string) error {
  expandedFilename, err := homedir.Expand(filename)
  if _, err := os.Stat(expandedFilename); os.IsNotExist(err) {
    return fmt.Errorf("file does not exist")
  }
  content, err := os.ReadFile(expandedFilename)
  if err != nil {
    return err
  }
  fmt.Println(string(content))
  return nil
}

