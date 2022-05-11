package kubectl

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli"

	"mini-kubernetes/tools/kubectl/command"
)

func Kubectl() {
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
	}
	for {
		fmt.Printf(">")
		cmdReader := bufio.NewReader(os.Stdin)
		if cmdStr, err := cmdReader.ReadString('\n'); err == nil {
			cmdStr = strings.Trim(cmdStr, "\r\n")
			if cmdStr == "quit" {
				return
			} else {
				if err := app.Run(strings.Split(cmdStr, " ")); err != nil {
					log.Fatal("[Fault] ", err)
				}
			}
		}
	}
}
