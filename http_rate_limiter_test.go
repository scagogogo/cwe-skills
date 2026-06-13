package cwe

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestNewRateLimiter(t *testing.T) {
	rl := NewRateLimiter(10.0, 5)
	if rl == nil {
		t.Fatal("NewRateLimiter returned nil")
	}
	if rl.GetRate() != 10.0 {
		t.Errorf("rate = %v, want 10.0", rl.GetRate())
	}
	if rl.GetBurst() != 5 {
		t.Errorf("burst = %d, want 5", rl.GetBurst())
	}
	if rl.GetInterval() != time.Duration(float64(time.Second)/10.0) {
		t.Errorf("interval = %v, want %v", rl.GetInterval(), time.Duration(float64(time.Second)/10.0))
	}
	// Initial tokens should equal burst
	if rl.Tokens() != 5.0 {
		t.Errorf("initial tokens = %v, want 5.0", rl.Tokens())
	}
}

func TestAllow_SuccessWhenTokensAvailable(t *testing.T) {
	rl := NewRateLimiter(1.0, 3)
	for i := 0; i < 3; i++ {
		if !rl.Allow() {
			t.Errorf("Allow() call %d should succeed, tokens available", i+1)
		}
	}
}

func TestAllow_FailWhenTokensDepleted(t *testing.T) {
	rl := NewRateLimiter(1.0, 2)
	rl.Allow()
	rl.Allow()
	if rl.Allow() {
		t.Error("Allow() should fail when tokens depleted")
	}
}

func TestAllow_RecoverAfterTime(t *testing.T) {
	rl := NewRateLimiter(1000.0, 1) // very fast rate, burst of 1
	rl.Allow()
	// Tokens should be 0 now
	if rl.Allow() {
		t.Error("Allow() should fail immediately after depletion")
	}
	// Wait for token refill
	time.Sleep(20 * time.Millisecond)
	if !rl.Allow() {
		t.Error("Allow() should succeed after waiting for refill")
	}
}

func TestWait_SuccessImmediatelyWhenTokensAvailable(t *testing.T) {
	rl := NewRateLimiter(1.0, 5)
	ctx := context.Background()
	err := rl.Wait(ctx)
	if err != nil {
		t.Errorf("Wait() should succeed immediately when tokens available, got error: %v", err)
	}
}

func TestWait_RespectsContextCancellation(t *testing.T) {
	rl := NewRateLimiter(1.0, 1)
	rl.Allow() // consume the only token

	ctx, cancel := context.WithCancel(context.Background())
	// Cancel context after a short delay
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	err := rl.Wait(ctx)
	if err == nil {
		t.Error("Wait() should return error when context is cancelled")
	}
	if ctx.Err() == nil {
		t.Error("context should be cancelled")
	}
}

func TestWaitForRequest_FirstCallReturnsImmediately(t *testing.T) {
	rl := NewRateLimiter(1.0, 5)
	start := time.Now()
	rl.WaitForRequest()
	elapsed := time.Since(start)
	if elapsed > 100*time.Millisecond {
		t.Errorf("first WaitForRequest() took %v, should be immediate", elapsed)
	}
}

func TestWaitForRequest_SubsequentCallsRespectInterval(t *testing.T) {
	rl := NewRateLimiter(100.0, 5) // 100 per second = 10ms interval
	rl.WaitForRequest()
	start := time.Now()
	rl.WaitForRequest()
	elapsed := time.Since(start)
	// Should have waited at least some portion of the interval
	// Being generous with the check due to timing
	if elapsed < time.Millisecond {
		t.Errorf("second WaitForRequest() took %v, expected some wait", elapsed)
	}
}

func TestGetInterval_SetInterval(t *testing.T) {
	rl := NewRateLimiter(1.0, 5)
	original := rl.GetInterval()
	newInterval := 500 * time.Millisecond
	rl.SetInterval(newInterval)
	got := rl.GetInterval()
	if got != newInterval {
		t.Errorf("GetInterval() after SetInterval = %v, want %v", got, newInterval)
	}
	// Verify that rate was also updated
	newRate := rl.GetRate()
	expectedRate := float64(time.Second) / float64(newInterval)
	if newRate != expectedRate {
		t.Errorf("GetRate() after SetInterval = %v, want %v", newRate, expectedRate)
	}
	_ = original
}

func TestGetRate(t *testing.T) {
	rl := NewRateLimiter(5.0, 10)
	if rl.GetRate() != 5.0 {
		t.Errorf("GetRate() = %v, want 5.0", rl.GetRate())
	}
}

func TestGetBurst(t *testing.T) {
	rl := NewRateLimiter(1.0, 7)
	if rl.GetBurst() != 7 {
		t.Errorf("GetBurst() = %d, want 7", rl.GetBurst())
	}
}

func TestResetLastRequest(t *testing.T) {
	rl := NewRateLimiter(1.0, 1)
	// Make a request
	rl.Allow()
	// Should be out of tokens
	if rl.Allow() {
		t.Error("should be out of tokens")
	}
	// Reset
	rl.ResetLastRequest()
	// Now Allow should work again because tokens are refilled
	if !rl.Allow() {
		t.Error("Allow() should succeed after ResetLastRequest()")
	}
}

func TestTokens(t *testing.T) {
	rl := NewRateLimiter(1.0, 5)
	tokens := rl.Tokens()
	if tokens < 4.99 || tokens > 5.01 {
		t.Errorf("Tokens() = %v, want approximately 5.0", tokens)
	}
	rl.Allow()
	tokens = rl.Tokens()
	if tokens < 3.99 || tokens > 4.01 {
		t.Errorf("Tokens() after one Allow = %v, want approximately 4.0", tokens)
	}
}

func TestTokens_RefillCap(t *testing.T) {
	rl := NewRateLimiter(1.0, 3)
	// Consume all tokens
	for i := 0; i < 3; i++ {
		rl.Allow()
	}
	// Wait for some refill (longer than burst capacity would need)
	time.Sleep(5 * time.Second)
	// Tokens should be capped at burst
	tokens := rl.Tokens()
	if tokens > 3.0 {
		t.Errorf("Tokens() should be capped at burst=3, got %v", tokens)
	}
}

func TestConcurrentAccess(t *testing.T) {
	rl := NewRateLimiter(1000.0, 100)
	var wg sync.WaitGroup
	successCount := 0
	var mu sync.Mutex

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if rl.Allow() {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	mu.Lock()
	count := successCount
	mu.Unlock()

	if count == 0 {
		t.Error("expected some successful Allow() calls from concurrent goroutines")
	}
	if count > 100 {
		t.Errorf("successCount = %d, should not exceed burst", count)
	}
}

func TestConcurrentWait(t *testing.T) {
	rl := NewRateLimiter(1000.0, 100)
	var wg sync.WaitGroup
	errCount := 0
	var mu sync.Mutex

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := rl.Wait(context.Background())
			if err != nil {
				mu.Lock()
				errCount++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	mu.Lock()
	count := errCount
	mu.Unlock()

	if count > 0 {
		t.Errorf("expected no errors from concurrent Wait(), got %d", count)
	}
}

func TestSetInterval_ZeroDoesNotPanic(t *testing.T) {
	rl := NewRateLimiter(1.0, 5)
	// Setting interval to 0 should not panic; rate should remain unchanged
	// because the implementation only updates rate if interval > 0
	originalRate := rl.GetRate()
	rl.SetInterval(0)
	// Rate should remain the same since interval is 0
	if rl.GetRate() != originalRate {
		t.Errorf("rate changed after SetInterval(0), was %v now %v", originalRate, rl.GetRate())
	}
}
