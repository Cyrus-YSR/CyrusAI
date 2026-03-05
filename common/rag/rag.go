package rag

import (
	"GopherAI/common/redis"
	redisPkg "GopherAI/common/redis"
	"GopherAI/config"
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	embeddingArk "github.com/cloudwego/eino-ext/components/embedding/ark"
	redisIndexer "github.com/cloudwego/eino-ext/components/indexer/redis"
	redisRetriever "github.com/cloudwego/eino-ext/components/retriever/redis"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/schema"
	redisCli "github.com/redis/go-redis/v9"
)

const (
	defaultVectorTopK = 5
	defaultBM25TopK   = 20
	defaultFuseTopK   = 5
	defaultRRFK       = 60
)

type RAGIndexer struct {
	embedding embedding.Embedder
	indexer   *redisIndexer.Indexer
}

type RAGQuery struct {
	embedding  embedding.Embedder
	retriever  retriever.Retriever
	rdb        *redisCli.Client
	indexName  string
	vectorTopK int
	bm25TopK   int
	fuseTopK   int
	rrfK       int
}

// 构建知识库索引器
// 专业说法：文本解析、文本切块、向量化、存储向量化结果
// 通俗理解：把“人能读的文档”，转换成“AI 能按语义搜索的格式”，并存起来
func NewRAGIndexer(filename, embeddingModel string) (*RAGIndexer, error) {

	// 用于控制整个初始化流程（超时 / 取消等），这里先用默认背景即 context.Background()
	ctx := context.Background()

	// 从环境变量中读取调用向量模型所需的 API Key
	apiKey := os.Getenv("OPENAI_API_KEY")

	// 向量的维度大小（等于向量模型输出的数字个数）
	// Redis 在创建向量索引时必须提前知道这个维度大小
	dimension := config.GetConfig().RagModelConfig.RagDimension

	// 1. 配置并创建“向量生成器”（Embedding 组件）
	// 可以理解为：找一个“翻译官”，
	// 专门负责把文本翻译成 AI 能理解的“向量表示”
	embedConfig := &embeddingArk.EmbeddingConfig{
		BaseURL: config.GetConfig().RagModelConfig.RagBaseUrl, // 向量模型服务地址
		APIKey:  apiKey,                                       // 鉴权信息
		Model:   embeddingModel,                               // 使用哪个向量模型
	}

	// 创建向量生成器实例
	// 后续所有文本的“向量化”都会通过它完成
	embedder, err := embeddingArk.NewEmbedder(ctx, embedConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create embedder: %w", err)
	}

	// ===============================
	// 2. 初始化 Redis 中的向量索引结构
	// ===============================
	// 可以理解为：先在 Redis 里建好“仓库”，
	// 告诉它以后要存向量，并且每个向量的维度是多少
	if err := redisPkg.InitRedisIndex(ctx, filename, dimension); err != nil {
		return nil, fmt.Errorf("failed to init redis index: %w", err)
	}

	// 获取 Redis 客户端，用于后续数据写入
	rdb := redisPkg.Rdb

	// ===============================
	// 3. 配置索引器（定义：文档如何被存进 Redis）
	// ===============================
	indexerConfig := &redisIndexer.IndexerConfig{
		Client:    rdb,                                     // Redis 客户端
		KeyPrefix: redis.GenerateIndexNamePrefix(filename), // 不同知识库使用不同前缀，避免冲突
		BatchSize: 10,                                      // 批量处理文档，提高写入效率

		// 定义：一段文档（Document）在 Redis 中该如何存储
		DocumentToHashes: func(ctx context.Context, doc *schema.Document) (*redisIndexer.Hashes, error) {

			// 从文档的元数据中取出来源信息（例如文件名、URL 等）
			source := ""
			if s, ok := doc.MetaData["source"].(string); ok {
				source = s
			}

			// 构建 Redis 中实际存储的数据结构（Hash 结构）
			return &redisIndexer.Hashes{
				// Redis Key，一般由“知识库名称 + 文档 ID”组成，用于唯一标识一个文档
				Key: fmt.Sprintf("%s:%s", filename, doc.ID),

				// Redis Hash 中的字段
				Field2Value: map[string]redisIndexer.FieldValue{
					// content：原始文本内容，用于后续的向量计算
					// EmbedKey 表示：该字段需要先做向量化处理，生成的向量会存入名为 "vector" 的字段中
					"content": {Value: doc.Content, EmbedKey: "vector"},

					// metadata：一些辅助信息，不参与向量计算，直接存储在 Redis 中
					"metadata": {Value: source},
				},
			}, nil
		},
	}

	// 将“向量生成器”交给索引器
	// 这样索引器在写入文本时，可以自动完成向量计算
	indexerConfig.Embedding = embedder

	// ===============================
	// 4. 创建最终可用的索引器实例
	// ===============================
	// 此时索引器已经具备：
	// - 文本 生成向量 的能力（通过 embedder 实现）
	// - 向量写入 Redis 的能力（通过 indexerConfig 配置）
	idx, err := redisIndexer.NewIndexer(ctx, indexerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create indexer: %w", err)
	}

	// 返回一个封装好的 RAGIndexer 实例
	// 后续只需要调用它，就可以把文档加入知识库
	return &RAGIndexer{
		embedding: embedder,
		indexer:   idx,
	}, nil
}

