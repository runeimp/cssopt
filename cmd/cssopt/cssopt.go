package main

import (
	"github.com/runeimp/cssopt"
	"github.com/runeimp/termlog"
)

func main() {
	tlog := termlog.New()
	optimizer := cssopt.GetOptimizer()
	css, err := optimizer.ProcessPath("css/base.css")
	if err != nil {
		tlog.Fatal(err.Error())
	}

	tlog.Info("cssopt CLI | css: %s", css)
}
