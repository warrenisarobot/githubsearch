package command

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/warrenisarobot/githubsearch/github"
)

func SearchCommand(cov *CommandOptionValues) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search for code",
		Run: func(cmd *cobra.Command, args []string) {
			SearchCommandRun(cmd, args, cov)
		},
	}
	cmd.Flags().StringVarP(&cov.SearchText, "query", "q", "", "Text for search query")
	cmd.Flags().StringVarP(&cov.SearchType, "searchtype", "s", "", "Search type, possible values are 'gopackage'")
	cmd.Flags().StringVarP(&cov.Organization, "organization", "o", "", "Limit search results to organization")
	return cmd
}

func SearchCommandRun(cmd *cobra.Command, args []string, cov *CommandOptionValues) {
	gh := github.NewAPI(cov.Token)
	fmt.Printf("Token: %s\nSearchText: %s\n", cov.Token, cov.SearchText)
	var err error
	switch cov.SearchType {
	case "gopackage":
		_, err = gh.GoSearch(cov.SearchText, cov.Organization, cov.MaxConcurrentRequests)
	default:
		_, err = gh.Search(cov.SearchText, cov.Organization, cov.MaxConcurrentRequests)

	}
	if err != nil {
		log.Error().Err(err).Msg("Search failed")
	}
}
