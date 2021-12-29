package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

const focusHosts = "/Users/vosjoppe/github/gofocus/hosts"
const swapFocusHosts = "/Users/vosjoppe/github/gofocus/hosts_swap"
const configHosts = "/Users/vosjoppe/github/gofocus/hosts_swap"

type HostsFile struct {
	IpAddress string      `json:"ip_address"`
	HostNames []HostNames `json:"host_names"`
}

type HostNames struct {
	HostName string `json:"host_name"`
}

func main() {
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

	// append blocked websites to focus hosts
	BlockedHosts := FormatHostFile(readJSON(configHosts))
	err = AppendToFile(focusHosts,BlockedHosts)
	if err != nil {
		panic(err)
	}

	// on any error, make sure we put back old hosts file
	defer func() {
		err := Copy(swapFocusHosts, focusHosts)
		if err != nil {
			panic(err)
		}
	}()

	countDown(finish)
	fmt.Println("\a") // \a is the bell literal.
	fmt.Printf("Go take a break")

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

		fmt.Fprint(os.Stdout, "Countdown: ", formatMinutes(timeRemaining), "   \r")
		os.Stdout.Sync()
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
	defer out.Close()

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

func FormatHostFile(file HostsFile) string{
	builder := strings.Builder{}
	for _, e := range file.HostNames{
		_, err := builder.WriteString(fmt.Sprintf("%s %s\n", file.IpAddress, e.HostName))
		if err != nil{
			panic(err)
		}
	}
	return builder.String()


}

func readJSON(file string) HostsFile{
	plan, _ := ioutil.ReadFile(file)
	data := HostsFile{}
	err := json.Unmarshal(plan, &data)
	if err != nil {
		panic(err)
	}
	return data

}