package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const service = "admin_system"

func main() {
	folder := fmt.Sprintf("rpc/%s/desc/proto", service)
	output := fmt.Sprintf("rpc/%s/%s.proto", service, service)

	files, err := os.ReadDir(folder)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	content := []byte(fmt.Sprintf(`syntax = "proto3";

package %s;
option go_package = "./%s";

`, service, service))

	var fileContents []byte
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".proto") {
			continue
		}

		name := filepath.Join(folder, file.Name())
		data, err := os.ReadFile(name)
		if err != nil {
			fmt.Printf("Error reading file %s: %s\n", name, err.Error())
			continue
		}

		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "syntax") ||
				strings.HasPrefix(trimmed, "import") ||
				strings.HasPrefix(trimmed, "//") {
				continue
			}
			fileContents = append(fileContents, line...)
			fileContents = append(fileContents, '\n')
		}
	}
	content = append(content, fileContents...)
	err = os.WriteFile(output, content, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}
	fmt.Println("Combine proto to", output)
}
