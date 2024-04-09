package templforms_test

import (
	"context"
	"io"
	"strconv"
	"testing"

	"github.com/monofone/templforms"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func Test_numberAttribute(t *testing.T) {
	type args struct {
		option templforms.Option
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "option, not selected",
			args: args{
				option: templforms.Option{
					Value:    "some-value",
					Label:    "some-label",
					Selected: false,
				},
			},
		},
		{
			name: "option, selected",
			args: args{
				option: templforms.Option{
					Value:    "some-value",
					Label:    "some-label",
					Selected: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, w := io.Pipe()
			go func() {
				_ = templforms.RawOption(tt.args.option).Render(context.Background(), w)
				_ = w.Close()
			}()
			doc, err := goquery.NewDocumentFromReader(r)
			if err != nil {
				t.Fatalf("failed to read template: %v", err)
			}
			// Expect the component to be present.
			if doc.Find(`option`).Length() == 0 {
				t.Error("expected option tag to be rendered, but it wasn't")
			}

			valueAttr, valueAttrExists := doc.Find("option").Attr("value")
			assert.True(t, valueAttrExists)
			assert.Equal(t, tt.args.option.Value, valueAttr)

			_, selectedAttrExists := doc.Find("option").Attr("selected")
			assert.Equal(t, tt.args.option.Selected, selectedAttrExists)
		})
	}
}

func TestRawSelect(t *testing.T) {
	type args struct {
		name          string
		selectOptions *templforms.SelectOptions
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "select no options",
			args: args{
				name:          "my-select",
				selectOptions: &templforms.SelectOptions{},
			},
		},
		{
			name: "select, size option",
			args: args{
				name: "my-select",
				selectOptions: &templforms.SelectOptions{
					Size: 5,
				},
			},
		},
		{
			name: "select, bool options",
			args: args{
				name: "my-select",
				selectOptions: &templforms.SelectOptions{
					GenericOptions: templforms.GenericOptions{
						ID:       "some-extra-id",
						Required: true,
						Disabled: true,
					},
					Multiple: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, w := io.Pipe()
			go func() {
				_ = templforms.RawSelect(tt.args.name, tt.args.selectOptions).Render(context.Background(), w)
				_ = w.Close()
			}()
			doc, err := goquery.NewDocumentFromReader(r)
			if err != nil {
				t.Fatalf("failed to read template: %v", err)
			}
			// Expect the component to be present.
			if doc.Find(`select`).Length() == 0 {
				t.Error("expected select tag to be rendered, but it wasn't")
			}

			nameAttr, nameAttrExists := doc.Find("select").Attr("name")
			assert.True(t, nameAttrExists)
			assert.Equal(t, tt.args.name, nameAttr)

			selectNode := doc.Find("select")

			assertStringAttribute(t, "size", selectNode, tt.args.selectOptions.Size > 1, strconv.Itoa(tt.args.selectOptions.Size))

			if len(tt.args.selectOptions.ID) > 0 {
				assertStringAttribute(t, "id", selectNode, true, tt.args.selectOptions.ID)
			} else {
				assertStringAttribute(t, "id", selectNode, true, tt.args.name)
			}

			if len(tt.args.selectOptions.Class) > 0 {
				assertStringAttribute(t, "class", selectNode, true, tt.args.selectOptions.Class)
			}

			assertBoolAttribute(t, "multiple", selectNode, tt.args.selectOptions.Multiple)
			assertBoolAttribute(t, "required", selectNode, tt.args.selectOptions.Required)
			assertBoolAttribute(t, "disabled", selectNode, tt.args.selectOptions.Disabled)
		})
	}
}

func assertBoolAttribute(t *testing.T, attributeName string, node *goquery.Selection, expectExists bool) {
	_, attrExists := node.Attr(attributeName)
	assert.Equal(t, expectExists, attrExists)

}

func assertStringAttribute(t *testing.T, attributeName string, node *goquery.Selection, expectExists bool, expectedContent string) {
	sizeAttr, sizteAttrExists := node.Attr(attributeName)
	if expectExists {
		assert.True(t, sizteAttrExists)
		assert.Equal(t, expectedContent, sizeAttr)
	} else {
		assert.False(t, sizteAttrExists)
	}
}
