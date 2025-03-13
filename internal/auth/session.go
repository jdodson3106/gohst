package auth

import (
	"context"
	"time"
)

type Session struct {
	Context    context.Context
	Provider   SessionProvider
	IsLoggedIn bool

	startTime      time.Time
	expirationTime time.Time
}

func NewSession(p SessionProvider) *Session {
	return &Session{
		Context:  context.Background(),
		Provider: p,
	}
}

func (s *Session) Profile() (p Profile, ok bool) {
	p, ok = s.Context.Value("profile").(Profile)
	if !ok {
		return Profile{}, false
	}

	return p, true
}

func (s *Session) ListProfiles() ([]string, error) {
	var profNames []string

	profs, err := s.Provider.ListProfiles()
	if err != nil {
		return profNames, err
	}

	for _, p := range profs {
		profNames = append(profNames, p.Name)
	}

	return profNames, nil
}
