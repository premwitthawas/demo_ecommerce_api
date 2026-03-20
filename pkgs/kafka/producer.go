package pkgs_kafka

import (
	"strings"

	"github.com/segmentio/kafka-go"
)

func NewKafkaProducer(urlsWithComma string) (*kafka.Writer, error) {
	urls := strings.Split(urlsWithComma, ",")
	writer := &kafka.Writer{
		Addr:         kafka.TCP(urls...),
		Balancer:     &kafka.LeastBytes{},
		Async:        false,
		RequiredAcks: kafka.RequireAll,
	}
	return writer, nil
}
