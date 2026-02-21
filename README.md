# jppgo - A JSON pretty printer written in Go

This is a simple command-line tool that takes JSON input and formats it in a more human-readable way. It supports various options for customizing the output, such as indentation levels, depth and gjson queries.

This tool was built using [gjson](https://github.com/tidwall/gjson) and [colorjson](https://github.com/TylerBrock/colorjson).

## Installation

### Using Go

```bash
go install github.com/KartoffelChipss/jppgo@latest
```

### Download Binary

Download from:
https://github.com/KartoffelChipss/jppgo/releases

After downloading, make the binary executable and move it to a directory in your PATH.

```bash
chmod +x jppgo
mv jppgo /usr/local/bin/
```

On macOS, you'll probably need to remove it from quarantine:

```bash
xattr -d com.apple.quarantine /usr/local/bin/jppgo
```

### Build from Source

You need to have Go installed on your system. Then you can clone the repository and build the tool:

```bash
git clone https://github.com/KartoffelChipss/jppgo
cd jppgo
go build
```

## Usage

```bash
jppgo [options] [file]
```

### Options

- `-h`, `--help`: Show help message
- `-i`, `--indent`: Set the indentation level (default: 2)
- `-d`, `--max-depth`: Set the maximum depth (-1 = unlimited) (default: -1)
- `-p`, `--path`: Set a gjson query path to filter the JSON output (e.g. "items.0.name")
- `-v`, `--version`: Show version information

### Examples

1. Pretty print JSON from a file:

```bash
jppgo data.json
```

2. Pretty print JSON from standard input:

```bash
cat data.json | jppgo
```

3. Pretty print JSON with a specific indentation level:

```bash
jppgo -i 4 data.json
```

4. Pretty print JSON with a maximum depth:

```bash
jppgo -d 3 data.json
```

5. Pretty print JSON with a gjson query path:

```bash
jppgo -p "items.0.name" data.json
```

You can find some more examples for gjson queires [here](https://gjson.dev).

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.