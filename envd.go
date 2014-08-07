package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/op/go-logging"
)

const AppVersion = "0.1.0"

var log = logging.MustGetLogger("envd")

func init() {
	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.Name}} [global options] command

VERSION:
   {{.Version}}

COMMANDS:
   {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
   {{end}}
GLOBAL OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}
`
}

func main() {
	app := cli.NewApp()
	app.Name = "envd"
	app.Usage = "Application launcher using environment variables from etcd."
	app.Version = AppVersion
	app.Action = runCommand
	app.Flags = []cli.Flag{
		cli.BoolFlag{"verbose, V", "Shows verbose logging."},
		cli.StringFlag{"environment, e", "default", "Environment name to watch."},
		cli.BoolFlag{"d", "Daemonize after launch."},
		cli.BoolFlag{"debug, D", "Turn on debug output."},
		cli.BoolFlag{"watch, w", "Watch for updates and restart if changed."},
		cli.StringFlag{"key", "/etc/envd/client.key", "Client key file."},
		cli.StringFlag{"cert", "/etc/envd/client.crt", "Client cert file."},
		cli.StringFlag{"cacert", "/etc/envd/cacert.crt", "Client CA cert file."},
		cli.StringFlag{"server", "http://localhost:4001", "Host to connect to etcd."},
	}
	app.Before = func(ctx *cli.Context) error {
		if ctx.GlobalBool("verbose") || ctx.GlobalBool("debug") {
			logging.SetLevel(logging.DEBUG, "envd")
		} else {
			logging.SetLevel(logging.INFO, "envd")
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("%s", err)
	}
}

func runCommand(ctx *cli.Context) {
	e := NewEnvironment(ctx.GlobalString("environment"))

	etcdHosts := strings.Split(ctx.GlobalString("server"), ",")

	if err := e.GetEnvironment(etcdHosts); err != nil {
		log.Fatalf("Unable to get environment: %s", err)
	}

	args := ctx.Args()
	if len(args) < 1 {
		cli.ShowAppHelp(ctx)
		log.Fatal("No command specified.")
	}

	log.Debug("Running command '%s %s'", args[0], strings.Join(args[1:], " "))
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = e.Env()

	stdout, outErr := cmd.StdoutPipe()
	if outErr != nil {
		log.Error("Unable to get StdoutPipe: %s", outErr)
	}

	stderr, errErr := cmd.StderrPipe()
	if errErr != nil {
		log.Error("Unable to get StderrPipe: %s", errErr)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal("Command %s failed: %s", strings.Join(args, " "), err)
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
