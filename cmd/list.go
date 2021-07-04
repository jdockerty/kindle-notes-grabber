/*
Copyright © 2020 Jack Dockerty

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
	"strings"

	"github.com/jdockerty/kindle-notes-grabber/notes"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the books which the application has ran against in your inbox.",
	Long: `Read from the defined YAML file for the books that have
been parsed from your inbox, these are deemed as 'completed' if they
have been seen by the application. The output from this file is simply
printed to the console in a numbered format.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		completedBooks, err := notes.LoadCompletedBooks()
		if err != nil {
			return err
		}

		fmt.Println("Completed books are:")
		var counter int = 1
		for book, _ := range *completedBooks {

			originalName := reverseFormat(book)
			fmt.Printf(" %d: %s\n", counter, originalName)
			counter += 1
		}

		return nil
	},
}

// reverseFormat is used to remove the dashes and notebook suffix to the
// names of books for the program, this is mainly used as a way to print
// out the books by their original name for clarity.
func reverseFormat(dashedNotebookName string) string {
	spacedNotebookName := strings.ReplaceAll(dashedNotebookName, "-", " ")
	bookName := strings.TrimSuffix(spacedNotebookName, "notebook")
	return bookName
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
