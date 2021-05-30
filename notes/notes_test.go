package notes_test

import (
	"testing"

	"github.com/emersion/go-imap"
	"github.com/jdockerty/kindle-notes-grabber/notes"
	"github.com/stretchr/testify/assert"
)

type mockClient struct{}

// Search is the mocked IMAP 'Search' function, this ensures that the mockClient satisfies
// the imapClient interface by implementing the correct method. 
func (mc mockClient) Search(search *imap.SearchCriteria) ([]uint32, error) {
	return []uint32{1, 2, 3}, nil
}

// Fetch is the mocked IMAP 'Fetch' function, this ensures that the mockClient satisfies
// the imapClient interface by implementing the correct method. 
func (mc mockClient) Fetch(seqset *imap.SeqSet, items []imap.FetchItem, ch chan *imap.Message) error {
	fakeMsg := &imap.Message{}
	ch <- fakeMsg
	return nil
}

func TestGetEmailIds(t *testing.T) {
	var m mockClient

	n := notes.New()

	// Can pass nil as search criteria as this takes a pointer to imap.SearchCriteria, which nil satisfies.
	ids := n.GetEmailIds(m, nil)

	var uint32Slice []uint32
	assert.IsType(t, uint32Slice, ids)
}

func TestGetAmazonMessages(t *testing.T) {
	assert := assert.New(t)
	var m mockClient

	var section imap.BodySectionName
	n := notes.New()
	fakeIds := []uint32{1, 2, 3}
	msgs := n.GetAmazonMessages(m, fakeIds, section)

	var receiveChannelType <-chan *imap.Message
	assert.IsType(receiveChannelType, msgs)
}

func TestGetNewNotesDefaults(t *testing.T) {
	assert := assert.New(t)
	testNotes := notes.New()

	var blankString string
	var blankNotes notes.Note

	assert.Empty(testNotes.Author, blankString)
	assert.Empty(testNotes.Title, blankString)
	assert.Empty(testNotes.Notes, blankNotes)

}
