package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"
)

const defaultDuration = 25 * time.Minute

func init() {
	flag.Parse()
}

func getDurationArg(start time.Time)(time.Time, error){
	arg := flag.Arg(0)
	if arg == "" {
		return start.Add(defaultDuration), nil
	}

	if n, err := strconv.Atoi(arg); err == nil {
		return start.Add(time.Duration(n)* time.Minute), nil
	}

	return time.Time{}, fmt.Errorf("could not parse time: %q", arg)
}