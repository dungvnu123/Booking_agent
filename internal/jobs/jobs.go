package jobs

import (
	"context"
	"log/slog"
)

// No business jobs are enabled in the scaffold. A single worker process is the
// intended initial owner. Add a DB advisory lock before running multiple workers.
type Definition struct{ Name, Schedule, Purpose string }

func Definitions() []Definition {
	return []Definition{
		{"slot-hold-cleanup", "every minute at second 0", "Expire holds"},
		{"booking-queue", "every minute at second 30", "Mark late and no-show"},
		{"monthly-roster", "every 5 minutes", "Open/lock periods and send reminders"},
		{"token-cleanup", "00:00 Asia/Ho_Chi_Minh daily", "Remove expired revoked tokens"},
	}
}
func Run(ctx context.Context) error {
	slog.Warn("worker is in scaffold mode; all business jobs are disabled", "jobs", len(Definitions()))
	<-ctx.Done()
	return nil
}
