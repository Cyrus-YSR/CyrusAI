package rag

import (
	"testing"

	"github.com/cloudwego/eino/schema"
)

func TestEscapeRediSearchQuery(t *testing.T) {
	in := `hello(world) a+b c:d "e" f*g`
	out := escapeRediSearchQuery(in)
	if out == in {
		t.Fatalf("expected escaped query, got unchanged: %q", out)
	}
	if len(out) <= len(in) {
		t.Fatalf("expected escaped query to be longer, in=%q out=%q", in, out)
	}
}

func TestParseFTSearchResult_WithScores(t *testing.T) {
	raw := []interface{}{
		int64(2),
		"doc1", "1.23", []interface{}{"content", "hello", "metadata", "m1"},
		"doc2", "0.50", []interface{}{"content", "world", "metadata", "m2"},
	}
	docs, err := parseFTSearchResult(raw)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(docs) != 2 {
		t.Fatalf("expected 2 docs, got %d", len(docs))
	}
	if docs[0].ID != "doc1" || docs[0].Content != "hello" {
		t.Fatalf("unexpected doc0: %#v", docs[0])
	}
	if v, ok := docs[0].MetaData["bm25_score"].(float64); !ok || v <= 0 {
		t.Fatalf("expected bm25_score float64 > 0, got %#v", docs[0].MetaData["bm25_score"])
	}
}

func TestFuseByRRF(t *testing.T) {
	vectorDocs := []*schema.Document{
		{ID: "a", Content: "va"},
		{ID: "b", Content: "vb"},
	}
	bm25Docs := []*schema.Document{
		{ID: "b", Content: "kb"},
		{ID: "c", Content: "kc"},
	}
	out := fuseByRRF(vectorDocs, bm25Docs, 3, 60)
	if len(out) != 3 {
		t.Fatalf("expected 3 docs, got %d", len(out))
	}
	if out[0].ID != "b" {
		t.Fatalf("expected 'b' first, got %q", out[0].ID)
	}
	if out[1].ID != "a" || out[2].ID != "c" {
		t.Fatalf("expected order b,a,c got %q,%q,%q", out[0].ID, out[1].ID, out[2].ID)
	}
	if _, ok := out[0].MetaData["rrf_score"]; !ok {
		t.Fatalf("expected rrf_score in metadata")
	}
}
