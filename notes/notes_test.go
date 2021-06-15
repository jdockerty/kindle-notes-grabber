package notes_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/emersion/go-imap"
	"github.com/jdockerty/kindle-notes-grabber/notes"
	"github.com/stretchr/testify/assert"
)

// mockClient is an empty struct used as a fake IMAP client for
// satisfying the respective interface of the 'notes' package.
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

func getFakeNotesData() *notes.Notes {

	fakeNotes := []notes.Note{
		{
			Annotation: "Annotation 1 used in test",
			Location:   "Page 1",
			Starred:    false,
			Type:       "Highlight",
		},
		{
			Annotation: "Annotation 2 used in test",
			Location:   "Page 2",
			Starred:    true,
			Type:       "Highlight",
		},
		{
			Annotation: "Annotation 3 used in test",
			Location:   "Page 50",
			Starred:    false,
			Type:       "Note",
		},
	}

	return &notes.Notes{
		Title: "test-book-title",
		Notes: fakeNotes,
	}

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

	assert.Empty(testNotes.Title, "A newly created Notes struct should contain a blank string for the 'Title'")
	assert.Empty(testNotes.Notes, "A newly created Notes struct should not contain any populated Note structs within the slice")

}

func TestWriteNoteFile(t *testing.T) {
	assert := assert.New(t)

	testNotes := getFakeNotesData()

	filename := fmt.Sprintf("%s-notes.txt", testNotes.Title)

	defer os.Remove(filename)

	i, err := notes.Write(testNotes)
	assert.Nil(err)
	assert.Greater(i, 0, "The number of bytes written should be greater than 0")

}
