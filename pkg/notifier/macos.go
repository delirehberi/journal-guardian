package notifier

import (
	"fmt"
	"os/exec"
)

type MacOSNotifier struct{}

func NewMacOSNotifier() *MacOSNotifier {
	return &MacOSNotifier{}
}

func (n *MacOSNotifier) Send(title, message string) error {
	if message == "" {
		return nil
	}
	// "display notification \"message\" with title \"title\""
	script := fmt.Sprintf("display notification \"%s\" with title \"%s\"", message, title)
	cmd := exec.Command("osascript", "-e", script)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error sending notification: %w", err)
	}
	return nil
}
