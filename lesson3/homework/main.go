package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type Options struct {
	From      string
	To        string
	Offset    uint
	Limit     uint
	Conv      string
	BlockSize uint
}

func ParseFlags() (*Options, error) {
	var opts Options

	flag.StringVar(&opts.From, "from", "", "file to read. by default - stdin")
	flag.StringVar(&opts.To, "to", "", "file to write. by default - stdout")
	flag.UintVar(&opts.Offset, "offset", 0, "bytes inside the input to be skipped when copying")
	flag.UintVar(&opts.Limit, "limit", 0, "maximum number of bytes to read")
	flag.StringVar(&opts.Conv, "conv", "", "do conv")
	flag.UintVar(&opts.BlockSize, "block-size", 0, "maximum number of bytes to read")
	flag.Parse()

	return &opts, nil
}

func main() {
	opts, err := ParseFlags()

	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "can not parse flags:", err)
		os.Exit(1)
	}
	if strings.Contains(opts.Conv, "lower_case") && strings.Contains(opts.Conv, "upper_case") {
		_, _ = fmt.Fprintln(os.Stderr, "error conv", err)
		os.Exit(1)
	}

	var buf []byte

	if opts.From == "" {
		if opts.Limit > 0 {
			lr := io.LimitReader(os.Stdin, int64(opts.Limit+opts.Offset))
			buf, _ = io.ReadAll(lr)
		} else {
			buf, _ = io.ReadAll(os.Stdin)
		}
	} else {
		r, err := os.Open(opts.From)

		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "error", err)
			os.Exit(1)
		}
		buf, _ = io.ReadAll(r)
		err = r.Close()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "error", err)
			os.Exit(1)
		}
	}

	if int(opts.Offset) > len(buf) {
		_, _ = fmt.Fprintln(os.Stderr, "Offset more file", err)
		os.Exit(1)
	}

	if opts.Offset != 0 {
		buf = buf[opts.Offset:]
	}

	str := string(buf)

	if opts.Conv != "" {
		if !strings.Contains(opts.Conv, "upper_case") && !strings.Contains(opts.Conv, "lower_case") && !strings.Contains(opts.Conv, "trim_spaces") {
			_, _ = fmt.Fprintln(os.Stderr, "unknown operation", err)
			os.Exit(1)
		}
		if strings.Contains(opts.Conv, "trim_spaces") {
			str = strings.TrimSpace(str)
		}
		if strings.Contains(opts.Conv, "lower_case") {
			str = strings.ToLower(str)
		}
		if strings.Contains(opts.Conv, "upper_case") {
			str = strings.ToUpper(str)
		}
	}

	if opts.To != "" {
		_, err := os.Stat(opts.To)
		if !os.IsNotExist(err) {
			_, _ = fmt.Fprintln(os.Stderr, "file exist", err)
			os.Exit(1)
		}
		to, err := os.Create(opts.To)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "error", err)
			os.Exit(1)
		}
		_, err = io.WriteString(to, str)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "error", err)
			os.Exit(1)
		}
		err = to.Close()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "error", err)
			os.Exit(1)
		}
	} else {
		if _, err = fmt.Fprint(os.Stdout, str); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "error", err)
			os.Exit(1)
		}
	}
}
