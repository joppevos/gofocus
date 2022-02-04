package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"
)

const defaultDuration = 8 * time.Second

func getDurationArg() (time.Duration, error) {
	flag.Parse()
	arg := flag.Arg(0)
	if arg == "" {
		return defaultDuration, nil
	}

	if n, err := strconv.Atoi(arg); err == nil {
		return time.Duration(n) * time.Minute, nil
	}

	return 0, fmt.Errorf("could not parse time: %q", arg)
}
