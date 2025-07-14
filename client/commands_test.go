package client

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

// mustEncodeRawJSON marshals v into json.RawMessage or panics.
func mustEncodeRawJSON(v interface{}) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

// mockClient implements Transport for unit testing.
type mockClient struct {
	responses []CommandResponse
	err       error
}

func (m *mockClient) Call(req CommandRequest, out interface{}) error {
	if m.err != nil {
		return m.err
	}
	*out.(*[]CommandResponse) = m.responses
	return nil
}

// newMockClient wraps mockClient in a real Client for testing.
func newMockClient(responses []CommandResponse, err error) *Client {
	return NewClient(&mockClient{responses: responses, err: err})
}

// TestBuildReportMulti verifies multi-response text extraction.
func TestBuildReportMulti(t *testing.T) {
	client := newMockClient([]CommandResponse{{Text: "a"}, {Text: "b"}}, nil)

	got, err := BuildReportMulti(client, "svc1", "svc2")
	want := []string{"a", "b"}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

// TestConfigGetMulti decodes multi-response JSON objects.
func TestConfigGetMulti(t *testing.T) {
	client := newMockClient([]CommandResponse{
		{Arguments: mustEncodeRawJSON(map[string]any{"val": 1})},
		{Arguments: mustEncodeRawJSON(map[string]any{"val": 2})},
	}, nil)

	type Config struct {
		Val int `json:"val"`
	}
	got, err := ConfigGetMulti[Config](client, "svc1", "svc2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 || got[0].Val != 1 || got[1].Val != 2 {
		t.Errorf("unexpected result: %+v", got)
	}
}

// TestListCommandsMulti extracts multi-service command lists.
func TestListCommandsMulti(t *testing.T) {
	client := newMockClient([]CommandResponse{
		{Arguments: mustEncodeRawJSON([]string{"foo", "bar"})},
		{Arguments: mustEncodeRawJSON([]string{"baz"})},
	}, nil)

	got, err := ListCommandsMulti(client, "svc1", "svc2")
	want := [][]string{{"foo", "bar"}, {"baz"}}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

// TestStatusGetMulti decodes boolean values from multiple status responses.
func TestStatusGetMulti(t *testing.T) {
	client := newMockClient([]CommandResponse{
		{Arguments: mustEncodeRawJSON(map[string]any{"ok": true})},
		{Arguments: mustEncodeRawJSON(map[string]any{"ok": false})},
	}, nil)

	type Status struct {
		OK bool `json:"ok"`
	}
	got, err := StatusGetMulti[Status](client, "svc1", "svc2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 || !got[0].OK || got[1].OK {
		t.Errorf("unexpected result: %+v", got)
	}
}

// TestVersionGetMulti verifies text version strings from multiple responses.
func TestVersionGetMulti(t *testing.T) {
	client := newMockClient([]CommandResponse{{Text: "v1"}, {Text: "v2"}}, nil)

	got, err := VersionGetMulti(client, "svc1", "svc2")
	want := []string{"v1", "v2"}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

// TestMultiHelpers_Error ensures all Multi helpers return errors on failure.
func TestMultiHelpers_Error(t *testing.T) {
	client := newMockClient(nil, errors.New("boom"))

	if _, err := BuildReportMulti(client, "svc"); err == nil {
		t.Error("expected error for BuildReportMulti")
	}
	if _, err := ConfigGetMulti[any](client, "svc"); err == nil {
		t.Error("expected error for ConfigGetMulti")
	}
	if _, err := ListCommandsMulti(client, "svc"); err == nil {
		t.Error("expected error for ListCommandsMulti")
	}
	if _, err := StatusGetMulti[any](client, "svc"); err == nil {
		t.Error("expected error for StatusGetMulti")
	}
	if _, err := VersionGetMulti(client, "svc"); err == nil {
		t.Error("expected error for VersionGetMulti")
	}
}
