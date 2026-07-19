package runner

// Running command and parse outputs

import (
	"fmt"
	"io"
	"log"
	"os/exec"
)

func Run(cmd string) (stdoutBytes, stderrBytes chan []byte) {
	cmdBuilder := exec.Command("pwsh", "-Command", cmd)
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
			log.Println(err)
		}
		stdoutBytes <- stdoutbytes
		close(stdoutBytes)
	}()

	go func() {
		stderrbytes, err := io.ReadAll(stderr)
		if err != nil {
			log.Println(err)
		}
		stderrBytes <- stderrbytes
		close(stderrBytes)
	}()

	go func() {
		if err := cmdBuilder.Wait(); err != nil {
			fmt.Println("Wait error:", err)
		}
	}()

	return stdoutBytes, stderrBytes
}
