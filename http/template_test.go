package http_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mrlyc/trumpet/http"
)

func TestRender(t *testing.T) {
	tmpl, err := http.NewTemplateFromString("test", `ok`)
	assert.NoError(t, err)

	output, err := tmpl.Render(nil)
	assert.NoError(t, err)
	assert.Equal(t, "ok", output)
}
