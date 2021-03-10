package types

type SensorData struct {
	EntityId   string   // `json:"entity_id" bson:"entity_id"`
	State      string   // `json:"state" bson:"state"`
	Attributes struct { // `json:"attributes" bson:"attributes"`
		UnitOfMeasurement string // `json:"unit_of_measurement" bson:"unit_of_measurement"`
		FriendlyName      string // `json:"friendly_name" bson:"friendly_name"`
		DeviceClass       string // `json:"device_class" bson:"device_class"`
	}
	LastChanged string   // `json:"last_changed" bson:"last_changed"`
	LastUpdated string   // `json:"last_updated" bson:"last_updated"`
	Context     struct { // `json:"context" bson:"context"`
		Id       string // `json:"id" bson:"id"`
		ParentId string // `json:"device_class" bson:"device_class"`
		UserId   string // `json:"user_id" bson:"user_id"`
	}
}

type SensorValue struct {
	Name  string
	Value string
}

type ArukeresoResult struct {
	Name  string
	Price int
}
