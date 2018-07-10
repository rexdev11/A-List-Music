package transcoder

import (
	"os"
	"os/exec"
)

type Transcoder interface {
	RunToMediaLibrary()(SoundFileMeta, err error)
	NewJob(file os.File, targetMime []string)(source SoundFileMeta)
	RunJobs()(data byte, err error)
	ExitChan() chan error
}

type TranscodeJob struct {
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
	transcodes chan
}

func BuildTranscodeClient() Transcoder {

}


func TransStore() {
	//initializeFFMPEG()

	//exec.Command("ffmpeg", )
	// catch STDOUT
}

func Transcode() {
	//initializeFFMPEG()
}
