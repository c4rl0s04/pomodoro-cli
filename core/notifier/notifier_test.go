package notifier

import (
	"testing"
)

// ensure BeeepNotifier implements Notifier at compile-time
var _ Notifier = (*BeeepNotifier)(nil)

func TestBeeepNotifier_ImplementsInterface(t *testing.T) {
	// A simple structural test to ensure our notifier implements the interface
	n := NewBeeepNotifier()
	if n == nil {
		t.Error("expected non-nil Notifier")
	}
}
