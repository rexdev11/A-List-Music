package transcoder

import (
	"os"
	"os/exec"
	"github.com/kjk/betterguid"
	"net/http"
	"fmt"
)

type Transcoder interface {
	RunToMediaLibrary()(SoundFileMeta, err error)
	NewJob(file os.File, targetMime []string)(source SoundFileMeta)
	RunJobs()(data byte, err error)
	ExitChan() chan error
}

type TranscodeJob struct {
	id string
	ready bool
	done bool
	sourceMeta SoundFileMeta
	targetMeta SoundFileMeta
	ffmpegCMD exec.Cmd
}

type SoundFileMeta struct {
	id string
	uri string
	encoding string
	codex string
	size int
}

func buildFFMPEGCMD(sourceMeta SoundFileMeta) *exec.Cmd {
	return exec.Command("ffmpeg", "-i", sourceMeta.uri, "-vn", "-ar 44100", "-ac 2", "-ab 192l", "-f mp3", sourceMeta.id + ".mp3")
}

type TranscoderClient struct {
	Transcodes chan TranscodeJob
}

func BuildTranscodeClient(transcoderClient chan TranscoderClient) TranscoderClient {

	client := make(chan Transcoder)
	transcoderClient = client
}


func TransStore() {
	//initializeFFMPEG()

	//exec.Command("ffmpeg", )
	// catch STDOUT
}

func Transcode() {
	//initializeFFMPEG()
}

//func (c *TranscoderClient) RunToMediaLibrary()(SoundFileMeta, err error) {
//	todo...
//}

func (c *TranscoderClient) NewJob(_file *os.File, targetMime []string)(source SoundFileMeta) {
	buffer := make([]byte, 512)
	_, err := _file.Read(buffer)
	if err != nil {
		panic(err)
	}

	encodeing := http.DetectContentType(buffer)
	fmt.Println(encodeing)
	// make a uri
	//uri := os.NewFile(file)
	// setMeta
	sourceMeta := SoundFileMeta{id: betterguid.New()}
	// buildCMD
	c.Transcodes <- TranscodeJob{
		id: betterguid.New(),
		ready: true,
		done: false,
		sourceMeta: sourceMeta,
		targetMeta: nil,
		ffmpegCMD: nil,
	}
	return SoundFileMeta{}
}
//func (c *TranscoderClient) RunJobs()(data byte, err error) {
//
//}
//func (c *TranscoderClient) ExitChan() chan error {
//
//}
