package main

import (
	"os/exec"
	"time"
)

func readAdb(out chan<- string) error {
	c := exec.Command("adb", "logcat")
	//c := exec.Command("ping", "www.google.com")
	o, err := c.StdoutPipe()
	if err != nil {
		return err
	}

	if err := c.Start(); err != nil {
		return err
	}

	go func() {
		b := make([]byte, 1, 1)
		var currentLine string
		for {
			n, err := o.Read(b)
			crashIf(err)

			if n == 0 {
				time.Sleep(10 * time.Millisecond)
				continue
			}
			if b[0] == '\n' {
				out <- currentLine
				currentLine = ""
			} else {
				currentLine += string(b[0])
			}
		}
	}()
	go func() {
		crashIf(c.Wait())
	}()

	return nil
}
