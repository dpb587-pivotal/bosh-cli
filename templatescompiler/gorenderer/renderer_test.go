package gorenderer_test

import (
	. "github.com/cloudfoundry/bosh-cli/templatescompiler/gorenderer"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	fakesys "github.com/cloudfoundry/bosh-utils/system/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("renderer", func() {
	var (
		fs       *fakesys.FakeFileSystem
		renderer TODORenderer
		context  TODOContext
	)

	BeforeEach(func() {
		logger := boshlog.NewLogger(boshlog.LevelNone)
		fs = fakesys.NewFakeFileSystem()
		context = TODOContext{
			ID: "fake-uuid",
			Properties: map[string]interface{}{
				"key": "fake-value",
			},
		}

		renderer = NewRenderer(fs, logger)
		fs.TempDirDir = "fake-temp-dir"
	})

	It("works happily", func() {
		fs.WriteFile("fake-src-path", []byte(`{{ .ID }} says {{ .Property "key" }}`))

		err := renderer.Render("fake-src-path", "fake-dst-path", context)
		Expect(err).ToNot(HaveOccurred())

		res, err := fs.ReadFile("fake-dst-path")
		Expect(err).ToNot(HaveOccurred())
		Expect(res).To(Equal([]byte("fake-uuid says fake-value")))
	})
})
