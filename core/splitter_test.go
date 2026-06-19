package core

import (
	"testing"
)

func makeRecords(n int) []Record {
	recs := make([]Record, n)
	for i := range n {
		recs[i] = Record{PrenomUtilisateur: "User"}
	}
	return recs
}

func TestSplitByLines_Exact(t *testing.T) {
	recs := makeRecords(100)
	chunks := SplitByLines(recs, 10)
	if len(chunks) != 10 {
		t.Fatalf("expected 10 chunks, got %d", len(chunks))
	}
	for i, chunk := range chunks {
		if len(chunk) != 10 {
			t.Fatalf("chunk %d: expected 10 rows, got %d", i, len(chunk))
		}
	}
}

func TestSplitByLines_WithRemainder(t *testing.T) {
	recs := makeRecords(105)
	chunks := SplitByLines(recs, 10)
	if len(chunks) != 11 {
		t.Fatalf("expected 11 chunks, got %d", len(chunks))
	}
	for i := 0; i < 10; i++ {
		if len(chunks[i]) != 10 {
			t.Fatalf("chunk %d: expected 10 rows, got %d", i, len(chunks[i]))
		}
	}
	if len(chunks[10]) != 5 {
		t.Fatalf("last chunk: expected 5 rows, got %d", len(chunks[10]))
	}
}

func TestSplitByLines_ZeroRecords(t *testing.T) {
	chunks := SplitByLines(nil, 10)
	if len(chunks) != 0 {
		t.Fatalf("expected 0 chunks, got %d", len(chunks))
	}
}

func TestSplitByLines_LinesLessThanOne(t *testing.T) {
	recs := makeRecords(5)
	chunks := SplitByLines(recs, 0)
	if len(chunks) != 5 {
		t.Fatalf("expected 5 chunks (clamped to 1), got %d", len(chunks))
	}
}

func TestSplitByFiles_Even(t *testing.T) {
	recs := makeRecords(100)
	chunks := SplitByFiles(recs, 5)
	if len(chunks) != 5 {
		t.Fatalf("expected 5 chunks, got %d", len(chunks))
	}
	for i, chunk := range chunks {
		if len(chunk) != 20 {
			t.Fatalf("chunk %d: expected 20 rows, got %d", i, len(chunk))
		}
	}
}

func TestSplitByFiles_Uneven(t *testing.T) {
	recs := makeRecords(103)
	chunks := SplitByFiles(recs, 5)
	if len(chunks) != 5 {
		t.Fatalf("expected 5 chunks, got %d", len(chunks))
	}
	// First 3 get 21, last 2 get 20
	expected := []int{21, 21, 21, 20, 20}
	for i, exp := range expected {
		if len(chunks[i]) != exp {
			t.Fatalf("chunk %d: expected %d rows, got %d", i, exp, len(chunks[i]))
		}
	}
}

func TestSplitByFiles_LessThanFiles(t *testing.T) {
	recs := makeRecords(3)
	chunks := SplitByFiles(recs, 5)
	if len(chunks) != 5 {
		t.Fatalf("expected 5 chunks, got %d", len(chunks))
	}
	if len(chunks[0]) != 1 {
		t.Fatalf("chunk 0: expected 1 row, got %d", len(chunks[0]))
	}
	for i := 3; i < 5; i++ {
		if len(chunks[i]) != 0 {
			t.Fatalf("chunk %d: expected 0 rows, got %d", i, len(chunks[i]))
		}
	}
}

func TestSplitByFiles_ZeroRecords(t *testing.T) {
	chunks := SplitByFiles(nil, 5)
	if len(chunks) != 5 {
		t.Fatalf("expected 5 empty chunks, got %d", len(chunks))
	}
}

func TestSplitByFiles_NumFilesLessThanOne(t *testing.T) {
	recs := makeRecords(10)
	chunks := SplitByFiles(recs, 0)
	if len(chunks) != 1 {
		t.Fatalf("expected 1 chunk (clamped), got %d", len(chunks))
	}
	if len(chunks[0]) != 10 {
		t.Fatalf("expected all 10 rows in single chunk, got %d", len(chunks[0]))
	}
}

func TestSplitByFiles_SingleFile(t *testing.T) {
	recs := makeRecords(50)
	chunks := SplitByFiles(recs, 1)
	if len(chunks) != 1 {
		t.Fatalf("expected 1 chunk, got %d", len(chunks))
	}
	if len(chunks[0]) != 50 {
		t.Fatalf("expected 50 rows, got %d", len(chunks[0]))
	}
}
