package services_test

import (
	"bytes"
	"context"
	"errors"
	"reflect"
	"testing"

	"megafon-buisness-reports/internal/domain/entities"
	"megafon-buisness-reports/internal/usecase/callstats"
	uc "megafon-buisness-reports/internal/usecase/services"
)

type mockRepo struct {
	gotFilter entities.CallFilter
	retCalls  []entities.Call
}

func (m *mockRepo) History(_ context.Context, f entities.CallFilter) ([]entities.Call, error) {
	m.gotFilter = f
	return m.retCalls, nil
}

type mockBuilder struct {
	lastSummary callstats.Summary
	buf         *bytes.Buffer
	fail        bool
}

func (m *mockBuilder) Build(s callstats.Summary) (*bytes.Buffer, error) {
	m.lastSummary = s
	if m.fail {
		return nil, errors.New("builder fail")
	}
	m.buf = bytes.NewBufferString("xlsx")
	return m.buf, nil
}

func TestReportService_GenerateCallReport(t *testing.T) {
	repo := &mockRepo{retCalls: []entities.Call{
		{Client: "111", Status: "success"},
	}}
	builder := &mockBuilder{}
	svc := uc.NewReportService(repo, builder)

	var filter entities.CallFilter

	buf, err := svc.GenerateCallReport(context.Background(), filter)
	if err != nil {
		t.Fatalf("GenerateCallReport: %v", err)
	}
	if buf == nil || buf.String() != "xlsx" {
		t.Fatalf("buf unexpected: %v", buf)
	}
	if !reflect.DeepEqual(repo.gotFilter, filter) {
		t.Errorf("filter not passed correctly")
	}

	summary := builder.lastSummary
	if summary[callstats.BucketAll].Count != 1 {
		t.Errorf("want 1 call in BucketAll, got %d", summary[callstats.BucketAll].Count)
	}
}

func TestReportService_BuilderError(t *testing.T) {
	repo := &mockRepo{}
	builder := &mockBuilder{fail: true}
	svc := uc.NewReportService(repo, builder)

	_, err := svc.GenerateCallReport(context.Background(), entities.CallFilter{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
