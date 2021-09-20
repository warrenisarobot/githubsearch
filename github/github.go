package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	defaultGithubHost        = "api.github.com"
	githubAPIsearchPath      = "/search/code"
	githubInvalidSearchChars = ".,:;/\\`'\"=*!?#$&+^|~<>(){}[]@"
	githubAPIMaxPages        = 100
)

var (
	goImportAlias     = regexp.MustCompile(`([^ ]*)[ ]*"[^"]*"`)
	goImportPathParts = regexp.MustCompile(`([^\/]*\/[^\.]*)\.(.*)`)
)

type GitHubAPI struct {
	Host   string
	Token  string
	client *http.Client
}

func NewAPI(token string) GitHubAPI {
	return GitHubAPI{
		Host:   defaultGithubHost,
		Token:  token,
		client: &http.Client{},
	}
}

func (g *GitHubAPI) Search(searchText []string, organization string, maxRequests int) ([]FileMatch, error) {
	page := 1
	per_page := githubAPIMaxPages
	u := g.url(githubAPIsearchPath)

	matchesToReturn := []FileMatch{}

	for {
		u.RawQuery = g.searchQuery(searchText, organization, page, per_page)

		log.Info().Str("url", u.String()).Interface("searchText", searchText).Str("queryParams", u.RawQuery).Msg("Creating search request")

		req, err := g.newRequest("GET", u.String())
		if err != nil {
			return nil, err
		}

		log.Info().Int("page", page).Msg("Getting search results")
		resp, err := g.client.Do(req)
		if err != nil {
			return nil, err
		}

		matches := &CodeSearchResults{}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("Could not read searh results: %w", err)
		}
		err = json.Unmarshal(data, matches)
		if err != nil {
			return nil, fmt.Errorf("Could not unmarshall search results: %w", err)
		}
		log.Debug().Int("results", len(matches.Items)).Int("totalResults", matches.TotalCount).Bool("incomplete", matches.IncompleteResults).Msg("Parsed search results")
		newres, err := g.searchResultsMatchText(searchText, matches.Items)
		if err != nil {
			return nil, fmt.Errorf("Failed matching search result text: %w", err)
		}
		matchesToReturn = append(matchesToReturn, newres...)
		if matches.TotalCount == 0 || ((page-1)*per_page)+len(matches.Items) >= matches.TotalCount {
			break
		}
		page += 1
	}
	return matchesToReturn, nil
}

//Search the github API for usage of an exported package resource.  This expects the searchText to be in the
//format of <ImportPath>.<Resource>, like:
//
//github.com/project/package/subpack.New
//
//This would find places that import the github.com/project/pacakge/subpack package, and have instances
//of the text subpack.New (adjusted if an alias is used in the import)
func (g *GitHubAPI) GoSearch(searchText, organization string, maxRequests int) ([]FileMatch, error) {
	matches := goImportPathParts.FindStringSubmatch(searchText)
	importPath := ""
	resource := ""
	if len(matches) > 1 {
		importPath = matches[1]
		resource = matches[2]
	}

	searchRes, err := g.Search([]string{importPath, resource}, organization, maxRequests)
	if err != nil {
		return nil, err
	}

	res := []FileMatch{}

	for _, item := range searchRes {
		alias := g.getImportAlias(string(item.Content()), importPath)
		searchMe := alias + "." + resource
		lineMatches := item.StringInLines(searchMe)
		if len(lineMatches) > 0 {
			item.LineMatches = lineMatches
			res = append(res, item)
		}
	}
	return res, nil
}

//check the given file context for the importPath and return if an alias is used.  Empty string means no alias
func (g *GitHubAPI) getImportAlias(fileContent, importPath string) string {
	for _, line := range strings.Split(string(fileContent), "\n") {
		if strings.Contains(line, importPath) {
			res := goImportAlias.FindStringSubmatch(line)
			if len(res) > 1 {
				alias := strings.TrimSpace(res[1])
				//there is text in front of the quoted import path
				if alias != "import" && alias != "" {
					return alias
				}
			}
			if len(res) > 0 {
				//only the import path matches, there is no alias
				parts := strings.Split(res[0], "/")
				//get the last part of the path, then remove the quote char
				lastPart := parts[len(parts)-1]
				if len(lastPart) < 2 {
					return ""
				}
				name := lastPart[:len(lastPart)-1]
				return name
			}
		}
	}
	return ""
}

func (g *GitHubAPI) getFile(url string) (*FileResults, error) {
	log.Info().Str("url", url).Msg("Downloading file from github")
	req, err := g.newRequest("GET", url)
	if err != nil {
		return nil, fmt.Errorf("Request creation failed: %w", err)
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Request failed: %w", err)
	}

	fr := &FileResults{}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to download body: %w", err)
	}
	err = json.Unmarshal(data, fr)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal body: %w", err)
	}
	return fr, nil

}

func (g *GitHubAPI) authHeader() (string, string) {
	return "Authorization", fmt.Sprintf("token %s", g.Token)
}

func (g *GitHubAPI) textMatchHeader() (string, string) {
	return "Accept", "application/vnd.github.v3.text-match+json"
}

func (g *GitHubAPI) url(path string) url.URL {
	return url.URL{
		Host:   g.Host,
		Scheme: "https",
		Path:   path,
	}
}

func (g *GitHubAPI) newRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add(g.authHeader())
	req.Header.Add(g.textMatchHeader())
	return req, nil
}

func (g *GitHubAPI) searchParts(org string) string {
	parts := []string{}
	if org != "" {
		parts = append(parts, fmt.Sprintf("org:%s", org))
	}
	return strings.Join(parts, " ")
}

func (g *GitHubAPI) searchQuery(searchText []string, organization string, page, per_page int) string {
	replacedString := strings.Join(searchText, " ")
	for _, ch := range githubInvalidSearchChars {
		replacedString = strings.ReplaceAll(replacedString, string(ch), " ")
	}
	v := url.Values{}
	v.Add("q", fmt.Sprintf("%s %s", replacedString, g.searchParts(organization)))
	v.Add("per_page", fmt.Sprintf("%d", per_page))
	v.Add("page", fmt.Sprintf("%d", page))
	return v.Encode()
}

//return a list of FileMatches that have all at least 1 exact match of the given searchText terms
func (g *GitHubAPI) searchResultsMatchText(searchText []string, results []CodeSearchMatch) ([]FileMatch, error) {
	res := []FileMatch{}
	for _, item := range results {
		localMatch := NewFileMatchFromCodeSearch(item)
		f, err := g.getFile(item.URL)
		if err != nil {
			return nil, fmt.Errorf("Could not get github file: %w", err)
		}
		content, err := f.DecodedContent()
		if err != nil {
			return nil, fmt.Errorf("Failed to decode content: %w", err)
		}
		localMatch.AddContent(string(content))
		found := true
		for _, searchTerm := range searchText {
			matches := localMatch.StringInLines(searchTerm)
			if len(matches) > 0 {
				localMatch.LineMatches = append(localMatch.LineMatches, matches...)
			} else {
				found = false
			}
		}
		if found {
			res = append(res, *localMatch)
		}
	}
	return res, nil
}
