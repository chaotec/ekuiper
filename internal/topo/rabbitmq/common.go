package rabbitmq

type RabbitMQRoutingKey string

const (
	RabbitMQRoutingKeyDeviceModel         RabbitMQRoutingKey = "device-model"
	RabbitMQRoutingKeyDevice              RabbitMQRoutingKey = "device"
	RabbitMQRoutingKeySceneLink           RabbitMQRoutingKey = "scene-link"
	RabbitMQRoutingKeyDeviceData          RabbitMQRoutingKey = "device-data"
	RabbitMQRoutingKeyDeviceAlarm         RabbitMQRoutingKey = "device-alarm"
	RabbitMQRoutingKeyDeviceStatus        RabbitMQRoutingKey = "device-status"
	RabbitMQRoutingKeyDeviceCommandStatus RabbitMQRoutingKey = "device-command-status"
	RabbitMQRoutingKeyServerSub           RabbitMQRoutingKey = "server-subscription"
)

type RabbitMQEvent string

const (
	RabbitMQEventAdd    RabbitMQEvent = "add"
	RabbitMQEventUpdate RabbitMQEvent = "update"
	RabbitMQEventDelete RabbitMQEvent = "delete"
)

type RabbitMQBody struct {
	MsgId     string             `json:"msgId,omitempty"`
	MsgTopic  RabbitMQRoutingKey `json:"msgTopic,omitempty"`
	MsgEvent  RabbitMQEvent      `json:"msgEvent,omitempty"`
	Timestamp string             `json:"timestamp,omitempty"`
	Body      interface{}        `json:"body,omitempty"`
}

type DeviceEvent struct {
	DeviceName      string     `json:"deviceName,omitempty"`
	DeviceModelName string     `json:"deviceModelName,omitempty"`
	Properties      []Property `json:"properties,omitempty"`
	Alarms          []Alarm    `json:"alarms,omitempty"`
	Status          string     `json:"status,omitempty"`
}

type Property struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	ValueType   string `json:"valueType,omitempty"`
	ReportValue string `json:"reportValue,omitempty"`
	DesireValue string `json:"desireValue,omitempty"`
}

type Alarm struct {
	AlarmId   string `json:"alarmId,omitempty"`
	AlarmName string `json:"alarmName,omitempty"`
}
