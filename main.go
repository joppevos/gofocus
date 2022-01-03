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
const configHosts = "hosts.json"

type HostsFile struct {
	IpAddress string      `json:"ip_address"`
	HostNames []HostNames `json:"host_names"`
}

type HostNames struct {
	HostName string `json:"host_name"`
}

// TODO: possibly no duplicate and simply append to a file.
// remove the appended lines when finished

func main() {
	go getCancelSignal()

	finish, err := getDurationArg()
	if err != nil {
		return
	}

	// Duplicate the current host file
	err = Copy(focusHosts, swapFocusHosts)
	if err != nil {
		panic(err)
	}
	//  put back old hosts file
	defer func() {
		err := Copy(swapFocusHosts, focusHosts)
		if err != nil {
			panic(err)
		}
		err = os.Remove(swapFocusHosts)
		if err != nil {
			panic(err)
		}
	}()

	HostFile, err := ReadJSON(configHosts)
	if err != nil {
		panic(err)
	}
	BlockedHosts := FormatHostFile(HostFile)
	// append blocked websites to focus hosts file
	err = AppendToFile(focusHosts, BlockedHosts)
	if err != nil {
		panic(err)
	}
	// flush cache to refresh hosts
	exec.Command("dscacheutil -flushcache\n")
	fmt.Printf("Go focus!\n")
	countDown(finish)
	fmt.Printf("\aGo take a break.\n") // \a is the bell system sound literal.
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

func ReadJSON(file string) (HostsFile, error) {
	plan, _ := ioutil.ReadFile(file)
	data := HostsFile{}
	err := json.Unmarshal(plan, &data)
	if err != nil {
		return HostsFile{}, fmt.Errorf("unable to read configuration file: '%s': %v", file, err)
	}
	return data, err

}

// getCancelSignal catch user input ctrl+c
// putting back the host file in its original state
func getCancelSignal() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan
	log.Println("Timer has been cancelled.")

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
