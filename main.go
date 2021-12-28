package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

const fileName = "hosts_duplicates"
const hostsFilePath = "/etc/hosts"

func main() {
	start := time.Now()
	finish, err := getDurationArg(start)
	if err != nil {
		return
	}

	fmt.Printf("Go focus!\n")

	// Duplicate the current hostfile
	if _, err := os.Stat(fileName); os.IsNotExist(err){
		err := Copy(hostsFilePath, fileName)
		if err != nil {
			panic(err)
		}
		fmt.Println("First Add the websites you want to block to:",fileName)
		return
	}

	// write it back
	err = Copy(fileName, hostsFilePath)
	if err != nil {
		panic(err)
	}

	countDown(finish)

	fmt.Println("\a") // \a is the bell literal.
	fmt.Printf("Go take a break")


}
func formatMinutes(t time.Duration) string{
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
	return out.Close()
}