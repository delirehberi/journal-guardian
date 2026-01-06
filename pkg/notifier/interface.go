package notifier

// Notifier interface for sending desktop notifications
type Notifier interface {
	Send(title, message string) error
}
