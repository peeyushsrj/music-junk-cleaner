![Build Status](https://circleci.com/gh/peeyushsrj/ui-music-junk-cleaner.png?style=shield)

# Music-junk-cleaner

It scans mp3 files in the current directory, and then recursively into subdirectories and it prompts a potential junk mp3 files. Once a junk is marked, it will handle same similar found patterns and delete the junk data from such files. 

[![Screenshot](https://s20.postimg.org/lrrw1qf65/Screenshot_from_2017-06-23_19-27-01.png)](https://postimg.org/image/66akhs37t/)

For library purpose or commandline usage, try [core-music-junk-cleaner](https://github.com/peeyushsrj/core-music-junk-cleaner).

## Installation

[Download](https://github.com/peeyushsrj/ui-music-junk-cleaner/releases) the binaries for your operating system, and run it in your music directory! Now open [http://localhost:7899/](http://localhost:7899/) in browser. Prompt-Mark-Clean!


## Build Instructions

- Install Go.
- Installing dependent packages.  
```
	go get github.com/gorilla/websocket
	github.com/skratchdot/open-golang/open
```
- Clone it `git clone https://github.com/peeyushsrj/ui-music-junk-cleaner/`
- Changed directory & build it `go build`.


## Future TODO

- [x] Launching a browser.
- [x] Setting up CI for X platform Binaries.
- [x] Run in current directory.
- [ ] Support other music formats.

## License

The MIT License (MIT) Copyright (c) 2017 Peeyush Singh