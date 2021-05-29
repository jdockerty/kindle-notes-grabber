package notes_test

import (
	"testing"

	"github.com/jdockerty/kindle-notes-grabber/notes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	// "github.com/stretchr/testify/suite"
)

type mockNotes struct {
	N notes.Notes
	mock.Mock
}

type MockNotes interface {
	GetEmailIds() []uint32
}

// WIP: Implement mocked function response for Ids?

func TestGetEmailIds(t *testing.T) {

	// notes := notes.New()
	// ids := mockNotes.GetEmailIds()
	// ids := mockNotes.GetEmailIds()
	ids := MockNotes.GetEmailIds()

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