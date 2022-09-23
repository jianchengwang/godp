package process

import (
	"encoding/json"
	"fmt"
	_pojo "godp/internal/pojo"
	"io/ioutil"
	"os"
)

func LoadProjectConfigJsonTpl() _pojo.ProjectConfig {
	packageJson, err := os.Open("tpl/project_config_tpl.json")
	if err != nil {
		fmt.Println(err)
	}

	defer packageJson.Close()
	byteValue, _ := ioutil.ReadAll(packageJson)
	fmt.Println(string(byteValue))

	projectConfig := _pojo.ProjectConfig{}
	err = json.Unmarshal(byteValue, &projectConfig)
	if err != nil {
		fmt.Println(err)
	}

	return projectConfig
}
