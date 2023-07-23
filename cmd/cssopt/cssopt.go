package main

import (
	"os"

	"github.com/runeimp/cssopt"
	"github.com/runeimp/cssopt/configuration"
	"github.com/runeimp/termlog"
)

func main() {
	tlog := termlog.New()

	if len(os.Args) == 1 {
		tlog.Fatal("Not enough arguments")
		os.Exit(1)
	}

	config := configuration.New()
	config.Comments.All = true // remove all comments
	config.Imports = true
	config.Newlines = "merge"
	optimizer := cssopt.GetOptimizer(config)

	css, err := optimizer.ProcessPath(os.Args[1])
	if err != nil {
		tlog.Fatal(err.Error())
	}
	tlog.Info("cssopt CLI | len(css): %d | css:\n%s", len(css), css)

}
