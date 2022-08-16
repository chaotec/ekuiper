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

type BaseUnion struct {
	Name string `json:"name,omitempty"`
	Id   string `json:"id,omitempty"`
}

type RabbitMQBody struct {
	MsgId     string             `json:"msgId,omitempty"`
	MsgTopic  RabbitMQRoutingKey `json:"msgTopic,omitempty"`
	MsgEvent  RabbitMQEvent      `json:"msgEvent,omitempty"`
	Timestamp string             `json:"timestamp,omitempty"`
	Body      []byte             `json:"body,omitempty"`
}

type DeviceData struct {
	Device      BaseUnion   `json:"device,omitempty"`
	DeviceModel BaseUnion   `json:"deviceModel,omitempty"`
	Properties  []Property  `json:"properties,omitempty"`
	Alarms      []BaseUnion `json:"alarms,omitempty"`
	Status      string      `json:"status,omitempty"`
}

type Property struct {
	BaseUnion
	ValueType   string `json:"valueType,omitempty"`
	ReportValue string `json:"reportValue,omitempty"`
	DesireValue string `json:"desireValue,omitempty"`
}
