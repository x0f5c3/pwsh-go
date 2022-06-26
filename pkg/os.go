package pkg

import (
	"bufio"
	"errors"
	"github.com/coreos/go-semver/semver"
	"github.com/go-ini/ini"
	"github.com/pterm/pterm"
	"io/fs"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

type SystemType = int

const (
	Debian SystemType = iota
	RPM
	Other
	MacOS
	Windows
)

func ReadOSRelease(configfile string) (SystemType, error) {
	if _, err := os.Stat(configfile); errors.Is(err, fs.ErrNotExist) {
		return Other, nil
	}
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

var OSType = func() SystemType {
	switch runtime.GOOS {
	case "windows":
		return Windows
	case "darwin":
		return MacOS
	}
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

func GetLocalVersion() (*semver.Version, error) {
	_, err := exec.LookPath("pwsh")
	if err != nil {
		return nil, errors.New("PowerShell is not installed or not in PATH")
	}
	cmd := exec.Command("pwsh", "-v")
	b, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	up := strings.Replace(string(b), "PowerShell ", "v", -1)
	return semver.NewVersion(up)
}
