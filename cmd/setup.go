/*
Copyright © 2021 Jack Dockerty

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/jdockerty/kindle-notes-grabber/notes"
	"github.com/spf13/cobra"
)

const (
	programDirectoryName = "kindle-notes"
	saveFile             = "completed-notebooks.yaml"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "A simple helper to create the relevant files and folders",
	Long: `A utility sub-command which creates the necessary directory and files for the 
usage of the application. This is relatively minor, but may save some headaches by removing 
the need to go back and forth with running the application, only to receive an error about 
a particular file or directory not being present.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		notesDir := fmt.Sprintf("%s/%s", homeDir, programDirectoryName)
		completedBooksPath := fmt.Sprintf("%s/%s", notesDir, saveFile)

		kindleNotesDirExists, err := notes.Exists(notesDir)
		if err != nil {
			return err
		}

		completedBooksFileExists, err := notes.Exists(completedBooksPath)
		if err != nil {
			return err
		}

		// Check both criteria is met, return if they both exist
		// Or create the save file if only the directory exists.
		if kindleNotesDirExists {
			fmt.Printf("%s directory exists\n", programDirectoryName)

			if completedBooksFileExists {
				fmt.Printf("%s file exists\n", saveFile)
				fmt.Println("No required actions")
				return nil
			} else {
				err := createSaveFile(notesDir, completedBooksPath)
				if err != nil {
					return err
				}
				return nil
			}
		}

		err = createDirectory(homeDir, notesDir)
		if err != nil {
			return err
		}

		err = createSaveFile(notesDir, completedBooksPath)
		if err != nil {
			return err
		}
		fmt.Println("Setup complete, you can now run `kng run'")
		return nil

	},
}

func createDirectory(homeDirectory, fullDirPath string) error {
	// If the directory doesn't exist, then the file cannot be contained within it either
	// so create it now
	fmt.Printf("Creating '%s' directory in '%s'\n", programDirectoryName, homeDirectory)
	err := os.Mkdir(fullDirPath, 0755)
	if err != nil {
		return err
	}

	return nil
}

func createSaveFile(dir, savePath string) error {
	fmt.Printf("Creating '%s' in '%s'\n", saveFile, dir)
	// We can discard the file pointer returned here, as we just need to create it
	_, err := os.Create(savePath)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
