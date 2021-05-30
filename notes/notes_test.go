package notes_test

import (
	"testing"

	"github.com/emersion/go-imap"
	"github.com/jdockerty/kindle-notes-grabber/notes"
	"github.com/stretchr/testify/assert"
)

var SearchMock func() ([]uint32, error)

type mockClient struct{}

func (mc mockClient) Search(search *imap.SearchCriteria) ([]uint32, error) {
	return SearchMock()
}

func TestGetEmailIds(t *testing.T) {
	var m mockClient

	n := notes.New()

	SearchMock = func() ([]uint32, error) {
		return []uint32{1, 2, 3}, nil
	}

	// Can pass nil as search criteria as this takes a pointer to imap.SearchCriteria
	// which nil satisfies.
	ids := n.GetEmailIds(m, nil)

	var uint32Slice []uint32
	assert.IsType(t, uint32Slice, ids)
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
