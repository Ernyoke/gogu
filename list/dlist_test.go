package list

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoublyLinkedList(t *testing.T) {
	assert := assert.New(t)

	list := InitDList(1)
	assert.Equal(1, list.data)

	// Removal of the first node is not permitted.
	err := list.Delete(&list.doubleNode)
	assert.Error(err)

	list.Pop()
	node := list.Unshift(2)
	list.InsertBefore(node, 3)

	n := 3
	list.Each(func(i int) {
		assert.Equal(n, i)
		n--
	})

	node = list.Append(4)
	assert.Equal(4, node.data)

	n = 0
	expected := []int{3, 2, 1, 4}
	list.Each(func(i int) {
		assert.Equal(expected[n], i)
		n++
	})

	n1, found := list.Find(4)
	assert.Equal(4, n1.data)
	assert.True(found)

	n2, found := list.Find(10)
	assert.Nil(n2)
	assert.False(found)

	err = list.InsertAfter(node, 6)
	assert.NoError(err)

	n = 0
	expected = []int{3, 2, 1, 4, 6}
	list.Each(func(i int) {
		assert.Equal(expected[n], i)
		n++
	})

	n3 := list.Unshift(7)
	assert.Equal(7, n3.data)

	n = 0
	expected = []int{7, 3, 2, 1, 4, 6}
	list.Each(func(i int) {
		assert.Equal(expected[n], i)
		n++
	})

	list.Shift()
	list.Shift()
	n = 0
	expected = []int{2, 1, 4, 6}
	list.Each(func(i int) {
		assert.Equal(expected[n], i)
		n++
	})

	list.Pop()
	n = 0
	expected = []int{2, 1, 4}
	list.Each(func(i int) {
		assert.Equal(expected[n], i)
		n++
	})

	err = list.Delete(node)
	assert.NoError(err)

	// node with value 4 has been already removed from the list
	err = list.Delete(node)
	assert.Error(err)

	err = list.InsertBefore(node, 11111)
	assert.Error(err)

	_, err = list.Replace(10, 0)
	assert.Error(err)

	list.Replace(2, 0)
	n = 0
	expected = []int{0, 1}
	list.Each(func(i int) {
		assert.Equal(expected[n], i)
		n++
	})

	for i := 0; i < 10; i++ {
		list.Append(i)
		assert.Equal(i, list.Last())
	}
	list.Pop()
	list.Pop()
	last := list.Pop()
	assert.Equal(6, last.data)
	assert.Equal(6, list.Last())
	assert.Equal(0, list.First())

	list.Clear()
	assert.Equal(0, list.First())
	assert.Equal(0, list.Last())
}

func Example_DoublyLinkedList() {
	list := InitDList(1)

	values := []int{2, 3, 4, 5, 6, 7, 8}
	for _, val := range values {
		list.Append(val)
	}
	sl := []int{}
	list.Each(func(val int) {
		sl = append(sl, val)
	})
	fmt.Println(sl)

	item := list.Pop()
	fmt.Println(item.data)

	sl = nil
	list.Each(func(val int) {
		sl = append(sl, val)
	})
	fmt.Println(sl)

	item = list.Shift()
	fmt.Println(item.data)

	sl = nil
	list.Each(func(val int) {
		sl = append(sl, val)
	})
	fmt.Println(sl)

	item, err := list.Replace(20, 10)
	fmt.Println(err)
	fmt.Println(item)

	item, err = list.Replace(7, 8)
	fmt.Println(item.data)
	item, err = list.Replace(8, 7)

	n := list.Unshift(1)
	fmt.Println(n.data)

	last := list.Append(8)
	item, _ = list.Find(8)
	fmt.Println(item.data)

	list.InsertAfter(last, 9)

	sl = nil
	list.Each(func(val int) {
		sl = append(sl, val)
	})
	fmt.Println(sl)

	list.Delete(last)

	sl = nil
	list.Each(func(val int) {
		sl = append(sl, val)
	})
	fmt.Println(sl)

	fmt.Println(list.First())
	fmt.Println(list.Last())

	// Output:
	// [1 2 3 4 5 6 7 8]
	// 7
	// [1 2 3 4 5 6 7]
	// 1
	// [2 3 4 5 6 7]
	// requested node does not exists
	// <nil>
	// 8
	// 1
	// 8
	// [1 2 3 4 5 6 7 8 9]
	// [1 2 3 4 5 6 7 9]
	// 1
	// 9
}
