package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

var envName string
var config string
var daemonize bool
var debug bool

func init() {
	flag.StringVar(&envName, "e", "default", "Environment name to watch.")
	flag.StringVar(&config, "c", "/etc/envd.yml", "envd config file.")
	flag.BoolVar(&daemonize, "d", false, "Daemonize after launch.")
	flag.BoolVar(&debug, "D", false, "Daemonize after launch.")
}

func runCommand(e *Environment, name string) {
	cmd := exec.Command(name)
	cmd.Env = e.Env()

	stdout, outErr := cmd.StdoutPipe()
	if outErr != nil {
		log.Println(outErr)
	}

	stderr, errErr := cmd.StderrPipe()
	if errErr != nil {
		log.Println(errErr)
	}

	err := cmd.Start()
	if err != nil {
		log.Fatal(fmt.Sprintf("Command %s failed: %s", name, err))
	}

	var outBuffer bytes.Buffer
	var errBuffer bytes.Buffer

	io.Copy(&outBuffer, stdout)
	if outBuffer.Len() > 0 {
		fmt.Fprintf(os.Stdout, "%s", outBuffer.String())
	}

	io.Copy(&errBuffer, stderr)
	if errBuffer.Len() > 0 {
		fmt.Fprintf(os.Stderr, "%s", errBuffer.String())
	}

	cmd.Wait()
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <command> <options>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(2)
	}

	flag.Parse()

	if debug {
		log.Println(fmt.Sprintf("Flag: envName: %s", envName))
		log.Println(fmt.Sprintf("Flag: config: %s", config))
		log.Println(fmt.Sprintf("Flag: daemonize: %t", daemonize))
	}

	if len(flag.Args()) > 0 {
		e := NewEnvironment(envName)
		err := e.GetEnvironment()
		if err != nil {
			log.Println(fmt.Sprintf("Error getting environment from etcd: %s", err))
			os.Exit(-1)
		}
		if debug {
			log.Println(fmt.Sprintf("Environment: %s", e.Env()))
		}
		runCommand(e, flag.Arg(0))
	} else {
		flag.Usage()
	}
}
