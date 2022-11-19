package runner

import (
  "github.com/yyamanoi1222/redrive-sqs/internal/sqs"
  "fmt"
  "os"
)

func Redrive(src string, dest string) error {
  cl := &sqs.SQS{}
  cl.Init(sqs.Config{})

  fmt.Printf("Start redriving to %s from %s \n", src, dest)

  for {
    msgs, err := cl.ReceiveMessage(src)

    if err != nil {
      return err
    }

    for _, msg := range msgs {
      redrive(cl, src, dest, msg)
    }
    if len(msgs) == 0 {
      break
    }
  }

  fmt.Printf("End redriving to %s from %s \n", dest, src)
  return nil
}

func redrive(sqs *sqs.SQS, src string, dest string, msg sqs.Message) {
  fmt.Printf("redriving... %s \n", *msg.MessageId)
  if err := sqs.SendMessage(dest, msg); err != nil {
    fmt.Fprintf(os.Stderr, "%s \n", err)
    os.Exit(1)
  }

  if err := sqs.DeleteMessage(src, *msg.ReceiptHandle); err != nil {
    fmt.Fprintf(os.Stderr, "%s \n", err)
    os.Exit(1)
  }
}
