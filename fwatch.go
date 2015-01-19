// This is a simple tool which watches a file for inactivity of a given
// time. If the file does not change for the given time, this process stops
// or executes an optional command.
package main

import (
	"fmt"
	"log"

	"os"

	"os/exec"
	"strconv"
	"time"

	"gopkg.in/fsnotify.v1"
)

// Simply execute the given command with the given args and print the output
// to the output of this process.
func execCommand(cmd string, args []string) error {
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		return err
	}
	return nil
}

func main() {

	if len(os.Args) < 4 {
		fmt.Println("usage: fwatch <file> <inactivetime> <retries> [<command>] ...")
		fmt.Println("retries=0 -> watch forever, but execute command on every interval")
		fmt.Println("If your nginx must be called every 10secs with a healthcheck you can")
		fmt.Println("execute a command if there are 3 drops:")
		fmt.Println("   fwatch /var/log/nginx/access.log 10s 3 echo \"no healthcheck\"")
		os.Exit(1)
	}

	filename := os.Args[1]
	dur, err1 := time.ParseDuration(os.Args[2])
	if err1 != nil {
		log.Fatal(err1)
	}
	retries, _ := strconv.Atoi(os.Args[3])
	var commando *string
	if len(os.Args) > 4 {
		commando = &os.Args[4]
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	err = watcher.Add(filename)
	if err != nil {
		log.Fatal(err)
	}

	t := time.NewTicker(dur)
	for {
		select {
		case <-t.C:
			if commando != nil {
				if cerr := execCommand(*commando, os.Args[5:]); cerr != nil {
					log.Fatal(cerr)
				}
			}
			// wenn retries == 0 do it forevery, otherwise only count down
			if retries > 0 {
				retries = retries - 1
				if retries == 0 {
					os.Exit(0)
				}

			}
		case ev := <-watcher.Events:
			if isModify(ev) {
				t.Stop()
				t = time.NewTicker(dur)
			}
		case err := <-watcher.Errors:
			log.Fatal("error:", err)
		}
	}

}

// Check if the target file was modified.
func isModify(e fsnotify.Event) bool {
	return e.Op&fsnotify.Write == fsnotify.Write ||
		e.Op&fsnotify.Rename == fsnotify.Rename ||
		e.Op&fsnotify.Remove == fsnotify.Remove ||
		e.Op&fsnotify.Chmod == fsnotify.Chmod
}
