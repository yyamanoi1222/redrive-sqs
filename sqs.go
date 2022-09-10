package main

import (
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/sqs"
  "regexp"
)

type SQS struct {
  cl *sqs.SQS
}

type Message sqs.Message

type Config struct {
  Endpoint string
}

const MAX_MESSAGE_RECEIVE_COUNT = 10

func (s *SQS) init(cfg Config) {
  sess := session.Must(session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable}))
  awsCfg := aws.Config{}

  if len(cfg.Endpoint) > 0 {
    awsCfg.Endpoint = aws.String(cfg.Endpoint)
  }
  s.cl = sqs.New(sess, &awsCfg)
}

func (s *SQS) receiveMessage(queueUrl string) ([]Message, error) {
  res, err := s.cl.ReceiveMessage(&sqs.ReceiveMessageInput{
    QueueUrl: aws.String(queueUrl),
    AttributeNames: []*string{
      aws.String("All"),
    },
    MessageAttributeNames: []*string{
      aws.String("All"),
    },
    MaxNumberOfMessages: aws.Int64(MAX_MESSAGE_RECEIVE_COUNT),
  })

  var msgs []Message

  for _, m := range res.Messages {
    msg := Message{
      MessageAttributes: m.MessageAttributes,
      Body: m.Body,
      Attributes: m.Attributes,
      ReceiptHandle:  m.ReceiptHandle,
      MessageId: m.MessageId,
    }
    msgs = append(msgs, msg)
  }

  return msgs, err
}

func (s *SQS) sendMessage(queueUrl string, msg Message) error {
  input := sqs.SendMessageInput{
    MessageAttributes: msg.MessageAttributes,
    MessageBody: msg.Body,
    QueueUrl: aws.String(queueUrl),
  }

  if isFifo(queueUrl) {
    input.MessageDeduplicationId = msg.Attributes[sqs.MessageSystemAttributeNameMessageDeduplicationId]
    input.MessageGroupId = msg.Attributes[sqs.MessageSystemAttributeNameMessageGroupId]
  }

  _, err := s.cl.SendMessage(&input)
  return err
}

func (s *SQS) deleteMessage(queueUrl string, receiptHandle string) error {
  _, err := s.cl.DeleteMessage(&sqs.DeleteMessageInput{
    QueueUrl: aws.String(queueUrl),
    ReceiptHandle: aws.String(receiptHandle),
  })
  return err
}

func isFifo(queueUrl string) bool {
  r := regexp.MustCompile(`\.fifo$`)
  return r.MatchString(queueUrl)
}
