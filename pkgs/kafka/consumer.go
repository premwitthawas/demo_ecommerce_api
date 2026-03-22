package pkgs_kafka

import (
	"strings"

	"github.com/segmentio/kafka-go"
)

func NewKafkaConsumer(urlsWithComma string, topic string, groupID string) *kafka.Reader {
	urls := strings.Split(urlsWithComma, ",")
	topicGroups := strings.Split(topic, ",")
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     urls,
		GroupTopics: topicGroups,
		GroupID:     groupID,
	})
	return reader
}
