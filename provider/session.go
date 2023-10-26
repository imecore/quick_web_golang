package provider

import (
	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"net/http"
	"quick_web_golang/config"
	"quick_web_golang/lib"
	"strconv"
	"time"
)

type Session struct {
	Manager *scs.SessionManager
}

func (s *Session) New() *Session {
	s.Manager = &scs.SessionManager{}
	return s
}

func (s *Session) Start() {
	if lib.IsDev() {
		s.Manager.Cookie.Secure = true
	}
	lifeTime := 24 * time.Hour
	if day, err := strconv.Atoi(config.Get(config.SessionLifeDay)); err == nil && day > 0 {
		lifeTime = time.Duration(day) * 24 * time.Hour
	}
	s.Manager.Cookie.Name = config.Get(config.CookieName)
	s.Manager.Cookie.SameSite = http.SameSiteNoneMode
	s.Manager.Cookie.Persist = true
	s.Manager.Lifetime = lifeTime
	s.Manager.Store = redisstore.New(Cache.Pool)
}

func (s *Session) Close() {
	return
}
