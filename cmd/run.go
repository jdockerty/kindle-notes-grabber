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
	"log"

	"github.com/jdockerty/kindle-notes-grabber/config"
	"github.com/jdockerty/kindle-notes-grabber/notes"
	"github.com/spf13/cobra"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

var conf *config.Config
var cfgFile string

var serviceName string

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "grab your kindle notes from a defined email inbox",
	Long: `Start the application and it will parse the emails which
match the relevant criteria from Amazon as kindle notes, placing them
into a .txt file with various metadata about the note or highlight,
such as page number it was taken and its type.`,
	Run: func(cmd *cobra.Command, args []string) {
		conf = initConfig()
		log.Println("Attempting to connect to providers server...")

		var im config.IMAPServer
		im.Populate(serviceName)

		c, err := client.DialTLS(im.Socket, nil)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Connected")

		// Don't forget to logout
		defer c.Logout()

		// Login
		if err := c.Login(conf.Email, conf.Password); err != nil {
			log.Fatal(err)
		}
		log.Println("Logged in")

			_, err = c.Select(mailbox, false)
	if err != nil {
		log.Fatal(err)
	}

	// NOTE: Hard-coded criteria for now.
	criteria := *imap.NewSearchCriteria()
	criteria.Body = []string{fromAmazon}

	ids := notes.GetEmailIds(c, &criteria)

	var notesCollection []*notes.Notes
	var section imap.BodySectionName

	completedBooks, err := notes.LoadCompletedBooks()
	if err != nil {
		log.Fatal(err)
	}

	for _, id := range ids {
		myNotes := notes.New()
		messages := myNotes.GetAmazonMessage(c, id, section)
		mailReaders := myNotes.GetMailReaders(messages, section)
		myNotes.Populate(mailReaders)

		// If the title exists in the map, skip it
		if _, ok := (*completedBooks)[myNotes.Title]; ok {
			log.Printf("%s already seen", myNotes.Title)
			continue
		}

		log.Printf("Adding %s to map", myNotes.Title)
		(*completedBooks)[myNotes.Title] = struct{}{}

		notes.Write(myNotes)
		notesCollection = append(notesCollection, myNotes)
	}

	err = notes.Save(notesCollection)
	if err != nil {
		log.Fatal(err)
	}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&serviceName, "service", "s", "gmail", "The service name of a particular provider.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() *config.Config {
	conf, err := config.New(cfgFile)
	if err != nil {
		log.Fatalf("Cannot read configuration: %s", err)
	}

	return conf
}
