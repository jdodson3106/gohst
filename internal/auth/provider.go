package auth

import (
	"os/exec"
	"time"
)

// SessionProvider interface defines the behavior
// that an auth session provider must adhere to
type SessionProvider interface {
	CreateProfile() error
	ListProfiles() ([]Profile, error)
	Login() error
	Logout() error
}

type GohstProvider struct {
	Profiles []Profile
	Created  time.Time
	LastUsed time.Time
}

func NewGohstProvider() *GohstProvider {
	return &GohstProvider{
		Created:  time.Now(),
		LastUsed: time.Now(),
	}
}

func (g *GohstProvider) CreateProfile() error {
	// TODO: Define how to store profile data
	return nil
}

func (g *GohstProvider) ListProfiles() ([]Profile, error) {
	if len(g.Profiles) > 0 {
		return g.Profiles, nil
	}

	return nil, nil
}

func (g *GohstProvider) Login() error {
	return nil
}

func (g *GohstProvider) Logout() error {
	return nil
}

type AwsProvider struct {
	Profiles []Profile
}

func (a *AwsProvider) CreateProfile() error {
	return nil
}

func (a *AwsProvider) ListProfiles() ([]Profile, error) {

	// if we have profiles, just return those
	if len(a.Profiles) > 0 {
		return a.Profiles, nil
	}

	// otherwise load the aws profiles
	cmd := exec.Command("aws", "configure", "list-profiles")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	profiles := make([]Profile, 0)
	line := make([]byte, 0)
	for _, v := range out {
		if v == '\n' {
			if len(line) > 0 {
				profiles = append(profiles, Profile{Name: string(line)})
				line = nil
			}
			continue
		}
		line = append(line, v)
	}

	return profiles, nil
}

func (a *AwsProvider) Login() error {
	return nil
}

func (a *AwsProvider) Logout() error {
	return nil
}
