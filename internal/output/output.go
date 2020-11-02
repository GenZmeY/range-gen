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

package output

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
)

var (
	endOfLine string      = "\n"
	devNull   *log.Logger = log.New(ioutil.Discard, "", 0)
	stdout    *log.Logger = log.New(os.Stdout, "", 0)
	stderr    *log.Logger = log.New(os.Stderr, "", 0)
)

func init() {
	SetEndOfLineNative()
}

func SetQuiet(enabled bool) {
	if enabled {
		stdout = devNull
		stderr = devNull
	}
}

func SetEndOfLineNative() {
	switch os := runtime.GOOS; os {
	case "windows":
		setEndOfLineWindows()
	default:
		setEndOfLineUnix()
	}
}

func EOL() string {
	return endOfLine
}

func setEndOfLineUnix() {
	endOfLine = "\n"
}

func setEndOfLineWindows() {
	endOfLine = "\r\n"
}

func Print(v ...interface{}) {
	stdout.Print(v...)
}

func Printf(format string, v ...interface{}) {
	stdout.Printf(format, v...)
}

func Println(v ...interface{}) {
	stdout.Print(fmt.Sprint(v...) + endOfLine)
}

func Error(v ...interface{}) {
	stderr.Print(v...)
}

func Errorf(format string, v ...interface{}) {
	stderr.Printf(format, v...)
}

func Errorln(v ...interface{}) {
	stderr.Print(fmt.Sprint(v...) + endOfLine)
}
