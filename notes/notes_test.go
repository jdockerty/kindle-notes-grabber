package notes_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/emersion/go-imap"
	"github.com/jdockerty/kindle-notes-grabber/notes"
	"github.com/stretchr/testify/assert"
)

var fakeNotes *notes.Notes = getFakeNotesData()
var m mockNotes

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

type mockNotes struct{}

// Save is the 
func (mn *mockNotes) Save(n []*notes.Notes) error {
	return nil
}

func getFakeNotesData() *notes.Notes {

	fakeNotes := []notes.Note{
		{
			Annotation: "Text that was highlighted 1",
			Location:   "Page 1",
			Starred:    false,
			Type:       "Highlight",
		},
		{
			Annotation: "Text that was highlighted 2",
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

	// Can pass nil as search criteria as this takes a pointer to imap.SearchCriteria, which nil satisfies.
	ids := notes.GetEmailIds(m, nil)

	var uint32Slice []uint32
	assert.IsType(t, uint32Slice, ids)
}

func TestGetAmazonMessage(t *testing.T) {
	assert := assert.New(t)
	var m mockClient

	var section imap.BodySectionName
	n := notes.New()

	var fakeId uint32 = 1
	msgs := n.GetAmazonMessage(m, fakeId, section)

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

	filename := fmt.Sprintf("%s.txt", fakeNotes.Title)

	defer os.Remove(filename)

	i, err := notes.Write(fakeNotes)
	assert.Nil(err)
	assert.Greater(i, 0, "The number of bytes written should be greater than 0")

}

func TestShouldReadAllNotes(t *testing.T) {
	numNotes := len(fakeNotes.Notes)
	assert.Equal(t, 3, numNotes)
}

func TestShouldWriteTitleToSaveFile(t *testing.T) {
	assert := assert.New(t)

	var multipleFakeNotes []*notes.Notes

	// Append duplicated notes to act as 2 items within the slice
	multipleFakeNotes = append(multipleFakeNotes, fakeNotes, fakeNotes)
	err := m.Save(multipleFakeNotes)
	assert.Nil(err)
}
