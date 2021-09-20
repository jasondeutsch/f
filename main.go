package main

import (
	"fmt"
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

	if len(os.Args) > 1 {
		_, err := file.WriteString(strings.Join(os.Args[1:], " "))
		if err != nil {
			panic(err)
		}
	}

	if err := runEditor(file.Name()); err != nil {
		panic(err)
	}

	b, err := os.ReadFile(file.Name())
	if err != nil {
		panic(err)
	}

	name, args := prepare(string(b))
	if len(name) == 0 {
		fmt.Println("noop")
		os.Exit(0)
	}

	execAndWait(name, args)
}

func prepare(s string) (string, []string) {
	s = collapse(s)
	s = clean(s)
	slc := strings.Split(s, " ")
	return slc[0], slc[1:]
}

// convert multiline commands (e.g. `\`) to single line
func collapse(s string) string {
	var sb strings.Builder
	for  i, v := range s {
		if (v == '\\'  && s[i+1] == '\n') || v == '\n' {
			continue
		}
		sb.WriteRune(v)
	}
	return sb.String()
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

func execAndWait(name string, args []string) {
	cmd := exec.Command(name, args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	defer cmd.Wait()
}
