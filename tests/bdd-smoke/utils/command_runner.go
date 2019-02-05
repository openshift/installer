package utils

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	expect "github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"
)

// RunOSCommandWithArgs executes a command in the operating system with arguments
func RunOSCommandWithArgs(command string, arguments []string, path string) string {
	cmd := exec.Command(command, arguments...)
	cmd.Dir = path
	cmd.Stdin = strings.NewReader("")
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error: failed running the command \"" + command + "\" with the arguments:")
		fmt.Println(arguments)
		fmt.Println(errb.String())
		log.Fatal(err)
	}
	return outb.String() + errb.String()
}

// RunInstallerWithSurvey run the installer given the command arguments, the path where will be launched and the command line values for survey
func RunInstallerWithSurvey(arguments []string, path string, commandLineArgs [][]string) string {
	openshiftInstallBin := os.Getenv("GOPATH") + "/src/github.com/openshift/installer/bin/openshift-install"

	c, state, errVt10x := vt10x.NewVT10XConsole()
	if errVt10x != nil {
		fmt.Println("Error: failed starting vt10x console")
		log.Fatal(errVt10x)
	}
	defer c.Close()

	donec := make(chan struct{})
	go func() {
		defer close(donec)

		for _, command := range commandLineArgs {
			c.ExpectString(command[0])
			c.SendLine(command[1])
		}
		c.ExpectEOF()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, openshiftInstallBin, arguments...)
	cmd.Dir = path
	cmd.Stdin = c.Tty()
	cmd.Stdout = c.Tty()
	cmd.Stderr = c.Tty()

	errCmd := cmd.Run()
	if errVt10x != nil {
		fmt.Println("Error: running the command : " + openshiftInstallBin + " with the arguments")
		fmt.Println(arguments)
		fmt.Println(commandLineArgs)
		log.Fatal(errCmd)
	}
	c.Tty().Close()
	<-donec

	return expect.StripTrailingEmptyLines(state.String())
}
