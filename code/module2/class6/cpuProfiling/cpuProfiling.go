package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"
)

var cpuProfiling *string = flag.String("cpuProfile", "./log/cpuProfile", "write cpu profile to file")

func main() {
	flag.Parse()

	f, err := os.Create(*cpuProfiling)
	if err != nil {
		panic(err)
	}

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	var result int
	for i := 0; i < 100000000; i++ {
		result += i
	}

	log.Println("result = ", result)
}
