package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"encoding/csv"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/jdockerty/kindle-notes-grabber/config"
)

type Note struct {
	Type       string
	Location   string
	Annotation string
	Starred    bool
}

type Notes struct {
	Author string
	Title  string
	Notes  []Note
}

const (

	// Index positions of the relevant records, these are the column headings in the CSV file.
	typeIndex       int = 0
	locationIndex   int = 1
	starredIndex    int = 2
	annotationIndex int = 3
)

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

func getMailReaders(messages <-chan *imap.Message, section imap.BodySectionName) []*mail.Reader {

	var mailReaders []*mail.Reader
	// Loop over the messages from the channel
	for m := range messages {

		messageBody := m.GetBody(&section)

		mailReader, err := mail.CreateReader(messageBody)
		if err != nil {
			log.Println("Using unknown charset for reading mail header.")
		}

		mailReaders = append(mailReaders, mailReader)

	}

	return mailReaders
}

// parseNotes is used to create a temporary CSV file which can be
// read from. This is done as it provides a simpler mechanicm than directly
// dealing with the emailAttachment directly, which is in a byte array, by
// instead cutting out the irrelevant rows and placing it into a temporary CSV file, we
// can leverage the csv package to handle the heavy lifting for us.
func parseNotes(title string, emailAttachment []byte) []Note {

	var notes []Note

	// CSV rows are separate by a newline character
	rows := bytes.Split(emailAttachment, []byte("\n"))
	tmpName := fmt.Sprintf("%s-tmp*.csv", title)
	tmpFile, err := ioutil.TempFile(".", tmpName)
	if err != nil {
		log.Fatal(err)
	}

	// Close and remove the temporary file once completed.
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	for lineNum, row := range rows {

		// First 8 lines (counting from 0) are generated by Amazon and aren't useful to us,
		// such as a book preview link etc.
		if lineNum < 7 {
			continue
		}

		// Write to tmp file
		// Parse CSV into struct
		// Append to Notes type

		tmpFile.Write(row)
		if err != nil {
			log.Fatal("Writing:", err)
		}

		// Re-add the newline for the next row
		// TODO Better way to do this, maybe around the splitting via \n is an issue?
		tmpFile.WriteString("\n")
	}

	// Move to the beginning of the file, as we've recently written to it
	// and thus moved the offset.
	tmpFile.Seek(0, 0)

	csvFile := csv.NewReader(tmpFile)
	csvFile.FieldsPerRecord = 4

	records, _ := csvFile.ReadAll()

	for index, record := range records {

		// Skip top record as these are the field labels
		if index == 0 {
			continue
		}

		var n Note
		n.Type = record[typeIndex]
		n.Location = record[locationIndex]
		n.Annotation = record[annotationIndex]

		n.Starred = false
		starred := record[starredIndex]
		if starred == "*" {
			n.Starred = true
		}

		notes = append(notes, n)
	}

	return notes
}

func main() {

	// TODO: Modify to take user argument later on, a flag with a default set if not specified.
	configPath := "kng-config.yaml"
	conf, err := config.New(configPath)
	if err != nil {
		log.Fatalf("Cannot read configuration: %s", err)
	}
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

	mailReaders := getMailReaders(messages, section)

	for _, mailReader := range mailReaders {
		header := mailReader.Header
		subject, err := header.Subject()

		if err == nil && strings.HasPrefix(subject, "Your Kindle Notes") {
			for {

				// Continue reading the parts until reaching EOF
				part, err := mailReader.NextPart()
				if err == io.EOF {
					break
				} else if err != nil {
					log.Println("Unknown charset in use, but usable part provided")
				}

				switch h := part.Header.(type) {

				// This is an attachment
				case *mail.AttachmentHeader:

					if filename, _ := h.Filename(); strings.HasSuffix(filename, ".csv") {
						log.Printf("Got attachment: %v\n", filename)

						_, params, err := h.ContentType()
						if err != nil {
							log.Fatal(err)
						}

						bookTitle := strings.TrimSuffix(params["name"], ".csv")

						log.Println(bookTitle)
						data, _ := ioutil.ReadAll(part.Body)
						parseNotes(bookTitle, data)
					}

				}
			}

		}

	}
	log.Println("Done!")
}
