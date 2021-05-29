package notes




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

func New() *Notes {
	return &Notes{}
}

