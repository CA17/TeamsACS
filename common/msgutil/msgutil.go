package msgutil

import (
	"encoding/json"

	"github.com/ca17/teamsacs/common/aes"
)

func _ekey() string {
	return string([]byte{49, 50, 51, 52, 53, 54, 119, 56, 49, 50, 51, 52, 53, 54, 106, 56, 49, 50, 51, 52, 53, 54, 120, 56})
}

func Marshal(v interface{}) ([]byte, error) {
	bs, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return aes.Encrypt(bs, _ekey())
}

func Unmarshal(data []byte, v interface{}) error {
	bs, err := aes.Decrypt(data, _ekey())
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, v)
}
