package cli

func ParseArgs(args []string) (string, string, string) {
	if len(args) != 4 {
		return "", "", ""
	}

	return args[1], args[2], args[3]
}
