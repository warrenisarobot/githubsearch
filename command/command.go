package command

import (
	"github.com/spf13/cobra"
)

type CommandOptionValues struct {
	Token                 string
	SearchType            string
	Organization          string
	MaxConcurrentRequests int
}

var searchCmd = &cobra.Command{}

func Generate(cov *CommandOptionValues) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "githubsearch",
		Short: "Search for a thing",
	}

	rootCmd.PersistentFlags().StringVarP(&cov.Token, "token", "t", "", "Github auth token to use for searching")
	rootCmd.PersistentFlags().IntVarP(&cov.MaxConcurrentRequests, "maxrequests", "m", 5, "Github auth token to use for searching")
	rootCmd.AddCommand(SearchCommand(cov))
	return rootCmd
}
