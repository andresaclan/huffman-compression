package cli

func ParseArgs(args []string) (string, string) {
	if len(args) != 3 {
		return "", ""
	}

	return args[1], args[2]
}
