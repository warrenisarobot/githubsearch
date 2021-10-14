package command

import (
	"os"

	"github.com/spf13/cobra"
)

type CommandOptionValues struct {
	TokenOverride         string
	SearchType            string
	Organization          string
	MaxConcurrentRequests int
	OutputFormt           string
}

func (cov *CommandOptionValues) Token() string {
	if cov.TokenOverride != "" {
		return cov.TokenOverride
	}
	val := os.Getenv("GH_TOKEN")
	if val != "" {
		return val
	}
	val = os.Getenv("GITHUB_TOKEN")
	return val
}

var searchCmd = &cobra.Command{}

func Generate(cov *CommandOptionValues) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "githubsearch",
		Short: "Search for a thing",
	}

	rootCmd.PersistentFlags().StringVarP(&cov.TokenOverride, "token", "t", "", "Github auth token to use for searching.  If this is not provided the env vars GH_TOKEN, or GITHUB_TOKEN will be used if available")
	rootCmd.PersistentFlags().IntVarP(&cov.MaxConcurrentRequests, "maxrequests", "m", 5, "Github auth token to use for searching")
	rootCmd.PersistentFlags().StringVarP(&cov.OutputFormt, "format", "f", "text", "output format.  Can be either 'text' or 'json'")
	rootCmd.AddCommand(SearchCommand(cov))
	return rootCmd
}
