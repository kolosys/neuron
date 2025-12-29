package mock_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kolosys/neuron/mock"
)

// MockValidator Tests
func TestMockValidator_ValidateSuccess(t *testing.T) {
	v := mock.NewMockValidator(nil)

	err := v.Validate([]byte(`{"valid":true}`), "application/json")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestMockValidator_InjectError(t *testing.T) {
	v := mock.NewMockValidator(nil)

	testErr := errors.New("invalid json")
	v.InjectValidationError(testErr)

	err := v.Validate([]byte(`invalid`), "application/json")

	if !errors.Is(err, testErr) {
		t.Errorf("expected injected error, got %v", err)
	}

	// Error should be cleared after one-shot
	err = v.Validate([]byte(`{}`), "application/json")
	if err != nil {
		t.Errorf("expected no error after one-shot, got %v", err)
	}
}

func TestMockValidator_RecordsCall(t *testing.T) {
	v := mock.NewMockValidator(&mock.MockValidatorOptions{
		RecordCalls: true,
	})

	data := []byte(`test`)
	v.Validate(data, "text/plain")

	calls := v.ValidateCalls()
	if len(calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(calls))
	}

	if string(calls[0].Data) != "test" {
		t.Errorf("expected data 'test', got %s", calls[0].Data)
	}

	if calls[0].ContentType != "text/plain" {
		t.Errorf("expected content type text/plain, got %s", calls[0].ContentType)
	}
}

func TestMockValidator_Reset(t *testing.T) {
	v := mock.NewMockValidator(&mock.MockValidatorOptions{
		RecordCalls: true,
	})

	testErr := errors.New("test")
	v.InjectValidationError(testErr)
	v.Validate([]byte{}, "")

	v.Reset()

	if len(v.ValidateCalls()) != 0 {
		t.Error("reset should clear calls")
	}

	err := v.Validate([]byte{}, "")
	if err != nil {
		t.Errorf("reset should clear errors, got %v", err)
	}
}

// MockRequestIDGenerator Tests
func TestMockRequestIDGenerator_DefaultEmpty(t *testing.T) {
	rig := mock.NewMockRequestIDGenerator(nil)

	id := rig.Generate()

	if id != "" {
		t.Errorf("expected empty string by default, got %s", id)
	}
}

func TestMockRequestIDGenerator_Sequence(t *testing.T) {
	rig := mock.NewMockRequestIDGenerator(&mock.MockRequestIDGeneratorOptions{
		IDs: []string{"id1", "id2", "id3"},
	})

	if id := rig.Generate(); id != "id1" {
		t.Errorf("expected id1, got %s", id)
	}

	if id := rig.Generate(); id != "id2" {
		t.Errorf("expected id2, got %s", id)
	}

	if id := rig.Generate(); id != "id3" {
		t.Errorf("expected id3, got %s", id)
	}

	// Past end - should return last
	if id := rig.Generate(); id != "id3" {
		t.Errorf("expected id3 (last), got %s", id)
	}
}

func TestMockRequestIDGenerator_CallCount(t *testing.T) {
	rig := mock.NewMockRequestIDGenerator(&mock.MockRequestIDGeneratorOptions{
		IDs: []string{"id1", "id2"},
	})

	rig.Generate()
	rig.Generate()
	rig.Generate()

	if rig.CallCount() != 3 {
		t.Errorf("expected 3 calls, got %d", rig.CallCount())
	}
}

func TestMockRequestIDGenerator_SetIDs(t *testing.T) {
	rig := mock.NewMockRequestIDGenerator(nil)

	rig.SetIDs([]string{"new1", "new2"})

	if id := rig.Generate(); id != "new1" {
		t.Errorf("expected new1, got %s", id)
	}

	if rig.CallCount() != 1 {
		t.Error("SetIDs should reset call count")
	}
}

func TestMockRequestIDGenerator_Reset(t *testing.T) {
	rig := mock.NewMockRequestIDGenerator(&mock.MockRequestIDGeneratorOptions{
		IDs: []string{"id1", "id2", "id3"},
	})

	rig.Generate()
	rig.Generate()

	rig.Reset()

	if rig.CurrentIndex() != 0 {
		t.Error("reset should reset index")
	}

	if rig.CallCount() != 0 {
		t.Error("reset should reset call count")
	}

	if id := rig.Generate(); id != "id1" {
		t.Errorf("expected id1 after reset, got %s", id)
	}
}

func TestMockRequestIDGenerator_CurrentIndex(t *testing.T) {
	rig := mock.NewMockRequestIDGenerator(&mock.MockRequestIDGeneratorOptions{
		IDs: []string{"id1", "id2", "id3"},
	})

	if rig.CurrentIndex() != 0 {
		t.Error("expected index 0 initially")
	}

	rig.Generate()
	if rig.CurrentIndex() != 1 {
		t.Error("expected index 1 after first generate")
	}

	rig.Generate()
	if rig.CurrentIndex() != 2 {
		t.Error("expected index 2 after second generate")
	}
}

