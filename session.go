package gohst

import (
	"context"
	"time"
)

type session struct {
	isLoggedIn     bool
	startTime      time.Time
	expirationTime time.Time
	context        context.Context
}

func (s *session) profile() (p string, ok bool) {
	p, ok = s.context.Value("profile").(string)
	if !ok || p == "" {
		return "", false
	}

	return p, true
}
