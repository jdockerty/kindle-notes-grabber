package notes_test

import (
	"testing"

	"github.com/jdockerty/kindle-notes-grabber/notes"

	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/suite"
)

// func TestGetEmailIds(t *testing.T) {

// 	notes := notes.New()
// }

func TestGetNewNotesDefaults(t *testing.T) {
	assert := assert.New(t)
	testNotes := notes.New()

	var blankString string
	var blankNotes notes.Note
	
	assert.Empty(testNotes.Author, blankString)
	assert.Empty(testNotes.Title, blankString)
	assert.Empty(testNotes.Notes, blankNotes)

	
}