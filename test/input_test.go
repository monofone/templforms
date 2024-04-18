package test

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/a-h/templ"
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
			name: "raw text input",
			args: args{
				option:    &templforms.InputOptions{},
				name:      "test-input",
				inputType: "text",
				value:     "some-test",
			},
		},
		{
			name: "raw number input",
			args: args{
				option:    &templforms.InputOptions{},
				name:      "test-input",
				inputType: "number",
				value:     "1",
			},
		},
		{
			name: "additional attr are rendered",
			args: args{
				option: &templforms.InputOptions{
					GenericOptions: templforms.GenericOptions{
						Attr: templ.Attributes{
							"data-test-id": "test",
							"class":        "border-2 border-dashed",
						},
					},
				},
				name:      "test-input",
				inputType: "number",
				value:     "1",
			},
		},
		{
			name: "additional attr overwrite existing ones",
			args: args{
				option: &templforms.InputOptions{
					GenericOptions: templforms.GenericOptions{
						Attr: templ.Attributes{
							"data-test-id":     "test",
							"aria-describedby": "test-aria-id",
						},
					},
					FieldError: errors.New("some error"),
				},
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
			nodeContent, err := inputNode.Html()
			t.Logf("html: %s %+v", nodeContent, err)
			assert.Equal(t, 1, inputNode.Length())
			// Expect the component to be present.
			assert.True(t, inputNode.Is("input"), "expected input tag to be rendered, but it wasn't")

			assertStringAttribute(t, "name", inputNode, tt.args.name)
			if len(tt.args.option.ID) > 0 {
				assertStringAttribute(t, "id", inputNode, tt.args.option.ID)
			} else {
				assertStringAttribute(t, "id", inputNode, tt.args.name)
			}

			// get value attribute from input field with goquery
			v, valueAttrExists := inputNode.Attr("value")
			assert.True(t, valueAttrExists, "value attribute is not present")
			assert.Equal(t, tt.args.value.(string), v)

			for k, v := range tt.args.option.Attr {
				if k != "aria-describedby" {
					assertStringAttribute(t, k, inputNode, v.(string))
				}
			}

			if tt.args.option.FieldError != nil {
				assertStringAttribute(t, "aria-describedby", inputNode, tt.args.name+"-error")
			}
		})
	}
}
