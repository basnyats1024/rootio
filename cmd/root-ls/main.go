// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// root-ls dumps the content of a ROOT file
package main // import "github.com/go-hep/rootio/cmd/root-ls"

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strings"

	"github.com/go-hep/rootio"
)

var (
	doProf = flag.String("profile", "", "filename of cpuprofile")
	doSI   = flag.Bool("sinfos", false, "print StreamerInfos")
	doTree = flag.Bool("t", false, "print Tree(s) (recursively)")
)

func main() {
	flag.Parse()

	if *doProf != "" {
		f, err := os.Create(*doProf)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if flag.NArg() <= 0 {
		fmt.Fprintf(os.Stderr, "**error** you need to give a ROOT file\n")
		flag.Usage()
		os.Exit(1)
	}

	for ii, fname := range flag.Args() {

		if ii > 0 {
			fmt.Printf("\n")
		}

		fmt.Printf("=== [%s] ===\n", fname)
		f, err := rootio.Open(fname)
		if err != nil {
			fmt.Fprintf(os.Stderr, "rootio: failed to open [%s]: %v\n", fname, err)
			os.Exit(1)
		}
		fmt.Printf("version: %v\n", f.Version())
		if *doSI {
			fmt.Printf("streamer-infos:\n")
			sinfos := f.StreamerInfo()
			for _, v := range sinfos {
				name := v.Name()
				fmt.Printf(" StreamerInfo for %q version=%d title=%q\n", name, v.ClassVersion(), v.Title())
				for _, elm := range v.Elements() {
					fmt.Printf("  %-15s %-20s offset=%3d type=%3d size=%3d %s\n", elm.TypeName(), elm.Name(), elm.Offset(), elm.Type(), elm.Size(), elm.Title())
				}
			}
		}

		for _, k := range f.Keys() {
			if *doTree {
				tree, ok := k.Value().(rootio.Tree)
				if ok {
					fmt.Printf("%-8s %-40s %s (entries=%d)\n", k.Class(), k.Name(), k.Title(), tree.Entries())
					displayBranches(tree, 2)
					continue
				}
			}
			fmt.Printf("%-8s %-40s %s (cycle=%d)\n", k.Class(), k.Name(), k.Title(), k.Cycle())
		}
	}
}

type brancher interface {
	Branches() []rootio.Branch
}

func displayBranches(bres brancher, indent int) {
	branches := bres.Branches()
	for _, b := range branches {
		fmt.Printf("%s%-20s %-20q %v\n", strings.Repeat(" ", indent), b.Name(), b.Title(), b.Class())
		displayBranches(b, indent+2)
	}
}
