package kubelet_utils

func IsStrInList(str string, list []string) bool {
	for _, str_in := range list {
		if str_in == str {
			return true
		}
	}
	return false
}
