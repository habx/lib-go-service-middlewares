package template_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	tpl "github.com/habx/lib-go-template"
)

func TestDoesSomething(t *testing.T) {
	a := assert.New(t)

	a.NoError(tpl.DoesSomething())

	t.Run("subtest", func(t *testing.T) {
		a := assert.New(t)
		a.NoError(tpl.DoesSomething())
	})
}
