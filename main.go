package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"
)

const HostsFilePath = "/etc/hosts"
const BackupHostsFilePath = "/tmp/hostsBackup"
const HostsConfigFilePath = "hosts"
const IPAddress = "127.0.0.1"

func main() {
	go getCancelSignal()

	durationArg, err := getDurationArg()
	if err != nil {
		return
	}

	// Duplicate the current host file.
	err = Copy(HostsFilePath, BackupHostsFilePath)
	if err != nil {
		log.Fatal(err)
	}

	defer rollBack()

	HostConfig := openAndRead(HostsConfigFilePath)
	BlockedHosts := formatHostsConfig(HostConfig)

	// Append blocked websites to focus hosts file.
	err = appendToHostsFile(HostsFilePath, BlockedHosts)
	if err != nil {
		log.Fatal(err)
	}
	flushCache()

	fmt.Printf("Go focus!\n")
	countDown(os.Stdout, durationArg)
	fmt.Printf("\aGo take a break.\n") // \a is the bell system sound literal.
}

// rollBack will place back the original hosts file,
// remove the backup and flush the cache.
func rollBack() {
	err := Copy(BackupHostsFilePath, HostsFilePath)
	if err != nil {
		log.Fatal(err)
	}
	err = os.Remove(BackupHostsFilePath)
	if err != nil {
		log.Fatal(err)
	}
	flushCache()
}

// flushCache on UNIX to refresh the hosts file.
func flushCache() {
	exec.Command("dscacheutil",  "-flushcache\n")
}

func formatMinutes(t time.Duration) string {
	minutes := int(t.Minutes())
	seconds := int(t.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func countDown(w io.Writer, duration time.Duration) {
	start := time.Now()
	c := start.Add(duration)
	for range time.Tick(1 * time.Second) {
		timeRemaining := -time.Since(c)
		_, err := fmt.Fprint(w, "", formatMinutes(timeRemaining), "   \r")
		if err != nil {
			panic(err)
		}
		if timeRemaining <= 0 {
			break
		}

	}
}

// Copy a source file to destination. Any existing file will be overwritten and will
// not copy file attributes.
func Copy(src, target string) error {

	in, err := os.Open(src)
	if err != nil {
		log.Println("unable to open source file")
		return err
	}
	defer in.Close()

	out, err := os.Create(target)
	if err != nil {
		log.Println("unable to create target file")
		return err
	}

	_, err = io.Copy(out, in)
	if err != nil {
		log.Println("unable to copy back original hosts file")
		return err
	}

	err = out.Close()
	if err != nil {
		return err

	}
	return nil
}

func appendToHostsFile(name string, data string) error {
	f, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(data)
	content, _ := ioutil.ReadFile(name)
	log.Println(string(content))
	if err != nil {
		return err
	}
	return nil
}

// formatHostsConfig creates a string block to append to /etc/hosts
func formatHostsConfig(HostUrls []string) string {
	builder := strings.Builder{}
	builder.WriteString("\n") // start on a newline
	for _, e := range HostUrls {
		_, err := builder.WriteString(fmt.Sprintf("%s %s\n", IPAddress, e))
		if err != nil {
			panic(err)
		}
	}
	return builder.String()
}

// openAndRead will read a file and return the content separated in a slice.
func openAndRead(name string) []string {
	content, err := ioutil.ReadFile(name)
	if err != nil {
		log.Fatalf("could not read file: %s", err)

	}
	return splitContent(content)
}

func splitContent(content []byte) []string {
	return strings.Split(string(content), "\n")

}

// getCancelSignal catch user input ctrl+c
// putting back the hosts file in its original state
func getCancelSignal() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("Are you sure?")
	<-quit
	log.Println("Timer has been cancelled.")

	err := Copy(BackupHostsFilePath, HostsFilePath)
	if err != nil {
		panic(err)
	}

	err = os.Remove(BackupHostsFilePath)
	if err != nil {
		panic(err)
	}
	flushCache()
	os.Exit(0)
}
