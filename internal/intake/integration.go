package intake

import "sort"

func Normalize(ids []string) []string {
	out := append([]string(nil), ids...)
	sort.Strings(out)
	return out
}
func Contains(ids []string, id string) bool {
	for _, v := range ids {
		if v == id {
			return true
		}
	}
	return false
}
