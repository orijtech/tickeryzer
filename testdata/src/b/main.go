package main

import (
	"fmt"
	"time"
)

func main() {
	go f()

	ticker := time.NewTicker(time.Second) // This is ok, main goroutine exits because the whole process quits when in main.
	for t := range ticker.C {
		fmt.Println("Current time: ", t)
	}
}

func f() {
	ticker := time.NewTicker(time.Second) // want `missing ticker.Stop\(\) call`
	for t := range ticker.C {
		fmt.Println("Current time in f: ", t)
	}
}
