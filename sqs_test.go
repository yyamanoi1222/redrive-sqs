package main

import (
  "testing"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/sqs"
  "crypto/rand"
  "math/big"
)

const endpoint = "http://localhost:9324"
const queueUrl = endpoint + "/000000000000/test"

func prepareInstance() *SQS {
  s := &SQS{}
  s.init(Config{Endpoint: endpoint})
  return s
}

func createDummyMessage() Message {
  n, err := rand.Int(rand.Reader, big.NewInt(100))
  if err != nil {
    panic(err)
  }
  return Message{
    Body: aws.String("test"),
    Attributes: map[string]*string{
      "MessageGroupId": aws.String(n.String()),
      "MessageDeduplicationId": aws.String(n.String()),
    },
  }
}

func TestMain(m *testing.M) {
  s := prepareInstance()
  s.cl.CreateQueue(&sqs.CreateQueueInput{
    QueueName: aws.String("test"),
  })

  m.Run()
}

func TestReceiveMessage(t *testing.T) {
  s := prepareInstance()
  msg := createDummyMessage()
  s.sendMessage(queueUrl, msg)
  msgs, err := s.receiveMessage(queueUrl)

  if err != nil {
    t.Errorf("err %s", err)
  }

  if len(msgs) == 0 {
    t.Errorf("err %s", "message not found")
  }
}

func TestSendMessage(t *testing.T) {
  s := prepareInstance()
  msg := createDummyMessage()
  if err := s.sendMessage(queueUrl, msg); err != nil {
    t.Errorf("err %s", err)
  }
}

func TestDeleteMessage(t *testing.T) {
  s := prepareInstance()
  msg := createDummyMessage()
  s.sendMessage(queueUrl, msg)
  msgs, _ := s.receiveMessage(queueUrl)

  if err := s.deleteMessage(queueUrl, *msgs[0].ReceiptHandle); err != nil {
    t.Errorf("err %s", err)
  }
}
