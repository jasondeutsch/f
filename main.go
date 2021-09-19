package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func main() {
	file, err := ioutil.TempFile("/tmp", "f.*.sh")
	if err != nil {
		panic(err)
	}
	defer os.Remove(file.Name())

	if err := runEditor(file.Name()); err != nil {
		panic(err)
	}

	b, err := os.ReadFile(file.Name())
	if b == nil {
		fmt.Println("noop")
		os.Exit(0)
	}
	if err != nil {
		panic(err)
	}

	name, args := prepare(string(b))

	execAndPipe(name, args)
}

func prepare(s string) (string, []string) {
	s = clean(s)
	slc := strings.Split(s, " ")
	return slc[0], slc[1:]
}

func clean(s string) string {
	return strings.TrimSpace(s)
}

func runEditor(path string) error {
	cmd := exec.Command("vi", path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func execAndPipe(name string, args []string) {
	cmd := exec.Command(name, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	defer cmd.Wait()

	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)
}