// IndexFile 读取文件内容并创建向量索引
func (r *RAGIndexer) IndexFile(ctx context.Context, filePath string) error {
	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// 将文件内容转换为文档
	// TODO: 这里可以根据需要进行文本切块，目前简单处理为一个文档一个向量
	doc := &schema.Document{
		ID:      "doc_1", // 可以使用 UUID 或其他唯一标识
		Content: string(content),
		MetaData: map[string]any{
			"source": filePath,
		},
	}

	// 使用 indexer 存储文档（会自动进行向量化）
	_, err = r.indexer.Store(ctx, []*schema.Document{doc})
	if err != nil {
		return fmt.Errorf("failed to store document: %w", err)
	}

	return nil
}

// DeleteIndex 删除指定文件的知识库索引（静态方法，不依赖实例）
func DeleteIndex(ctx context.Context, filename string) error {
	if err := redisPkg.DeleteRedisIndex(ctx, filename); err != nil {
		return fmt.Errorf("failed to delete redis index: %w", err)
	}
	return nil
}

// NewRAGQuery 创建 RAG 查询器（用于向量检索和问答）
func NewRAGQuery(ctx context.Context, username string) (*RAGQuery, error) {
	cfg := config.GetConfig()
	apiKey := os.Getenv("OPENAI_API_KEY")

	// 创建 embedding 模型
	embedConfig := &embeddingArk.EmbeddingConfig{
		BaseURL: cfg.RagModelConfig.RagBaseUrl,
		APIKey:  apiKey,
		Model:   cfg.RagModelConfig.RagEmbeddingModel,
	}
	embedder, err := embeddingArk.NewEmbedder(ctx, embedConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create embedder: %w", err)
	}

	// 获取用户上传的文件名（假设每个用户只有一个文件）
	// 这里需要从用户目录读取文件，假设文件名为 "knowledge.txt"
	userDir := fmt.Sprintf("uploads/%s", username)
	files, err := os.ReadDir(userDir)
	if err != nil || len(files) == 0 {
		return nil, fmt.Errorf("no uploaded file found for user %s", username)
	}

	var filename string
	for _, f := range files {
		if !f.IsDir() {
			filename = f.Name()
			break
		}
	}

	if filename == "" {
		return nil, fmt.Errorf("no valid file found for user %s", username)
	}

	// 创建 retriever
	rdb := redisPkg.Rdb
	indexName := redis.GenerateIndexName(filename)

	retrieverConfig := &redisRetriever.RetrieverConfig{
		Client:       rdb,
		Index:        indexName,
		Dialect:      2,
		ReturnFields: []string{"content", "metadata", "distance"},
		TopK:         5,
		VectorField:  "vector",
		DocumentConverter: func(ctx context.Context, doc redisCli.Document) (*schema.Document, error) {
			resp := &schema.Document{
				ID:       doc.ID,
				Content:  "",
				MetaData: map[string]any{},
			}
			for field, val := range doc.Fields {
				if field == "content" {
					resp.Content = val
				} else {
					resp.MetaData[field] = val
				}
			}
			return resp, nil
		},
	}
	retrieverConfig.Embedding = embedder

	rtr, err := redisRetriever.NewRetriever(ctx, retrieverConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create retriever: %w", err)
	}

	return &RAGQuery{
		embedding:  embedder,
		retriever:  rtr,
		rdb:        rdb,
		indexName:  indexName,
		vectorTopK: defaultVectorTopK,
		bm25TopK:   defaultBM25TopK,
		fuseTopK:   defaultFuseTopK,
		rrfK:       defaultRRFK,
	}, nil
}

// RetrieveDocuments 检索相关文档
func (r *RAGQuery) RetrieveDocuments(ctx context.Context, query string) ([]*schema.Document, error) {
	vectorDocs, vecErr := r.retriever.Retrieve(ctx, query)
	bm25Docs, bm25Err := r.retrieveBM25(ctx, query, r.bm25TopK)

	if vecErr != nil && bm25Err != nil {
		return nil, fmt.Errorf("failed to retrieve documents: vector=%v bm25=%v", vecErr, bm25Err)
	}

	fused := fuseByRRF(vectorDocs, bm25Docs, r.fuseTopK, r.rrfK)
	if len(fused) > 0 {
		return fused, nil
	}

	if vecErr == nil {
		return vectorDocs, nil
	}
	return bm25Docs, nil
}

