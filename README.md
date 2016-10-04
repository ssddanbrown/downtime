# Downtime

A really basic, simple Go app to log http downtime. Tracks the time it cannot make a request for.

### Usage

Either run the script directly with go (`go run downtime.go http://example.com`) or compile it (with `go install` or `go build`) and run the binary.

``` bash
# Command format
downtime <options> <url>

# Example
downtime -p 1 -f output.log https://google.com
```

#### Options

* `-f output.txt` - Specify the output log. Outputs to stdout if ommitted.
* `-p 1` - Specify the ping frequency in seconds.
