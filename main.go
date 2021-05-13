package main

import (
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Email    string `yaml:"email" env:"KNG_EMAIL"`
	Password string `yaml:"password" env:"KNG_PASSWORD"`
}

func readConfig() *Config {

	var cfg Config

	err := cleanenv.ReadConfig("kng-config.yml", &cfg)
	if err != nil {
		log.Fatalf("Configuration is not set")
	}

	return &cfg
}

func kindleMessageIds(c *client.Client) []uint32 {
	_, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}

	criteria := imap.NewSearchCriteria()

	// TODO: Implement a customisable time range for when to check for.
	// twoDaysAgo := time.Now().AddDate(0, 0, -2)
	// criteria.SentSince = twoDaysAgo

	// TODO: Look into searching via IMAP, this doesn't seem to work
	// as expected when looking for value in the email subject, will parse
	// subject manually for now.
	// subjSearch := "OR SUBJECT \"Your Kindle Notes\""

	fromAmazon := "FROM no-reply@amazon.com"
	criteria.Body = []string{fromAmazon}

	ids, err := c.Search(criteria)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Got Ids:", ids)
	return ids

}

// func readEmails(c *client.Client) {

// 	ids := kindleMessageIds(c)

// 	kindleNoteMessages := findNotes(c, ids)

// 	for _, email := range kindleNoteMessages {
// 		bodySt, err := imap.ParseBodySectionName("RFC822")
// 		f := email.GetBody(bodySt)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		log.Println(f)
// 	}

// }

// func findNotes(c *client.Client, ids []uint32) []*imap.Message {

// 	var kindleNoteMessages []*imap.Message

// 	if len(ids) > 0 {
// 		log.Println("Parsing emails...")

// 		// Create a set of UIDs for the emails, each email has a specific ID associated with it
// 		seqset := new(imap.SeqSet)
// 		seqset.AddNum(ids...)

// 		messages := make(chan *imap.Message, 1)
// 		done := make(chan error, 1)

// 		var section imap.BodySectionName
// 		items := []imap.FetchItem{section.FetchItem()}
// 		log.Println("fetching")
// 		go func() {
// 			if err := c.Fetch(seqset, items, messages); err != nil {
// 				log.Fatal(err)
// 			}
// 			log.Println("fetching...")
// 		}()

// 		for msg := range messages {

// 			if subj := msg.Envelope.Subject; strings.HasPrefix(subj, "Your Kindle Notes") {
// 				log.Println(subj)
// 				kindleNoteMessages = append(kindleNoteMessages, msg)
// 			}
// 		}

// 		if err := <-done; err != nil {
// 			log.Fatal(err)
// 		}

// 		log.Printf("%d emails gathered", len(kindleNoteMessages))
// 		return kindleNoteMessages
// 	} else {
// 		log.Println("No Kindle Note emails to parse.")
// 		return kindleNoteMessages
// 	}
// }

func main() {

	// var wg sync.WaitGroup

	conf := readConfig()

	log.Println("Connecting to server...")

	// Connect to server
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

	ids := kindleMessageIds(c)

	seqSet := new(imap.SeqSet)
	seqSet.AddNum(ids...)

	// Get the whole message body
	var section imap.BodySectionName
	items := []imap.FetchItem{section.FetchItem()}

	// Bufferred channel for the last 10 messages
	// NOTE: Could make this user configurable in the future?
	messages := make(chan *imap.Message, 10)
	go func() {
		if err := c.Fetch(seqSet, items, messages); err != nil {
			log.Fatal(err)
		}
	}()

	// Loop over the messages from the channel
	for m := range messages {

		messageBody := m.GetBody(&section)

		mailReader, err := mail.CreateReader(messageBody)
		if err != nil {
			log.Fatal(err)
		}

		// If the email has a subject, continue through processing
		// All the Amazon Kindle emails should have this.
		header := mailReader.Header
		if subject, err := header.Subject(); err == nil {

			// Common prefix for subject header in the emails
			// TODO: Maybe a more efficient way to do this?
			// Could start with consolidating the above if statement, 
			// since all Amazon emails have a subject?
			if strings.HasPrefix(subject, "Your Kindle Notes") {
				

				for {

					// Continue reading the parts until reaching EOF
					part, err := mailReader.NextPart()
					if err == io.EOF {
						break
					} else if err != nil {
						log.Fatal(err)
					}

					switch h := part.Header.(type) {

					// This is an attachment
					case *mail.AttachmentHeader:

						if filename, _ := h.Filename(); strings.HasSuffix(filename, ".csv") {
							log.Printf("Got attachment: %v\n", filename)
							
							contentType, params, err := h.ContentType()
							if err != nil {
								log.Fatal(err)
							}

							if contentType != "text/csv" {
								continue
							}
							log.Println(contentType, params)
							rawData, _ := ioutil.ReadAll(part.Body)
							ioutil.WriteFile("test.csv", rawData, 0644)

						}

					}
				}

				mailReader.Close()
			}
		}

	}

	log.Println("Done!")
}
