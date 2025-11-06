package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
	"zpwoot/internal/model"
	"zpwoot/pkg/logger"
)

// SessionCache gerencia cache de informações de sessões
type SessionCache struct {
	cache *cache.Cache
}

// SessionInfo representa as informações de uma sessão no cache
type SessionInfo struct {
	data map[string]string
}

// NewSessionInfo cria uma nova instância de SessionInfo
func NewSessionInfo() *SessionInfo {
	return &SessionInfo{
		data: make(map[string]string),
	}
}

// Get obtém um valor do SessionInfo
func (s *SessionInfo) Get(key string) string {
	if s.data == nil {
		return ""
	}
	return s.data[key]
}

// Set define um valor no SessionInfo
func (s *SessionInfo) Set(key, value string) {
	if s.data == nil {
		s.data = make(map[string]string)
	}
	s.data[key] = value
}

// GetData retorna todos os dados
func (s *SessionInfo) GetData() map[string]string {
	return s.data
}

var (
	// Cache global de sessões (adaptado do wmial.bak para sessões)
	sessionInfoCache *SessionCache
)

// InitSessionCache inicializa o cache de sessões
func InitSessionCache() {
	sessionInfoCache = &SessionCache{
		cache: cache.New(cache.NoExpiration, 10*time.Minute),
	}
	logger.Log.Info().Msg("Session cache initialized")
}

// GetSessionCache retorna a instância do cache
func GetSessionCache() *SessionCache {
	if sessionInfoCache == nil {
		InitSessionCache()
	}
	return sessionInfoCache
}

// Set armazena informações da sessão no cache
func (sc *SessionCache) Set(sessionID string, sessionInfo *SessionInfo) {
	sc.cache.Set(sessionID, sessionInfo, cache.NoExpiration)
}

// Get obtém informações da sessão do cache
func (sc *SessionCache) Get(sessionID string) (*SessionInfo, bool) {
	if item, found := sc.cache.Get(sessionID); found {
		if sessionInfo, ok := item.(*SessionInfo); ok {
			return sessionInfo, true
		}
	}
	return nil, false
}

// Delete remove informações da sessão do cache
func (sc *SessionCache) Delete(sessionID string) {
	sc.cache.Delete(sessionID)
}

// UpdateSessionInfo atualiza uma informação específica da sessão
func (sc *SessionCache) UpdateSessionInfo(sessionID, key, value string) *SessionInfo {
	sessionInfo, found := sc.Get(sessionID)
	if !found {
		sessionInfo = NewSessionInfo()
	}
	sessionInfo.Set(key, value)
	sc.Set(sessionID, sessionInfo)
	return sessionInfo
}

// CreateSessionInfoFromModel cria SessionInfo a partir do modelo Session
func CreateSessionInfoFromModel(session *model.Session) *SessionInfo {
	sessionInfo := NewSessionInfo()
	sessionInfo.Set("Id", session.ID)
	sessionInfo.Set("Name", session.Name)
	sessionInfo.Set("DeviceJID", session.DeviceJID)
	sessionInfo.Set("Status", session.Status)
	sessionInfo.Set("QRCode", session.QRCode)
	
	// Webhook config
	if session.WebhookConfig != nil {
		sessionInfo.Set("WebhookEnabled", fmt.Sprintf("%t", session.WebhookConfig.Enabled))
		sessionInfo.Set("WebhookURL", session.WebhookConfig.URL)
		if len(session.WebhookConfig.Events) > 0 {
			sessionInfo.Set("Events", strings.Join(session.WebhookConfig.Events, ","))
		}
		sessionInfo.Set("WebhookToken", session.WebhookConfig.Token)
	}
	
	// Proxy config
	if session.ProxyConfig != nil {
		sessionInfo.Set("ProxyEnabled", fmt.Sprintf("%t", session.ProxyConfig.Enabled))
		sessionInfo.Set("ProxyProtocol", session.ProxyConfig.Protocol)
		sessionInfo.Set("ProxyHost", session.ProxyConfig.Host)
		sessionInfo.Set("ProxyPort", fmt.Sprintf("%d", session.ProxyConfig.Port))
		sessionInfo.Set("ProxyUsername", session.ProxyConfig.Username)
		sessionInfo.Set("ProxyPassword", session.ProxyConfig.Password)
	}

	return sessionInfo
}
