package domain

import (
	"testing"
	"thop/internal/adapters/log"
	"thop/internal/domain/template"
	"thop/internal/domain/thop"
	"thop/test/gen/mock"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_Create(t *testing.T) {
	const templateVersion = template.V1

	// given
	defaultDirectory := "default/directory"
	for _, testCase := range []struct {
		name         string
		expectedName string
		path         string
		expectedPath string
	}{
		{name: "", path: "", expectedName: defaultDirectory, expectedPath: defaultDirectory},
		{name: "test name", path: "", expectedName: "test name", expectedPath: defaultDirectory},
		{name: "", path: "test/directory", expectedName: "test/directory", expectedPath: "test/directory"},
		{name: "test name", path: "test/directory", expectedName: "test name", expectedPath: "test/directory"},
	} {
		testName := "creates new default template with options \"" + testCase.name + "\" and \"" + testCase.path + "\""
		t.Run(testName, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLogger := log.NewNoopLogger()
			mockSelector := mock.NewMockSelector(ctrl)
			mockFileStorage := mock.NewMockFileStorage(ctrl)
			mockMultiplexer := mock.NewMockMultiplexer(ctrl)

			app := thop.New(
				mockLogger,
				mockSelector,
				mockFileStorage,
				mockMultiplexer,
				defaultDirectory,
			)

			createTemplate := thop.CreateTemplate{
				Name: testCase.name,
				Path: testCase.path,
			}
			savedTemplate := template.NewTemplate(
				"foo/bar/default_directory.yaml",
				templateVersion,
				"default",
				"default",
				"",
				[]string{},
				[]*template.Window{},
			)

			var capturedTemplate *template.Template
			mockFileStorage.
				EXPECT().
				Save(gomock.Any()).
				DoAndReturn(func(templ *template.Template) (*template.Template, error) {
					capturedTemplate = templ
					return savedTemplate, nil
				})

			// when
			templ, err := app.Create(createTemplate)

			// then
			assert.Nil(t, err)
			assert.Equal(t, template.FilePath(""), capturedTemplate.FilePath())
			assert.Equal(t, testCase.expectedName, capturedTemplate.Name())
			assert.Equal(t, testCase.expectedPath, capturedTemplate.Path())
			assert.Equal(t, templateVersion, capturedTemplate.Version())
			assert.Equal(t, savedTemplate, templ)
		})
	}
}
