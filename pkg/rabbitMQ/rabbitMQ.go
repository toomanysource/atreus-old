package rabbitmq

import (
	"fmt"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

/*
	交换思路是点对点模式，交换机名和队列名相同，
	user1向user2发消息的过程中，将消息发送到名为user1_user2的交换机
	user2获取user1消息时，须声明user1_user2的队列并获取消息
	由于uid唯一，因此不会出现冲突问题
*/

// MQ Client,既可作为生产者，又作为消费者
type Client struct {
	conn    *amqp.Connection // mq connection
	channel *amqp.Channel    // mq channel
	url     string           // mq url
	ack     bool             // 消费完后是否删除消息
}

// 根据MQURL生成客户端
func NewMQClient(url string, ack bool) (*Client, error) {
	// 连接到 RabbitMQ 服务器
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	fmt.Println("... rabbitMQ conn init successfully")

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	fmt.Println("... rabbitMQ channel init successfully")
	return &Client{conn: conn, channel: ch, url: url, ack: ack}, nil
}

// 断开连接
func (client *Client) CleanUp() {
	client.channel.Close()
	client.conn.Close()
}

// 对发送消息进行了封装，隐藏了用户申明exchange的过程，只需关心谁对谁发送消息
func (client *Client) Publish(userId, toUserId int, content string) error {
	/*
		生产者采取的思路是先声明exchange，若该exchange不存在，则生成该exchange，同时需要额外声明同名队列并将其bind到该exchange上
	*/
	exchangeName := strconv.Itoa(userId) + "_" + strconv.Itoa(toUserId)

	// 不存在交换机则声明交换机和队列，并绑定
	if !client.exsitExchange(exchangeName) {
		if err := client.declareExchange(exchangeName); err != nil {
			return err
		}
	}

	err := client.channel.Publish(
		exchangeName,
		"", // Routing key在fanout类型的Exchange中无用
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // 将消息持久化
			ContentType:  "text/plain",
			Body:         []byte(content),
			// 在消息头部携带用户信息
			Headers: amqp.Table{
				// "UserID":    userId,
				// "ToUserID":  toUserId,
				"Timestamp": time.Now().Format(time.RFC3339),
			},
			//Expiration: "15000", // 15s消息过期
		})

	return err
}

// 对接收消息进行了封装，隐藏了用户申明queue的过程，只需关心谁从谁获取消息(userId <---message--- fromUserId,the messages were beginning on prevTime)
func (client *Client) Consume(userId, fromUserId int, prevTime int64) ([]string, error) {
	// 消费者需事先申明queue
	queueName := strconv.Itoa(fromUserId) + "_" + strconv.Itoa(userId)
	q, err := client.channel.QueueDeclarePassive(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	// 消息数为0则不接收消息
	if q.Messages == 0 {
		return nil, nil
	}
	// 解析时间格式
	layout := time.RFC3339
	// 获取消息
	msgs, err := client.channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	messages := make([]string, 0, q.Messages)
	// 获取一定数量的messages（不获取即时消息）
	for i := 0; i < q.Messages; i++ {
		select {
		case <-time.After(time.Millisecond * 200):
			// 超时
			return nil, fmt.Errorf("timeOut in waiting messages")
		case msg := <-msgs:
			// 获取消息的时间戳
			timestamp := msg.Headers["Timestamp"].(string)
			praseTime, err := time.Parse(layout, timestamp)
			if err != nil {
				return nil, fmt.Errorf("parse timestamp err:%s", err)
			}
			if praseTime.Unix() >= prevTime {
				// 将满足时间条件的message并入
				messages = append(messages, string(msg.Body))
			}
			// 确认消息处理完成并删除
			if client.ack {
				msg.Ack(true)
			}
		}
	}

	return messages, nil
}

// 是否存在指定的交换机
func (client *Client) exsitExchange(exchangeName string) bool {
	return client.channel.ExchangeDeclarePassive(
		exchangeName,
		"fanout", // 使用fanout类型的Exchange，广播消息给所有绑定的Queue
		true,
		false,
		false,
		false,
		nil,
	) == nil
}

// 创建交换机和队列，并进行绑定
func (client *Client) declareExchange(exchangeName string) error {
	// 创建交换机和队列
	if err := client.channel.ExchangeDeclare(
		exchangeName,
		"fanout", // 使用fanout类型的Exchange，广播消息给所有绑定的Queue
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	} else {
		fmt.Println("success in declare exchange", exchangeName)
	}
	// fmt.Printf("交换机 %s 不存在，已默认创建\n", exchangeName)
	// 声明新的队列
	if _, err := client.channel.QueueDeclare(
		exchangeName, // 队列名与交换机同名
		true,         // 持久化Queue
		false,        // 非自动删除Queue
		false,        // 非排他Queue
		false,        // 不等待服务器完成
		nil,          // 额外参数
	); err != nil {
		return err
	} else {
		fmt.Println("success in declare queue", exchangeName)
	}

	// 将队列绑定至该交换机
	if err := client.channel.QueueBind(
		exchangeName,
		"",
		exchangeName,
		false,
		nil,
	); err != nil {
		return err
	} else {
		fmt.Println("success in bind queue", exchangeName, "to exchange", exchangeName)
	}

	return nil
}
