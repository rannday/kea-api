package client

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

// TestCallCommand_Success verifies a successful call returns a valid response.
func TestCallCommand_Success(t *testing.T) {
	client := newMockClient([]CommandResponse{
		{Result: ResultSuccess, Arguments: mustEncodeRawJSON(map[string]any{"ok": true})},
	}, nil)

	got, err := CallCommand(client, "status-get", "svc")
	if err != nil {
		t.Fatalf("CallCommand() error = %v", err)
	}
	if len(got) != 1 || got[0].Result != ResultSuccess {
		t.Errorf("CallCommand() = %+v", got)
	}
}

// TestCallCommand_Error verifies error is returned on transport failure.
func TestCallCommand_Error(t *testing.T) {
	client := newMockClient(nil, errors.New("transport failure"))
	_, err := CallCommand(client, "status-get", "svc")
	if err == nil || err.Error() != "status-get failed: transport failure" {
		t.Errorf("expected transport error, got %v", err)
	}
}

// TestCallCommand_Empty verifies error is returned when no response is received.
func TestCallCommand_Empty(t *testing.T) {
	client := newMockClient([]CommandResponse{}, nil)
	_, err := CallCommand(client, "empty-reply", "svc")
	if err == nil || err.Error() != "empty-reply returned empty response" {
		t.Errorf("expected empty response error, got %v", err)
	}
}

// TestCallCommand_BadResult verifies error is returned on non-success result code.
func TestCallCommand_BadResult(t *testing.T) {
	client := newMockClient([]CommandResponse{
		{Result: ResultUnsupported, Text: "bad command"},
	}, nil)
	_, err := CallCommand(client, "unsupported", "svc")
	if err == nil || err.Error() != "unsupported command: bad command" {
		t.Errorf("expected unsupported error, got %v", err)
	}
}

// TestCallAndDecode verifies arguments are decoded into a typed slice.
func TestCallAndDecode(t *testing.T) {
	client := newMockClient([]CommandResponse{
		{Result: ResultSuccess, Arguments: mustEncodeRawJSON(map[string]string{"foo": "bar"})},
	}, nil)

	type Foo struct {
		Foo string `json:"foo"`
	}
	got, err := CallAndDecode[Foo](client, "config-get", "svc")
	if err != nil {
		t.Fatalf("CallAndDecode() error = %v", err)
	}
	if len(got) != 1 || got[0].Foo != "bar" {
		t.Errorf("CallAndDecode() = %+v", got)
	}
}

// TestCallAndDecode_Empty verifies error when response is empty.
func TestCallAndDecode_Empty(t *testing.T) {
	client := newMockClient([]CommandResponse{}, nil)
	_, err := CallAndDecode[map[string]string](client, "config-get", "svc")
	if err == nil || !strings.Contains(err.Error(), "returned empty response") {
		t.Errorf("expected empty response error, got %v", err)
	}
}

// TestCallAndDecode_BadJSON verifies decoding error on invalid JSON.
func TestCallAndDecode_BadJSON(t *testing.T) {
	client := newMockClient([]CommandResponse{
		{Result: ResultSuccess, Arguments: json.RawMessage(`{not valid}`)},
	}, nil)
	_, err := CallAndDecode[map[string]string](client, "config-get", "svc")
	if err == nil || !strings.Contains(err.Error(), "decode config-get arguments") {
		t.Errorf("expected decode error, got %v", err)
	}
}

// TestDecodeFirst verifies DecodeFirst returns the first decoded response.
func TestDecodeFirst(t *testing.T) {
	client := newMockClient([]CommandResponse{
		{Result: ResultSuccess, Arguments: mustEncodeRawJSON(map[string]string{"val": "x"})},
	}, nil)

	type V struct {
		Val string `json:"val"`
	}
	got, err := DecodeFirst[V](client, "config-get", "svc")
	if err != nil {
		t.Fatalf("DecodeFirst() error = %v", err)
	}
	if got.Val != "x" {
		t.Errorf("DecodeFirst() = %+v", got)
	}
}

