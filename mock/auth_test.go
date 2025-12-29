package mock_test

import (
	"context"
	"errors"
	"testing"

	"github.com/kolosys/neuron/mock"
)

func TestMockAuthProvider_DefaultToken(t *testing.T) {
	ap := mock.NewMockAuthProvider(nil)

	token, err := ap.GetToken(context.Background())

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if token != "" {
		t.Errorf("expected empty token by default, got %s", token)
	}
}

func TestMockAuthProvider_InitialToken(t *testing.T) {
	ap := mock.NewMockAuthProvider(&mock.MockAuthProviderOptions{
		InitialToken: "test-token",
	})

	token, _ := ap.GetToken(context.Background())

	if token != "test-token" {
		t.Errorf("expected test-token, got %s", token)
	}
}

func TestMockAuthProvider_GetAuthHeader(t *testing.T) {
	ap := mock.NewMockAuthProvider(nil)

	header := ap.GetAuthHeader("test-token")

	expected := "Bearer test-token"
	if header != expected {
		t.Errorf("expected %s, got %s", expected, header)
	}
}

func TestMockAuthProvider_CustomHeaderFormat(t *testing.T) {
	ap := mock.NewMockAuthProvider(&mock.MockAuthProviderOptions{
		HeaderFormat: "X-API-Key {}",
	})

	header := ap.GetAuthHeader("secret-key")

	expected := "X-API-Key secret-key"
	if header != expected {
		t.Errorf("expected %s, got %s", expected, header)
	}
}

func TestMockAuthProvider_TokenRotation(t *testing.T) {
	ap := mock.NewMockAuthProvider(&mock.MockAuthProviderOptions{
		InitialToken: "token1",
		Tokens:       []string{"token2", "token3"},
	})

	// First call should return token1
	token1, _ := ap.GetToken(context.Background())
	if token1 != "token1" {
		t.Errorf("expected token1, got %s", token1)
	}

	// Manual rotation
	ap.RotateToken()

	// Second call should return token2
	token2, _ := ap.GetToken(context.Background())
	if token2 != "token2" {
		t.Errorf("expected token2, got %s", token2)
	}

	// Rotate again
	ap.RotateToken()

	// Third call should return token3
	token3, _ := ap.GetToken(context.Background())
	if token3 != "token3" {
		t.Errorf("expected token3, got %s", token3)
	}

	// Rotate past end - should stay at last token
	ap.RotateToken()
	token4, _ := ap.GetToken(context.Background())
	if token4 != "token3" {
		t.Errorf("expected token3 (last), got %s", token4)
	}
}

func TestMockAuthProvider_RecordsGetTokenCalls(t *testing.T) {
	ap := mock.NewMockAuthProvider(&mock.MockAuthProviderOptions{
		RecordCalls:  true,
		InitialToken: "test-token",
	})

	ctx := context.Background()
	ap.GetToken(ctx)
	ap.GetToken(ctx)

	calls := ap.GetTokenCalls()
	if len(calls) != 2 {
		t.Fatalf("expected 2 calls, got %d", len(calls))
	}

	if calls[0].Result != "test-token" {
		t.Errorf("first call result: expected test-token, got %s", calls[0].Result)
	}

	if calls[1].Result != "test-token" {
		t.Errorf("second call result: expected test-token, got %s", calls[1].Result)
	}
}

func TestMockAuthProvider_RecordsGetHeaderCalls(t *testing.T) {
	ap := mock.NewMockAuthProvider(&mock.MockAuthProviderOptions{
		RecordCalls: true,
	})

	ap.GetAuthHeader("token1")
	ap.GetAuthHeader("token2")

	calls := ap.GetHeaderCalls()
	if len(calls) != 2 {
		t.Fatalf("expected 2 calls, got %d", len(calls))
	}

	if calls[0].Token != "token1" {
		t.Errorf("first call token: expected token1, got %s", calls[0].Token)
	}

	if calls[1].Token != "token2" {
		t.Errorf("second call token: expected token2, got %s", calls[1].Token)
	}
}

func TestMockAuthProvider_InjectTokenError(t *testing.T) {
	ap := mock.NewMockAuthProvider(&mock.MockAuthProviderOptions{
		InitialToken: "valid-token",
	})

	testErr := errors.New("auth failed")
	ap.InjectTokenError(testErr)

	_, err := ap.GetToken(context.Background())

	if !errors.Is(err, testErr) {
		t.Errorf("expected injected error, got %v", err)
	}

	// Error should be cleared after one-shot consumption
	token, err := ap.GetToken(context.Background())
	if err != nil {
		t.Errorf("expected no error after one-shot, got %v", err)
	}
	if token != "valid-token" {
		t.Errorf("expected valid-token after error cleared, got %s", token)
	}
}

func TestMockAuthProvider_ClearInjectedErrors(t *testing.T) {
	ap := mock.NewMockAuthProvider(&mock.MockAuthProviderOptions{
		InitialToken: "test-token",
	})

	testErr := errors.New("auth failed")
	ap.InjectTokenError(testErr)
	ap.ClearInjectedErrors()

	token, err := ap.GetToken(context.Background())

	if err != nil {
		t.Errorf("expected no error after clear, got %v", err)
	}

	if token != "test-token" {
		t.Errorf("expected test-token, got %s", token)
	}
}

