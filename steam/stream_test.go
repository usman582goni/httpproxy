package stream

import (
	"bytes"
	"io"
	"io/ioutil"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStream(t *testing.T) {
	const numReaders = 16
	const dataSize = 1024 * 1024 // 1 MiB

	stream := New(NewBuffer())

	wg := new(sync.WaitGroup)
	wg.Add(numReaders)

	r := rand.New(rand.NewSource(123))
	buffer := new(bytes.Buffer)
	io.CopyN(buffer, r, dataSize)
	content := buffer.Bytes()

	for i := 0; i < numReaders; i++ {
		go func(i int) {
			r, err := stream.NewReader()
			defer wg.Done()
			time.Sleep(time.Duration(i) * 50 * time.Millisecond)
			require.Nil(t, err)
			data, err := ioutil.ReadAll(r)
			require.Nil(t, err)
			assert.Equal(t, data, content)
		}(i)
	}
	io.Copy(stream, bytes.NewBuffer(content))
	err := stream.CloseWrite()
	assert.Nil(t, err)
	wg.Wait()
}
