package quranenc

import "github.com/urfave/cli/v2"

func Command() *cli.Command {
	return &cli.Command{
		Name:   "quranenc",
		Usage:  "download data from QuranEnc.com",
		Flags:  flags,
		Action: cliAction,
	}
}

var flags = []cli.Flag{
	&cli.StringFlag{
		Name:    "dst",
		Aliases: []string{"d"},
		Usage:   "destination directory",
		Value:   ".",
	},
	&cli.BoolFlag{
		Name:    "clear-cache",
		Aliases: []string{"cc"},
		Usage:   "clear download cache",
	},
}
