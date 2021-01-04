package p

import (
	"fmt"
	"time"
)

func f() {
	ticker := time.NewTicker(time.Second) // want `missing ticker.Stop\(\) call`
	for t := range ticker.C {
		fmt.Println("Current time: ", t)
	}
}

func g(ticker *time.Ticker) { // BAD: this should be reported
	for t := range ticker.C {
		fmt.Println("Current time: ", t)
	}
}
