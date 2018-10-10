package http_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mrlyc/trumpet/http"
	"github.com/spf13/afero"
)

func TestRender(t *testing.T) {
	fileName := "test"
	setupFs(afero.NewMemMapFs())
	defer tearDownFs()
	writeFile(fileName, []byte(`ok`))

	tmpl, err := http.NewTemplate(fileName)
	assert.NoError(t, err)

	output, err := tmpl.Render(nil)
	assert.NoError(t, err)
	assert.Equal(t, "ok", output)
}

func TestFHasAttr(t *testing.T) {
	fileName := "test"
	setupFs(afero.NewMemMapFs())
	defer tearDownFs()
	writeFile(fileName, []byte(`{{ if hasattr . "OK" }}ok{{ end }}`))

	var output string
	tmpl, err := http.NewTemplate(fileName)
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

func TestFLookup(t *testing.T) {
	fileName := "test"
	setupFs(afero.NewMemMapFs())
	defer tearDownFs()

	var (
		tmpl   *http.Template
		output string
		err    error
		data   = make(map[string]interface{}, 0)
	)
	err = json.Unmarshal([]byte(`
		{
			"store": {
				"book": [
					{
						"category": "reference",
						"author": "Nigel Rees",
						"title": "Sayings of the Century",
						"price": 8.95
					},
					{
						"category": "fiction",
						"author": "Evelyn Waugh",
						"title": "Sword of Honour",
						"price": 12.99
					},
					{
						"category": "fiction",
						"author": "Herman Melville",
						"title": "Moby Dick",
						"isbn": "0-553-21311-3",
						"price": 8.99
					},
					{
						"category": "fiction",
						"author": "J. R. R. Tolkien",
						"title": "The Lord of the Rings",
						"isbn": "0-395-19395-8",
						"price": 22.99
					}
				],
				"bicycle": {
					"color": "red",
					"price": 19.95
				}
			},
			"expensive": 10
		}
	`), &data)
	assert.NoError(t, err)

	cases := []struct{ input, excepted string }{
		{"$.expensive", "10"},
		{"$.store.book[0].price", "8.95"},
		{"$.store.book[-1].isbn", "0-395-19395-8"},
		{"$.store.book[0,1].price", "[8.95 12.99]"},
		{"$.store.book[0:2].price", "[8.95 12.99 8.99]"},
		{"$.store.book[?(@.isbn)].price", "[8.99 22.99]"},
		{"$.store.book[?(@.price > 10)].title", "[Sword of Honour The Lord of the Rings]"},
		// {"$.store.book[?(@.price < $.expensive)].price", "[8.95 8.99]"},  // bug
		{"$.store.book[:].price", "[8.95 12.99 8.99 22.99]"},
		{"$.store.book[?(@.author =~ /(?i).*REES/)].author[0]", "Nigel Rees"},
	}

	for _, c := range cases {
		writeFile(fileName, []byte(fmt.Sprintf(`{{ lookup . "%v" "!!!" }}`, c.input)))
		tmpl, err = http.NewTemplate(fileName)
		assert.NoError(t, err)

		output, err = tmpl.Render(data)
		if assert.NoError(t, err) {
			assert.Equal(t, c.excepted, output)
		}
	}

}
