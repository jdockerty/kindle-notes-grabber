package notes

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-message/mail"
	"gopkg.in/yaml.v3"
)

// imapClient interface which satisfies the required methods defined in
// emersion/go-imap/client, this enables pluggability when testing as
// the external calls to an email account and their return values can
// be mocked.
type imapClient interface {
	Search(criteria *imap.SearchCriteria) ([]uint32, error)
	Fetch(seqset *imap.SeqSet, items []imap.FetchItem, ch chan *imap.Message) error
}

const (

	// Index positions of the relevant records, these are the column headings in the CSV file.
	typeIndex       int = 0
	locationIndex   int = 1
	starredIndex    int = 2
	annotationIndex int = 3

	programDirectoryName = "kindle-notes"
)

// Note is a struct which contains a singular record about a note or highlight
// from a particular book. The difference between a note and a highlight is that
// a highlight has no annotation.
type Note struct {
	Type       string
	Location   string
	Annotation string
	Starred    bool
}

// Notes are the overarching struct for the program. This encapsulates a slice of 'Note', which
// is the various records of information that a person has jotted down about a book, including
// other metadata pertaining to it.
type Notes struct {
	Title string
	Notes []Note
}

func GetEmailIds(c imapClient, sc *imap.SearchCriteria) []uint32 {

	// TODO: Look into searching via IMAP, this doesn't seem to work
	// as expected when looking for value in the email subject, will parse
	// subject manually for now.
	// subjSearch := "OR SUBJECT \"Your Kindle Notes\""

	// fromAmazon := "FROM no-reply@amazon.com"
	// criteria.Body = []string{fromAmazon}

	ids, err := c.Search(sc)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Got Ids:", ids)
	return ids
}

func (n *Notes) GetAmazonMessage(c imapClient, id uint32, section imap.BodySectionName) <-chan *imap.Message {
	// Create a set of UIDs for the emails, each email has a specific ID associated with it
	seqSet := new(imap.SeqSet)

	// Add the ids of the Amazon messages which can be parsed for Kindle note emails later
	seqSet.AddNum(id)

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

func (n *Notes) GetMailReaders(messages <-chan *imap.Message, section imap.BodySectionName) []*mail.Reader {

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
// dealing with the emailAttachment, which is in a byte array, by
// instead cutting out the irrelevant rows and placing it into a temporary CSV file, we
// can leverage the csv package to handle the heavy lifting for us.
func parseNotes(title string, emailAttachment []byte) []Note {

	var parsedNotes []Note

	// CSV rows are separate by a newline character
	rows := bytes.Split(emailAttachment, []byte("\n"))
	tmpName := fmt.Sprintf("%s-tmp*.csv", title)
	tmpFile, err := ioutil.TempFile(".", tmpName)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Created temporary file for", title)
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
		log.Printf("Row [%d]: Written for %s\n", lineNum, title)
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

		parsedNotes = append(parsedNotes, n)
	}

	return parsedNotes
}

func (n *Notes) Populate(mailReaders []*mail.Reader) {
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
						log.Println("Got file:", filename)

						_, params, err := h.ContentType()
						if err != nil {
							log.Fatal(err)
						}

						bookTitle := strings.TrimSuffix(params["name"], ".csv")

						// Change the title to lower case and replace spaces with dashes for consistency
						adjustedTitle := strings.ReplaceAll(strings.ToLower(bookTitle), " ", "-")

						log.Println("Adjusted title:", adjustedTitle)
						data, _ := ioutil.ReadAll(part.Body)

						myNotes := parseNotes(adjustedTitle, data)
						n.Notes = append(n.Notes, myNotes...)

						n.Title = adjustedTitle
						log.Println("Set title for notebook", n.Title)
						log.Println("Notes in notebook", len(n.Notes))
					}

				}
			}

		}

	}
}

// Write is used to write the Notes struct, for a given book, into a text file.
// This creates a file with the name of <book-title>-notes.txt and writes each
// Note struct into it, separating each entry with a newline.
// TODO: Sort before writing so that notes appear before highlights etc?
func Write(n *Notes) (int, error) {

	log.Printf("Writing notes for %s\n", n.Title)
	var totalBytes int

	// The 'notebook' prefix is automatically added by Amazon to the CSV file, we can just use .txt as an extension
	notesFilename := fmt.Sprintf("%s.txt", n.Title)
	log.Println("Note file:", notesFilename)
	f, err := os.Create(notesFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for _, note := range n.Notes {
		log.Println("Got note:", note)
		noteEntry := fmt.Sprintf("Annotation: %s\nLocation: %s\nType: %s\nStarred: %t\n\n",
			note.Annotation, note.Location, note.Type, note.Starred)

		bytesWritten, err := f.WriteString(noteEntry)
		if err != nil {
			log.Fatal(err)
		}
		totalBytes += bytesWritten
	}

	return totalBytes, nil
}

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Save will write the titles which have been written into a 'completed notebooks' file of
// key-value pairs, this is done in order to avoid parsing the same titles multiple times
// when they have already been processed.
func Save(n []*Notes) error {

	booksCompleted := make(map[string]bool, len(n))

	for _, completedBook := range n {

		// Skip any erronious issues of the title being blank
		if completedBook.Title == "" {
			continue
		}

		booksCompleted[completedBook.Title] = true
	}

	userHomeDirectory, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	saveDirectory := fmt.Sprintf("%s/%s", userHomeDirectory, programDirectoryName)
	saveDirectoryExists, err := exists(saveDirectory)
	if err != nil {
		return err
	}

	// Save the completed notebooks file if the home directory exists,
	// otherwise we need to return an error
	if saveDirectoryExists {
		savePath := fmt.Sprintf("%s/completed-notebooks.yaml", saveDirectory)
		f, err := os.Create(savePath)
		if err != nil {
			return err
		}
		defer f.Close()

		enc := yaml.NewEncoder(f)
		defer enc.Close()

		enc.Encode(booksCompleted)
		log.Println("Written notes to", savePath)

	} else {
		return fmt.Errorf("a 'kindle-notes' directory does not at '%s' to write the completed notebooks save file", userHomeDirectory)
	}

	return nil
}

// New returns a default Notes struct with none of the fields populated, this is
// ready to be used throughout the program.
func New() *Notes {
	return &Notes{}
}
