package main

import (
	"fmt"
	"time"

	"github.com/atalkowski/go-rpc/myutils"
)

func testTimers() {
	fmt.Println("Hello World")
	timer := myutils.NewPhaseTimer()
	timer.Log("Init-A")
	timer.LogN(1, "Init-B")
	time.Sleep(50 * time.Millisecond)
	timer.Log("Part1-A")
	time.Sleep(30 * time.Millisecond)
	timer.Log("Last-A")
	timer.LogN(1, "Last-B")
	time.Sleep(170 * time.Millisecond)
	timer.AllDone()
	fmt.Printf("testTimers completed:\n%v\n", timer.ToString())

}

func main() {
	testTimers()
}
