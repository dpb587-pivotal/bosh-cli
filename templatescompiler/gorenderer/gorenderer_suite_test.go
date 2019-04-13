package gorenderer_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGorenderer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gorenderer Suite")
}
