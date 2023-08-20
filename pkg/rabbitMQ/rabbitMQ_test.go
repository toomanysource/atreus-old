package rabbitmq

import (
	"fmt"
	"log"
	"testing"
	"time"
)

const (
	endpoint = "amqp://guest:guest@localhost:5672/"
)

func _assert(condition bool, content string, a ...any) {
	if condition {
		log.Fatalf(content, a...)

	}
}

func TestPublish(t *testing.T) {
	client, err := NewMQClient(endpoint, true)
	_assert(err != nil, "MQ client init failed err %s", err)
	defer client.CleanUp()
	for i := 1; i < 5; i++ {
		err = client.Publish(101, 102, fmt.Sprintf("hi 102 seeId %d", i))
		_assert(err != nil, "publish failed err:%s", err)
	}
}

func TestConsume(t *testing.T) {
	client, err := NewMQClient(endpoint, true) // 消费完后删除消息
	_assert(err != nil, "MQ client init failed err %s", err)
	defer client.CleanUp()
	var msgs []string
	msgs, err = client.Consume(102, 101, time.Now().AddDate(0, 0, -1).Unix()) // 获取昨天之后的消息
	_assert(err != nil, "consume failed err %s", err)
	fmt.Printf("%d get msgs from %d: %v", 102, 101, msgs)
}

func TestDelete(t *testing.T) {
	client, err := NewMQClient(endpoint, true) // 消费完后删除消息
	_assert(err != nil, "MQ client init failed err %s", err)
	defer client.CleanUp()
	exchangeName := "101_102"
	queueName := exchangeName
	// 删除exchange
	err = client.channel.ExchangeDelete(exchangeName, false, false)
	if err != nil {
		log.Fatalf("failed to delete exchange %s", exchangeName)
	} else {
		fmt.Println("success in delete exchange", exchangeName)
	}

	// 删除queue
	_, err = client.channel.QueueDelete(exchangeName, false, false, false)
	if err != nil {
		log.Fatalf("failed to delete queue %s", queueName)
	} else {
		fmt.Println("success in delete queue", queueName)
	}
}
