package main

import (
	"fmt"
	"time"
)


func main() {

	start := time.Now()
	finish, err := getDurationArg(start)
	if err != nil {
		return
	}
	fmt.Printf("Go focus for %s!\n", finish.Sub(start))
	countDown(finish)
	fmt.Println("\a") // \a is the bell literal.
	fmt.Printf("Go take a break")


}

func countDown(duration time.Time) {

	for range time.Tick(1 * time.Second) {
		timeRemaining := -time.Since(duration)

		if timeRemaining <= 0 {
			break
		}
	}
}