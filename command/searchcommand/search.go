package searchcommand

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/warrenisarobot/githubsearch/command"
)

func Generate(cov *command.CommandOptionValues) *cobra.Command {
	return &cobra.Command{
		Use:   "search",
		Short: "Search for code",
		Run: func(cmd *cobra.Command, args []string) {
			SearchCmd(cmd, args, cov)
		},
	}
}

func SearchCmd(cmd *cobra.Command, args []string, cov *command.CommandOptionValues) {
	fmt.Printf("Token: %s", cov.Token)
}
