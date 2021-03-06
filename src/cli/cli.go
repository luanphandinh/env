package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"
)

// Executable is a callback for give commmand
// It recieve CLI object as a parameter
type Executable func(context *CLI)

type Argument struct {
	Name        string
	Exec        Executable
	Usage       string
	Description string
}

// Command declare basic struct of a simple command
// Command can be nested inside command.
// If Command have other commands within it, Exec will be ignored.
type Command struct {
	Name        string
	Usage       string
	Description string
	Exec        Executable
	Commands    []Command
}

// CLI main object for CLI
// Contain args and cfgs bag
type CLI struct {
	Commands    []Command
	Arguments   []Argument
	Init        func(cli *CLI)
	HandlePanic func(cli *CLI, e interface{})
	args        []string
	configs     interface{}
	cmd         *Command
}

// ShiftStrictArg shift cli.args by 1
// Return the shifted arguments
// Will panic if there is no more arguments in args bag
func (cli *CLI) ShiftStrictArg() string {
	if len(cli.args) < 1 {
		panic("Too few arguments")
	}

	return cli.ShiftArg()
}

// ShiftArg return the shifted first arg
// If there is no arguments left, return empty string
func (cli *CLI) ShiftArg() string {
	if len(cli.args) < 1 {
		return ""
	}

	arg := cli.args[0]
	cli.args = cli.args[1:]

	return arg
}

func (cli *CLI) init() {
	cli.args = os.Args[1:]
	if cli.Init != nil {
		cli.Init(cli)
	}
}

// Run CLI command base on cli.args
func (cli *CLI) Run() {
	cli.init()

	defer func() {
		if r := recover(); r != nil {
			cli.HandlePanic(cli, r)
			cli.help()
		}
	}()

	// Apply cli level arguments
	for len(cli.args) > 0 && strings.HasPrefix(cli.args[0], "--") {
		arg := cli.ShiftArg()
		for _, executable := range cli.Arguments {
			if strings.EqualFold(executable.Name, arg[2:]) {
				executable.Exec(cli)
			}
		}
	}

	cli.invokeCommand(cli.Commands)
}

func (cli *CLI) invokeCommand(cmds []Command) {
	name := cli.ShiftStrictArg()

	for _, cmd := range cmds {
		if matchCmd(name, cmd.Name) {
			cli.cmd = &cmd
			if len(cmd.Commands) > 0 {
				cli.invokeCommand(cmd.Commands)
				return
			}

			if cmd.Exec != nil {
				cmd.Exec(cli)
				return
			}
		}
	}

	panic("")
}

func matchCmd(name, cmds string) bool {
	cases := strings.Split(cmds, "|")
	for _, c := range cases {
		if name == c {
			return true
		}
	}

	return false
}

func (cli *CLI) help() {
	if cli.cmd == nil {
		fmt.Print(
			`
CLI control your local development environment
version:  2.1

`)
		fmt.Println("Cli's arguments:")
		printArguments(cli.Arguments)
		fmt.Println()

		fmt.Println("Commands:")
		printCommands(cli.Commands)
		fmt.Println()
	} else {
		fmt.Println(fmt.Sprintf("%s %s\n", cli.cmd.Name, cli.cmd.Description))
		fmt.Printf("Available sub commands:\n\n")
		printCommands(cli.cmd.Commands)
		fmt.Println()
	}
}

func printArguments(args []Argument) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 4, ' ', 0)
	for _, arg := range args {
		fmt.Fprintln(w, fmt.Sprintf("\t%s\t%s", arg.Name, arg.Description))
		fmt.Fprintln(w, fmt.Sprintf("\t%s\t%s", "", arg.Usage))
	}
	w.Flush()
}

func printCommands(cmds []Command) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 4, ' ', 0)
	for _, cmd := range cmds {
		fmt.Fprintln(w, fmt.Sprintf("\t%s\t%s", cmd.Name, cmd.Description))
	}
	w.Flush()
}

// GetLastExecutedCommand return last command that cli run.
func (cli *CLI) GetLastExecutedCommand() *Command {
	return cli.cmd
}

// GetConfigs return interface{} config, that being set via SetConfigs.
func (cli *CLI) GetConfigs() interface{} {
	return cli.configs
}

// SetConfigs allowed to put custom configuration to cli object.
func (cli *CLI) SetConfigs(cfg interface{}) {
	cli.configs = cfg
}

// GetArgs returns all available args
func (cli *CLI) GetArgs() []string {
	return cli.args
}

// ExecCmd execute "c" command with given args..
func (cli *CLI) ExecCmd(c string, args ...string) error {
	cmd := exec.Command(c, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