// BuildRAGPrompt 构建包含检索文档的提示
func BuildRAGPrompt(query string, docs []*schema.Document) string {
	if len(docs) == 0 {
		return query
	}

	contextText := ""
	for i, doc := range docs {
		contextText += fmt.Sprintf("[文档 %d]: %s\n\n", i+1, doc.Content)
	}

	prompt := fmt.Sprintf(`基于以下参考文档回答用户的问题。如果文档中没有相关信息，请说明无法找到相关信息。

参考文档：
%s

用户问题：%s

请提供准确、完整的回答：`, contextText, query)

	return prompt
}

func (r *RAGQuery) retrieveBM25(ctx context.Context, query string, topK int) ([]*schema.Document, error) {
	if r.rdb == nil {
		return nil, fmt.Errorf("redis client is nil")
	}
	if topK <= 0 {
		return []*schema.Document{}, nil
	}

	q := strings.TrimSpace(query)
	if q == "" {
		q = "*"
	} else {
		q = escapeRediSearchQuery(q)
	}

	args := []interface{}{
		"FT.SEARCH", r.indexName, q,
		"WITHSCORES",
		"LIMIT", "0", strconv.Itoa(topK),
		"RETURN", "2", "content", "metadata",
		"DIALECT", "2",
	}

	raw, err := r.rdb.Do(ctx, args...).Result()
	if err != nil {
		return nil, err
	}
	return parseFTSearchResult(raw)
}

func escapeRediSearchQuery(q string) string {
	var b strings.Builder
	b.Grow(len(q) + 8)
	for _, ch := range q {
		switch ch {
		case '\\', '-', '[', ']', '{', '}', '(', ')', '<', '>', '~', '*', '"', '\'', ':', ';', '!', '?', '@', '|', '&', '=', '+', '%', ',':
			b.WriteByte('\\')
			b.WriteRune(ch)
		default:
			b.WriteRune(ch)
		}
	}
	return b.String()
}

func parseFTSearchResult(raw interface{}) ([]*schema.Document, error) {
	arr, ok := raw.([]interface{})
	if !ok || len(arr) < 1 {
		return []*schema.Document{}, nil
	}

	docs := make([]*schema.Document, 0, max(0, (len(arr)-1)/2))

	i := 1
	for i < len(arr) {
		docID, ok := toString(arr[i])
		if !ok {
			return nil, fmt.Errorf("invalid FT.SEARCH response doc id at %d", i)
		}
		i++

		var scoreStr string
		if i < len(arr) {
			if s, ok := toString(arr[i]); ok {
				scoreStr = s
				i++
			}
		}

		if i >= len(arr) {
			break
		}

		fieldsArr, ok := arr[i].([]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid FT.SEARCH response fields at %d", i)
		}
		i++

		var content string
		meta := map[string]any{}
		if scoreStr != "" {
			if v, err := strconv.ParseFloat(scoreStr, 64); err == nil {
				meta["bm25_score"] = v
			} else {
				meta["bm25_score_raw"] = scoreStr
			}
		}

		for j := 0; j+1 < len(fieldsArr); j += 2 {
			k, ok := toString(fieldsArr[j])
			if !ok {
				continue
			}
			v, _ := toString(fieldsArr[j+1])
			if k == "content" {
				content = v
			} else {
				meta[k] = v
			}
		}

		docs = append(docs, &schema.Document{
			ID:       docID,
			Content:  content,
			MetaData: meta,
		})
	}

	return docs, nil
}

func toString(v interface{}) (string, bool) {
	switch t := v.(type) {
	case string:
		return t, true
	case []byte:
		return string(t), true
	default:
		return "", false
	}
}

func fuseByRRF(vectorDocs, bm25Docs []*schema.Document, topK int, rrfK int) []*schema.Document {
	if topK <= 0 {
		return []*schema.Document{}
	}
	if rrfK <= 0 {
		rrfK = defaultRRFK
	}

	type agg struct {
		doc   *schema.Document
		score float64
	}
	byID := map[string]*agg{}

	add := func(d *schema.Document, rank int) {
		if d == nil || d.ID == "" {
			return
		}
		a := byID[d.ID]
		if a == nil {
			a = &agg{doc: d}
			byID[d.ID] = a
		}
		a.score += 1.0 / float64(rrfK+rank)
	}

	for i, d := range vectorDocs {
		add(d, i+1)
	}
	for i, d := range bm25Docs {
		add(d, i+1)
	}

	all := make([]*agg, 0, len(byID))
	for _, a := range byID {
		all = append(all, a)
	}
	sort.Slice(all, func(i, j int) bool {
		if all[i].score == all[j].score {
			return all[i].doc.ID < all[j].doc.ID
		}
		return all[i].score > all[j].score
	})

	if len(all) > topK {
		all = all[:topK]
	}

	out := make([]*schema.Document, 0, len(all))
	for _, a := range all {
		if a.doc.MetaData == nil {
			a.doc.MetaData = map[string]any{}
		}
		a.doc.MetaData["rrf_score"] = a.score
		out = append(out, a.doc)
	}
	return out
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
