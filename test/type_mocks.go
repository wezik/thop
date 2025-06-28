package test

import (
	"os/exec"
	"thop/internal/types/project"
	"thop/internal/types/template"

	"github.com/stretchr/testify/mock"
)

type MockExecutor struct {
	mock.Mock
}

func (m *MockExecutor) Execute(cmd *exec.Cmd) (string, int, error) {
	args := m.Called(cmd)
	return args.String(0), args.Int(1), args.Error(2)
}

func (m *MockExecutor) ExecuteInteractive(cmd *exec.Cmd) (int, error) {
	args := m.Called(cmd)
	return args.Int(0), args.Error(1)
}

type MockMultiplexer struct {
	mock.Mock
}

func (m *MockMultiplexer) AttachProject(p project.Project) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *MockMultiplexer) ListActiveSessions() ([]project.Project, error) {
	args := m.Called()
	return args.Get(0).([]project.Project), args.Error(1)
}

func (m *MockMultiplexer) KillSession(p project.Project) error {
	args := m.Called(p)
	return args.Error(0)
}

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) List() ([]project.Project, error) {
	args := m.Called()
	return args.Get(0).([]project.Project), args.Error(1)
}

func (m *MockStorage) ListWithInvalid() ([]project.Project, error) {
	args := m.Called()
	return args.Get(0).([]project.Project), args.Error(1)
}

func (m *MockStorage) Find(name project.Name) (project.Project, error) {
	args := m.Called(name)
	return args.Get(0).(project.Project), args.Error(1)
}

func (m *MockStorage) Save(p *project.Project) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *MockStorage) Delete(uuid project.UUID) error {
	args := m.Called(uuid)
	return args.Error(0)
}

func (m *MockStorage) PrepareTemplateFile(uuid project.UUID) (string, error) {
	args := m.Called(uuid)
	return args.String(0), args.Error(1)
}

type MockService struct {
	mock.Mock
}

func (m *MockService) CreateProject(root template.Root, name project.Name) error {
	args := m.Called(root, name)
	return args.Error(0)
}

func (m *MockService) OpenProject(name project.Name) error {
	args := m.Called(name)
	return args.Error(0)
}

func (m *MockService) DeleteProject(name project.Name) error {
	args := m.Called(name)
	return args.Error(0)
}

func (m *MockService) EditProject(name project.Name) error {
	args := m.Called(name)
	return args.Error(0)
}

func (m *MockService) KillSession(name project.Name) error {
	args := m.Called(name)
	return args.Error(0)
}

type MockProjectSelector struct {
	mock.Mock
}

func (s *MockProjectSelector) SelectFrom(items []project.Project, prompt string) (*project.Project, error) {
	args := s.Called(items, prompt)
	return args.Get(0).(*project.Project), args.Error(1)
}
