package callstats

import "megafon-buisness-reports/internal/domain/entities"

type Bucket int

const (
	BucketAll Bucket = iota + 1
	BucketMissed
	BucketClientCalledBack
	BucketWeCalledBackOK
	BucketWeCalledBackFail
	BucketNoCallBack
)

type Stat struct {
	Count   int
	Numbers []string
}

type Summary map[Bucket]*Stat

// Aggregate группирует звонки по шести категориям.
func Aggregate(calls []entities.Call) Summary {
	sum := Summary{
		BucketAll:              &Stat{},
		BucketMissed:           &Stat{},
		BucketClientCalledBack: &Stat{},
		BucketWeCalledBackOK:   &Stat{},
		BucketWeCalledBackFail: &Stat{},
		BucketNoCallBack:       &Stat{},
	}

	for _, c := range calls {
		add(sum[BucketAll], c.Client)

		if entities.CallType(c.Status) == entities.CallMissed {
			add(sum[BucketMissed], c.Client)
		}

		if entities.CallType(c.Status) == entities.CallIn && c.Duration <= 5 {
			add(sum[BucketNoCallBack], c.Client)
			continue
		}

		if entities.CallType(c.Status) == entities.CallOut && c.Duration <= 5 {
			add(sum[BucketWeCalledBackFail], c.Client)
			continue
		}

		switch c.MissedStatus {
		case entities.MissedClientCalledBack:
			add(sum[BucketClientCalledBack], c.Client)
		case entities.MissedWeCalledBackOK:
			add(sum[BucketWeCalledBackOK], c.Client)
		case entities.MissedWeCalledBackFail:
			add(sum[BucketWeCalledBackFail], c.Client)
		case entities.MissedNoCallBack:
			add(sum[BucketNoCallBack], c.Client)
		}
	}
	return sum
}

func add(s *Stat, num string) {
	s.Count++
	s.Numbers = append(s.Numbers, num)
}
