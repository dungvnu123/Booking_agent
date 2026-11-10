package clock

import "time"

// Clock is injected into time-sensitive booking, roster and token use cases.
type Clock interface{ Now() time.Time }
type System struct{ Location *time.Location }

func (s System) Now() time.Time { return time.Now().In(s.Location) }

type Fixed struct{ Time time.Time }

func (f Fixed) Now() time.Time { return f.Time }
