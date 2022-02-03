package iot

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type SessionCtx struct {
	ID        string
	CreatedAt time.Time
}

type GidInSession struct {
	Gid   string
	Scope string
}

type SessionManager struct {
	cookieName  string
	lock        sync.Mutex
	sessionsWS  map[string]*SessionCtx
	gidPerUser  map[string]GidInSession
	maxlifetime int64
}

func (manager *SessionManager) sessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "default"
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (mg *SessionManager) StoreGid(user, gid, scope string) {
	mg.lock.Lock()
	defer mg.lock.Unlock()
	mg.gidPerUser[user] = GidInSession{Gid: gid, Scope: scope}
}

func (mg *SessionManager) GetGid(user string) (string, bool) {
	mg.lock.Lock()
	defer mg.lock.Unlock()
	if v, ok := mg.gidPerUser[user]; ok {
		return v.Gid, true
	}
	return "", false
}

func (mg *SessionManager) GetJwtScope(user string) (string, bool) {
	mg.lock.Lock()
	defer mg.lock.Unlock()
	if v, ok := mg.gidPerUser[user]; ok {
		return v.Scope, true
	}
	return "", false
}

func (mg *SessionManager) GetSession(w http.ResponseWriter, r *http.Request) (*SessionCtx, error) {
	mg.lock.Lock()
	defer mg.lock.Unlock()
	var sn string
	cookie, err := r.Cookie(mg.cookieName)
	if err != nil || cookie.Value == "" {
		sn = mg.sessionId()
		log.Printf("Set a new session %s", sn)
		cookie := http.Cookie{Name: mg.cookieName, Value: url.QueryEscape(sn), Path: "/", HttpOnly: true, MaxAge: int(mg.maxlifetime)}
		http.SetCookie(w, &cookie)
	} else {
		sn, _ = url.QueryUnescape(cookie.Value)
		//fmt.Println("Using cookie value ", cookie.Value)
	}
	var session *SessionCtx
	var ok bool
	if session, ok = mg.sessionsWS[sn]; !ok {
		session = &SessionCtx{
			CreatedAt: time.Now(),
		}
		mg.sessionsWS[sn] = session
		//fmt.Println("Insert the new session data ", sn)
	}
	return session, nil
}

func (mg *SessionManager) GC() {
	mg.lock.Lock()
	defer mg.lock.Unlock()
	//log.Printf("Session GC call after %d seconds", manager.maxlifetime)
	destrKeys := []string{}
	now := time.Now()
	var max float64
	max = float64(mg.maxlifetime)
	for k, v := range mg.sessionsWS {
		elapsed := now.Sub(v.CreatedAt)
		if elapsed.Seconds() > max {
			destrKeys = append(destrKeys, k)
		}
	}
	for _, k := range destrKeys {
		log.Printf("Delete session %s", k)
		delete(mg.sessionsWS, k)
	}
	seconds, _ := time.ParseDuration(fmt.Sprintf("%ds", mg.maxlifetime))
	time.AfterFunc(time.Duration(seconds), func() {
		mg.GC()
	})
}

func InitSession() {
	go sessMgr.GC()
}
