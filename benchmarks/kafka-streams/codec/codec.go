package codec

import (
	"encoding/json"

	"github.com/Telefonica/redis-vs-nats/model"
)

type MessageCodec struct{}

func (c *MessageCodec) Encode(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

func (c *MessageCodec) Decode(data []byte) (interface{}, error) {
	var message model.Message
	if err := json.Unmarshal([]byte(data), &message); err != nil {
		return nil, err
	}
	return &message, nil
}
