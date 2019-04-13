package gorenderer

import (
	"bytes"
	"text/template"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

// TODO not renderer-specific
type TODORenderer interface {
	Render(srcPath, dstPath string, context TODOContext) error
}

type renderer struct {
	fs     boshsys.FileSystem
	logger boshlog.Logger
	logTag string
}

func NewRenderer(
	fs boshsys.FileSystem,
	logger boshlog.Logger,
) TODORenderer {
	return &renderer{
		fs:     fs,
		logger: logger,
		logTag: "goRenderer",
	}
}

func (r renderer) Render(srcPath, dstPath string, context TODOContext) error {
	r.logger.Debug(r.logTag, "Rendering template %s", dstPath)

	srcBytes, err := r.fs.ReadFile(srcPath)
	if err != nil {
		return bosherr.WrapError(err, "Reading template")
	}

	tmpl := template.New(srcPath)
	_, err = tmpl.Parse(string(srcBytes))
	if err != nil {
		return bosherr.WrapError(err, "Parsing template")
	}

	dstBuffer := bytes.NewBuffer(nil)

	err = tmpl.Execute(dstBuffer, context)
	if err != nil {
		return bosherr.WrapError(err, "Rendering template")
	}

	err = r.fs.WriteFile(dstPath, dstBuffer.Bytes())
	if err != nil {
		return bosherr.WrapError(err, "Writing rendered template")
	}

	return nil
}
