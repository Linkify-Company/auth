package memory

import (
	"sync"
	"time"
)

type EmailMemory struct {
	// Email / *authData
	emailItems map[string]*authData
	mx         sync.RWMutex
}

const (
	// Количество времени, за которое должны подтвердить почту
	timeout time.Duration = time.Minute * 10
	// Максимальное допустимое количество попыток дозволеных на подтверждение email
	defaultNumberAttempts int = 3
)

type authData struct {
	// Время отправки кода на почту
	timeCreate time.Time
	// Количество неправильных попыток, затраченных на введение кода
	numberAttempts    int
	authorizationCode int
}

func NewEmailMemory() *EmailMemory {
	return &EmailMemory{
		emailItems: make(map[string]*authData),
		mx:         sync.RWMutex{},
	}
}

func (m *EmailMemory) Set(email string, authorizationCode int) {
	m.mx.Lock()
	defer m.mx.Unlock()

	m.emailItems[email] = &authData{
		timeCreate:        time.Now(),
		authorizationCode: authorizationCode,
	}
}

func (m *EmailMemory) IsValid(email string, authorizationCode int) bool {
	m.mx.Lock()
	defer m.mx.Unlock()

	if m.emailItems[email] == nil {
		return false
	}

	if m.emailItems[email].timeCreate.Sub(time.Now()) >= timeout {
		m.remove(email)
		return false
	}
	if m.emailItems[email].authorizationCode != authorizationCode {
		m.emailItems[email].numberAttempts++

		if m.emailItems[email].numberAttempts >= defaultNumberAttempts {
			m.remove(email)
		}
		return false
	}
	m.remove(email)

	return true
}

func (m *EmailMemory) IsExist(email string) bool {
	m.mx.RLock()
	defer m.mx.RUnlock()

	return m.emailItems[email] != nil
}

func (m *EmailMemory) remove(email string) {
	delete(m.emailItems, email)
}
