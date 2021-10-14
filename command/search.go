package command

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/warrenisarobot/githubsearch/github"
)

type SearchResult interface {
	String() (string, error)
	JSON() (string, error)
}

func SearchCommand(cov *CommandOptionValues) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search [text to search]",
		Short: "Search for code",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			SearchCommandRun(cmd, args, cov)
		},
	}
	cmd.Flags().StringVarP(&cov.SearchType, "searchtype", "s", "", "Search type, possible values are 'gopackage'")
	cmd.Flags().StringVarP(&cov.Organization, "organization", "o", "", "Limit search results to organization")
	return cmd
}

func SearchCommandRun(cmd *cobra.Command, args []string, cov *CommandOptionValues) {
	gh := github.NewAPI(cov.Token())
	searchText := strings.Join(args, " ")
	var err error
	var res SearchResult
	switch cov.SearchType {
	case "gopackage":
		res, err = gh.GoSearch(searchText, cov.Organization, cov.MaxConcurrentRequests)
	default:
		res, err = gh.Search([]string{searchText}, cov.Organization, cov.MaxConcurrentRequests)

	}
	if err != nil {
		log.Error().Err(err).Msg("Search failed")
	}

	var out string
	if cov.OutputFormt == "json" {
		out, err = res.JSON()
	} else {
		out, err = res.String()
	}
	if err != nil {
		log.Err(err).Msg("Error displaying return value")
		return
	}

	fmt.Print(out)
}
