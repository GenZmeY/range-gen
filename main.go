package main

import (
	"bufio"
	"image"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"

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
	Version string = "dev"

	hashes map[string]string
	names  []string
)

func main() {
	closeHandler() // Ctrl+C

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

	if ArgJobs == 0 {
		ArgJobs = runtime.NumCPU()
	}
	runtime.GOMAXPROCS(ArgJobs)

	hashes, names := calcHashes(files)

	ranges := calcRanges(hashes, names, Threshold)

	writeRanges(ranges)

	wg := new(sync.WaitGroup)

	//wg.Add(1)

	wg.Wait()

	os.Exit(EXIT_SUCCESS)
}

func closeHandler() {
	interrupt := make(chan os.Signal, 2)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-interrupt
		output.Println("Closed by SIGINT")
		os.Exit(EXIT_SUCCESS)
	}()
}

func calcHashes(files []os.FileInfo) (map[string]string, []string) {
	var hashes map[string]string
	var keys []string
	var err error

	hashes = make(map[string]string)

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".png") {
			keys = append(keys, file.Name())
			hashes[file.Name()], err = calcHash(filepath.Join(ArgInput, file.Name()))

			if err != nil {
				output.Errorln(err)
			}
		}
	}
	sort.Strings(keys)
	return hashes, keys
}

func calcHash(filepath string) (string, error) {
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

func calcRanges(hashes map[string]string, names []string, Threshold int) string {
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

	return ranges.String()
}

func writeRanges(ranges string) {
	mode := os.FileMode(int(0644))
	targetFile, err := os.OpenFile(ArgOutput, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	defer targetFile.Close()
	if err == nil {
		datawriter := bufio.NewWriter(targetFile)
		_, err = datawriter.WriteString(ranges)
		if err == nil {
			err = datawriter.Flush()
		}
	}

	if err != nil {
		output.Errorln(err)
		os.Exit(EXIT_FILE_WRITE_ERR)
	}
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
