package common

import  (
	"encoding/json"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func InterfaceToMap(data interface{}) (map[string]interface{}, error) {
	var dataToMap map[string]interface{}
	dataByte, err := json.Marshal(data)
	if nil != err {
		return nil, err
	}

	if err := json.Unmarshal(dataByte, &dataToMap); nil != err {
		return nil, err
	}

	return dataToMap, nil
}

func Encrypt(data string, keyString string) (string) {
	key := []byte(keyString)

	sig := hmac.New(sha256.New, key)
	sig.Write([]byte(data))

	return hex.EncodeToString(sig.Sum(nil))
}
