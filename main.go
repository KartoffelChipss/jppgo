package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/TylerBrock/colorjson"
	flag "github.com/spf13/pflag"
	"github.com/tidwall/gjson"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	showVersion := flag.BoolP("version", "v", false, "print version and exit")
	indent := flag.IntP("indent", "i", 2, "indent level")
	maxDepth := flag.IntP("max-depth", "d", -1, "maximum depth (-1 = unlimited)")
	path := flag.StringP("path", "p", "", "gjson path (e.g. hits.0.analytics)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s [options] [file]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  cat file.json | %s\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s data.json\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *showVersion {
		fmt.Printf("Version: %s\nCommit: %s\nBuilt: %s\n", version, commit, date)
		os.Exit(0)
	}

	if *indent < 0 {
		fmt.Fprintln(os.Stderr, "Indent must be non-negative.")
		os.Exit(1)
	}

	var input []byte
	args := flag.Args()

	if len(args) > 1 {
		fmt.Fprintln(os.Stderr, "Only one input file may be specified.")
		os.Exit(1)
	}

	if len(args) == 1 {
		fileBytes, err := os.ReadFile(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading file:", err)
			os.Exit(1)
		}
		input = fileBytes
	} else {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			fmt.Fprintln(os.Stderr, "No input provided. Pipe JSON or specify a file.")
			os.Exit(1)
		}

		stdinBytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		input = stdinBytes
	}

	if !json.Valid(input) {
		fmt.Fprintln(os.Stderr, "Invalid JSON input.")
		os.Exit(1)
	}

	if *path != "" {
		result := gjson.GetBytes(input, *path)
		if !result.Exists() {
			fmt.Fprintln(os.Stderr, "Path not found.")
			os.Exit(1)
		}
		input = []byte(result.Raw)
	}

	var data interface{}
	if err := json.Unmarshal(input, &data); err != nil {
		panic(err)
	}

	if *maxDepth >= 0 {
		data = truncateDepth(data, 0, *maxDepth)
	}

	formatter := colorjson.NewFormatter()
	formatter.Indent = *indent

	coloredJSON, err := formatter.Marshal(data)
	if err != nil {
		panic(err)
	}

	os.Stdout.Write(coloredJSON)
}

func truncateDepth(v interface{}, current, max int) interface{} {
	// If on max depth, return a placeholder instead of the actual value
	if current >= max {
		switch v.(type) {
		case map[string]interface{}:
			return "{...}"
		case []interface{}:
			return "[...]"
		default:
			return v
		}
	}

	switch val := v.(type) {
	case map[string]interface{}: // Recursively truncate nested maps
		newMap := make(map[string]interface{})
		for k, v2 := range val {
			newMap[k] = truncateDepth(v2, current+1, max)
		}
		return newMap

	case []interface{}: // Recursively truncate nested arrays
		newArr := make([]interface{}, len(val))
		for i, v2 := range val {
			newArr[i] = truncateDepth(v2, current+1, max)
		}
		return newArr

	default:
		return val
	}
}
