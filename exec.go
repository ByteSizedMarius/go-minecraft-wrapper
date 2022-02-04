package wrapper

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
)

type JavaExec interface {
	Stdout() io.ReadCloser
	Stdin() io.WriteCloser
	Start() error
	Kill() error
}

type defaultJavaExec struct {
	cmd *exec.Cmd
}

func (j *defaultJavaExec) Stdout() io.ReadCloser {
	r, _ := j.cmd.StdoutPipe()
	return r
}

func (j *defaultJavaExec) Stdin() io.WriteCloser {
	w, _ := j.cmd.StdinPipe()
	return w
}

func (j *defaultJavaExec) Start() error {
	return j.cmd.Start()
}

func (j *defaultJavaExec) Kill() error {
	return j.cmd.Process.Kill()
}

func javaExecCmd(serverPath string, initialHeapSize, maxHeapSize int) *defaultJavaExec {
	initialHeapFlag := fmt.Sprintf("-Xms%dM", initialHeapSize)
	maxHeapFlag := fmt.Sprintf("-Xmx%dM", maxHeapSize)
	cmd := exec.Command("java", initialHeapFlag, maxHeapFlag, "-jar", serverPath, "nogui")

	// normalize path
	serverPathClean := path.Clean(serverPath)

	// check if path is a dir
	file, err := os.Open(serverPathClean)
	if err != nil {
		log.Print(err)
		return &defaultJavaExec{}
	}

	info, err := file.Stat()
	if err != nil {
		log.Print(err)
		return &defaultJavaExec{}
	}

	// if given path is a directory (aka if this program is executed in different directory),
	// change the directory, the jar will be executed from
	if !info.IsDir() {
		cmd.Dir = path.Dir(serverPath)
	}

	return &defaultJavaExec{cmd: cmd}
}
