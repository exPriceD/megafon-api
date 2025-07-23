package entities

import "time"

type CallType string   // in, out, missed…
type CallStatus string // success, missed…

type MissedStatusCode int

const (
	MissedNone             MissedStatusCode = 0 // звонок принят сразу
	MissedClientCalledBack MissedStatusCode = 1 // клиент перезвонил
	MissedWeCalledBackOK   MissedStatusCode = 2 // мы перезвонили, дозвонились
	MissedNoCallBack       MissedStatusCode = 3 // клиенту не перезвонили
	MissedWeCalledBackFail MissedStatusCode = 4 // перезвонили, но не дозвонились
)

const (
	CallAll    CallType = "all"
	CallIn     CallType = "in"
	CallOut    CallType = "out"
	CallMissed CallType = "missed"
)

// Call — доменная сущность «звонок».
type Call struct {
	UID         string
	Type        CallType
	Status      CallStatus
	Client      string
	Diversion   string
	TelnumName  string
	Destination string
	User        string
	UserName    string
	GroupName   string
	Start       time.Time
	Wait        int
	Duration    int
	Record      string
	Rating      int
	Note        string

	MissedStatus MissedStatusCode
}
