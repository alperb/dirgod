package arguments

var AVAILABLE_ARGS = []string{
	"-d",
	"--debug",
	"-v",
	"--verbose",
}

type Arguments struct {
	RawArgs      []string
	DEBUG_MODE   bool
	VERBOSE_MODE bool
	Filename     string
}

func NewArguments(args []string) *Arguments {
	a := &Arguments{args, false, false, args[len(args)-1]}
	a.Init()
	return a
}

func (a *Arguments) Init() {
	a.setDebugMode()
	a.validate()
}

func (a *Arguments) validate() {
	for idx, arg := range a.RawArgs {
		if len(a.RawArgs)-1 == idx {
			// we are at the last argument, aka the filename
			break
		}

		if !contains(AVAILABLE_ARGS, arg) {
			panic("Invalid argument: " + arg + " at index: " + string(idx))
		}
	}
}

func (a *Arguments) setDebugMode() {
	for _, arg := range a.RawArgs {
		if arg == "-d" || arg == "--debug" {
			a.DEBUG_MODE = true
		} else if arg == "-v" || arg == "--verbose" {
			a.VERBOSE_MODE = true
		}
	}
}

// helpers

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
