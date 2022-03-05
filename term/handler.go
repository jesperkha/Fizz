package term

// Implementation of the ArgList interface

type ArgHandler struct {
	flags   []string
	options []string
	args    []string
	subcmd  string
}

func (a *ArgHandler) HasFlag(flag string) bool {
	for _, f := range a.flags {
		if f == flag {
			return true
		}
	}
	return false
}

func (a *ArgHandler) HasOption(option string) bool {
	for _, o := range a.options {
		if o == option {
			return true
		}
	}
	return false
}

func (a *ArgHandler) SubCommand() string {
	return a.subcmd
}

func (a *ArgHandler) Args() []string {
	return a.args
}
