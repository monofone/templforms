package templforms

import "github.com/a-h/templ"

type GenericOptions struct {
	ID       string
	Disabled bool
	Required bool
	Class    string
	Attr     templ.Attributes
}

type Option struct {
	Value    string
	Label    string
	Selected bool
}

type OptionGroup struct {
	Label    string
	Disabled bool
}
