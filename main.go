package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/warrenisarobot/githubsearch/command"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cov := &command.CommandOptionValues{}
	cmd := command.Generate(cov)
	cmd.Execute()
}
