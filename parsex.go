package parsex

type Program struct {
	// Pointer to a [var ... struct{}] containing all program options.
	//
	// All input flags will be stored in the struct,
	// --help menu will be automatically generated.
	//
	// Use the following field tags:
	//
	// `alt:"<single letter alternative use>" desc:"<description of the option>`
	Data any
	// The name of the executable / command
	Name string
	// Will be displayed in the --help menu
	Desc string
}
