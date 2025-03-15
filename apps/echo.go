package apps

func (a *Apps) Echo(args ...string) string {
	result := ""
	for i, _ := range args {
		result += args[i] + " "
	}
	return result
}
