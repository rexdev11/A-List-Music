package transcoder

import (
	"os"
	"io/ioutil"
	"encoding/json"
	"fmt"
)


type Transcoder interface {
	TransStore(file os.File)(url string, err error)
	Transcode(file os.File, toMime string)(data byte, err error)
}

type LocalVars struct {
	paths map[string] string `json:"paths"`
}

type Env interface {
	GetVars() (LocalVars, error)
	ExitChan() chan error
}

type SoundFileMeta struct {
	uri string
	encoding string
	codexLib string
	size int
}

//type FFMPEG interface {
//	Trancode(file os.File, meta SoundFileMeta) (os.File)
//
//}

func initializeFFMPEG() () {
	fmt.Print("starting transcoder")

	// Get ENVIRONMENT VARS
	EnvPath := "./local.json"
	localVars := make(map[string]string)
	local, err := ioutil.ReadFile(EnvPath)
	if (err != nil) {
		panic(err)
	}
	if err = json.Unmarshal(local, &localVars); err == nil {
		fmt.Print(localVars)
	}

}

func ProcessEnv(envChan chan Env, envPath string) {
		for {
			fmt.Println("\n env!")
			envVars := <- envChan
			vars := make(map[string] string)
			content, err := ioutil.ReadFile(envPath)
			if err == nil {
				if err = json.Unmarshal(content, &vars); err == nil {
					//vals, err := envVars.GetVars()
					//if err != nil {
					//	envVars.ExitChan() <- err
					//}
					//	fmt.Println("printing Env")
					//	fmt.Println(vals)
				}

			}
			envVars.ExitChan() <- err
		}
}

type ReadEnvVars struct {
	ec chan LocalVars
	exitChan chan error
}

func NewReadEnvVars() *ReadEnvVars {
	return &ReadEnvVars{
		ec: make(chan LocalVars, 1),
		exitChan: make(chan error, 1),
	}
}

func (ev ReadEnvVars) ExitChan() chan error {
	return ev.exitChan
}

func (ev ReadEnvVars) GetVars(localVars LocalVars) (LocalVars, error) {

	ev.ec <- localVars

	// ??
	return LocalVars{map[string]string{}}, nil
}

func TransStore() {
	initializeFFMPEG()

	//exec.Command("ffmpeg", )
	// catch STDOUT
}

func Transcode() {
	initializeFFMPEG()
}

func mimeOf(file os.File) {

}

type EnvClient struct {
	envs chan Env
}

//func (c *EnvClient) GetLocalVars() ([]LocalVars, error) {
	//arr := make([]LocalVars, 0)
	//
	//localVar := NewReadEnvVars()
	//c.envs <- localVar
//

//}