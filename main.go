package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	fromStrPtr := flag.String("f", "", "milliseconds to convert to ISO 8601")
	toStrPtr := flag.String("t", "", "ISO 8601 date to convert")
	flag.Parse()
	
	switch {
	case *fromStrPtr == "" && *toStrPtr == "":
		exit(errors.New("no date given"))
	case *fromStrPtr != "" && *toStrPtr != "":
		exit(errors.New("either -t or -f can be specified not both"))
	case *fromStrPtr != "":
		run(*fromStrPtr, func(s string) (interface{}, error) {
			millis, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return nil, err
			}
			unixTime := time.Unix(0, millis*int64(time.Millisecond))
			return unixTime.Format(time.RFC3339Nano), nil
		})
	case *toStrPtr != "":
		run(*toStrPtr, func(s string) (interface{}, error) {
			t, err := time.Parse(time.RFC3339Nano, s)
			if err != nil {
				return nil, err
			}
			return t.UnixNano() / int64(time.Millisecond), nil
		})
	}
}

func run(in string, conv func(string) (interface{}, error)) {
	arg, err := readArg(in)
	if err != nil {
		exit(err)
	}
	out, err := conv(arg)
	if err != nil {
		exit(err)
	}
	fmt.Println(out)
	exit(nil)
}

func readArg(arg string) (string, error) {
	// we interpret - as indication to read from std in
	if arg != "-" {
		return arg, nil
	}
	// assume pipe
	reader := bufio.NewReader(os.Stdin)
	s, err := reader.ReadString('\n')
	if err != nil {
		return s, err
	}
	return strings.TrimSuffix(s, "\n"), nil
}

func exit(err error) {
	if err != nil {
		println(err.Error())
		flag.Usage()
		os.Exit(1)
	}
	os.Exit(0)
}
