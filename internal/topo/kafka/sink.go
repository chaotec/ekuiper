package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lf-edge/ekuiper/internal/conf"
	"github.com/lf-edge/ekuiper/pkg/api"
	kafka "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

type sink struct {
	Username     string
	Password     string
	URL          string
	Topic        string
	Acks         int
	Retries      int
	BatchSize    int
	LingerMs     int
	BufferMemory int
	writer       *kafka.Writer
}

func (s *sink) Configure(props map[string]interface{}) error {
	if i, ok := props["username"]; ok {
		if u, ok := i.(string); ok {
			s.Username = u
		} else {
			return fmt.Errorf("Not valid username %v.", i)
		}
	}

	if i, ok := props["password"]; ok {
		if p, ok := i.(string); ok {
			s.Password = p
		} else {
			return fmt.Errorf("Not valid password %v.", i)
		}
	}

	if i, ok := props["url"]; ok {
		if u, ok := i.(string); ok {
			s.URL = u
		} else {
			return fmt.Errorf("Not valid url %v.", i)
		}
	}

	if i, ok := props["topic"]; ok {
		if t, ok := i.(string); ok {
			s.Topic = t
		} else {
			return fmt.Errorf("Not valid topic %v.", i)
		}
	}

	if i, ok := props["acks"]; ok {
		if a, ok := i.(int); ok {
			s.Acks = a
		} else {
			return fmt.Errorf("Not valid acks %v.", i)
		}
	}

	if i, ok := props["retries"]; ok {
		if r, ok := i.(int); ok {
			s.Retries = r
		} else {
			return fmt.Errorf("Not valid retries %v.", i)
		}
	}

	if i, ok := props["batchSize"]; ok {
		if b, ok := i.(int); ok {
			s.BatchSize = b
		} else {
			return fmt.Errorf("Not valid batchSize %v.", i)
		}
	}

	if i, ok := props["lingerMs"]; ok {
		if l, ok := i.(int); ok {
			s.LingerMs = l
		} else {
			return fmt.Errorf("Not valid lingerMs %v.", i)
		}
	}

	if i, ok := props["bufferMemory"]; ok {
		if l, ok := i.(int); ok {
			s.BufferMemory = l
		} else {
			return fmt.Errorf("Not valid bufferMemory %v.", i)
		}
	}

	conf.Log.Debugf("Initialized with configurations %#v.", s)
	return nil
}

func (s *sink) Open(ctx api.StreamContext) error {
	var w *kafka.Writer
	if s.Username != "" {
		w = &kafka.Writer{
			Addr:         kafka.TCP(s.URL),
			Topic:        s.Topic,
			Balancer:     &kafka.LeastBytes{},
			BatchSize:    s.BatchSize,
			BatchBytes:   int64(s.BufferMemory),
			BatchTimeout: time.Duration(s.LingerMs),
			MaxAttempts:  s.Retries,
			RequiredAcks: kafka.RequiredAcks(s.Acks),
			Transport:    &kafka.Transport{SASL: plain.Mechanism{Username: s.Username, Password: s.Password}},
		}
	} else {
		w = &kafka.Writer{
			Addr:         kafka.TCP(s.URL),
			Topic:        s.Topic,
			Balancer:     &kafka.LeastBytes{},
			BatchSize:    s.BatchSize,
			BatchBytes:   int64(s.BufferMemory),
			BatchTimeout: time.Duration(s.LingerMs),
			MaxAttempts:  s.Retries,
			RequiredAcks: kafka.RequiredAcks(s.Acks),
		}
	}
	s.writer = w
	return nil
}

func (s *sink) Collect(ctx api.StreamContext, data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = s.writer.WriteMessages(context.Background(), kafka.Message{Value: bytes})
	if err != nil {
		return err
	}
	return nil
}

func (s *sink) Close(ctx api.StreamContext) error {
	s.writer.Close()
	return nil
}

func GetSink() *sink {
	return &sink{}
}
