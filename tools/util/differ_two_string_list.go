package util

import (
	"fmt"
	"sort"
)

func ConvertIntListToStringList(intList []int) []string {
	var stringList []string
	for _, string_ := range intList {
		stringList = append(stringList, fmt.Sprintf("%d", string_))
	}
	return stringList
}

func DifferTwoStringList(olds []string, news []string) (added []string, deleted []string) {
	sort.Strings(news)
	for _, old := range olds {
		if sort.SearchStrings(news, old) == len(news) {
			deleted = append(deleted, old)
		}
	}
	sort.Strings(olds)
	for _, new_ := range news {
		if sort.SearchStrings(olds, new_) == len(olds) {
			added = append(added, new_)
		}
	}
	return
}
