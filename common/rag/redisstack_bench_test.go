package rag

import (
	"context"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"testing"
	"time"

	redisCli "github.com/redis/go-redis/v9"
)

func TestRedisVectorLatencyProfile(t *testing.T) {
	if os.Getenv("REDIS_PROFILE_RUN") == "" {
		t.Skip("set REDIS_PROFILE_RUN=1 to run")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	opts := redisBenchOptionsFromEnv()
	rdb := newBenchRedisClient(opts)

	defer func() { _ = rdb.Close() }()

	indexName, prefix, cleanup, err := prepareVectorIndex(ctx, rdb, opts)
	if err != nil {
		t.Fatalf("prepare index failed: %v", err)
	}
	defer cleanup()

	queryVec := make([]float32, opts.dim)
	fillRandomFloat32(queryVec, rand.New(rand.NewSource(7)))
	queryBlob := packFloat32LE(queryVec)

	durations := make([]time.Duration, 0, opts.profileIters)
	for i := 0; i < opts.profileIters; i++ {
		start := time.Now()
		if err := doKNNQuery(ctx, rdb, indexName, queryBlob, opts.topK); err != nil {
			t.Fatalf("query failed: %v", err)
		}
		durations = append(durations, time.Since(start))
	}

	sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
	p50 := percentileDuration(durations, 0.50)
	p95 := percentileDuration(durations, 0.95)
	p99 := percentileDuration(durations, 0.99)

	t.Logf("redis=%s index=%s prefix=%s mode=%s dim=%d n=%d topK=%d", opts.addr, indexName, prefix, opts.mode, opts.dim, opts.n, opts.topK)
	t.Logf("latency: p50=%s p95=%s p99=%s iters=%d", p50, p95, p99, opts.profileIters)
}

func BenchmarkRedisVectorSearch(b *testing.B) {
	if os.Getenv("REDIS_BENCH_RUN") == "" {
		b.Skip("set REDIS_BENCH_RUN=1 to run")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	opts := redisBenchOptionsFromEnv()
	rdb := newBenchRedisClient(opts)
	defer func() { _ = rdb.Close() }()

	indexName, _, cleanup, err := prepareVectorIndex(ctx, rdb, opts)
	if err != nil {
		b.Fatalf("prepare index failed: %v", err)
	}
	defer cleanup()

	queryVec := make([]float32, opts.dim)
	fillRandomFloat32(queryVec, rand.New(rand.NewSource(7)))
	queryBlob := packFloat32LE(queryVec)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := doKNNQuery(ctx, rdb, indexName, queryBlob, opts.topK); err != nil {
			b.Fatalf("query failed: %v", err)
		}
	}
}

type redisBenchOptions struct {
	addr         string
	password     string
	db           int
	mode         string
	dim          int
	n            int
	topK         int
	batch        int
	hnswM        int
	hnswEFConstr int
	profileIters int
	allowLargeN  bool
	dialTimeout  time.Duration
	readTimeout  time.Duration
	writeTimeout time.Duration
	poolTimeout  time.Duration
}

func redisBenchOptionsFromEnv() redisBenchOptions {
	return redisBenchOptions{
		addr:         envString("REDIS_ADDR", "127.0.0.1:6379"),
		password:     envString("REDIS_PASSWORD", ""),
		db:           envInt("REDIS_DB", 0),
		mode:         envString("REDIS_VECTOR_MODE", "FLAT"),
		dim:          envInt("REDIS_VECTOR_DIM", 1024),
		n:            envInt("REDIS_VECTOR_N", 50000),
		topK:         envInt("REDIS_VECTOR_TOPK", 5),
		batch:        envInt("REDIS_VECTOR_BATCH", 500),
		hnswM:        envInt("REDIS_HNSW_M", 16),
		hnswEFConstr: envInt("REDIS_HNSW_EF_CONSTRUCTION", 200),
		profileIters: envInt("REDIS_PROFILE_ITERS", 500),
		allowLargeN:  os.Getenv("REDIS_BENCH_ALLOW_LARGE_N") == "1",
		dialTimeout:  envDurationMs("REDIS_DIAL_TIMEOUT_MS", 10_000),
		readTimeout:  envDurationMs("REDIS_READ_TIMEOUT_MS", 600_000),
		writeTimeout: envDurationMs("REDIS_WRITE_TIMEOUT_MS", 600_000),
		poolTimeout:  envDurationMs("REDIS_POOL_TIMEOUT_MS", 600_000),
	}
}

func newBenchRedisClient(opts redisBenchOptions) *redisCli.Client {
	return redisCli.NewClient(&redisCli.Options{
		Addr:         opts.addr,
		Password:     opts.password,
		DB:           opts.db,
		Protocol:     2,
		DialTimeout:  opts.dialTimeout,
		ReadTimeout:  opts.readTimeout,
		WriteTimeout: opts.writeTimeout,
		PoolTimeout:  opts.poolTimeout,
	})
}

func prepareVectorIndex(ctx context.Context, rdb *redisCli.Client, opts redisBenchOptions) (string, string, func(), error) {
	if err := ensureRediSearch(ctx, rdb); err != nil {
		return "", "", func() {}, err
	}

	if opts.n > 5_000_000 && !opts.allowLargeN {
		return "", "", func() {}, fmt.Errorf("REDIS_VECTOR_N=%d too large; set REDIS_BENCH_ALLOW_LARGE_N=1 to proceed", opts.n)
	}

	now := time.Now().UnixNano()
	indexName := fmt.Sprintf("bench_vec_idx_%d", now)
	prefix := fmt.Sprintf("bench_vec_doc_%d:", now)

	if err := createVectorIndex(ctx, rdb, indexName, prefix, opts); err != nil {
		return "", "", func() {}, err
	}

	if err := populateVectors(ctx, rdb, prefix, opts); err != nil {
		_ = dropIndexDD(ctx, rdb, indexName)
		return "", "", func() {}, err
	}

	cleanup := func() { _ = dropIndexDD(context.Background(), rdb, indexName) }
	return indexName, prefix, cleanup, nil
}

func ensureRediSearch(ctx context.Context, rdb *redisCli.Client) error {
	if err := rdb.Do(ctx, "FT._LIST").Err(); err != nil {
		return fmt.Errorf("redis does not support RediSearch commands: %w", err)
	}
	return nil
}

func createVectorIndex(ctx context.Context, rdb *redisCli.Client, indexName, prefix string, opts redisBenchOptions) error {
	mode := stringsToUpper(opts.mode)
	if mode != "FLAT" && mode != "HNSW" {
		return fmt.Errorf("unsupported REDIS_VECTOR_MODE=%s", opts.mode)
	}

	args := []interface{}{
		"FT.CREATE", indexName,
		"ON", "HASH",
		"PREFIX", "1", prefix,
		"SCHEMA",
		"content", "TEXT",
		"vector", "VECTOR", mode,
	}

	if mode == "FLAT" {
		args = append(args,
			"6",
			"TYPE", "FLOAT32",
			"DIM", strconv.Itoa(opts.dim),
			"DISTANCE_METRIC", "COSINE",
		)
	} else {
		args = append(args,
			"10",
			"TYPE", "FLOAT32",
			"DIM", strconv.Itoa(opts.dim),
			"DISTANCE_METRIC", "COSINE",
			"M", strconv.Itoa(opts.hnswM),
			"EF_CONSTRUCTION", strconv.Itoa(opts.hnswEFConstr),
		)
	}

	return rdb.Do(ctx, args...).Err()
}

func populateVectors(ctx context.Context, rdb *redisCli.Client, prefix string, opts redisBenchOptions) error {
	if opts.n <= 0 {
		return nil
	}
	if opts.batch <= 0 {
		opts.batch = 200
	}

	rng := rand.New(rand.NewSource(1))
	vec := make([]float32, opts.dim)

	for start := 0; start < opts.n; start += opts.batch {
		end := start + opts.batch
		if end > opts.n {
			end = opts.n
		}
		pipe := rdb.Pipeline()
		for i := start; i < end; i++ {
			fillRandomFloat32(vec, rng)
			blob := packFloat32LE(vec)
			key := fmt.Sprintf("%s%d", prefix, i)
			pipe.HSet(ctx, key, "content", fmt.Sprintf("doc-%d", i), "vector", blob)
		}
		if _, err := pipe.Exec(ctx); err != nil {
			return err
		}
	}

	return nil
}

func doKNNQuery(ctx context.Context, rdb *redisCli.Client, indexName string, vecBlob []byte, topK int) error {
	if topK <= 0 {
		topK = 5
	}
	query := fmt.Sprintf("*=>[KNN %d @vector $vec AS score]", topK)
	args := []interface{}{
		"FT.SEARCH", indexName, query,
		"PARAMS", "2", "vec", vecBlob,
		"SORTBY", "score",
		"RETURN", "2", "content", "score",
		"DIALECT", "2",
	}
	return rdb.Do(ctx, args...).Err()
}

func dropIndexDD(ctx context.Context, rdb *redisCli.Client, indexName string) error {
	return rdb.Do(ctx, "FT.DROPINDEX", indexName, "DD").Err()
}

func fillRandomFloat32(dst []float32, rng *rand.Rand) {
	for i := range dst {
		dst[i] = rng.Float32()
	}
}

func packFloat32LE(v []float32) []byte {
	b := make([]byte, 4*len(v))
	for i, x := range v {
		binary.LittleEndian.PutUint32(b[i*4:], math.Float32bits(x))
	}
	return b
}

func percentileDuration(sorted []time.Duration, p float64) time.Duration {
	if len(sorted) == 0 {
		return 0
	}
	if p <= 0 {
		return sorted[0]
	}
	if p >= 1 {
		return sorted[len(sorted)-1]
	}
	idx := int(float64(len(sorted)-1) * p)
	return sorted[idx]
}

func envString(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func envInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

func envDurationMs(key string, defMs int) time.Duration {
	ms := envInt(key, defMs)
	if ms <= 0 {
		return 0
	}
	return time.Duration(ms) * time.Millisecond
}

func stringsToUpper(s string) string {
	if s == "" {
		return s
	}
	b := []byte(s)
	for i := range b {
		if b[i] >= 'a' && b[i] <= 'z' {
			b[i] = b[i] - 'a' + 'A'
		}
	}
	return string(b)
}
