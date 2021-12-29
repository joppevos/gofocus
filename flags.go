package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"
)

const defaultDuration = 8 * time.Second

func getDurationArg() (time.Time, error) {
	flag.Parse()
	start := time.Now()
	arg := flag.Arg(0)
	if arg == "" {
		return start.Add(defaultDuration), nil
	}

	if n, err := strconv.Atoi(arg); err == nil {
		return start.Add(time.Duration(n) * time.Minute), nil
	}

	return time.Time{}, fmt.Errorf("could not parse time: %q", arg)
}
