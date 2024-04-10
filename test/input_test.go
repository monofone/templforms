package test

import (
	"context"
	"io"
	"testing"

	"github.com/monofone/templforms"
	"github.com/stretchr/testify/assert"

	"github.com/PuerkitoBio/goquery"
)

func Test_RawInput(t *testing.T) {
	type args struct {
		option    *templforms.InputOptions
		name      string
		inputType string
		value     interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "raw input",
			args: args{
				option:    &templforms.InputOptions{},
				name:      "test-input",
				inputType: "test",
				value:     "some-test",
			},
		},
		{
			name: "raw input",
			args: args{
				option:    &templforms.InputOptions{},
				name:      "test-input",
				inputType: "number",
				value:     "1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, w := io.Pipe()
			go func() {
				_ = templforms.RawInput(tt.args.name, tt.args.inputType, tt.args.value, tt.args.option).Render(context.Background(), w)
				_ = w.Close()
			}()
			doc, err := goquery.NewDocumentFromReader(r)
			if err != nil {
				t.Fatalf("failed to read template: %v", err)
			}

			inputNode := doc.Find(`input`)

			// Expect the component to be present.
			assert.Equal(t, 1, inputNode.Length(), "expected input tag to be rendered, but it wasn't")

			assertStringAttribute(t, "name", inputNode, tt.args.name)
			v, _ := inputNode.Attr("value")
			assert.Equal(t, tt.args.value.(string), v)
		})
	}
}
