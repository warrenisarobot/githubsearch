package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/warrenisarobot/githubsearch/command"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Str("foo", "bar").Msg("Hello world")
	cov := &command.CommandOptionValues{}
	cmd := command.Generate(cov)
	cmd.Execute()
}
