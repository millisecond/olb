package picker

import (
	"time"
	"github.com/millisecond/olb/model"
)

// Picker selects a target from a list of targets.
type Picker func(targets []*model.Target) *model.Target

// Pickers contains the available Picker functions.
// Update config/load.go#load after updating.
var Pickers = map[string]Picker{
	"rnd": RndPicker,
	//"rr":  rrPicker,
}

// rndPicker picks a random target from the list of targets.
func RndPicker(targets []*model.Target) *model.Target {
	return targets[randIntn(len(targets))]
}

//// rrPicker picks the next target from a list of targets using round-robin.
//func rrPicker(targets []*model.Target) *model.Target {
//	u := targets[r.total%uint64(len(targets))]
//	atomic.AddUint64(&r.total, 1)
//	return u
//}

// stubbed out for testing
// we implement the randIntN function using the nanosecond time counter
// since it is 15x faster than using the pseudo random number generator
// (12 ns vs 190 ns) Most HW does not seem to provide clocks with ns
// resolution but seem to be good enough for µs resolution. Since
// requests are usually handled within several ms we should have enough
// variation. Within 1 ms we have 1000 µs to distribute among a smaller
// set of entities (<< 100)
var randIntn = func(n int) int {
	if n == 0 {
		return 0
	}
	return int(time.Now().UnixNano()/int64(time.Microsecond)) % n
}
