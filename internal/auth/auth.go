package auth

import (
	"errors"
	"sync"
)

type Auth struct {
	mu       sync.RWMutex
	password string
}

func New(password string) *Auth {
	return &Auth{password: password}
}

func (a *Auth) Authenticate(pass string) error {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.password != "" && pass != a.password {
		return errors.New("ERR invalid password")
	}
	return nil
}

func (a *Auth) UpdatePassword(newPass string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.password = newPass
}
