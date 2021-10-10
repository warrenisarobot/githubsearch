package command

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/warrenisarobot/githubsearch/github"
)

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
	gh := github.NewAPI(cov.Token)
	searchText := strings.Join(args, " ")
	fmt.Printf("Token: %s\nSearchText: %s\n", cov.Token, searchText)
	var err error
	var res []github.FileMatch
	switch cov.SearchType {
	case "gopackage":
		res, err = gh.GoSearch(searchText, cov.Organization, cov.MaxConcurrentRequests)
	default:
		res, err = gh.Search([]string{searchText}, cov.Organization, cov.MaxConcurrentRequests)

	}
	if err != nil {
		log.Error().Err(err).Msg("Search failed")
	}

	fmt.Print("Results:\n")
	for _, item := range res {
		fmt.Printf("%s\n", item.String())
	}
}
