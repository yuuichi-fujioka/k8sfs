package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"runtime/pprof"
)

func startHandlingSignal() {

	cpuprof()

	ch := make(chan os.Signal)

	signal.Notify(
		ch,
		syscall.SIGQUIT)

	go func() {
		for {
			s := <-ch
			switch s {
			case syscall.SIGQUIT:
				memprof()
				restartCpuprof()
			}
		}
	}()
}

func restartCpuprof() {
	pprof.StopCPUProfile()
	f, err := os.Create("cpu.prof")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
}
func cpuprof() {
	f, err := os.Create("cpu.prof")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
}

func memprof() {
	f, err := os.Create("mem.prof")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	pprof.WriteHeapProfile(f)
}
