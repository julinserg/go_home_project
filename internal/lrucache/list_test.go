package lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})

	t.Run("complexMy", func(t *testing.T) {
		l := NewList()
		l.PushBack(10)
		l.PushBack(20)
		l.PushBack(30)
		l.Remove(l.Back())
		l.Remove(l.Front())
		require.Equal(t, 1, l.Len())
		require.Equal(t, 20, l.Front().Value)
		l.Remove(l.Front())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
		require.Equal(t, 0, l.Len())

		l.PushFront(30)
		l.PushFront(20)
		l.PushFront(10)

		elems1 := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems1 = append(elems1, i.Value.(int))
		}
		require.Equal(t, []int{10, 20, 30}, elems1)

		l.MoveToFront(l.Back())
		require.Nil(t, l.Front().Prev)
		require.Nil(t, l.Back().Next)
		l.MoveToFront(l.Back())
		require.Nil(t, l.Front().Prev)
		require.Nil(t, l.Back().Next)

		elems2 := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems2 = append(elems2, i.Value.(int))
		}
		require.Equal(t, []int{20, 30, 10}, elems2)

		l.Remove(l.Back())
		l.Remove(l.Front())

		require.Equal(t, 30, l.Back().Value)
		require.Equal(t, 1, l.Len())

		l.MoveToFront(l.Back())
		require.Nil(t, l.Front().Prev)
		require.Nil(t, l.Back().Next)

		require.Equal(t, 30, l.Back().Value)

		l.PushBack(40)
		l.PushBack(50)

		l.MoveToFront(l.Front().Next)
		require.Nil(t, l.Front().Prev)
		require.Nil(t, l.Back().Next)

		elems3 := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems3 = append(elems3, i.Value.(int))
		}
		require.Equal(t, []int{40, 30, 50}, elems3)
	})
}
