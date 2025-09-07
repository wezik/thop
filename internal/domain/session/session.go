package session

type Session struct {
	name string
}

func (s *Session) Name() string {
	return s.name
}

func NewSession(name string) *Session {
	return &Session{
		name: name,
	}
}
