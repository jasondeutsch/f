package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

/*
   TODOS:
   2. support other editors
   4. figure out how to paste to prompt instead of just executing
   5. remove panics, handle errors
   6. enforce opts before text
   7. main() is getting bloated...
*/

const (
	logPath = "/usr/local/var/f"
	logFile = "flog"
)

type options struct {
	DryRun     bool
	CmdLogName string
}

func main() {
	setup()
	opts := flags()

	file, err := ioutil.TempFile(logPath, "temp.*.sh")
	if err != nil {
		panic(err)
	}
	defer func() {
		// not working
		err = os.Remove(file.Name())
		if err != nil {
			log.Println(fmt.Sprintf("error removing temp file: %v", err))
		}
	}()

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

	name, args := prepare(string(b))
	if len(name) == 0 {
		fmt.Println("noop")
		os.Exit(0)
	}

	err = writeLog(opts.CmdLogName, string(b))
	if err != nil {
		log.Fatal(err)
	}

	if opts.DryRun {
		fmt.Println(string(b))
		os.Exit(0)
	}
	execAndWait(name, args)
}

func writeLog(name, cmd string) error {
	// If the file doesn't exist, create it, or append to the file
	file, err := os.OpenFile(logPath+"/"+logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString("--- " + name + "\n" + cmd)
	return err
}

func setup() {
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		if err := os.Mkdir(logPath, os.ModePerm); err != nil {
			panic(err)
		}
	}
}

func flags() options {
	// get flags
	helpFlag := flag.Bool("help", false, "print usage information")
	helpFlagShort := flag.Bool("h", false, "print usage information (shorthand)")
	dryRunFlag := flag.Bool("dry-run", false, "print command without execution")
	historyFlag := flag.Bool("history", false, "print command log")
	cmdLogNameFlag := flag.String("log", "", "label command and save to f's history")
	cmdLogNameFlagShort := flag.String("l", "", "label command and save to f's history (shorthand)")

	flag.Parse()

	// handle flags
	if *helpFlag || *helpFlagShort {
		fmt.Print(help)
		os.Exit(0)
	}

	if *historyFlag {
		// lazy...
		// What? this is not enterprise software.
		// This will error if logFile not exists.
		cmd := exec.Command("cat", logPath+"/"+logFile)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
		os.Exit(0)
	}

	var opts options
	max := func(a, b string) string {
		if len(b) > len(a) {
			return b
		}
		return a
	}
	if n := max(*cmdLogNameFlag, *cmdLogNameFlagShort); n != "" {
		opts.CmdLogName = n
	}

	if *dryRunFlag {
		opts.DryRun = true
		opts.CmdLogName = strings.Join([]string{opts.CmdLogName, "(dry-run)"}, " ")
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
  -l, --label   Label command and write to log.	
  --history     Print command history (only named commands are saved).
  --dry-run     Print command without execution.
`
