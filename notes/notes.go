package notes

import (
	"log"

	"github.com/emersion/go-imap"
	_ "github.com/emersion/go-imap/client"
)

type imapClient interface {
	Search(criteria *imap.SearchCriteria) ([]uint32, error)
	Fetch(seqset *imap.SeqSet, items []imap.FetchItem, ch chan *imap.Message) error
}

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

func (n *Notes) GetEmailIds(c imapClient, sc *imap.SearchCriteria) []uint32 {
	// criteria := imap.NewSearchCriteria()

	// TODO: Implement a customisable time range for when to check for.
	// twoDaysAgo := time.Now().AddDate(0, 0, -2)
	// criteria.SentSince = twoDaysAgo

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

func (n *Notes) GetAmazonMessages(c imapClient, ids []uint32, section imap.BodySectionName) <-chan *imap.Message {
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

func New() *Notes {
	return &Notes{}
}
