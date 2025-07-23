package entities

import "time"

type Period string

const (
	PeriodToday     Period = "today"
	PeriodYesterday Period = "yesterday"
	PeriodThisWeek  Period = "this_week"
	PeriodLastWeek  Period = "last_week"
	PeriodThisMonth Period = "this_month"
	PeriodLastMonth Period = "last_month"
)

type CallFilter struct {
	Start         *time.Time
	End           *time.Time
	Period        Period
	Type          CallType
	Limit         int
	User          string
	Diversion     string
	Client        string
	Groups        []string
	FirstAnswered bool
	ProcessMissed bool
	MissedStatus  []int
}
