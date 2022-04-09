package properties

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
)

var checkConfigs = []string{
	"CompileTimeLimit",
	"CompileMemoryLimit",
	"Username",
	"Password",
	"ContainerName",
}

type Error struct {
	Info string
}

func (e Error) Error() string {
	return e.Info
}

var config = make(map[string]string)

func LoadProperties() error {
	buffer, err := ioutil.ReadFile("config.json")
	if err != nil {
		return Error{
			Info: "config.json文件不存在" + err.Error(),
		}
	}
	err = json.Unmarshal(buffer, &config)
	if err != nil {
		return Error{
			Info: "config.json文件格式错误" + err.Error(),
		}
	}
	for _, key := range checkConfigs {
		if _, ok := config[key]; !ok {
			return Error{
				Info: "配置:" + key + "不存在",
			}
		}
	}
	return err
}

func Get(key string) (string, error) {
	if v, ok := config[key]; ok {
		return v, nil
	} else {
		return "", Error{Info: "配置:" + key + "不存在"}
	}
}

func GetAnyway(key string) string {
	if v, ok := config[key]; ok {
		return v
	}
	return ""
}

func GetInt(key string) (int, error) {
	str, err := Get(key)
	if err != nil {
		return 0, err
	}
	value, err := strconv.Atoi(str)
	return value, err
}

func Set(key, value string) error {
	config[key] = value
	buf, err := json.Marshal(config)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("config.json", buf, 0777)
	return err
}
