package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type EditorDetails struct {
	VsCodePath           string
	FileType             string
	EditorName           []string
	FolderToOpenInEditor string
}

type EditorDetailsForVSCode struct {
	*EditorDetails
	FileManagerForVSCode *FileManager
}

type IEditorHandler interface {
	LaunchVSCode()
	SetOrUpdatePath()
}

func (editor *EditorDetailsForVSCode) SetOrUpdatePath() error {
	// Open the .zshrc file
	filePath := os.ExpandEnv("$HOME/.zshrc")
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()
	// Read the existing content of the file
	scanner := bufio.NewScanner(file)
	exists, existing_entry := editor.FileManagerForVSCode.CheckIfEntryExists(scanner, editor.FileManagerForVSCode.ProjectPathKey)
	if exists {
		rootProjectDirectory := strings.Split(existing_entry, "=")[1]
		editor.FileManagerForVSCode.UpdateRoots(rootProjectDirectory)
		// If entry is same as the project path value then return and do nothing.
		if rootProjectDirectory == editor.FileManagerForVSCode.ProjectPathValue {
			fmt.Printf("%s[message:] We are not updating the path as it is already set to: %s \n", Cyan,rootProjectDirectory)
			return nil
		}
		//entry exists but its different from the project path value then update the entry
		if user_ans := editor.FileManagerForVSCode.AskIfUserUpdatesToNewPath(); user_ans {
			editor.FileManagerForVSCode.UpdateRoots(editor.FileManagerForVSCode.ProjectPathValue)
			editor.WriteContentToProfile(scanner, filePath, file)
		}
	} else {
		// If entry does not exist then check for the project path
		var path string
		if editor.FileManagerForVSCode.ProjectPathValue != "" {
			path = editor.FileManagerForVSCode.ProjectPathValue
		} else {
			path = editor.FileManagerForVSCode.GetFilePathInput()
		}
		if path == "" {
			return fmt.Errorf("[+]Error : You have failed to provide the project path. Please try again")
		} else {
			editor.FileManagerForVSCode.UpdateRoots(path)
			editor.WriteContentToProfile(scanner, filePath, file)
		}
	}
	return nil
}

func (editor *EditorDetailsForVSCode) WriteContentToProfile(scanner *bufio.Scanner, filePath string, file *os.File) {
	str_content := editor.FileManagerForVSCode.NewContentForZshFile(scanner)
	editor.FileManagerForVSCode.WriteContentToFile(filePath, str_content)
	file.Sync()

}

func (editor *EditorDetailsForVSCode) LaunchVSCode() {
	filePath, err := editor.FileManagerForVSCode.SearchPackage(editor.FolderToOpenInEditor)

	if err != nil {
		fmt.Printf(BgRed+"Error launching VS Code instance: %v\n", err)
		return
	}
	var cmd *exec.Cmd

	if filePath != "" {
		cmd = exec.Command(editor.VsCodePath, filePath)
	} else {
		cmd = exec.Command(editor.VsCodePath, filePath)
	}
	err = cmd.Run()
	if err != nil {
		fmt.Printf(BgRed+"Error launching VS Code instance: %v\n", err)
	} else {
		fmt.Println(Yellow + "...Launching VS Code")
	}
}
