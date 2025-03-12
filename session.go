package gohst

import (
	"context"
	"os/exec"
	"time"
)

type profile struct {
	name string
}

type session struct {
	isLoggedIn     bool
	startTime      time.Time
	expirationTime time.Time
	context        context.Context
	profiles       []profile
}

func (s *session) profile() (p string, ok bool) {
	p, ok = s.context.Value("profile").(string)
	if !ok || p == "" {
		return "", false
	}

	return p, true
}

func (s *session) listProfiles() ([]string, error) {
	var profNames []string

	// if there are not profiles, then load them
	if len(s.profiles) == 0 {
		profiles, err := s.loadProfiles()
		if err != nil {
			// TODO: I should probably log sys errors somewhere??
			return nil, err
		}
		s.profiles = profiles
	}

	for _, p := range s.profiles {
		profNames = append(profNames, p.name)
	}
	return profNames, nil
}

func (s *session) loadProfiles() ([]profile, error) {
	cmd := exec.Command("aws", "configure", "list-profiles")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	profiles := make([]profile, 0)
	line := make([]byte, 0)
	for _, v := range out {
		if v == '\n' {
			if len(line) > 0 {
				profiles = append(profiles, profile{name: string(line)})
				line = nil
			}
			continue
		}
		line = append(line, v)
	}

	return profiles, nil
}
