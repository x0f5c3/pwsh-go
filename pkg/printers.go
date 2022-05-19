package pkg

import (
	"bufio"
	"github.com/pterm/pterm"
	"os/exec"
	"sync"
)

var (
	Out = pterm.PrefixPrinter{
		MessageStyle: &pterm.ThemeDefault.InfoMessageStyle,
		Prefix: pterm.Prefix{
			Style: &pterm.ThemeDefault.InfoPrefixStyle,
			Text:  "STDOUT",
		},
	}
	Err = pterm.PrefixPrinter{
		Prefix: pterm.Prefix{
			Style: &pterm.ThemeDefault.ErrorPrefixStyle,
			Text:  "STDERR",
		},
		MessageStyle: &pterm.ThemeDefault.ErrorMessageStyle,
	}
)

type Cmd struct {
	PrintStdout bool
	PrintStderr bool
	cmd         *exec.Cmd
}

func (c *Cmd) SetPrintStdout(PrintStdout bool) {
	c.PrintStdout = PrintStdout
}

func (c *Cmd) SetPrintStderr(PrintStderr bool) {
	c.PrintStderr = PrintStderr
}

func Command(name string, arg ...string) *Cmd {
	cmd := exec.Command(name, arg...)
	return &Cmd{
		cmd: cmd,
	}
}

func (c *Cmd) Run() error {
	wg := &sync.WaitGroup{}
	if c.PrintStdout {
		wg.Add(1)
		out, err := c.cmd.StdoutPipe()
		if err != nil {
			return err
		}
		go func() {
			defer wg.Done()
			scanner := bufio.NewScanner(out)
			for scanner.Scan() {
				t := scanner.Text()
				Out.Println(t)
			}
		}()
	}
	if c.PrintStderr {
		wg.Add(1)
		out, err := c.cmd.StderrPipe()
		if err != nil {
			return err
		}
		go func() {
			defer wg.Done()
			scanner := bufio.NewScanner(out)
			for scanner.Scan() {
				t := scanner.Text()
				Err.Println(t)
			}
		}()
	}
	err := c.cmd.Start()
	if err != nil {
		return err
	}
	err = c.cmd.Wait()
	wg.Wait()
	return err
}
