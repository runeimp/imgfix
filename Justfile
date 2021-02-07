
alias _build-mac := _build-macos
alias _build-win := _build-windows

@_default:
	cls
	just --list


# Build linux, macos, or windows app (default: current OS)
build os='':
	#!/bin/sh
	cls
	if [ "{{os}}" = '' ]; then
		os='{{os()}}'
	else
		os='{{os}}'
	fi
	just "_build-${os}"
	ls -al "bin/${os}/"

# Build macOS app
_build-linux:
	rm -rf bin/linux
	mkdir -p bin/linux
	GOOS=linux GOARCH=amd64 go build -o 'bin/linux/imgfix' ./...

# Build macOS app
_build-macos:
	rm -rf bin/macos
	mkdir -p bin/macos
	GOOS=darwin GOARCH=amd64 go build -o 'bin/macos/imgfix' ./...

# Build Windows app
_build-windows:
	rm -rf bin/windows
	mkdir -p bin/windows
	GOOS=windows GOARCH=amd64 go build -o 'bin/windows/imgfix.exe' ./...


# Install macOS app
install:
	rm -rf bin/{{os()}}
	mkdir -p bin/{{os()}}
	GOARCH=amd64 go build -o 'bin/{{os()}}/imgfix' ./...
	cp 'bin/{{os()}}/imgfix' "${GOBIN}/"


# Run the app
run +args='-dry-run -verbose is-dir/* is-*':
	@cls
	@just _file-setup
	go run ./... {{args}}


@_file-setup:
	# mv is-empty-text.txt is-empty-text.jpeg
	mv is-gif-animated.gif is-gif-animated.jpg 2>/dev/null; exit 0
	# mv is-jpeg-1.jpg is-jpeg-1.jpeg
	cp is-jpeg-2.jpg is-jpeg-2.png 2>/dev/null; exit 0
	mv is-png.png is-png 2>/dev/null; exit 0
	mv is-png.png is-png.gif 2>/dev/null; exit 0
