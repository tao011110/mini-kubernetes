package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli"

	"mini-kubernetes/tools/kubectl/command"
)

func Initial() *cli.App {
	app := cli.NewApp()
	app.Name = "kubectl"
	app.Version = "0.0.0"
	app.Usage = "Command line tool to communicate with apiserver."
	app.Flags = []cli.Flag{}
	app.Commands = []cli.Command{
		command.NewHelloCommand(),
		command.NewCreateCommand(),
		command.NewGetCommand(),
		command.NewDeleteCommand(),
		command.NewDescribeCommand(),
		command.NewUpdateCommand(),
	}
	return app
}

func main() {
	app := Initial()
	for {
		fmt.Printf(">")
		cmdReader := bufio.NewReader(os.Stdin)
		if cmdStr, err := cmdReader.ReadString('\n'); err == nil {
			cmdStr = strings.Trim(cmdStr, "\r\n")
			if cmdStr == "quit" {
				return
			} else {
				_ = ParseArgs(app, cmdStr)
			}
		}
	}
}

func ParseArgs(app *cli.App, cmdStr string) error {
	err := app.Run(strings.Split(cmdStr, " "))
	if err != nil {
		log.Fatal("[Fault] ", err)
	}
	return err
}
