package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type FileManager struct {
	Root             string
	ProfileType      string
	FileType         string
	ProjectPathValue string
	ProjectPathKey   string
}

type FileHandler interface {
	CheckIfEntryExists(scanner *bufio.Scanner, key string) (result bool, str string)
}

type IFileHandlerForZsh interface {
	FileHandler
	UpdateContentInZsh(scanner *bufio.Scanner) string
}

func (fm *FileManager) CheckIfEntryExists(scanner *bufio.Scanner, key string) (exists bool, str string) {
	for scanner.Scan() {
		str = scanner.Text()
		if strings.Contains(str, key) {
			fmt.Println(Green+"Entry already exists for path variable in the zsh file", str)
			return true, str
		}
	}
	return exists, str
}

func (fm *FileManager) WriteContentToFile(filePath, str_content string) {
	err := os.WriteFile(filePath, []byte(str_content), 0644)
	if err != nil {
		fmt.Println(BoldHiRed+"Error writing to zshrc file:", err)

	}
	fmt.Println(BgHiGreen + "Successfully wrote variable to zshrc file!")
}

func (fm *FileManager) NewContentForZshFile(scanner *bufio.Scanner) string {
	var content []string
	var str_content string
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("Line:", line)
		if strings.Contains(line, fmt.Sprintf("%s=", fm.ProjectPathKey)) {
			line = fmt.Sprintf("%s=%s", fm.ProjectPathKey, fm.ProjectPathValue)
		}
		content = append(content, line, "\n")
	}
	joinedString := strings.Join(content, " ")

	newLine := fmt.Sprintf("export %s=%s\n", fm.ProjectPathKey, fm.ProjectPathValue)
	str_content = joinedString + newLine
	return str_content
}
func (fm *FileManager) GetFilePathInput() (answer string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(UnderlineCyan + "Please provide the path to your project directory: ")
	answer, _ = reader.ReadString('\n')
	answer = answer[:len(answer)-1] // Remove newline character
	return answer
}
func (fm *FileManager) AskIfUserUpdatesToNewPath() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(UnderlineCyan + "Do you wish to update it with the current path (y/n): ")
	answer, _ := reader.ReadString('\n')
	answer = answer[:len(answer)-1] // Remove newline character
	res := strings.ToLower(answer)
	if res == "y" || res == "yes" {
		fmt.Println("You answered:", "Yes")
		return true
	}
	fmt.Println("You answered:", "No")
	return false
}
func (fm *FileManager) UpdateRoots(folderName string) {
	fm.Root = folderName
}

func (editor *FileManager) SearchPackage(packageName string) (string, error) {
	// Walk through the directory structure
	res := ""
	err := filepath.Walk(editor.Root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				fmt.Printf(Red+"Permission is denied for the path: %s\n", path)
				return filepath.SkipDir
			}
			return err
		}
		if info.IsDir() && strings.Contains(path, packageName) {
			res = path
			return filepath.SkipDir
		}
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("%s Not able to find the provided package or path so exiting please try again: %v", UnderlineRed, err)
	}
	return res, err
}
