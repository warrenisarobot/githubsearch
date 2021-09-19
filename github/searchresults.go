package github

type CodeSearchResults struct {
	TotalCount        int               `json:"total_count"`
	IncompleteResults bool              `json:"incomplete_results"`
	Items             []CodeSearchMatch `json:"items"`
}

type CodeSearchMatch struct {
	Name        string              `json:"name"`
	Path        string              `json:"path"`
	SHA         string              `json:"sha"`
	URL         string              `json:"url"`
	GitURL      string              `json:"git_url"`
	HTMLURL     string              `json:"html_url"`
	Score       float32             `json:"score"`
	Repository  *RepositoryInfo     `json:"repository"`
	TextMatches []CodeTextMatchInfo `json:"text_matches"`
}

type RepositoryInfo struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	HTMLURL  string `json:"html_url"`
	URL      string `json:"url"`
}

type CodeTextMatchInfo struct {
	ObjectURL  string            `json:"object_url"`
	ObjectType string            `json:"object_type"`
	Fragement  string            `json:"fragment"`
	Matches    []CodeSearchMatch `json:"matches"`
}

type CodeTextMatch struct {
	Text    string `json:"text"`
	Indices []int  `json:"indices"`
}

func (ctm *CodeTextMatch) RowNum() int {
	if ctm == nil || len(ctm.Indices) != 2 {
		return 0
	}
	return ctm.Indices[0]
}

func (ctm *CodeTextMatch) ColNum() int {
	if ctm == nil || len(ctm.Indices) != 2 {
		return 0
	}
	return ctm.Indices[1]
}
