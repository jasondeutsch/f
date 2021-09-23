package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

/*
   TODOS:
   1. write to files in /usr/local/var/f/
      a. tempfile for current command
      b. log of commands in single file
   2. support other editors
   4. figure out how to paste to prompt instead of just executing
*/

func main() {
	opts := flags()

	file, err := ioutil.TempFile("/tmp", "f.*.sh")
	if err != nil {
		panic(err)
	}
	defer os.Remove(file.Name())

	if len(os.Args) > 1 {
		skipFlags := ""
		for _, v := range os.Args[1:] {
			if v[0] == '-' {
				continue
			}
			skipFlags += " " + v
		}
		skipFlags = strings.TrimLeft(skipFlags, " ")
		_, err := file.WriteString(skipFlags)
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

	if opts.DryRun {
		fmt.Println(string(b))
		os.Exit(0)
	}

	name, args := prepare(string(b))
	if len(name) == 0 {
		fmt.Println("noop")
		os.Exit(0)
	}

	execAndWait(name, args)
}

type options struct {
	DryRun bool
}

func flags() options {
	// get flags
	helpFlag := flag.Bool("help", false, "print usage information")
	helpFlagShort := flag.Bool("h", false, "print usage information (shorthand)")
	dryrunFlag := flag.Bool("dry-run", false, "print command without execution")

	flag.Parse()

	// handle flags
	if (*helpFlag) ||
		(helpFlagShort != nil && *helpFlagShort) {
		fmt.Print(help)
		os.Exit(0)
	}

	var opts options

	if *dryrunFlag {
		opts.DryRun = true
	}

	return opts
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
	for i, v := range s {
		if (v == '\\' && s[i+1] == '\n') || v == '\n' {
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

const help = `
Usage: f  
   or: f [options] [text] Open editor with text provided e.g. f !! to open with last command

Options:
  -n, --name    Save command by name.	
  --history     Print command history (only named commands are saved).
  --dry-run     Print command without execution.
`
