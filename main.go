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

func amazonEmailIds(c *client.Client) []uint32 {

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


func getAmazonMessages(c *client.Client, ids []uint32, section imap.BodySectionName) <-chan *imap.Message {
	// Create a set of UIDs for the emails, each email has a specific ID associated with it
	seqSet := new(imap.SeqSet)

	// Add the ids of the Amazon messages which can be parsed for Kindle note emails later
	seqSet.AddNum(ids...)

	// Get the whole message body
	items := []imap.FetchItem{section.FetchItem()}

	// Bufferred channel for the last 10 messages
	// NOTE: Could make this user configurable in the future?
	messages := make(chan *imap.Message, 10)

	// Run separate goroutine for fetching messages, these are
	// passed back over the channel defined above
	go func() {
		if err := c.Fetch(seqSet, items, messages); err != nil {
			log.Fatal(err)
		}
	}()

	return messages
}

func main() {

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

	_, err = c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}

	ids := amazonEmailIds(c)

	var section imap.BodySectionName
	messages := getAmazonMessages(c, ids, section)

	// Loop over the messages from the channel
	for m := range messages {

		messageBody := m.GetBody(&section)

		mailReader, err := mail.CreateReader(messageBody)
		if err != nil {
			log.Println("Using unknown charset for reading mail header.")
		}

		// If the email has a subject, continue through processing
		// All the Amazon Kindle emails should have this.
		header := mailReader.Header

		subject, err := header.Subject()

		if err == nil && strings.HasPrefix(subject, "Your Kindle Notes") {

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

	log.Println("Done!")
}
