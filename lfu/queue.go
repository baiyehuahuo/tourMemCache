package lfu

import "container/heap"

type queue []*entry

var _ heap.Interface = (*queue)(nil)

func (q *queue) Len() int {
	return len(*q)
}

func (q *queue) Less(i, j int) bool {
	return (*q)[i].freq < (*q)[j].freq
}

func (q *queue) Swap(i, j int) {
	(*q)[i], (*q)[j] = (*q)[j], (*q)[i]
	(*q)[i].index = i
	(*q)[j].index = j
}

func (q *queue) Push(x any) {
	*q = append(*q, x.(*entry))
}

func (q *queue) Pop() any {
	length := q.Len()
	ans := (*q)[length-1]
	*q = (*q)[:length-1]
	return ans
}

func (q *queue) Update(en *entry, value any, freq int) {
	en.value = value
	en.freq = freq
	heap.Fix(q, en.index)
}
