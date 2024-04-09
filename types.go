package templforms

type GenericOptions struct {
	ID       string
	Disabled bool
	Required bool
	Class    string
}

type Option struct {
	Value    string
	Label    string
	Selected bool
	Children []Option
}
