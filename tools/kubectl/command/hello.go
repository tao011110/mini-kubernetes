package command

import (
	"fmt"

	"github.com/urfave/cli"
)

func NewHelloCommand() cli.Command {
	return cli.Command{
		Name:  "hello",
		Usage: "Test",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) error {
			fmt.Println("Hi! This is kubectl for minik8s.")
			return nil
		},
	}
}
