package api

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func parsePDF(name string, path string) string {
	newName := []string{}
	newName = append(newName, name)
	newName = append(newName, "Parsed")
	pdfName := strings.Join(newName, "")

	cmd := "pdftotext"
	args := []string{path, "tmp/" + pdfName}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("Successfully created parsed file")

	return "./tmp/" + pdfName
}
