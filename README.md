# octlango
CLI to get statistics on languages used on GitHub.

## Install
```
go install github.com/korosuke613/octlango@latest
```

or from the [Release](https://github.com/korosuke613/octlango/releases/latest).

## Usage

```
❯ octlango -h            
NAME:
   octlango - CLI to get statistics on languages used on GitHub.

USAGE:
   main [global options] command [command options] [arguments...]

VERSION:
   v0.0.4

COMMANDS:
   version  Print octlango version
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --username GITHUB_USERNAME, -u GITHUB_USERNAME  your GITHUB_USERNAME [$OCTLANGO_GH_USERNAME]
   --token GITHUB_TOKEN, -t GITHUB_TOKEN           your GITHUB_TOKEN [$OCTLANGO_GH_TOKEN, $GITHUB_TOKEN]
   --sort-by-size, -s                              if true, the order is by size. (default: true)
   --reverse-order, -r                             if true, reverse the result. (default: false)
   --help, -h                                      show help (default: false)
   --version, -v                                   print the version (default: false)
```

## Example
input
```
./octlango -u korosuke613 -t YOUR_GITHUB_TOKEN
```

output
```json5
{
  "updated_range": {
    "oldest": "2020-08-02T16:43:48Z",
    "latest": "2021-07-11T13:01:20Z"
  },
  "language_sizes": [
    {
      "name": "TypeScript",
      "size": 537091,
      "percentage": 50.01
    },
    {
      "name": "Vue",
      "size": 103000,
      "percentage": 9.59
    },
    {
      "name": "JavaScript",
      "size": 93888,
      "percentage": 8.74
    },
    {
      "name": "HCL",
      "size": 89233,
      "percentage": 8.31
    },
    {
      "name": "HTML",
      "size": 80865,
      "percentage": 7.53
    },
    {
      "name": "Go",
      "size": 65508,
      "percentage": 6.1
    },
    // ...
  ]
}
```
