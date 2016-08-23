package utils

import (
	"errors"
	"time"
)

// DeadlineActionFunc is a function called by enforcer
type DeadlineActionFunc func(success chan<- struct{}, terminate <-chan struct{})

// DeadlineEnforcer encforces deadlines
type DeadlineEnforcer struct {
	// Action that enforcec end
	Action DeadlineActionFunc
}

// DoWithin do with duration
func (deadlineEnforcer DeadlineEnforcer) DoWithin(duration time.Duration) error {
	success := make(chan struct{})
	terminate := make(chan struct{})

	go deadlineEnforcer.Action(success, terminate)
	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-success:
		return nil
	case <-timer.C:
		close(terminate)
		return errors.New("timeout")
	}
}
