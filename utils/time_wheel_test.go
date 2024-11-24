package utils

import (
	"fmt"
	"testing"
)

func TestTimingWheel_Start(t *testing.T) {
	Instance().Start()
	go func() {
		for t := range Instance().TickChan() {
			fmt.Printf("time wheel run one round: %s", t.Format("02 15:04:05.999"))
		}
	}()
}
