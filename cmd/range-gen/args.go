/*/
 * range-gen creates a list of scene ranges based on a set of frames from the video.
 * Copyright (C) 2020 GenZmeY
 * mailto: genzmey@gmail.com
 *
 * This file is part of range-gen.
 *
 * range-gen is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
/*/

package main

import (
	"github.com/juju/gnuflag"

	"range-gen/internal/output"

	"errors"
	"os"
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

	ArgVersion bool = false
	ArgHelp    bool = false
)

func printHelp() {
	output.Println("Сreates a list of scene ranges based on a set of frames from the video")
	output.Println("")
	output.Println("Usage: range-gen [option]... <input_dir> <output_file> <threshold>")
	output.Println("input_dir          Directory with png images")
	output.Println("output_file        Range list file")
	output.Println("threshold          Image similarity threshold (0-1024)")
	output.Println("")
	output.Println("Options:")
	output.Println("  -j, --jobs N     Allow N jobs at once")
	output.Println("  -n, --noise      Default noise level for each range")
	output.Println("  -h, --help       Show this page")
	output.Println("  -v, --version    Show version")
}

func printVersion() {
	output.Println("range-gen ", Version)
}

func init() {
	gnuflag.IntVar(&ArgJobs, "jobs", 0, "")
	gnuflag.IntVar(&ArgJobs, "j", 0, "")
	gnuflag.IntVar(&ArgDefaultNoiseLevel, "noise", -1, "")
	gnuflag.IntVar(&ArgDefaultNoiseLevel, "n", -1, "")
	gnuflag.BoolVar(&ArgVersion, "version", false, "")
	gnuflag.BoolVar(&ArgVersion, "v", false, "")
	gnuflag.BoolVar(&ArgHelp, "help", false, "")
	gnuflag.BoolVar(&ArgHelp, "h", false, "")
}

func parseArgs() error {
	gnuflag.Parse(false)

	switch {
	case ArgHelp:
		printHelp()
		os.Exit(EXIT_SUCCESS)
	case ArgVersion:
		printVersion()
		os.Exit(EXIT_SUCCESS)
	}

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
