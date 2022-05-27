package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
)

type StatusJSON struct {
	last_parsed_height_cosmos uint64
	last_parsed_height_native uint64
}

func NewStatusJSON(homePath string, cosmos_height uint64, native_height uint64) {

	jsonStruct := StatusJSON{}
	jsonStruct.last_parsed_height_native = native_height
	jsonStruct.last_parsed_height_cosmos = cosmos_height

	content, err := json.Marshal(jsonStruct)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile(filepath.Join(homePath, "status.json"), content, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func GetHeightStatus(homePath string) (c_height, n_height uint64) {
	content, err := ioutil.ReadFile(filepath.Join(homePath, "status.json"))
	if err != nil {
		log.Fatal(err)
	}
	heightStatus := StatusJSON{}
	err = json.Unmarshal(content, &heightStatus)
	if err != nil {
		log.Fatal(err)
	}
	return heightStatus.last_parsed_height_cosmos, heightStatus.last_parsed_height_native
}
