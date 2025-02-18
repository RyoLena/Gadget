package syncx

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestLimitPool(t *testing.T) {

	expectedMaxAttempts := 3
	expectedVal := []byte("A")

	pool := NewLimitPool(expectedMaxAttempts, func() []byte {
		var buffer bytes.Buffer
		buffer.Write(expectedVal)
		return buffer.Bytes()
	})

	var wg sync.WaitGroup
	bufChan := make(chan []byte, expectedMaxAttempts)

	// 从Pool中并发获取缓冲区
	for i := 0; i < expectedMaxAttempts; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			buf, ok := pool.Get()
			assert.True(t, ok)
			assert.NotZero(t, buf)
			assert.Equal(t, string(expectedVal), string(buf))

			bufChan <- buf
		}()
	}

	wg.Wait()
	close(bufChan)

	// 超过最大申请次数返回零值
	val, ok := pool.Get()
	assert.False(t, ok)
	assert.Zero(t, val)

	// 归还一个
	pool.Put(<-bufChan)

	// 再次申请仍可以拿到非零值缓冲区
	val, ok = pool.Get()
	assert.True(t, ok)
	assert.NotZero(t, string(expectedVal), string(val))

	// 超过最大申请次数返回零值
	val, ok = pool.Get()
	assert.False(t, ok)
	assert.Zero(t, val)
}
