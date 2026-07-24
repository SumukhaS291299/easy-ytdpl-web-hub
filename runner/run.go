package runner

// Running command and parse outputs

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
)

type Shell struct {
	Program string
	Args    []string
}

var (
	cachedShell Shell
	once        sync.Once
)

func candidates() []Shell {
	switch runtime.GOOS {
	case "windows":
		return []Shell{
			{Program: "pwsh", Args: []string{"-Command"}},
			{Program: "powershell", Args: []string{"-Command"}},
			{Program: "cmd", Args: []string{"/C"}},
		}

	case "darwin": // macOS
		return []Shell{
			{Program: "zsh", Args: []string{"-c"}},
			{Program: "bash", Args: []string{"-c"}},
			{Program: "sh", Args: []string{"-c"}},
		}

	default: // Linux and other Unix-like OSes
		return []Shell{
			{Program: "bash", Args: []string{"-c"}},
			{Program: "sh", Args: []string{"-c"}},
		}
	}
}

func DetectShell() (Shell, error) {
	var err error

	once.Do(func() {
		for _, shell := range candidates() {
			if _, e := exec.LookPath(shell.Program); e == nil {
				cachedShell = shell
				return
			}
		}
		err = errors.New("no supported shell found")
	})
	log.Infof("Using shell: %s %v", cachedShell.Program, cachedShell.Args)

	return cachedShell, err
}

func Run(cmd string) (stdoutBytes, stderrBytes chan []byte) {
	shell, err := DetectShell()
	if err != nil {
		log.Error("[Error]:Running the code\n", err)
		return nil, nil
	}

	args := append(shell.Args, cmd)

	cmdBuilder := exec.Command(shell.Program, args...)
	log.Info("Running command", "["+strings.ToUpper(shell.Program)+"]:\t", cmd)
	stdout, err := cmdBuilder.StdoutPipe()
	if err != nil {
		return
	}
	stderr, err := cmdBuilder.StderrPipe()
	if err != nil {
		return
	}
	if err := cmdBuilder.Start(); err != nil {
		return
	}
	fmt.Printf("Ran the command %s...\n", cmd)

	stdoutBytes = make(chan []byte)
	stderrBytes = make(chan []byte)

	go func() {
		stdoutbytes, err := io.ReadAll(stdout)
		if err != nil {
			log.Error(err)
		}
		stdoutBytes <- stdoutbytes
		close(stdoutBytes)
	}()

	go func() {
		stderrbytes, err := io.ReadAll(stderr)
		if err != nil {
			log.Error(err)
		}
		stderrBytes <- stderrbytes
		close(stderrBytes)
	}()

	go func() {
		if err := cmdBuilder.Wait(); err != nil {
			log.Error("Wait error:", err)
		}
	}()

	return stdoutBytes, stderrBytes
}
