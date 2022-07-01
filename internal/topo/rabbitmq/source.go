package rabbitmq

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	"github.com/lf-edge/ekuiper/internal/conf"
	"github.com/lf-edge/ekuiper/pkg/api"
	"github.com/mitchellh/mapstructure"
	"github.com/streadway/amqp"
)

var once sync.Once

type source struct {
	Username     string
	Password     string
	URL          string
	Exchange     string
	ExchangeType string
	RoutingKey   string
	conn         *amqp.Connection
	channel      *amqp.Channel
	msgs         <-chan amqp.Delivery
}

func (s *source) Configure(_ string, props map[string]interface{}) error {
	conf.Log.Infof("configurations map %#v.", props)
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
			return fmt.Errorf("Not valid addr %v.", i)
		}
	}

	if i, ok := props["exchange"]; ok {
		if e, ok := i.(string); ok {
			s.Exchange = e
		} else {
			return fmt.Errorf("Not valid exchange %v.", i)
		}
	}

	if i, ok := props["exchangetype"]; ok {
		if e, ok := i.(string); ok {
			s.ExchangeType = e
		} else {
			return fmt.Errorf("Not valid exchangeType %v.", i)
		}
	}

	if i, ok := props["routingkey"]; ok {
		if r, ok := i.(string); ok {
			s.RoutingKey = r
		} else {
			return fmt.Errorf("Not valid routingKey %v.", i)
		}
	}

	conf.Log.Infof("Initialized with configurations %#v.", s)
	return nil
}

func (s *source) Open(ctx api.StreamContext, consumer chan<- api.SourceTuple, errCh chan<- error) {
	logger := ctx.GetLogger()
	ruleId := ctx.GetRuleId()

	logger.Infof("Connet to rabbitmq.")
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/", s.Username, s.Password, s.URL))
	if err != nil {
		logger.Infof("Failed to connet to rabbitmq.")
	}
	s.conn = conn

	logger.Infof("Declare a channel.")
	channel, err := s.conn.Channel()
	if err != nil {
		logger.Infof("Failed to declare a channel.")
		return
	}
	s.channel = channel

	logger.Infof("Declare a exchange.")
	if err := s.channel.ExchangeDeclarePassive(
		s.Exchange,
		s.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		logger.Infof("Failed to declare a exchange.")
		return
	}

	logger.Infof("Declare a queue.")
	queue, err := s.channel.QueueDeclare(
		ruleId,
		false,
		true,
		true,
		false,
		nil,
	)
	if err != nil {
		logger.Infof("Failed to declare a queue.")
		return
	}

	logger.Infof("Bind a queue.")
	if err = s.channel.QueueBind(queue.Name, string(s.RoutingKey), s.Exchange, false, nil); err != nil {
		logger.Infof("Failed bind a queue.")
		return
	}

	logger.Infof("Declare a consumer.")
	msgs, err := s.channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Infof("Failed to declare a consumer.")
		return
	}
	s.msgs = msgs

	err = s.Connect(ctx, consumer)
	if err != nil {
		errCh <- err
		return
	}
	return
}

func (s *source) Connect(ctx api.StreamContext, consumer chan<- api.SourceTuple) error {
	logger := ctx.GetLogger()
	// Send data to data channel
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-s.msgs:
			var (
				body         RabbitMQBody
				deviceEvents []DeviceEvent
			)
			logger.Infof("body = %s\n", msg.Body)
			body.Body = []DeviceEvent{}
			if err := json.Unmarshal(msg.Body, &body); err != nil {
				continue
			}
			logger.Infof("after body = %s\n", body)
			if err := mapstructure.Decode(body.Body, &deviceEvents); err != nil {
				continue
			}
			logger.Infof("device data = %s\n", deviceEvents)
			switch RabbitMQRoutingKey(s.RoutingKey) {
			case RabbitMQRoutingKeyDeviceData:
				logger.Infoln("In RabbitMQRoutingKeyDeviceData")
				for _, d := range deviceEvents {
					result := make(map[string]interface{})
					meta := make(map[string]interface{})
					logger.Infof("d = %s", d)
					for _, p := range d.Properties {
						logger.Infof("p = %s", p)
						switch p.ValueType {
						case "int":
							v, err := strconv.ParseInt(p.ReportValue, 10, 64)
							if err != nil {
								logger.Errorf("Failed to convert to int type.")
								continue
							}
							result[p.Name] = v
						case "float":
							v, err := strconv.ParseFloat(p.ReportValue, 64)
							if err != nil {
								logger.Errorf("Failed to convert to float type.")
								continue
							}
							result[p.Name] = v
						case "string":
							result[p.Name] = p.ReportValue
						case "bool":
							v, err := strconv.ParseBool(p.ReportValue)
							if err != nil {
								logger.Errorf("Failed to convert to bool type.")
								continue
							}
							result[p.Name] = v
						}
					}
					meta["deviceName"] = d.DeviceName
					meta["deviceModelName"] = d.DeviceModelName
					logger.Infof("result = %s", result)
					logger.Infof("meta = %s", meta)
					consumer <- api.NewDefaultSourceTuple(result, meta)
				}
			case RabbitMQRoutingKeyDeviceStatus:
				for _, d := range deviceEvents {
					result := make(map[string]interface{})
					meta := make(map[string]interface{})
					result["status"] = d.Status
					meta["deviceName"] = d.DeviceName
					meta["deviceModelName"] = d.DeviceModelName
					consumer <- api.NewDefaultSourceTuple(result, meta)
				}
			case RabbitMQRoutingKeyDeviceCommandStatus:
			}
		}
	}
}

func (s *source) Close(ctx api.StreamContext) error {
	ctx.GetLogger().Infof("Close rabbitmq source")
	s.channel.Close()
	return nil
}

func GetSource() *source {
	return &source{}
}
