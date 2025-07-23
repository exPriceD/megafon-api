package request

import "time"

type Period string

const (
	Today     Period = "today"
	Yesterday Period = "yesterday"
	ThisWeek  Period = "this_week"
	LastWeek  Period = "last_week"
	ThisMonth Period = "this_month"
	LastMonth Period = "last_month"
)

type CallType string

const (
	All    CallType = "all"
	In     CallType = "in"
	Out    CallType = "out"
	Missed CallType = "missed"
)

type HistoryParams struct {
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
