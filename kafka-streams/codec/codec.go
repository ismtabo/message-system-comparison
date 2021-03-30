package codec

import (
	"encoding/json"

	"github.com/Telefonica/redis-vs-nats/model"
)

type MessageCodec struct{}

func (c *MessageCodec) Encode(value interface{}) (data []byte, err error) {
	return json.Marshal(value)
}

func (c *MessageCodec) Decode(data []byte) (value interface{}, err error) {
	var message model.Message
	err = json.Unmarshal([]byte(data), &message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}
