package pkg

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"fmt"
	"github.com/jinzhu/copier"
	ar "github.com/mkrautz/goar"
	"github.com/pterm/pterm"
	"github.com/smira/go-xz"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
)

type MD5Sums = map[string]string

func indirect(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

type MessageType = int

const (
	Stdout MessageType = iota
	Stderr
)

type Message struct {
	typ MessageType
	msg string
}

func Install(path string) error {
	name := filepath.Base(path)
	sp, err := pterm.DefaultSpinner.Start(fmt.Sprintf("Installing %s", name))
	if err != nil {
		return err
	}
	cmd := exec.Command("dpkg", "-i", path)
	msgChan := make(chan Message)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		msgType := Stdout
		for scanner.Scan() {
			msg := Message{typ: msgType, msg: scanner.Text()}
			msgChan <- msg
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stderr)
		msgType := Stderr
		for scanner.Scan() {
			msg := Message{typ: msgType, msg: scanner.Text()}
			msgChan <- msg
		}
	}()
	err = cmd.Start()
	if err != nil {
		return err
	}
	go func() {
		pterm.Error.PrintOnError(cmd.Wait())
		wg.Wait()
		close(msgChan)
	}()
	err = cmd.Start()
	if err != nil {
		return err
	}
	for v := range msgChan {
		switch v.typ {
		case Stdout:
			sp.UpdateText(v.msg)
		case Stderr:
			pterm.Error.Println(v.msg)

		}
	}
	return nil
}

func archiveLen(in io.Reader) (int, error) {
	res := 0
	archive := tar.NewReader(in)
	for {
		hdr, err := archive.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return 0, err
		}
		if hdr.Name == "./" {
			continue
		}
		res += 1
	}
	return res, nil
}

func UnpackDeb() error {
	file, err := os.Open("test_data/test.deb")
	pterm.Error.PrintOnError(err)
	if err != nil {
		return err
	}
	arc := ar.NewReader(file)
	for {
		header, err := arc.Next()
		if err == io.EOF {
			break
		}
		pterm.Error.PrintOnError(err)
		if strings.HasPrefix(header.Name, "control.tar") {
			bufReader := bufio.NewReader(arc)
			var tarInput io.Reader
			switch header.Name {
			case "control.tar":
				tarInput = bufReader
			case "control.tar.gz":
				ungzip, err := gzip.NewReader(bufReader)
				pterm.Fatal.PrintOnError(err)
				defer ungzip.Close()
				tarInput = ungzip
			case "control.tar.xz":
				unxz, err := xz.NewReader(bufReader)
				pterm.Fatal.PrintOnError(err)
				defer unxz.Close()
				tarInput = unxz

			}
			var toCount tar.Reader
			untar := tar.NewReader(tarInput)
			pterm.Fatal.PrintOnError(copier.CopyWithOption(&toCount, &untar, copier.Option{IgnoreEmpty: false, DeepCopy: true}))
			total, err := archiveLen(&toCount)
			pterm.Fatal.PrintOnError(err)
			if err != nil {
				return err
			}
			pb, err := pterm.DefaultProgressbar.WithTotal(total).WithTitle(fmt.Sprintf("Unpacking %s", header.Name)).Start()
			pterm.Fatal.PrintOnError(err)
			if err != nil {
				return err
			}
			for {
				hdr, err := untar.Next()
				if err == io.EOF {
					break
				}
				pterm.Fatal.PrintOnError(err)
				if hdr.Name == "./" {
					continue
				}
				pb.UpdateTitle("Unpacking " + hdr.Name)
				unpackPath := filepath.Join("./test_data", hdr.Name)
				f, err := os.Create(unpackPath)
				pterm.Fatal.PrintOnError(err)
				if _, err = io.Copy(f, untar); err != nil {
					pterm.Fatal.PrintOnError(err)
				}
				pterm.Success.Println(fmt.Sprintf("Unpacked %s", hdr.Name))
				pb.Increment()
			}
			_, err = pb.Stop()
			pterm.Fatal.PrintOnError(err)
		}

	}
	return err
}