func TestMockAuthProvider_ContextCancellation(t *testing.T) {
	ap := mock.NewMockAuthProvider(nil)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := ap.GetToken(ctx)

	if err != context.Canceled {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

func TestMockAuthProvider_SetTokens(t *testing.T) {
	ap := mock.NewMockAuthProvider(nil)

	ap.SetTokens([]string{"new1", "new2", "new3"})

	token1, _ := ap.GetToken(context.Background())
	if token1 != "new1" {
		t.Errorf("expected new1, got %s", token1)
	}

	ap.RotateToken()
	token2, _ := ap.GetToken(context.Background())
	if token2 != "new2" {
		t.Errorf("expected new2, got %s", token2)
	}
}

func TestMockAuthProvider_CurrentTokenIndex(t *testing.T) {
	ap := mock.NewMockAuthProvider(&mock.MockAuthProviderOptions{
		InitialToken: "token1",
		Tokens:       []string{"token2", "token3"},
	})

	if ap.CurrentTokenIndex() != 0 {
		t.Error("expected index 0 initially")
	}

	ap.RotateToken()
	if ap.CurrentTokenIndex() != 1 {
		t.Error("expected index 1 after rotation")
	}

	ap.RotateToken()
	if ap.CurrentTokenIndex() != 2 {
		t.Error("expected index 2 after second rotation")
	}
}

func TestMockAuthProvider_Reset(t *testing.T) {
	ap := mock.NewMockAuthProvider(&mock.MockAuthProviderOptions{
		RecordCalls:  true,
		InitialToken: "token1",
		Tokens:       []string{"token2", "token3"},
	})

	ap.GetToken(context.Background())
	ap.RotateToken()
	ap.GetToken(context.Background())
	ap.InjectTokenError(errors.New("test"))

	if len(ap.GetTokenCalls()) == 0 {
		t.Error("expected recorded calls before reset")
	}

	ap.Reset()

	if ap.CurrentTokenIndex() != 0 {
		t.Error("reset should return to first token")
	}

	if len(ap.GetTokenCalls()) != 0 {
		t.Error("reset should clear recorded calls")
	}

	token, err := ap.GetToken(context.Background())
	if err != nil {
		t.Errorf("reset should clear errors, got %v", err)
	}

	if token != "token1" {
		t.Errorf("expected token1 after reset, got %s", token)
	}
}

func TestMockAuthProvider_ClearRecorded(t *testing.T) {
	ap := mock.NewMockAuthProvider(&mock.MockAuthProviderOptions{
		RecordCalls:  true,
		InitialToken: "test-token",
	})

	ap.GetToken(context.Background())
	ap.GetAuthHeader("test-token")

	if len(ap.GetTokenCalls()) == 0 || len(ap.GetHeaderCalls()) == 0 {
		t.Error("expected recorded calls")
	}

	ap.ClearRecorded()

	if len(ap.GetTokenCalls()) != 0 || len(ap.GetHeaderCalls()) != 0 {
		t.Error("expected all calls to be cleared")
	}
}

func TestMockAuthProvider_SetHeaderFormat(t *testing.T) {
	ap := mock.NewMockAuthProvider(nil)

	ap.SetHeaderFormat("Authorization: {}") // e.g., for custom format
	header := ap.GetAuthHeader("token123")

	if header != "Authorization: token123" {
		t.Errorf("expected 'Authorization: token123', got %s", header)
	}
}

func TestMockAuthProvider_RecordErrorCalls(t *testing.T) {
	ap := mock.NewMockAuthProvider(&mock.MockAuthProviderOptions{
		RecordCalls: true,
	})

	testErr := errors.New("auth failed")
	ap.InjectTokenError(testErr)

	ap.GetToken(context.Background())

	calls := ap.GetTokenCalls()
	if len(calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(calls))
	}

	if !errors.Is(calls[0].Err, testErr) {
		t.Errorf("expected error to be recorded")
	}
}

func TestMockAuthProvider_DisabledRecording(t *testing.T) {
	ap := mock.NewMockAuthProvider(&mock.MockAuthProviderOptions{
		RecordCalls: false,
	})

	ap.GetToken(context.Background())
	ap.GetAuthHeader("token")

	if len(ap.GetTokenCalls()) != 0 {
		t.Error("expected no recorded calls when recording disabled")
	}

	if len(ap.GetHeaderCalls()) != 0 {
		t.Error("expected no recorded header calls when recording disabled")
	}
}

func TestMockAuthProvider_ThreadSafety(t *testing.T) {
	ap := mock.NewMockAuthProvider(&mock.MockAuthProviderOptions{
		RecordCalls:  true,
		InitialToken: "token1",
		Tokens:       []string{"token2", "token3"},
	})

	ctx := context.Background()

	// Concurrent GetToken and RotateToken calls
	done := make(chan bool, 5)
	for i := 0; i < 5; i++ {
		go func() {
			for j := 0; j < 5; j++ {
				ap.GetToken(ctx)
				ap.RotateToken()
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 5; i++ {
		<-done
	}

	calls := ap.GetTokenCalls()
	if len(calls) != 25 {
		t.Errorf("expected 25 calls, got %d", len(calls))
	}
}
