package notifier

import (
	"github.com/gen2brain/beeep"
)

// Notifier defines the interface for sending desktop notifications
type Notifier interface {
	Notify(title, message string) error
}

// BeeepNotifier implements Notifier using the gen2brain/beeep cross-platform library
type BeeepNotifier struct{}

// NewBeeepNotifier creates a new BeeepNotifier
func NewBeeepNotifier() *BeeepNotifier {
	return &BeeepNotifier{}
}

// Notify sends a desktop notification and plays the default system sound
func (n *BeeepNotifier) Notify(title, message string) error {
	return beeep.Alert(title, message, "")
}
