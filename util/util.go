package util

import (
	"cmp"
	"sort"
)

func SortMapValues[V cmp.Ordered, K cmp.Ordered](dict map[K]V) map[K]V {
	keys := make([]K, 0, len(dict))
	for key := range dict {
		keys = append(keys, key)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		ival, iok := dict[keys[i]]
		jval, jok := dict[keys[j]]
		return iok && jok && ival < jval
	})
	sorted := make(map[K]V)
	for _, key := range keys[:100] {
		sorted[key] = dict[key]
	}
	return sorted
}
