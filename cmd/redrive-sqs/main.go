package main

import (
  "github.com/urfave/cli/v2"
  "github.com/yyamanoi1222/redrive-sqs/internal/runner"
  "os"
  "fmt"
)

func main() {
  cli.AppHelpTemplate = `NAME:
   {{.Name}}
USAGE:
   {{.HelpName}} {{if .Commands}} command [command options]{{end}}
   {{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
COMMANDS:
{{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .Copyright }}
COPYRIGHT:
   {{.Copyright}}
   {{end}}{{if .Version}}
VERSION:
   {{.Version}}
   {{end}}
`
  app := cli.App{
    Name: "redrive-sqs",
    Version: "v0.0.0",
    Action: func(cCtx *cli.Context) error {
      return runner.Redrive(
        cCtx.String("src"),
        cCtx.String("dest"),
      )
    },
    Flags: []cli.Flag{
      &cli.StringFlag{
        Name: "src",
        Usage: "Source Queue URL",
        Required: true,
      },
      &cli.StringFlag{
        Name: "dest",
        Usage: "Destination Queue URL",
        Required: true,
      },
    },
  }
  if err := app.Run(os.Args); err != nil {
    fmt.Printf("Error %s", err)
  }
}
