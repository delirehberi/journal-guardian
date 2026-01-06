package notifier

import (
	"fmt"
	"os/exec"
)

type LinuxNotifier struct{}

func NewLinuxNotifier() *LinuxNotifier {
	return &LinuxNotifier{}
}

func (n *LinuxNotifier) Send(title, message string) error {
	if message == "" {
		return nil
	}
	cmd := exec.Command("notify-send", title, message)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error sending notification: %w", err)
	}
	return nil
}
