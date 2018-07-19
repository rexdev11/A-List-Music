package configs

import (
	"io/ioutil"
	"encoding/json"
)

type Paths struct {
	BIN string `json:"bin"`
	FFMPEG string `json:"ffmpeg"`
}

type LocalVars struct {
	Paths Paths `json:"paths"`
	Error error
}

func GetEnvironmentVars() (LocalVars) {
	localVars := LocalVars{}
	EnvPath := "./local.env.json"
	// Get ENVIRONMENT VARS
	local, err := ioutil.ReadFile(EnvPath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(local, &localVars)
	if err != nil {
		panic(err)
	}
	return localVars
}

