package request

import (
	"fmt"
	"net/url"
	"strings"
)

func (p HistoryParams) ToQuery() url.Values {
	q := url.Values{}

	if p.Start != nil {
		q.Set("start", p.Start.UTC().Format("20060102T150405Z"))
	}
	if p.End != nil {
		q.Set("end", p.End.UTC().Format("20060102T150405Z"))
	}
	if p.Period != "" {
		q.Set("period", string(p.Period))
	}
	if p.Type != "" {
		q.Set("type", string(p.Type))
	}
	if p.Limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", p.Limit))
	}
	if p.User != "" {
		q.Set("user", p.User)
	}
	if p.Diversion != "" {
		q.Set("diversion", p.Diversion)
	}
	if p.Client != "" {
		q.Set("client", p.Client)
	}
	if len(p.Groups) > 0 {
		q.Set("groups", strings.Join(p.Groups, ","))
	}
	if p.FirstAnswered {
		q.Set("first_answered", "true")
	}
	if p.ProcessMissed {
		q.Set("processMissed", "true")
	}
	if len(p.MissedStatus) > 0 {
		var s []string
		for _, m := range p.MissedStatus {
			s = append(s, fmt.Sprintf("%d", m))
		}
		q.Set("missedStatus", strings.Join(s, ","))
	}
	return q
}
