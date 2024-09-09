package redis

import (
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)

type ConnectionRegister func() redis.UniversalClient

type Manager struct {
	connectionRegisters map[string]ConnectionRegister
	connections         map[string]redis.UniversalClient
	defaultConnection   string
	mu                  sync.RWMutex
}

func (m *Manager) Register(name string, cr ConnectionRegister) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.connectionRegisters[name] = cr
}

func (m *Manager) GetConnection(name ...string) (redis.UniversalClient, error) {
	var con string
	if len(name) > 0 {
		con = name[0]
	} else {
		con = m.defaultConnection
	}

	m.mu.RLock()

	if c, ok := m.connections[con]; ok {
		m.mu.RUnlock()
		return c, nil
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()

	if cr, ok := m.connectionRegisters[con]; ok {
		m.connections[con] = cr()
		return m.connections[con], nil
	}

	return nil, fmt.Errorf("unregisterd redis connection [%s]", con)
}

func (m *Manager) Use(name string) *Manager {
	m.defaultConnection = name

	return m
}

func NewManager() *Manager {
	return &Manager{
		connectionRegisters: map[string]ConnectionRegister{},
		connections:         map[string]redis.UniversalClient{},
		defaultConnection:   "default",
		mu:                  sync.RWMutex{},
	}
}
