package template

import (
	"errors"
	"testing"
	"thop/internal/domain/template"
	"thop/test/gen/mock"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_TemplateCreation(t *testing.T) {
	t.Run("creates template", func(t *testing.T) {
		// given
		defaultDirectory := "default/directory"
		savedPath := template.FilePath("some/saved/path")
		expectedVersion := template.V1
		for _, testCase := range []struct {
			name         string
			expectedName string
			path         string
			expectedPath string
		}{
			// in short name should default to path, path should default to pre-set directory
			{name: "", path: "", expectedName: defaultDirectory, expectedPath: defaultDirectory},
			{name: "test name", path: "", expectedName: "test name", expectedPath: defaultDirectory},
			{name: "", path: "test/directory", expectedName: "test/directory", expectedPath: "test/directory"},
			{name: "test name", path: "test/directory", expectedName: "test name", expectedPath: "test/directory"},
		} {

			t.Run("with name \""+testCase.name+"\" and path \""+testCase.path+"\"", func(t *testing.T) {
				// and given
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				mockFileStorage := mock.NewMockFileStorage(ctrl)

				mockFileStorage.EXPECT().
					Save(gomock.Any()).
					DoAndReturn(func(templ *template.Template) (*template.Template, error) {
						saved := template.NewTemplate(
							savedPath,
							templ.Version(),
							templ.Name(),
							templ.SessionName(),
							templ.Path(),
							templ.Commands(),
							templ.Windows(),
						)
						return saved, nil
					})

				service := template.NewTemplateService(
					mockFileStorage,
					defaultDirectory,
				)

				// when
				templ, err := service.Create(template.CreateTemplate{
					Name: testCase.name,
					Path: testCase.path,
				})

				// then
				assert.Nil(t, err)
				assert.Equal(t, savedPath, templ.FilePath())
				assert.Equal(t, testCase.expectedName, templ.Name())
				assert.Equal(t, testCase.expectedPath, templ.Path())
				assert.Equal(t, expectedVersion, templ.Version())
			})

		}
	})
}

func Test_TemplateFinding(t *testing.T) {
	t.Run("finds template", func(t *testing.T) {
		// given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFileStorage := mock.NewMockFileStorage(ctrl)
		expectedTemplate := template.NewTemplate(
			"foo/bar/1",
			template.V1,
			"1",
			"",
			"",
			[]string{},
			[]*template.Window{},
		)

		mockFileStorage.EXPECT().
			Find("foo/bar/1").
			Return(expectedTemplate, nil)

		service := template.NewTemplateService(
			mockFileStorage,
			"",
		)

		// when
		templ, err := service.Find("foo/bar/1")

		// then
		assert.Nil(t, err)
		assert.Equal(t, expectedTemplate, templ)
	})

	t.Run("propagates file finding error", func(t *testing.T) {
		// given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFileStorage := mock.NewMockFileStorage(ctrl)
		expectedErr := errors.New("not found")

		mockFileStorage.EXPECT().
			Find("foo/bar/1").
			Return(nil, expectedErr)

		service := template.NewTemplateService(
			mockFileStorage,
			"",
		)

		// when
		templ, err := service.Find("foo/bar/1")

		// then
		assert.Nil(t, templ)
		assert.Equal(t, expectedErr, err)
	})
}

func Test_TemplateListing(t *testing.T) {
	t.Run("lists templates", func(t *testing.T) {
		// given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFileStorage := mock.NewMockFileStorage(ctrl)

		expectedTemplateFiles := []*template.File{
			template.NewFile("1", "foo/bar/1"),
		}

		expectedTemplates := []*template.Template{
			template.NewTemplate(
				"foo/bar/1",
				template.V1,
				"1",
				"",
				"",
				[]string{},
				[]*template.Window{},
			),
		}

		mockFileStorage.EXPECT().List().Return(expectedTemplateFiles, nil)
		mockFileStorage.EXPECT().LoadTemplate(expectedTemplateFiles[0].Path()).Return(expectedTemplates[0], nil)

		service := template.NewTemplateService(
			mockFileStorage,
			"",
		)

		// when
		templates, err := service.List()

		// then
		assert.Nil(t, err)
		assert.Equal(t, expectedTemplates, templates)
	})

	t.Run("ignores template loading errors", func(t *testing.T) {
		// given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFileStorage := mock.NewMockFileStorage(ctrl)

		expectedTemplateFiles := []*template.File{
			template.NewFile("1", "foo/bar/1"),
			template.NewFile("2", "foo/bar/2"),
		}

		expectedTemplates := []*template.Template{
			template.NewTemplate(
				"foo/bar/1",
				template.V1,
				"1",
				"",
				"",
				[]string{},
				[]*template.Window{},
			),
		}

		mockFileStorage.EXPECT().
			List().
			Return(expectedTemplateFiles, nil)

		mockFileStorage.EXPECT().
			LoadTemplate(expectedTemplateFiles[0].Path()).
			Return(expectedTemplates[0], nil)

		mockFileStorage.EXPECT().
			LoadTemplate(expectedTemplateFiles[1].Path()).
			Return(nil, errors.New("parse error"))

		service := template.NewTemplateService(
			mockFileStorage,
			"",
		)

		// when
		templates, err := service.List()

		// then
		assert.Nil(t, err)
		assert.Equal(t, expectedTemplates, templates)
	})

	t.Run("propagates file listing error", func(t *testing.T) {
		// given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFileStorage := mock.NewMockFileStorage(ctrl)

		expectedErr := errors.New("read error")

		mockFileStorage.EXPECT().
			List().
			Return(nil, expectedErr)

		service := template.NewTemplateService(
			mockFileStorage,
			"",
		)

		// when
		templates, err := service.List()

		// then
		assert.Nil(t, templates)
		assert.Equal(t, err, expectedErr)
	})
}

func Test_FileListing(t *testing.T) {
	t.Run("lists files", func(t *testing.T) {
		// given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFileStorage := mock.NewMockFileStorage(ctrl)

		expectedTemplateFiles := []*template.File{
			template.NewFile("1", "foo/bar/1"),
			template.NewFile("2", "foo/bar/2"),
		}

		mockFileStorage.EXPECT().List().Return(expectedTemplateFiles, nil)

		service := template.NewTemplateService(
			mockFileStorage,
			"",
		)

		// when
		files, err := service.ListFiles()

		// then
		assert.Nil(t, err)
		assert.Equal(t, expectedTemplateFiles, files)
	})

	t.Run("propagates file listing error", func(t *testing.T) {
		// given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFileStorage := mock.NewMockFileStorage(ctrl)

		expectedErr := errors.New("read error")

		mockFileStorage.EXPECT().
			List().
			Return(nil, expectedErr)

		service := template.NewTemplateService(
			mockFileStorage,
			"",
		)

		// when
		files, err := service.ListFiles()

		// then
		assert.Nil(t, files)
		assert.Equal(t, err, expectedErr)
	})
}

func Test_TemplateLoading(t *testing.T) {
	t.Run("loads template", func(t *testing.T) {
		// given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		filePath := template.FilePath("foo/bar/1")

		mockFileStorage := mock.NewMockFileStorage(ctrl)
		expectedTemplate := template.NewTemplate(
			filePath,
			template.V1,
			"1",
			"",
			"",
			[]string{},
			[]*template.Window{},
		)

		mockFileStorage.EXPECT().
			LoadTemplate(filePath).
			Return(expectedTemplate, nil)

		service := template.NewTemplateService(
			mockFileStorage,
			"",
		)

		// when
		templ, err := service.Load(filePath)

		// then
		assert.Nil(t, err)
		assert.Equal(t, expectedTemplate, templ)
	})

	t.Run("propagates file loading error", func(t *testing.T) {
		// given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		filePath := template.FilePath("foo/bar/1")

		mockFileStorage := mock.NewMockFileStorage(ctrl)
		expectedErr := errors.New("not found")

		mockFileStorage.EXPECT().
			LoadTemplate(filePath).
			Return(nil, expectedErr)

		service := template.NewTemplateService(
			mockFileStorage,
			"",
		)

		// when
		templ, err := service.Load(filePath)

		// then
		assert.Nil(t, templ)
		assert.Equal(t, expectedErr, err)
	})
}
