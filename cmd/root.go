/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	folderToOpenInEditor string
	projectPathValue          string
)

var rootCmd = &cobra.Command{
	Use:   "vs-code-cli",
	Short: "Its a cli tool to open your work space with vs code",
	Long: `It help us to open your work space with vs code.
	Its highly extendable and we are planning several features such as 
	extension manager, auto update, etc based on workspace.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		editor := EditorDetailsForVSCode{
			EditorDetails: &EditorDetails{
				VsCodePath:           "/Applications/Visual Studio Code.app/Contents/MacOS/Electron",
				FileType:             "zsh",
				EditorName:           []string{"Visual Studio Code"},
				FolderToOpenInEditor: folderToOpenInEditor,
			},
			FileManagerForVSCode: &FileManager{
				Root:             projectPathValue,
				ProjectPathValue: projectPathValue,
				ProjectPathKey:   "PROJECT_PATH",
			},
		}
		err := editor.SetOrUpdatePath()
		if err == nil {
			editor.LaunchVSCode()
		} else {
			fmt.Printf("%v", err)
		}

	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&folderToOpenInEditor, "folder", "f", "", "folder name")
	rootCmd.Flags().StringVarP(&projectPathValue, "project", "p", "", "project path")
}
