package sessions

import (
	"net/http"
	"time"
)

var (
	sid         string        = "go.sid"
	maxlifetime time.Duration = time.Minute * 20
)

func SetMaxLifeTime(duration time.Duration) time.Duration {
	if duration > time.Second {
		maxlifetime = duration
	}

	return maxlifetime
}

func SetSid(name string) string {
	if name != "" {
		sid = name
	}

	return sid
}
func GetSid() string {
	return sid
}

/*
	получить актуальную сессию
*/
func Get(w http.ResponseWriter, r *http.Request) *Session {
	var (
		err    error
		cookie *http.Cookie
		ses    *Session
		ok     bool
		ses_id string
	)

	cookie, err = r.Cookie(sid)
	if err != nil {
		ses_id, ses = registry.new()

		cookie = &http.Cookie{Name: sid, Value: ses_id, Path: "/", HttpOnly: true, MaxAge: int(maxlifetime.Seconds())}
		http.SetCookie(w, cookie)

		return ses
	}

	ses, ok = registry.get(cookie.Value)
	if !ok {
		ses = registry.create(cookie.Value)

		cookie.MaxAge = int(maxlifetime.Seconds())
		cookie.Path = "/"
		http.SetCookie(w, cookie)

		return ses
	}

	cookie.MaxAge = int(maxlifetime.Seconds())
	http.SetCookie(w, cookie)

	return ses
}

/*
	удаление сессии и неактулизируем куки
*/
func Delete(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie(sid)
	if err == nil {
		registry.delete(cookie.Value)
	}

	cookie.MaxAge = -1
	cookie.Expires = time.Now()
	http.SetCookie(w, cookie)

	return nil
}
