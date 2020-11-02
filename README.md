# Range-Gen

![build release](https://github.com/GenZmeY/range-gen/workflows/build%20release/badge.svg)
[![GitHub top language](https://img.shields.io/github/languages/top/GenZmeY/range-gen)](https://golang.org)
[![GitHub](https://img.shields.io/github/license/genzmey/range-gen)](https://www.gnu.org/licenses/gpl-3.0.en.html)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/GenZmeY/range-gen)](https://github.com/GenZmeY/range-gen/releases)

*range-gen creates a list of scene ranges based on a set of frames from the video. This program is used in the [video2d-2x](https://github.com/GenZmeY/video2d-2x) project.*

***

# Build & Install
**Note:** You can get the compiled version for your platform on the [release page](https://github.com/GenZmeY/range-gen/releases).

Dependencies:  
- linux distro
- git
- golang 1.13
- make  

Get the source:  
`git clone https://github.com/GenZmeY/range-gen`  

Build:  
`cd range-gen && make`  

Install:  
`make install`  

**Build versions for all plaforms:**  
`make -j $(nproc) compile`  
(executables will be in `range-gen/bin` folder)

# Usage
```
Usage: range-gen [option]... <input_dir> <output_file> <threshold>
input_dir          Directory with png images
output_file        Range list file
threshold          Image similarity threshold (0-1024)

Options:
  -j, --jobs N     Allow N jobs at once
  -n, --noise      Default noise level for each range (0-3)
  -h, --help       Show this page
  -v, --version    Show version
```

# License
Range-gen is licensed under the [GNU GPLv3](https://www.gnu.org/licenses/gpl-3.0.en.html), but uses a [go-perceptualhash](https://github.com/dsoprea/go-perceptualhash) ([BSD 3-Clause License](https://github.com/dsoprea/go-perceptualhash/blob/master/LICENSE)) to calculate hashes of images.
