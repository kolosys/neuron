package neuron

import (
	"io"
	"testing"
)

type benchData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Data string `json:"data"`
}

func BenchmarkSerializeBody_Internal(b *testing.B) {
	data := benchData{
		ID:   1,
		Name: "test",
		Data: "some random data for benchmark to make it realistic",
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r := serializeBody(data)
			if r != nil {
				if rc, ok := r.(io.ReadCloser); ok {
					io.Copy(io.Discard, rc)
					rc.Close()
				} else {
					io.Copy(io.Discard, r)
				}
			}
		}
	})
}
