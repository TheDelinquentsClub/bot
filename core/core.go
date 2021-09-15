package core

import (
	"github.com/diamondburned/arikawa/v3/state"
	"time"
)

var State *state.State

var TimeNow time.Time

func init() {
	TimeNow = time.Now()
}
