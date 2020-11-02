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

	"range-gen/internal/output"
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

	allFiles, err := ioutil.ReadDir(ArgInput)
	if err != nil {
		output.Errorln("Read dir error")
		os.Exit(EXIT_DIR_READ_ERR)
	}

	if ArgJobs == 0 {
		ArgJobs = runtime.NumCPU()
	}
	runtime.GOMAXPROCS(ArgJobs)

	pngFiles := pngList(allFiles)
	hashes := calcHashes(pngFiles)
	ranges := calcRanges(hashes, Threshold)
	writeRanges(ranges)

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

func pngList(files []os.FileInfo) []os.FileInfo {
	var pngFiles []os.FileInfo

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".png") {
			pngFiles = append(pngFiles, file)
		}
	}

	return pngFiles
}

func calcHashes(files []os.FileInfo) map[string]string {
	var (
		hashes          map[string]string = make(map[string]string)
		jobChan         chan int          = make(chan int, ArgJobs)
		mutex           sync.Mutex
		activeJobs      int = 0
		index           int = 0
		complete        int = 0
		size            int = len(files)
		progress        int = 0
		currentProgress int = 0
	)

	for complete < size {
		for activeJobs < ArgJobs && index < size {
			go calcHash(ArgInput, files[index], &hashes, &mutex, jobChan)
			activeJobs++
			index++
		}
		<-jobChan
		activeJobs--
		complete++
		currentProgress = complete * 100 / size
		if currentProgress != progress {
			output.Print("\r[" + strconv.Itoa(currentProgress) + "%] " + strconv.Itoa(complete) + "/" + strconv.Itoa(size))
			progress = currentProgress
		}
	}

	return hashes
}

func calcHash(path string, fileinfo os.FileInfo, hashes *map[string]string, mutex *sync.Mutex, jobChan chan int) {
	file, err := os.Open(filepath.Join(path, fileinfo.Name()))
	if err != nil {
		output.Errorln(err)
	}
	defer file.Close()

	image, _, err := image.Decode(file)
	if err != nil {
		output.Errorln(err)
	}
	hash := blockhash.NewBlockhash(image, 16).Hexdigest()
	mutex.Lock()
	(*hashes)[fileinfo.Name()] = hash
	mutex.Unlock()
	jobChan <- 1
}

func calcRanges(hashes map[string]string, Threshold int) string {
	var ranges strings.Builder
	var prevHash string = ""
	var dist int = 0
	var startName = ""
	var rangeNoise string = ""

	if ArgDefaultNoiseLevel >= 0 {
		rangeNoise = "\t" + string(ArgDefaultNoiseLevel)
	}

	names := []string{}
	for key := range hashes {
		names = append(names, key)
	}
	sort.Strings(names)

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
