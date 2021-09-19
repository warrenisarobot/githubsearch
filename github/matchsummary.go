package github

import "strings"

type FileMatch struct {
	Name string
	Path string
	//raw JSON file via api
	URL string
	//human readable version
	HTMLURL     string
	content     string
	lines       []string
	LineMatches []LineMatch
}

type LineMatch struct {
	Row         int
	Col         int
	MatchedText string
}

func NewFileMatchFromCodeSearch(csm CodeSearchMatch) *FileMatch {
	return &FileMatch{
		Name:    csm.Name,
		Path:    csm.Path,
		URL:     csm.URL,
		HTMLURL: csm.HTMLURL,
	}
}

func (fm *FileMatch) AddContent(c string) {
	fm.content = c
	fm.lines = strings.Split(fm.content, "\n")
}

func (fm *FileMatch) Content() string {
	return fm.content
}

func (fm *FileMatch) Lines() []string {
	return fm.lines
}

func (fm *FileMatch) StringInLines(searchText string) []LineMatch {
	res := []LineMatch{}
	for row, line := range fm.Lines() {
		col := strings.Index(line, searchText)
		if col >= 0 {
			res = append(res, LineMatch{Row: row + 1, Col: col + 1, MatchedText: searchText})
		}
	}
	return res
}