// MockBodyProvider Tests
func TestMockBodyProvider_DefaultValues(t *testing.T) {
	bp := mock.NewMockBodyProvider(nil)

	if bp.ContentType() != "application/octet-stream" {
		t.Errorf("expected default content type, got %s", bp.ContentType())
	}

	body, err := bp.Body()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	data, _ := io.ReadAll(body)
	if len(data) != 0 {
		t.Errorf("expected empty body by default, got %d bytes", len(data))
	}
}

func TestMockBodyProvider_CustomValues(t *testing.T) {
	bp := mock.NewMockBodyProvider(&mock.MockBodyProviderOptions{
		ContentType: "application/json",
		Body:        []byte(`{"test":true}`),
	})

	if bp.ContentType() != "application/json" {
		t.Errorf("expected application/json, got %s", bp.ContentType())
	}

	body, _ := bp.Body()
	data, _ := io.ReadAll(body)

	if string(data) != `{"test":true}` {
		t.Errorf("expected json body, got %s", data)
	}
}

func TestMockBodyProvider_SetBody(t *testing.T) {
	bp := mock.NewMockBodyProvider(nil)

	bp.SetBody([]byte(`new body`))

	body, _ := bp.Body()
	data, _ := io.ReadAll(body)

	if string(data) != "new body" {
		t.Errorf("expected 'new body', got %s", data)
	}
}

func TestMockBodyProvider_SetContentType(t *testing.T) {
	bp := mock.NewMockBodyProvider(nil)

	bp.SetContentType("text/plain")

	if bp.ContentType() != "text/plain" {
		t.Errorf("expected text/plain, got %s", bp.ContentType())
	}
}

func TestMockBodyProvider_InjectError(t *testing.T) {
	bp := mock.NewMockBodyProvider(nil)

	testErr := errors.New("io error")
	bp.InjectBodyError(testErr)

	_, err := bp.Body()

	if !errors.Is(err, testErr) {
		t.Errorf("expected injected error, got %v", err)
	}

	// Error should be cleared after one-shot
	body, err := bp.Body()
	if err != nil {
		t.Errorf("expected no error after one-shot, got %v", err)
	}

	if body == nil {
		t.Error("expected body after error cleared")
	}
}

func TestMockBodyProvider_CallCount(t *testing.T) {
	bp := mock.NewMockBodyProvider(nil)

	bp.Body()
	bp.Body()
	bp.Body()

	if bp.CallCount() != 3 {
		t.Errorf("expected 3 calls, got %d", bp.CallCount())
	}
}

func TestMockBodyProvider_Reset(t *testing.T) {
	bp := mock.NewMockBodyProvider(nil)

	bp.InjectBodyError(errors.New("test"))
	bp.Body()

	bp.Reset()

	if bp.CallCount() != 0 {
		t.Error("reset should clear call count")
	}

	_, err := bp.Body()
	if err != nil {
		t.Errorf("reset should clear errors, got %v", err)
	}
}

func TestMockBodyProvider_MultipleReads(t *testing.T) {
	bp := mock.NewMockBodyProvider(&mock.MockBodyProviderOptions{
		Body: []byte(`test data`),
	})

	// First read
	body1, _ := bp.Body()
	data1, _ := io.ReadAll(body1)

	// Second read
	body2, _ := bp.Body()
	data2, _ := io.ReadAll(body2)

	if string(data1) != "test data" {
		t.Errorf("first read: expected 'test data', got %s", data1)
	}

	if string(data2) != "test data" {
		t.Errorf("second read: expected 'test data', got %s", data2)
	}

	if bp.CallCount() != 2 {
		t.Errorf("expected 2 calls, got %d", bp.CallCount())
	}
}

func TestMockValidator_ClearRecorded(t *testing.T) {
	v := mock.NewMockValidator(&mock.MockValidatorOptions{
		RecordCalls: true,
	})

	v.Validate([]byte(`test`), "application/json")

	if len(v.ValidateCalls()) == 0 {
		t.Error("expected calls before clear")
	}

	v.ClearRecorded()

	if len(v.ValidateCalls()) != 0 {
		t.Error("expected calls to be cleared")
	}
}

func TestMockValidator_DisabledRecording(t *testing.T) {
	v := mock.NewMockValidator(&mock.MockValidatorOptions{
		RecordCalls: false,
	})

	v.Validate([]byte(`test`), "application/json")

	if len(v.ValidateCalls()) != 0 {
		t.Error("expected no calls when recording disabled")
	}
}
