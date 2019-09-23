package main

import (
	"flag"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/sampaioletti/wagoja/pkg/jswasm"
)

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatal("must specify js file to run")
	}
	wd, _ := os.Getwd()
	//look for file
	file := flag.Arg(0)
	if !path.IsAbs(file) {
		file = filepath.Join(wd, file)
	}
	if !strings.HasSuffix(file, ".js") {
		file = file + ".js"
	}
	if _, err := os.Stat(file); err != nil {
		log.Fatal("File Not found")
	}

	jw := jswasm.New()
	jw.RunFile(file)

}
