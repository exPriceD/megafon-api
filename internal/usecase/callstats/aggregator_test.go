package callstats_test

import (
	"testing"
	"time"

	"megafon-buisness-reports/internal/domain/entities"
	"megafon-buisness-reports/internal/usecase/callstats"
)

func TestAggregate(t *testing.T) {
	now := time.Now()

	calls := []entities.Call{
		{Client: "111", Status: "success", MissedStatus: entities.MissedNone},
		{Client: "222", Status: "missed", MissedStatus: entities.MissedClientCalledBack},
		{Client: "333", Status: "missed", MissedStatus: entities.MissedWeCalledBackOK},
		{Client: "444", Status: "missed", MissedStatus: entities.MissedWeCalledBackFail},
		{Client: "555", Status: "missed", MissedStatus: entities.MissedNoCallBack},
	}

	calls = append(calls, entities.Call{Client: "666", Status: "success", MissedStatus: entities.MissedNone, Start: now})

	got := callstats.Aggregate(calls)

	expect := map[callstats.Bucket]int{
		callstats.BucketAll:              6,
		callstats.BucketMissed:           4,
		callstats.BucketClientCalledBack: 1,
		callstats.BucketWeCalledBackOK:   1,
		callstats.BucketWeCalledBackFail: 1,
		callstats.BucketNoCallBack:       1,
	}

	for bucket, want := range expect {
		if got[bucket].Count != want {
			t.Errorf("bucket %v: want %d, got %d",
				bucket, want, got[bucket].Count)
		}
	}

	has := func(nums []string, target string) bool {
		for _, n := range nums {
			if n == target {
				return true
			}
		}
		return false
	}
	if !has(got[callstats.BucketNoCallBack].Numbers, "555") {
		t.Errorf("‘555’ не найден в BucketNoCallBack")
	}
	if !has(got[callstats.BucketWeCalledBackFail].Numbers, "444") {
		t.Errorf("‘444’ не найден в BucketWeCalledBackFail")
	}

	allNums := map[string]struct{}{}
	for _, s := range got {
		for _, n := range s.Numbers {
			allNums[n] = struct{}{}
		}
	}
	if len(allNums) != 6 {
		t.Errorf("ожидалось 6 уникальных номеров, получили %d", len(allNums))
	}

	for bucket, stat := range got {
		if stat == nil {
			t.Fatalf("stat for bucket %v is nil", bucket)
		}
		if stat.Numbers == nil {
			t.Fatalf("bucket %v: Numbers slice is nil", bucket)
		}
	}
}

func TestAggregate_Deduplicate(t *testing.T) {
	calls := []entities.Call{
		{Client: "111", Status: "missed", MissedStatus: entities.MissedNoCallBack},
		{Client: "111", Status: "missed", MissedStatus: entities.MissedNoCallBack}, // дубль
	}

	sum := callstats.Aggregate(calls)

	if sum[callstats.BucketNoCallBack].Count != 2 {
		t.Fatalf("Count want 2, got %d",
			sum[callstats.BucketNoCallBack].Count)
	}
	if len(sum[callstats.BucketNoCallBack].Numbers) != 1 {
		t.Fatalf("unique Numbers want 1, got %d",
			len(sum[callstats.BucketNoCallBack].Numbers))
	}
	if sum[callstats.BucketAll].Count != 2 {
		t.Fatalf("BucketAll Count want 2, got %d",
			sum[callstats.BucketAll].Count)
	}
}
