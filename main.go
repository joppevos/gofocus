package main

import (
	"encoding/json"
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

const focusHosts = "/etc/hosts"
const swapFocusHosts = "/etc/hosts_swap"
const configHosts = "/Users/vosjoppe/github/gofocus/hostsfile.json" // TODO

type HostsFile struct {
	IpAddress string      `json:"ip_address"`
	HostNames []HostNames `json:"host_names"`
}

type HostNames struct {
	HostName string `json:"host_name"`
}

func main() {
	go getCancelSignal()

	finish, err := getDurationArg()
	if err != nil {
		return
	}
	fmt.Printf("Go focus!\n")

	// Duplicate the current host file
	err = Copy(focusHosts, swapFocusHosts)
	if err != nil {
		panic(err)
	}

	// remove swap file
	defer func() {
		err := os.Remove(swapFocusHosts)
		if err != nil {
			panic(err)
		}
	}()

	// append blocked websites to focus hosts file
	BlockedHosts := FormatHostFile(readJSON(configHosts))
	err = AppendToFile(focusHosts, BlockedHosts)

	// flush cache to refresh hosts
	exec.Command("dscacheutil -flushcache\n")
	if err != nil {
		panic(err)
	}

	//  put back old hosts file
	defer func() {
		err := Copy(swapFocusHosts, focusHosts)
		if err != nil {
			panic(err)
		}
	}()

	countDown(finish)

	fmt.Printf("\aGo take a break") // \a is the bell system sound literal.
}

func formatMinutes(t time.Duration) string {
	minutes := int(t.Minutes())
	seconds := int(t.Seconds()) % 60

	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

func countDown(duration time.Time) {
	for range time.Tick(1 * time.Second) {
		timeRemaining := -time.Since(duration)
		if timeRemaining <= 0 {
			break
		}
		_, err := fmt.Fprint(os.Stdout, "Countdown: ", formatMinutes(timeRemaining), "   \r")
		if err != nil {
			panic(err)
		}
		err = os.Stdout.Sync()
		if err != nil {
			panic(err)
		}
	}
}

// Copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func Copy(src, target string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(target)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	err = os.Chmod(target, 0777)
	if err != nil {
		log.Fatal(err)
	}

	return out.Close()
}

func AppendToFile(file string, data string) error {

	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(data)
	if err != nil {
		return err
	}
	return nil
}

func FormatHostFile(file HostsFile) string {
	builder := strings.Builder{}
	builder.WriteString("\n") // start on a newline
	for _, e := range file.HostNames {
		_, err := builder.WriteString(fmt.Sprintf("%s %s\n", file.IpAddress, e.HostName))
		if err != nil {
			panic(err)
		}
	}
	return builder.String()
}

func readJSON(file string) HostsFile {
	plan, _ := ioutil.ReadFile(file)
	data := HostsFile{}
	err := json.Unmarshal(plan, &data)
	if err != nil {
		panic(err)
	}
	return data

}

// getCancelSignal catch user input ctrl+c
// putting back the host file in its original state
func getCancelSignal() {
	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan
	log.Println("Program killed !")

	err := Copy(swapFocusHosts, focusHosts)
	if err != nil {
		panic(err)
	}

	err = os.Remove(swapFocusHosts)
	if err != nil {
		panic(err)
	}
	os.Exit(0)
}
