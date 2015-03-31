package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/tcnksm/go-latest"
)

type CLI struct {
	// out/err stream is the stdout and stderr
	// to write message from CLI
	outStream, errStream io.Writer
}

// Run executes CLI and return its exit code
func (c *CLI) Run(args []string) int {
	var githubTag latest.GithubTag

	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(c.errStream)

	flags.StringVar(&githubTag.Repository, "repo", "", "Repository name")
	flags.StringVar(&githubTag.Owner, "owner", "", "Repository owner name")

	flgNew := flags.Bool("new", false, "Check TAG is new and greater")
	flgVersion := flags.Bool("version", false, "Print version information")
	flgHelp := flags.Bool("help", false, "Print this message and quit")
	flgDebug := flags.Bool("debug", false, "Print verbose(debug) output")

	if err := flags.Parse(args[1:]); err != nil {
		fmt.Fprint(c.errStream, "Failed to parse flag")
		return 1
	}

	if *flgVersion {
		fmt.Fprintf(c.errStream, "%s Version v%s\n", Name, Version)
		return 0
	}

	if *flgHelp {
		fmt.Fprintf(c.errStream, helpText)
		return 0
	}

	if os.Getenv(envDebug) != "" {
		*flgDebug = true
	}

	parsedArgs := flags.Args()
	if len(parsedArgs) != 1 {
		fmt.Fprintf(c.errStream, "Invalid arguments\n")
		return 1
	}

	f := latest.DeleteFrontV()
	githubTag.FixVersionStrFunc = f
	target := f(parsedArgs[0])
	res, err := latest.Check(target, &githubTag)

	if err != nil {
		fmt.Fprintf(c.errStream, "Failed to check Tags on GitHub: %s\n", err.Error())
		return 1
	}

	if *flgNew && !res.New {
		if *flgDebug {
			fmt.Fprintf(c.outStream, "%s is not new\n", target)
		}
		return 1
	}

	if *flgNew && res.New {
		if *flgDebug {
			fmt.Fprintf(c.outStream, "%s is new\n", target)
		}
		return 0
	}

	if !res.Latest {
		if *flgDebug {
			fmt.Fprintf(c.outStream,
				"%s is not latest (%s)\n", target, res.Current)
		}
		return 1
	}

	if *flgDebug {
		fmt.Fprintf(c.outStream, "%s is latest\n", target)
	}
	return 0
}

const helpText = `Usage: latest [options] TAG(VERSION)

    latest command check TAG(VERSION) is latest. If is not latest,
    it returns non-zero value. It try to compare version by Semantic
    Versioning. By default, it tries to check tags on GitHub.

Options:

    -owner=NAME    Set GitHub repository owner name.

    -repo=NAME     Set Github repository name.

    -simple=URL    Try to check HTML on provided URL. HTML must be included
                   specific meta tag. See more detail spec on GitHub. 
                   https://github.com/tcnksm/go-latest

    -new           Check TAG(VERSION) is new. 'new' means TAG(VERSION)
                   is not exist and greater than others.

    -help          Print this message and quit.

    -debug         Print verbose(debug) output.

Example:

    $ latest -debug 0.2.0
`
