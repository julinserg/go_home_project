package lrucache

import (
	"crypto/rand"
	"math/big"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(3)

		c.Set("a", 1)
		c.Set("b", 2)
		c.Set("c", 3)
		c.Set("d", 4)

		val, ok := c.Get("a")
		require.False(t, ok)
		require.Nil(t, val)
		val, ok = c.Get("b")
		require.True(t, ok)
		require.Equal(t, 2, val)
		val, ok = c.Get("c")
		require.True(t, ok)
		require.Equal(t, 3, val)
		val, ok = c.Get("d")
		require.True(t, ok)
		require.Equal(t, 4, val)

		c.Clear()

		c.Set("a", 1)
		c.Set("b", 2)
		c.Set("c", 3)

		c.Set("a", 10)
		c.Set("b", 20)
		c.Set("d", 4)

		val, ok = c.Get("c")
		require.False(t, ok)
		require.Nil(t, val)
		val, ok = c.Get("a")
		require.True(t, ok)
		require.Equal(t, 10, val)
		val, ok = c.Get("b")
		require.True(t, ok)
		require.Equal(t, 20, val)
		val, ok = c.Get("d")
		require.True(t, ok)
		require.Equal(t, 4, val)

		c.Clear()

		c.Set("a", 1)
		c.Set("b", 2)
		c.Set("c", 3)

		c.Get("a")
		c.Get("c")
		c.Set("d", 4)

		val, ok = c.Get("b")
		require.False(t, ok)
		require.Nil(t, val)
		val, ok = c.Get("a")
		require.True(t, ok)
		require.Equal(t, 1, val)
		val, ok = c.Get("c")
		require.True(t, ok)
		require.Equal(t, 3, val)
		val, ok = c.Get("d")
		require.True(t, ok)
		require.Equal(t, 4, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			nBig, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
			if err != nil {
				return
			}
			c.Get(Key(strconv.FormatInt(nBig.Int64(), 10)))
		}
	}()

	wg.Wait()
}
