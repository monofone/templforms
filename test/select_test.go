package test

import (
	"context"
	"io"
	"strconv"
	"testing"

	"github.com/monofone/templforms"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func Test_RawOption(t *testing.T) {
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
			assert.Equal(t, 1, doc.Find(`option`).Length(), "expected option tag to be rendered")

			valueAttr, valueAttrExists := doc.Find("option").Attr("value")
			assert.True(t, valueAttrExists)
			assert.Equal(t, tt.args.option.Value, valueAttr)

			_, selectedAttrExists := doc.Find("option").Attr("selected")
			assert.Equal(t, tt.args.option.Selected, selectedAttrExists)
		})
	}
}

func Test_RawOptionGroup(t *testing.T) {
	type args struct {
		optionGroup templforms.OptionGroup
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "option, not selected",
			args: args{
				optionGroup: templforms.OptionGroup{
					Label:    "some-label",
					Disabled: false,
				},
			},
		},
		{
			name: "option, selected",
			args: args{
				optionGroup: templforms.OptionGroup{
					Label:    "some-label",
					Disabled: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, w := io.Pipe()
			go func() {
				_ = templforms.RawOptionGroup(tt.args.optionGroup).Render(context.Background(), w)
				_ = w.Close()
			}()
			doc, err := goquery.NewDocumentFromReader(r)
			if err != nil {
				t.Fatalf("failed to read template: %v", err)
			}
			// Expect the component to be present.
			if doc.Find(`optgroup`).Length() == 0 {
				t.Error("expected option tag to be rendered, but it wasn't")
			}

			optgroupNode := doc.Find(`optgroup`)

			assertStringAttribute(t, "label", optgroupNode, tt.args.optionGroup.Label)

			assertBoolAttribute(t, "disabled", optgroupNode, tt.args.optionGroup.Disabled)
		})
	}
}

func TestOptionsGroupIntegration(t *testing.T) {
	r, w := io.Pipe()

	go func() {
		options := []templforms.Option{
			templforms.Option{
				Label: "Option 1",
				Value: "option-1",
			},
			templforms.Option{
				Label: "Option 2",
				Value: "option-2",
			},
			templforms.Option{
				Label:    "Option 3",
				Value:    "option-3",
				Selected: true,
			},
		}

		_ = optionGroupTest(options).Render(context.Background(), w)
		_ = w.Close()
	}()
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		t.Fatalf("failed to read template: %v", err)
	}
	// Expect the component to be present.
	assert.Equal(t, 3, doc.Find(`optgroup > option`).Length())

	optgroupNode := doc.Find(`optgroup`)

	assertStringAttribute(t, "label", optgroupNode, "some-label")

	assertBoolAttribute(t, "disabled", optgroupNode, false)
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

			if tt.args.selectOptions.Size > 1 {
				assertStringAttribute(t, "size", selectNode, strconv.Itoa(tt.args.selectOptions.Size))
			}

			if len(tt.args.selectOptions.ID) > 0 {
				assertStringAttribute(t, "id", selectNode, tt.args.selectOptions.ID)
			} else {
				assertStringAttribute(t, "id", selectNode, tt.args.name)
			}

			if len(tt.args.selectOptions.Class) > 0 {
				assertStringAttribute(t, "class", selectNode, tt.args.selectOptions.Class)
			}

			assertBoolAttribute(t, "multiple", selectNode, tt.args.selectOptions.Multiple)
			assertBoolAttribute(t, "required", selectNode, tt.args.selectOptions.Required)
			assertBoolAttribute(t, "disabled", selectNode, tt.args.selectOptions.Disabled)
		})
	}
}
