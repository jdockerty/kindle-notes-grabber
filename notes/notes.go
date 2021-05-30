package notes

import (
	"log"

	"github.com/emersion/go-imap"
	_ "github.com/emersion/go-imap/client"
)

type imapClient interface {
	Search(criteria *imap.SearchCriteria) ([]uint32, error)
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

func New() *Notes {
	return &Notes{}
}
