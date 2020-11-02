package main

import (
	"bufio"
	"image"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"

	_ "image/png"

	"github.com/dsoprea/go-perceptualhash"

	"range-gen/output"
)

const (
	EXIT_SUCCESS        int = 0
	EXIT_ARG_ERR        int = 1
	EXIT_FILE_READ_ERR  int = 2
	EXIT_DIR_READ_ERR   int = 3
	EXIT_FILE_WRITE_ERR int = 4
)

var (
	hashes map[string]string
	names  []string
)

func main() {
	output.SetEndOfLineNative()

	if err := parseArgs(); err != nil {
		output.Errorln(err)
		os.Exit(EXIT_ARG_ERR)
	}

	Threshold, err := strconv.Atoi(ArgThreshold)
	if err != nil {
		output.Errorln("Can't convert threshold to integer")
		os.Exit(EXIT_ARG_ERR)
	}

	files, err := ioutil.ReadDir(ArgInput)
	if err != nil {
		output.Errorln("Read dir error")
		os.Exit(EXIT_DIR_READ_ERR)
	}

	hashes = make(map[string]string)
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".png") {
			names = append(names, file.Name())
			hashes[file.Name()], err = calchash(ArgInput + file.Name())

			if err != nil {
				output.Errorln(err)
			}
		}
	}
	sort.Strings(names)

	var ranges strings.Builder
	var prevHash string = ""
	var dist int = 0
	var startName = ""
	var rangeNoise string = ""

	if ArgDefaultNoiseLevel >= 0 {
		rangeNoise = "\t" + string(ArgDefaultNoiseLevel)
	}

	for i := 0; i < len(names); i++ {
		name := names[i]
		if startName == "" {
			startName = name
		}
		dist = hammingDistance(prevHash, hashes[name])
		if dist >= Threshold {
			ranges.WriteString(startName + "\t" + names[i-1] + rangeNoise + output.EOL())
			startName = name
		}
		prevHash = hashes[name]
	}
	ranges.WriteString(startName + "\t" + names[len(names)-1] + rangeNoise + output.EOL())

	mode := os.FileMode(int(0644))
	targetFile, err := os.OpenFile(ArgOutput, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	defer targetFile.Close()
	if err == nil {
		datawriter := bufio.NewWriter(targetFile)
		_, err = datawriter.WriteString(ranges.String())
		if err == nil {
			err = datawriter.Flush()
		}
	}

	if err != nil {
		output.Errorln(err)
		os.Exit(EXIT_FILE_WRITE_ERR)
	}

	os.Exit(EXIT_SUCCESS)
}

func calchash(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	image, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	return blockhash.NewBlockhash(image, 16).Hexdigest(), nil
}

func hammingDistance(prev, cur string) int {
	var dist int = 0
	var p, c int64

	if prev == "" || cur == "" {
		return 0
	}

	for i := 0; i < len(cur); i++ {
		p, _ = strconv.ParseInt(string([]rune(prev)[i]), 16, 64)
		c, _ = strconv.ParseInt(string([]rune(cur)[i]), 16, 64)
		if p > c {
			dist += int(p - c)
		} else {
			dist += int(c - p)
		}
	}

	return dist
}
