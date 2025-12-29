package mock_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kolosys/neuron"
	"github.com/kolosys/neuron/mock"
)

func TestMockRateLimiter_AllowByDefault(t *testing.T) {
	rl := mock.NewMockRateLimiter(nil)

	result := rl.Allow(context.Background(), "GET", "/api/users")

	if !result {
		t.Error("expected allow to return true by default")
	}
}

func TestMockRateLimiter_SetAllow(t *testing.T) {
	rl := mock.NewMockRateLimiter(nil)
	rl.SetAllow(false)

	result := rl.Allow(context.Background(), "GET", "/api/users")

	if result {
		t.Error("expected allow to return false after SetAllow(false)")
	}
}

func TestMockRateLimiter_SetAllowForEndpoint(t *testing.T) {
	rl := mock.NewMockRateLimiter(nil)

	rl.SetAllowForEndpoint("GET", "/api/users", false)
	rl.SetAllowForEndpoint("GET", "/api/posts", true)

	if rl.Allow(context.Background(), "GET", "/api/users") {
		t.Error("expected allow false for /api/users")
	}

	if !rl.Allow(context.Background(), "GET", "/api/posts") {
		t.Error("expected allow true for /api/posts")
	}
}

func TestMockRateLimiter_RecordsAllowCalls(t *testing.T) {
	rl := mock.NewMockRateLimiter(&mock.MockRateLimiterOptions{
		RecordCalls: true,
	})

	ctx := context.Background()
	rl.Allow(ctx, "GET", "/api/users")
	rl.Allow(ctx, "POST", "/api/users")

	calls := rl.AllowCalls()
	if len(calls) != 2 {
		t.Fatalf("expected 2 recorded calls, got %d", len(calls))
	}

	if calls[0].Method != "GET" || calls[0].Endpoint != "/api/users" {
		t.Error("first call details incorrect")
	}

	if calls[1].Method != "POST" || calls[1].Endpoint != "/api/users" {
		t.Error("second call details incorrect")
	}
}

func TestMockRateLimiter_DisabledRecording(t *testing.T) {
	rl := mock.NewMockRateLimiter(&mock.MockRateLimiterOptions{
		RecordCalls: false,
	})

	rl.Allow(context.Background(), "GET", "/api/users")

	calls := rl.AllowCalls()
	if len(calls) != 0 {
		t.Errorf("expected no recorded calls, got %d", len(calls))
	}
}

func TestMockRateLimiter_Wait(t *testing.T) {
	rl := mock.NewMockRateLimiter(nil)

	err := rl.Wait(context.Background(), "GET", "/api/users")

	if err != nil {
		t.Errorf("expected no error from Wait, got %v", err)
	}
}

func TestMockRateLimiter_WaitDuration(t *testing.T) {
	rl := mock.NewMockRateLimiter(nil)
	rl.SetWaitDuration(10 * time.Millisecond)

	start := time.Now()
	rl.Wait(context.Background(), "GET", "/api/users")
	elapsed := time.Since(start)

	if elapsed < 10*time.Millisecond {
		t.Errorf("expected at least 10ms wait, got %v", elapsed)
	}
}

