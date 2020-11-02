package main

import (
	"github.com/juju/gnuflag"

	"errors"
)

var (
	ArgInput          string
	ArgInputIsSet     bool = false
	ArgOutput         string
	ArgOutputIsSet    bool = false
	ArgThreshold      string
	ArgThresholdIsSet bool = false

	ArgJobs              int = 0
	ArgDefaultNoiseLevel int = 0
)

func Use(vals ...interface{}) {
	for _, val := range vals {
		_ = val
	}
}

func init() {
	gnuflag.IntVar(&ArgJobs, "jobs", 0, "")
	gnuflag.IntVar(&ArgJobs, "j", 0, "")
	gnuflag.IntVar(&ArgDefaultNoiseLevel, "noise", -1, "")
	gnuflag.IntVar(&ArgDefaultNoiseLevel, "n", -1, "")
}

func parseArgs() error {
	gnuflag.Parse(false)

	for i := 0; i < 3 && i < gnuflag.NArg(); i++ {
		switch i {
		case 0:
			ArgInput = gnuflag.Arg(0)
			ArgInputIsSet = true
		case 1:
			ArgOutput = gnuflag.Arg(1)
			ArgOutputIsSet = true
		case 2:
			ArgThreshold = gnuflag.Arg(2)
			ArgThresholdIsSet = true
		}
	}

	if !ArgInputIsSet {
		return errors.New("Input directory not specified")
	}

	if !ArgOutputIsSet {
		return errors.New("Output file not specified")
	}

	if !ArgThresholdIsSet {
		return errors.New("Threshold not specified")
	}

	return nil
}
