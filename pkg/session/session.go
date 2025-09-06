package session

type Session struct {
	name      string
	sessionId string
}

func (s *Session) Name() string {
	return s.name
}

func (s *Session) SessionId() string {
	return s.sessionId
}

func NewSession(name string, tmuxId string) *Session {
	return &Session{
		name:      name,
		sessionId: tmuxId,
	}
}
