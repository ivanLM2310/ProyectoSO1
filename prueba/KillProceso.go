package main

import (
	"fmt"
	"log"
	"os"
	"syscall"
)

func main() {
	pid := os.Getpid()
	fmt.Println("PID:", pid)
	r, w, err := os.Pipe()

	if err != nil {
		panic(err)
	}

	defer r.Close()

	process, err := os.StartProcess("/bin/ps", []string{"-ef"}, &os.ProcAttr{Files: []*os.File{nil, w, os.Stderr}})

	if err != nil {
		panic(err)
	}

	processState, err := process.Wait()

	if err != nil {
		panic(err)
	}

	err = process.Release()

	if err != nil {
		panic(err)
	}

	fmt.Println("PID proceso? : ", 6960)
	fmt.Println("Proceso terminado? : ", processState.Success())

	err = process.Signal(syscall.SIGKILL)
	proc, err := os.FindProcess(6960)
	if err != nil {
		log.Println(err)
	}
	// Kill el proceso
	proc.Kill()

	if err != nil {
		fmt.Println(err)
		return
	}

	w.Close()

}
