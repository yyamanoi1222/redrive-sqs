package main

import (
  "flag"
  "fmt"
  "os"
)

func main() {
  var (
    src string
    dest string
  )

  flag.StringVar(&src, "s", "", "src dead letter queue url")
  flag.StringVar(&dest, "d", "", "destination queue url")
  flag.Parse()

  if len(src) == 0 || len(dest) == 0 {
    fmt.Fprintf(os.Stderr, "invalid arguments \n")
    os.Exit(1)
  }

  sqs := &SQS{}
  sqs.init(Config{})

  fmt.Printf("Start redriving to %s from %s \n", dest, src)

  for {
    msgs, err := sqs.receiveMessage(src)

    if err != nil {
      fmt.Fprintf(os.Stderr, "%s \n", err)
      os.Exit(1)
    }

    for _, msg := range msgs {
      redrive(sqs, src, dest, msg)
    }
    if len(msgs) == 0 {
      break
    }
  }

  fmt.Printf("End redriving to %s from %s \n", dest, src)
}

func redrive(sqs *SQS, src string, dest string, msg Message) {
  fmt.Printf("redriving... %s \n", *msg.MessageId)
  if err := sqs.sendMessage(dest, msg); err != nil {
    fmt.Fprintf(os.Stderr, "%s \n", err)
    os.Exit(1)
  }

  if err := sqs.deleteMessage(src, *msg.ReceiptHandle); err != nil {
    fmt.Fprintf(os.Stderr, "%s \n", err)
    os.Exit(1)
  }
}
