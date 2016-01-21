package sessions

import (
	"errors"
	"net/http"
	"time"
)

var (
	sid         string        = "go.sid"
	maxlifetime time.Duration = time.Minute * 20
)

func init() {
	ticker := time.NewTicker(time.Minute)
	go func() {
		var time_now time.Time

		for {
			time_now = <-ticker.C
			registry.remove_old(time_now)
		}
	}()
}

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

/*
	получить актуальную сессию
*/
func Get(w http.ResponseWriter, r *http.Request) (*Session, error) {
	var (
		ses *Session
		err error
	)

	cookie, err := r.Cookie(sid)
	if err != nil {
		// нет куки, поэтому создаем новую сессию
		return sessionStart(w, "")
	} else {
		ses, err = registry.get(cookie.Value)

		if err != nil {
			// если сессии нет, то создадим новую
			return sessionStart(w, cookie.Value)
		} else if !ses.is_actual(time.Now()) {
			// если сессия устарела, то очистим и создадим новую
			registry.delete(cookie.Value)
			return sessionStart(w, "")
		}

		// существует актуальная сессия. Обновим время актуальности
		ses.update_last_time()

		cookie.MaxAge = int(maxlifetime.Seconds())
		//cookie
		http.SetCookie(w, cookie)

		return ses, nil
	}
}

/*
	новая сессия
*/
func sessionStart(w http.ResponseWriter, ses_id string) (*Session, error) {
	session := &Session{
		values:    make(map[string]interface{}),
		last_time: time.Now(),
	}

	if ses_id == "" {
		var err error

		for i := 0; i < 100; i++ {
			ses_id = registry.create_ses_id()

			if _, err = registry.get(ses_id); err != nil { // если отсутствует, то ок
				break
			}
		}

		if err == nil { // если существует сессия
			return nil, errors.New("not available")
		}
	}

	cookie := http.Cookie{Name: sid, Value: ses_id, Path: "/", HttpOnly: true, MaxAge: int(maxlifetime.Seconds())}
	http.SetCookie(w, &cookie)

	registry.add(ses_id, session)

	return session, nil
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