// TestDecodeFirst_SingleItem verifies DecodeFirst handles a map response.
func TestDecodeFirst_SingleItem(t *testing.T) {
	client := newMockClient([]CommandResponse{
		{Result: ResultSuccess, Arguments: mustEncodeRawJSON(map[string]any{"hello": "world"})},
	}, nil)

	got, err := DecodeFirst[map[string]any](client, "config-get", "svc")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if got["hello"] != "world" {
		t.Errorf("expected 'world', got: %+v", got)
	}
}

// TestDecodeFirstWithText verifies text and arguments are decoded.
func TestDecodeFirstWithText(t *testing.T) {
	client := newMockClient([]CommandResponse{
		{Result: ResultSuccess, Text: "hello", Arguments: mustEncodeRawJSON(map[string]string{"val": "y"})},
	}, nil)

	type V struct {
		Val string `json:"val"`
	}
	text, got, err := DecodeFirstWithText[V](client, "status-get", "svc")
	if err != nil {
		t.Fatalf("DecodeFirstWithText() error = %v", err)
	}
	if text != "hello" || got.Val != "y" {
		t.Errorf("DecodeFirstWithText() = (%q, %+v)", text, got)
	}
}

// TestDecodeFirstWithText_SingleItem verifies DecodeFirstWithText returns map and text.
func TestDecodeFirstWithText_SingleItem(t *testing.T) {
	client := newMockClient([]CommandResponse{
		{Result: ResultSuccess, Text: "text-ok", Arguments: mustEncodeRawJSON(map[string]string{"foo": "bar"})},
	}, nil)

	text, got, err := DecodeFirstWithText[map[string]string](client, "config-get", "svc")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if text != "text-ok" {
		t.Errorf("unexpected text: %v", text)
	}
	if got["foo"] != "bar" {
		t.Errorf("expected bar, got: %+v", got)
	}
}

// TestDecodeFirstWithText_BadJSON verifies error on bad JSON with text.
func TestDecodeFirstWithText_BadJSON(t *testing.T) {
	client := newMockClient([]CommandResponse{
		{Result: ResultSuccess, Text: "oops", Arguments: json.RawMessage(`{bad}`)},
	}, nil)

	_, _, err := DecodeFirstWithText[map[string]string](client, "status-get", "svc")
	if err == nil || !strings.Contains(err.Error(), "decode arguments") {
		t.Errorf("expected decode error, got %v", err)
	}
}

// TestDecodeFirst_Error verifies DecodeFirst returns error when CallAndDecode fails.
func TestDecodeFirst_Error(t *testing.T) {
	client := newMockClient([]CommandResponse{}, nil) // triggers "empty response" error
	_, err := DecodeFirst[map[string]string](client, "config-get", "svc")
	if err == nil || !strings.Contains(err.Error(), "empty response") {
		t.Errorf("expected decode error, got %v", err)
	}
}

// TestDecodeFirstWithText_Error verifies DecodeFirstWithText returns error when CallCommand fails.
func TestDecodeFirstWithText_Error(t *testing.T) {
	client := newMockClient([]CommandResponse{}, nil) // triggers "empty response" error
	_, _, err := DecodeFirstWithText[map[string]string](client, "status-get", "svc")
	if err == nil || !strings.Contains(err.Error(), "empty response") {
		t.Errorf("expected decode error, got %v", err)
	}
}

// TestCallAndExtractText verifies .Text is extracted without decoding .Arguments.
func TestCallAndExtractText(t *testing.T) {
	client := newMockClient([]CommandResponse{
		{Result: ResultSuccess, Text: "build-report text"},
	}, nil)

	got, err := CallAndExtractText(client, "build-report", "svc")
	if err != nil {
		t.Fatalf("CallAndExtractText() error = %v", err)
	}
	if got != "build-report text" {
		t.Errorf("CallAndExtractText() = %q, want %q", got, "build-report text")
	}
}

// TestCallAndExtractText_Error verifies error handling.
func TestCallAndExtractText_Error(t *testing.T) {
	client := newMockClient([]CommandResponse{}, nil) // triggers "empty response"
	_, err := CallAndExtractText(client, "build-report", "svc")
	if err == nil || !strings.Contains(err.Error(), "returned empty response") {
		t.Errorf("expected empty response error, got %v", err)
	}
}
