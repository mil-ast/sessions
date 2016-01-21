package sessions

import (
	"sync"
	"time"
)

type Session struct {
	values    map[string]interface{} // int, float64, bool, string
	last_time time.Time
	mu        sync.RWMutex
}

func (this Session) Get(key string) interface{} {
	this.mu.RLock()
	value := this.values[key]

	this.mu.RUnlock()
	return value
}

func (this Session) GetString(key string) string {
	this.mu.RLock()
	defer this.mu.RUnlock()

	if _, ok := this.values[key]; !ok {
		return ""
	}

	value := this.values[key]
	switch value.(type) {
	case string:
		return value.(string)
	}

	return ""
}

func (this Session) GetInt(key string) int {
	this.mu.RLock()
	defer this.mu.RUnlock()

	if _, ok := this.values[key]; !ok {
		return 0
	}

	value := this.values[key]
	switch value.(type) {
	case int:
		return value.(int)
	}

	return 0
}

func (this Session) GetUint64(key string) uint64 {
	this.mu.RLock()
	defer this.mu.RUnlock()

	if _, ok := this.values[key]; !ok {
		return 0
	}

	value := this.values[key]
	switch value.(type) {
	case uint64:
		return value.(uint64)
	}

	return 0
}

func (this Session) GetBool(key string) bool {
	this.mu.RLock()
	defer this.mu.RUnlock()

	if _, ok := this.values[key]; !ok {
		return false
	}

	value := this.values[key]
	switch value.(type) {
	case bool:
		return value.(bool)
	}

	return false
}

func (this Session) GetFloat64(key string) float64 {
	this.mu.RLock()
	defer this.mu.RUnlock()

	if _, ok := this.values[key]; !ok {
		return 0.0
	}

	value := this.values[key]
	switch value.(type) {
	case float64:
		return value.(float64)
	}

	return 0.0
}

/*
	проверка существование ключа
*/
func (this Session) Exists(key string) bool {
	this.mu.RLock()
	_, ok := this.values[key]
	this.mu.RUnlock()

	return ok
}

func (this *Session) Set(key string, value interface{}) {
	this.mu.Lock()
	this.values[key] = value
	this.mu.Unlock()
}

func (this *Session) SetMap(values map[string]interface{}) {
	this.mu.Lock()

	for key, value := range values {
		this.values[key] = value
	}

	this.mu.Unlock()
}

func (this *Session) Delete(key string) {
	this.mu.Lock()
	delete(this.values, key)
	this.mu.Unlock()
}

func (this *Session) update_last_time() {
	this.mu.Lock()
	this.last_time = time.Now()
	this.mu.Unlock()
}

func (this Session) is_actual(time_now time.Time) bool {
	this.mu.RLock()
	actual := time_now.Before(this.last_time.Add(maxlifetime))
	this.mu.RUnlock()

	return actual
}
