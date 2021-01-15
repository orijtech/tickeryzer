package p

import (
	"fmt"
	"sync"
	"time"
)

var m sync.Mutex

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

func h() {
	ticker1 := time.NewTicker(time.Second) // want `missing ticker1.Stop\(\) call`
	defer Stop()
	ticker2 := time.NewTicker(time.Second) // want `missing ticker2.Stop\(\) call`
	defer Stop()
	go func() {
		for range ticker2.C {
		}
	}()

	for t := range ticker1.C {
		fmt.Println("Current time: ", t)
	}
}

func i() {
	ticker := time.NewTicker(time.Second) // want `missing ticker.Stop\(\) call`
	m.Lock()
	defer m.Unlock()

	for t := range ticker.C {
		fmt.Println("Current time: ", t)
	}
}

func Stop() {}
