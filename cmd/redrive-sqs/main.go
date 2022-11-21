package main

import (
  "github.com/yyamanoi1222/redrive-sqs/internal/runner"
  "flag"
  "os"
  "fmt"
)

func main() {
  var (
    src string
    dest string
  )

  flag.StringVar(&src, "src", "", "src dead letter queue url")
  flag.StringVar(&dest, "dest", "", "destination queue url")
  flag.Parse()


  if len(src) == 0 || len(dest) == 0 {
    fmt.Fprintf(os.Stderr, "invalid arguments \n")
    os.Exit(1)
  }

  if err := runner.Redrive(src, dest); err != nil {
    fmt.Fprintf(os.Stderr, "%s \n", err)
  }
}
