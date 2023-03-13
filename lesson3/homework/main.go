package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type Options struct {
	From   string
	To     string
	Offset uint
	Limit  uint
	Conv   string
}

func ParseFlags() (*Options, error) {
	var opts Options

	flag.StringVar(&opts.From, "from", "", "file to read. by default - stdin")
	flag.StringVar(&opts.To, "to", "", "file to write. by default - stdout")
	flag.UintVar(&opts.Offset, "offset", 0, "bytes inside the input to be skipped when copying")
	flag.UintVar(&opts.Limit, "limit", 0, "maximum number of bytes to read")
	flag.StringVar(&opts.Conv, "conv", "", "do conv")

	flag.Parse()

	return &opts, nil
}

func main() {
	opts, err := ParseFlags()

	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "can not parse flags:", err)
		os.Exit(1)
	}
	var buf []byte
	if opts.From == "" {
		/*if opts.Limit == 0 && opts.Offset == 0 {
			r := io.Reader(os.Stdin)
			buf, _ = io.ReadAll(r)
		} else {
			r := io.Reader(os.Stdin)
			lr := io.LimitReader(r, int64(opts.Limit+opts.Offset))
			buf, _ = io.ReadAll(lr)
		}*/
		if opts.Limit > 0 {
			r := io.Reader(os.Stdin)
			lr := io.LimitReader(r, int64(opts.Limit+opts.Offset))
			buf, _ = io.ReadAll(lr)
		} else {
			r := io.Reader(os.Stdin)
			buf, _ = io.ReadAll(r)
		}
	} else {
		r, err := os.Open(opts.From)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "error", err)
			os.Exit(1)
		}
		buf, _ = io.ReadAll(r) //io.ReadCloser
		//fmt.Println(string(buf))
		r.Close()
	}

	if opts.Offset != 0 {
		if int(opts.Offset) > len(buf) || opts.Offset < 0 {
			_, _ = fmt.Fprintln(os.Stderr, "error", err)
			os.Exit(1)
		}
		buf = buf[opts.Offset:]
	}

	/*if opts.Limit != 0 {
		if int(opts.Limit) < len(buf) {
			buf = buf[:opts.Limit]
		}
	}*/

	str := string(buf)

	if opts.Conv != "" {
		if strings.Contains(opts.Conv, "lower_case") && strings.Contains(opts.Conv, "upper_case") {
			_, _ = fmt.Fprintln(os.Stderr, "error", err)
			os.Exit(1)
		} else if strings.Contains(opts.Conv, "trim_spaces") {
			str = strings.TrimSpace(str)
		} else if strings.Contains(opts.Conv, "lower_case") {
			str = strings.ToLower(str)
		} else if strings.Contains(opts.Conv, "upper_case") {
			str = strings.ToUpper(str)
		} else {
			_, _ = fmt.Fprintln(os.Stderr, "error", err)
			os.Exit(1)
		}
	}

	if opts.To != "" {
		_, err := os.Stat(opts.To)
		if !os.IsNotExist(err) {
			_, _ = fmt.Fprintln(os.Stderr, "error", err)
			os.Exit(1)
		}
		to, err := os.Create(opts.To)
		if os.IsNotExist(err) {
			_, _ = fmt.Fprintln(os.Stderr, "error", err)
			os.Exit(1)
		}
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "error", err)
			os.Exit(1)
		}
		//to.WriteString(string(buf))
		to.WriteString(str)
		to.Close()
	} else {
		if _, err = fmt.Fprint(os.Stdout, str); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "error", err)
			os.Exit(1)
		}
		/*if _, err = fmt.Fprint(os.Stdout, string(buf)); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "error", err)
			os.Exit(1)
		}*/
		//fmt.Print(string(buf))
	}
}