func TestMockRateLimiter_WaitContextCancellation(t *testing.T) {
	rl := mock.NewMockRateLimiter(nil)
	rl.SetWaitDuration(100 * time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := rl.Wait(ctx, "GET", "/api/users")

	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

func TestMockRateLimiter_RecordsWaitCalls(t *testing.T) {
	rl := mock.NewMockRateLimiter(&mock.MockRateLimiterOptions{
		RecordCalls: true,
	})

	ctx := context.Background()
	rl.Wait(ctx, "GET", "/api/users")
	rl.Wait(ctx, "POST", "/api/posts")

	calls := rl.WaitCalls()
	if len(calls) != 2 {
		t.Fatalf("expected 2 recorded calls, got %d", len(calls))
	}

	if calls[0].Method != "GET" || calls[0].Endpoint != "/api/users" {
		t.Error("first wait call details incorrect")
	}
}

func TestMockRateLimiter_InjectWaitError(t *testing.T) {
	rl := mock.NewMockRateLimiter(nil)

	testErr := errors.New("rate limit exceeded")
	rl.InjectWaitError(testErr)

	err := rl.Wait(context.Background(), "GET", "/api/users")

	if !errors.Is(err, testErr) {
		t.Errorf("expected injected error, got %v", err)
	}

	// Error should be cleared after one-shot consumption
	err = rl.Wait(context.Background(), "GET", "/api/users")
	if err != nil {
		t.Errorf("expected no error after one-shot, got %v", err)
	}
}

func TestMockRateLimiter_ClearInjectedErrors(t *testing.T) {
	rl := mock.NewMockRateLimiter(nil)

	testErr := errors.New("rate limit exceeded")
	rl.InjectWaitError(testErr)
	rl.ClearInjectedErrors()

	err := rl.Wait(context.Background(), "GET", "/api/users")

	if err != nil {
		t.Errorf("expected no error after clear, got %v", err)
	}
}

func TestMockRateLimiter_ClearRecorded(t *testing.T) {
	rl := mock.NewMockRateLimiter(&mock.MockRateLimiterOptions{
		RecordCalls: true,
	})

	ctx := context.Background()
	rl.Allow(ctx, "GET", "/api/users")
	rl.Wait(ctx, "GET", "/api/users")

	if len(rl.AllowCalls()) == 0 || len(rl.WaitCalls()) == 0 {
		t.Error("expected recorded calls")
	}

	rl.ClearRecorded()

	if len(rl.AllowCalls()) != 0 || len(rl.WaitCalls()) != 0 {
		t.Error("expected all calls to be cleared")
	}
}

func TestMockRateLimiter_Reset(t *testing.T) {
	rl := mock.NewMockRateLimiter(nil)

	rl.SetAllow(false)
	rl.SetAllowForEndpoint("GET", "/api/users", true)
	rl.InjectWaitError(errors.New("test"))

	rl.Reset()

	if !rl.Allow(context.Background(), "GET", "/api/endpoint") {
		t.Error("reset should restore default allow=true")
	}

	err := rl.Wait(context.Background(), "GET", "/api/users")
	if err != nil {
		t.Errorf("reset should clear injected errors, got %v", err)
	}
}

func TestMockRateLimitHandler_UpdateFromHeaders(t *testing.T) {
	rlh := mock.NewMockRateLimitHandler(nil)

	info := &neuron.RateLimitInfo{
		Limit:     100,
		Remaining: 50,
	}

	err := rlh.UpdateFromHeaders("GET", "/api/users", info)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestMockRateLimitHandler_RecordsUpdates(t *testing.T) {
	rlh := mock.NewMockRateLimitHandler(nil)

	info1 := &neuron.RateLimitInfo{Limit: 100, Remaining: 90}
	info2 := &neuron.RateLimitInfo{Limit: 100, Remaining: 80}

	rlh.UpdateFromHeaders("GET", "/api/users", info1)
	rlh.UpdateFromHeaders("GET", "/api/users", info2)

	updates := rlh.Updates()
	if len(updates) != 2 {
		t.Fatalf("expected 2 updates, got %d", len(updates))
	}

	if updates[0].Info.Remaining != 90 {
		t.Errorf("first update remaining: expected 90, got %d", updates[0].Info.Remaining)
	}

	if updates[1].Info.Remaining != 80 {
		t.Errorf("second update remaining: expected 80, got %d", updates[1].Info.Remaining)
	}
}

func TestMockRateLimitHandler_WasExhausted(t *testing.T) {
	rlh := mock.NewMockRateLimitHandler(nil)

	info := &neuron.RateLimitInfo{Limit: 100, Remaining: 0}
	rlh.UpdateFromHeaders("GET", "/api/users", info)

	if !rlh.WasExhausted() {
		t.Error("expected WasExhausted to return true")
	}
}

func TestMockRateLimitHandler_LastUpdate(t *testing.T) {
	rlh := mock.NewMockRateLimitHandler(nil)

	if rlh.LastUpdate() != nil {
		t.Error("expected LastUpdate to be nil initially")
	}

	info := &neuron.RateLimitInfo{Limit: 100, Remaining: 50}
	rlh.UpdateFromHeaders("GET", "/api/users", info)

	last := rlh.LastUpdate()
	if last == nil {
		t.Fatal("expected LastUpdate to return an update")
	}

	if last.Info.Remaining != 50 {
		t.Errorf("expected remaining 50, got %d", last.Info.Remaining)
	}
}

func TestMockRateLimitHandler_UpdatesForEndpoint(t *testing.T) {
	rlh := mock.NewMockRateLimitHandler(nil)

	info := &neuron.RateLimitInfo{Limit: 100, Remaining: 50}
	rlh.UpdateFromHeaders("GET", "/api/users", info)
	rlh.UpdateFromHeaders("GET", "/api/users", info)
	rlh.UpdateFromHeaders("POST", "/api/users", info)

	updates := rlh.UpdatesForEndpoint("GET", "/api/users")
	if len(updates) != 2 {
		t.Errorf("expected 2 updates for GET /api/users, got %d", len(updates))
	}
}

func TestMockRateLimitHandler_ClearRecordedUpdates(t *testing.T) {
	rlh := mock.NewMockRateLimitHandler(nil)

	info := &neuron.RateLimitInfo{Limit: 100}
	rlh.UpdateFromHeaders("GET", "/api/users", info)

	if len(rlh.Updates()) == 0 {
		t.Error("expected updates before clear")
	}

	rlh.ClearRecordedUpdates()

	if len(rlh.Updates()) != 0 {
		t.Error("expected updates to be cleared")
	}
}

func TestMockRateLimitHandler_Reset(t *testing.T) {
	rlh := mock.NewMockRateLimitHandler(nil)

	info := &neuron.RateLimitInfo{Limit: 100}
	rlh.UpdateFromHeaders("GET", "/api/users", info)

	if len(rlh.Updates()) == 0 {
		t.Error("expected updates before reset")
	}

	rlh.Reset()

	if len(rlh.Updates()) != 0 {
		t.Error("reset should clear updates")
	}
}

func TestMockRateLimiter_String(t *testing.T) {
	rl := mock.NewMockRateLimiter(nil)
	rl.SetAllow(false)
	rl.SetAllowForEndpoint("GET", "/api/users", true)

	str := rl.String()
	if str == "" {
		t.Error("expected non-empty string representation")
	}
}

func TestMockRateLimiter_ThreadSafety(t *testing.T) {
	rl := mock.NewMockRateLimiter(&mock.MockRateLimiterOptions{
		RecordCalls: true,
	})

	ctx := context.Background()

	// Concurrent Allow and Wait calls
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				rl.Allow(ctx, "GET", "/api/users")
				rl.Wait(ctx, "GET", "/api/users")
			}
		}()
	}

	time.Sleep(100 * time.Millisecond)

	if len(rl.AllowCalls()) != 100 {
		t.Errorf("expected 100 allow calls, got %d", len(rl.AllowCalls()))
	}
	if len(rl.WaitCalls()) != 100 {
		t.Errorf("expected 100 wait calls, got %d", len(rl.WaitCalls()))
	}
}
