package write

import (
	"github.com/Shopify/sarama"
)

type KafkaWriter struct {
	cli   sarama.SyncProducer
	topic string
}

func newConfig() *sarama.Config {
	// kafka配置
	config := sarama.NewConfig()
	// 等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 随机的分区类型
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	// 是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Version = sarama.V0_11_0_2
	return config
}

func NewKafkaWriter(addr string, topic string) (*KafkaWriter, error) {
	producer, err := sarama.NewSyncProducer([]string{addr}, newConfig())
	if err != nil {
		return nil, err
	}
	weiter := KafkaWriter{
		cli:   producer,
		topic: topic,
	}
	return &weiter, nil
}

func (w *KafkaWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	msg := &sarama.ProducerMessage{}
	msg.Topic = w.topic
	msg.Value = sarama.ByteEncoder(p)
	_, _, err = w.cli.SendMessage(msg)
	return n, err
}

func (w *KafkaWriter) CLose() error {
	return w.cli.Close()
}
