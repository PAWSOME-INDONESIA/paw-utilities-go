package example

import "testing"

func TestEsQuery(t *testing.T) {
	esQuery()
}

func BenchmarkEsQuery(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		esQuery()
	}
}
