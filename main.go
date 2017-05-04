package main

import (
	"os"

	"github.com/mitchellh/cli"
	"github.com/wj24021040/etcd_2v3/command"
)

var (
	Commands = map[string]cli.CommandFactory{
		"update": command.UpdateCommandFactory,
	}
)

func main() {
	cli := &cli.CLI{
		Args:     os.Args[1:],
		Commands: Commands,
		HelpFunc: cli.BasicHelpFunc("etcd_2v3"),
	}

	_, _ = cli.Run()

}
