package response

type CallDTO struct {
	UID          string `json:"uid"`
	Type         string `json:"type"`
	Status       string `json:"status"`
	Client       string `json:"client"`
	Diversion    string `json:"diversion"`
	TelnumName   string `json:"telnum_name"`
	Destination  string `json:"destination"`
	User         string `json:"user"`
	UserName     string `json:"user_name"`
	GroupName    string `json:"group_name"`
	StartRaw     string `json:"start"`
	Wait         int    `json:"wait"`
	Duration     int    `json:"duration"`
	Record       string `json:"record"`
	Rating       int    `json:"rating"`
	Note         string `json:"note"`
	MissedStatus int    `json:"missedStatus"`
}
