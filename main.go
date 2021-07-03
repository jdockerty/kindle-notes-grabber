package main

import (
	"log"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/jdockerty/kindle-notes-grabber/config"
	"github.com/jdockerty/kindle-notes-grabber/notes"
)

const (
	// TODO: Modify to take user argument later on, a flag with a default set if not specified.
	configPath = "kng-config.yaml"
	fromAmazon = "FROM no-reply@amazon.com"
	mailbox    = "INBOX"
)

func main() {

	conf, err := config.New(configPath)
	if err != nil {
		log.Fatalf("Cannot read configuration: %s", err)
	}
	log.Println("Connecting to server...")

	// Connect to server
	// TODO: Implement other providers as mapping format, e.g. gmail : imap.gmail.com:993
	c, err := client.DialTLS("imap.gmail.com:993", nil)
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
	log.Println("Done!")
}
