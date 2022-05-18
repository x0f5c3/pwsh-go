package pkg

import (
	"bufio"
	"github.com/go-ini/ini"
	"github.com/pterm/pterm"
	"os/exec"
	"strings"
	"sync"
)

type LinuxType int

const (
	Debian LinuxType = iota
	RPM
	Other
)

func ReadOSRelease(configfile string) (LinuxType, error) {
	cfg, err := ini.Load(configfile)
	if err != nil {
		pterm.Error.WithShowLineNumber(true).PrintOnError(err)
		return Other, err
	}
	idLike := cfg.Section("").Key("ID_LIKE").Value()
	if strings.Contains(idLike, "ubuntu") || strings.Contains(idLike, "debian") {
		return Debian, nil
	} else if strings.Contains(idLike, "rhel") || strings.Contains(idLike, "fedora") {
		return RPM, nil
	}
	return Other, nil
}

var OSType = func() LinuxType {
	res, err := ReadOSRelease("/etc/os-release")
	pterm.Error.WithFatal(true).WithShowLineNumber(true).PrintOnError(err)
	return res
}()

func RunInstaller(path string) (string, string, error) {
	cmd := exec.Command("sudo", "apt", "install", "-y", path)
	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return "", "", err
	}
	stdErr, err := cmd.StderrPipe()
	if err != nil {
		return "", "", err
	}
	err = cmd.Start()
	if err != nil {
		return "", "", err
	}
	stdOutString := ""
	stdErrString := ""
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		b := bufio.NewScanner(stdOut)
		for b.Scan() {
			t := b.Text()
			pterm.Info.WithPrefix(pterm.Prefix{
				Text:  "STDOUT",
				Style: pterm.NewStyle(pterm.FgCyan),
			}).Println(t)
			stdOutString += t
			stdOutString += "\n"
		}
	}()
	go func() {
		defer wg.Done()
		b := bufio.NewScanner(stdErr)
		for b.Scan() {
			t := b.Text()
			pterm.Info.WithPrefix(pterm.Prefix{
				Text:  "STDERR",
				Style: pterm.NewStyle(pterm.FgCyan),
			}).Println(t)
			stdErrString += t
			stdErrString += "\n"
		}
	}()
	err = cmd.Wait()
	if err != nil {
		return "", "", err
	}
	wg.Wait()
	return stdOutString, stdErrString, nil
}
