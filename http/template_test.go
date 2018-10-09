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

func TestFHasAttr(t *testing.T) {
	var output string
	tmpl, err := http.NewTemplateFromString("test", `{{ if hasAttr . "OK" }}ok{{ end }}`)
	assert.NoError(t, err)

	dict := map[string]interface{}{
		"OK": 1,
	}
	assert.True(t, tmpl.FHasAttr(dict, "OK"))
	assert.False(t, tmpl.FHasAttr(dict, "_"))

	output, err = tmpl.Render(dict)
	assert.NoError(t, err)
	assert.Equal(t, "ok", output)

	object := struct{ OK int }{OK: 1}
	assert.True(t, tmpl.FHasAttr(object, "OK"))
	assert.False(t, tmpl.FHasAttr(object, "_"))

	output, err = tmpl.Render(object)
	assert.NoError(t, err)
	assert.Equal(t, "ok", output)
}
