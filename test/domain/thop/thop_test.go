package thop

import (
	"errors"
	"slices"
	"strings"
	"testing"
	"thop/internal/adapters/log"
	"thop/internal/domain/selector"
	"thop/internal/domain/session"
	"thop/internal/domain/template"
	"thop/internal/domain/thop"
	"thop/test/gen/mock"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func someTemplate() *template.Template {
	return template.NewTemplate(
		"foo/bar/default_directory.yaml",
		template.V1,
		"default",
		"default",
		"",
		[]string{},
		[]*template.Window{},
	)
}

func Test_Create(t *testing.T) {
	t.Run("creates new template", func(t *testing.T) {
		// given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockLogger := log.NewNoopLogger()
		mockSelector := mock.NewMockSelector(ctrl)
		mockTemplateService := mock.NewMockTemplateService(ctrl)
		mockMultiplexer := mock.NewMockMultiplexer(ctrl)

		app := thop.New(
			mockLogger,
			mockSelector,
			mockTemplateService,
			mockMultiplexer,
		)

		createTemplate := template.CreateTemplate{
			Name: "test name",
			Path: "test/directory",
		}

		expectedTemplate := someTemplate()

		mockTemplateService.EXPECT().Create(createTemplate).Return(expectedTemplate, nil)

		// when
		templ, err := app.Create(createTemplate)

		// then
		assert.Nil(t, err)
		assert.Equal(t, expectedTemplate, templ)
	})

	t.Run("propagates template creation error", func(t *testing.T) {
		// given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockLogger := log.NewNoopLogger()
		mockSelector := mock.NewMockSelector(ctrl)
		mockTemplateService := mock.NewMockTemplateService(ctrl)
		mockMultiplexer := mock.NewMockMultiplexer(ctrl)

		app := thop.New(
			mockLogger,
			mockSelector,
			mockTemplateService,
			mockMultiplexer,
		)

		createTemplate := template.CreateTemplate{
			Name: "test name",
			Path: "test/directory",
		}

		expectedErr := errors.New("not found")

		mockTemplateService.EXPECT().Create(createTemplate).Return(nil, expectedErr)

		// when
		templ, err := app.Create(createTemplate)

		// then
		assert.Nil(t, templ)
		assert.Equal(t, expectedErr, err)
	})
}

func Test_OpenSelect(t *testing.T) {
	t.Run("opens existing template", func(t *testing.T) {
		// given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockLogger := log.NewNoopLogger()
		mockSelector := mock.NewMockSelector(ctrl)
		mockTemplateService := mock.NewMockTemplateService(ctrl)
		mockMultiplexer := mock.NewMockMultiplexer(ctrl)

		chosenTemplate := someTemplate()
		templates := []*template.Template{chosenTemplate}

		mockTemplateService.EXPECT().
			List().
			Return(templates, nil)

		mockSelector.EXPECT().
			SelectFrom(gomock.Any(), selector.OperationOpen).
			DoAndReturn(func(entries []*selector.Entry, op selector.Operation) (*selector.Entry, error) {
				return entries[0], nil
			})

		mockMultiplexer.EXPECT().
			ListSessions().
			Return(nil, nil)

		mockMultiplexer.EXPECT().
			AttachTemplate(chosenTemplate).
			Return(nil)

		app := thop.New(
			mockLogger,
			mockSelector,
			mockTemplateService,
			mockMultiplexer,
		)

		// when
		err := app.OpenSelect()

		// then
		assert.Nil(t, err)
	})

	t.Run("opens active session", func(t *testing.T) {
		// given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockLogger := log.NewNoopLogger()
		mockSelector := mock.NewMockSelector(ctrl)
		mockTemplateService := mock.NewMockTemplateService(ctrl)
		mockMultiplexer := mock.NewMockMultiplexer(ctrl)

		chosenSession := session.NewSession("foo")

		mockTemplateService.EXPECT().
			List().
			Return(nil, nil)

		mockSelector.EXPECT().
			SelectFrom(gomock.Any(), selector.OperationOpen).
			DoAndReturn(func(entries []*selector.Entry, op selector.Operation) (*selector.Entry, error) {
				return entries[0], nil
			})

		mockMultiplexer.EXPECT().
			ListSessions().
			Return([]*session.Session{chosenSession}, nil)

		mockMultiplexer.EXPECT().
			AttachSession(chosenSession).
			Return(nil)

		app := thop.New(
			mockLogger,
			mockSelector,
			mockTemplateService,
			mockMultiplexer,
		)

		// when
		err := app.OpenSelect()

		// then
		assert.Nil(t, err)
	})

	t.Run("merges active sessions and templates", func(t *testing.T) {
		// given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockLogger := log.NewNoopLogger()
		mockSelector := mock.NewMockSelector(ctrl)
		mockTemplateService := mock.NewMockTemplateService(ctrl)
		mockMultiplexer := mock.NewMockMultiplexer(ctrl)

		templates := []*template.Template{
			template.NewTemplate(
				"foo/bar/template1.yaml",
				template.V1,
				"existing_session_name",
				"session_name",
				"",
				[]string{},
				[]*template.Window{},
			),
			template.NewTemplate(
				"foo/bar/template2.yaml",
				template.V1,
				"template_name",
				"",
				"",
				[]string{},
				[]*template.Window{},
			),
			template.NewTemplate(
				"foo/bar/template3.yaml",
				template.V1,
				"not_active_name",
				"not_active_name",
				"",
				[]string{},
				[]*template.Window{},
			),
		}

		sessions := []*session.Session{
			session.NewSession("session_name"),
			session.NewSession("template_name"),
			session.NewSession("non_template_session"),
		}

		entrySort := func(a, b *selector.Entry) int {
			return strings.Compare(a.Name(), b.Name())
		}

		expectedEntries := []*selector.Entry{
			selector.NewEntry(templates[0].Name(), sessions[0].Name(), selector.TagActiveTemplate),
			selector.NewEntry(templates[1].Name(), sessions[1].Name(), selector.TagActiveTemplate),
			selector.NewEntry(sessions[2].Name(), sessions[2].Name(), selector.TagActiveSession),
			selector.NewEntry(templates[2].Name(), string(templates[2].FilePath()), selector.TagTemplate),
		}
		slices.SortFunc(expectedEntries, entrySort)

		var capturedEntries []*selector.Entry

		mockTemplateService.EXPECT().
			List().
			Return(templates, nil)

		mockSelector.EXPECT().
			SelectFrom(gomock.Any(), selector.OperationOpen).
			DoAndReturn(func(entries []*selector.Entry, op selector.Operation) (*selector.Entry, error) {
				capturedEntries = entries
				slices.SortFunc(capturedEntries, entrySort) // Sort for easier comparison
				return entries[0], nil
			})

		mockMultiplexer.EXPECT().
			ListSessions().
			Return(sessions, nil)

		mockMultiplexer.EXPECT().
			AttachTemplate(gomock.Any()).
			MinTimes(0).
			Return(nil)

		mockMultiplexer.EXPECT().
			AttachSession(gomock.Any()).
			MinTimes(0).
			Return(nil)

		app := thop.New(
			mockLogger,
			mockSelector,
			mockTemplateService,
			mockMultiplexer,
		)

		// when
		err := app.OpenSelect()

		// then
		assert.Nil(t, err)
		assert.Equal(t, expectedEntries, capturedEntries)
	})
}
