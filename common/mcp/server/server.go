package mcp

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

//wttr.in JSON 响应结构

type WttrResponse struct {
	CurrentCondition []struct {
		TempC         string `json:"temp_C"`
		Humidity      string `json:"humidity"`
		WindspeedKmph string `json:"windspeedKmph"`
		WeatherDesc   []struct {
			Value string `json:"value"`
		} `json:"weatherDesc"`
	} `json:"current_condition"`

	NearestArea []struct {
		AreaName []struct {
			Value string `json:"value"`
		} `json:"areaName"`
	} `json:"nearest_area"`
}

//统一对外天气结构

type WeatherResponse struct {
	Location    string  `json:"location"`
	Temperature float64 `json:"temperature"`
	Condition   string  `json:"condition"`
	Humidity    int     `json:"humidity"`
	WindSpeed   float64 `json:"windSpeed"`
}

//Weather API Client

type WeatherAPIClient struct{}

func NewWeatherAPIClient() *WeatherAPIClient {
	return &WeatherAPIClient{}
}

func (c *WeatherAPIClient) GetWeather(ctx context.Context, city string) (*WeatherResponse, error) {
	apiURL := fmt.Sprintf(
		"https://wttr.in/%s?format=j1&lang=zh",
		city,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %w", err)
	}

	var wttrResp WttrResponse
	if err := json.Unmarshal(body, &wttrResp); err != nil {
		return nil, fmt.Errorf("json parse failed: %w", err)
	}

	if len(wttrResp.CurrentCondition) == 0 {
		return nil, fmt.Errorf("no weather data")
	}

	cc := wttrResp.CurrentCondition[0]

	temp, _ := strconv.ParseFloat(cc.TempC, 64)
	humidity, _ := strconv.Atoi(cc.Humidity)
	wind, _ := strconv.ParseFloat(cc.WindspeedKmph, 64)

	location := city
	if len(wttrResp.NearestArea) > 0 &&
		len(wttrResp.NearestArea[0].AreaName) > 0 {
		location = wttrResp.NearestArea[0].AreaName[0].Value
	}

	condition := "未知"
	if len(cc.WeatherDesc) > 0 {
		condition = cc.WeatherDesc[0].Value
	}

	return &WeatherResponse{
		Location:    location,
		Temperature: temp,
		Condition:   condition,
		Humidity:    humidity,
		WindSpeed:   wind,
	}, nil
}

type DDGInstantAnswer struct {
	AbstractText  string `json:"AbstractText"`
	AbstractURL   string `json:"AbstractURL"`
	RelatedTopics []struct {
		Text     string `json:"Text"`
		FirstURL string `json:"FirstURL"`
	} `json:"RelatedTopics"`
}

type SearchAPIClient struct{}

func NewSearchAPIClient() *SearchAPIClient {
	return &SearchAPIClient{}
}

type SearchResult struct {
	Title   string
	URL     string
	Snippet string
}

func (c *SearchAPIClient) ddgHTMLSearch(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	base := "https://duckduckgo.com/html/"
	u, err := url.Parse(base)
	if err != nil {
		return nil, fmt.Errorf("parse ddg html url failed: %w", err)
	}
	q := u.Query()
	q.Set("q", query)
	q.Set("kl", "cn-zh") // prefer Chinese when possible
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create ddg html request failed: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; MCP-Search/1.0)")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	resp, err := (&http.Client{Timeout: 10 * time.Second}).Do(req)
	if err != nil {
		return nil, fmt.Errorf("ddg html request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ddg html read failed: %w", err)
	}
	html := string(body)

	// Very simple extraction: anchors with class result__a
	re := regexp.MustCompile(`<a[^>]*class="result__a"[^>]*href="([^"]+)"[^>]*>(.*?)</a>`)
	matches := re.FindAllStringSubmatch(html, -1)
	results := make([]SearchResult, 0, limit)
	for _, m := range matches {
		if len(m) < 3 {
			continue
		}
		urlStr := htmlUnescape(m[1])
		title := stripTags(htmlUnescape(m[2]))
		if title == "" || urlStr == "" {
			continue
		}
		results = append(results, SearchResult{Title: title, URL: urlStr})
		if len(results) >= limit {
			break
		}
	}
	return results, nil
}

func (c *SearchAPIClient) baiduHTMLSearch(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	base := "https://www.baidu.com/s"
	u, err := url.Parse(base)
	if err != nil {
		return nil, fmt.Errorf("parse baidu url failed: %w", err)
	}
	q := u.Query()
	q.Set("wd", query)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create baidu request failed: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; MCP-Search/1.0)")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	resp, err := (&http.Client{Timeout: 10 * time.Second}).Do(req)
	if err != nil {
		return nil, fmt.Errorf("baidu request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("baidu read failed: %w", err)
	}
	html := string(body)

	if strings.Contains(html, "百度安全验证") || strings.Contains(html, "ppui-static-wap") || strings.Contains(html, "mkdjump") {
		return nil, fmt.Errorf("baidu verification required")
	}

	patterns := []*regexp.Regexp{
		regexp.MustCompile(`<h3[^>]*class="t"[^>]*>\s*<a[^>]*href="([^"]+)"[^>]*>(.*?)</a>`),
		regexp.MustCompile(`<h3[^>]*class="[^"]*c-title[^"]*"[^>]*>\s*<a[^>]*href="([^"]+)"[^>]*>(.*?)</a>`),
	}
	var matches [][]string
	for _, re := range patterns {
		matches = re.FindAllStringSubmatch(html, -1)
		if len(matches) > 0 {
			break
		}
	}
	results := make([]SearchResult, 0, limit)
	for _, m := range matches {
		if len(m) < 3 {
			continue
		}
		urlStr := htmlUnescape(m[1])
		title := stripTags(htmlUnescape(m[2]))
		if title == "" || urlStr == "" {
			continue
		}
		results = append(results, SearchResult{Title: title, URL: urlStr})
		if len(results) >= limit {
			break
		}
	}
	return results, nil
}

func (c *SearchAPIClient) bingHTMLSearch(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	base := "https://www.bing.com/search"
	u, err := url.Parse(base)
	if err != nil {
		return nil, fmt.Errorf("parse bing url failed: %w", err)
	}
	q := u.Query()
	q.Set("q", query)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create bing request failed: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; MCP-Search/1.0)")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	resp, err := (&http.Client{Timeout: 10 * time.Second}).Do(req)
	if err != nil {
		return nil, fmt.Errorf("bing request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("bing read failed: %w", err)
	}
	html := string(body)

	re := regexp.MustCompile(`(?s)<li[^>]*class="b_algo"[^>]*>.*?<h2>\s*<a[^>]*href="([^"]+)"[^>]*>(.*?)</a>.*?(?:<p>(.*?)</p>)?`)
	matches := re.FindAllStringSubmatch(html, -1)
	results := make([]SearchResult, 0, limit)
	for _, m := range matches {
		if len(m) < 3 {
			continue
		}
		urlStr := htmlUnescape(m[1])
		title := stripTags(htmlUnescape(m[2]))
		snippet := ""
		if len(m) >= 4 {
			snippet = stripTags(htmlUnescape(m[3]))
			snippet = strings.TrimSpace(snippet)
		}
		if title == "" || urlStr == "" {
			continue
		}
		results = append(results, SearchResult{Title: title, URL: urlStr, Snippet: snippet})
		if len(results) >= limit {
			break
		}
	}
	return results, nil
}

func (c *SearchAPIClient) bingRSSSearch(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	base := "https://www.bing.com/search"
	u, err := url.Parse(base)
	if err != nil {
		return nil, fmt.Errorf("parse bing url failed: %w", err)
	}
	q := u.Query()
	q.Set("q", query)
	q.Set("format", "rss")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create bing rss request failed: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; MCP-Search/1.0)")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	resp, err := (&http.Client{Timeout: 10 * time.Second}).Do(req)
	if err != nil {
		return nil, fmt.Errorf("bing rss request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("bing rss read failed: %w", err)
	}

	type rss struct {
		Channel struct {
			Items []struct {
				Title       string `xml:"title"`
				Link        string `xml:"link"`
				Description string `xml:"description"`
			} `xml:"item"`
		} `xml:"channel"`
	}

	var feed rss
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("bing rss parse failed: %w", err)
	}

	results := make([]SearchResult, 0, min(limit, len(feed.Channel.Items)))
	for _, it := range feed.Channel.Items {
		title := strings.TrimSpace(it.Title)
		link := strings.TrimSpace(it.Link)
		snippet := strings.TrimSpace(stripTags(it.Description))
		if title == "" || link == "" {
			continue
		}
		results = append(results, SearchResult{Title: title, URL: link, Snippet: snippet})
		if len(results) >= limit {
			break
		}
	}
	return results, nil
}

func htmlUnescape(s string) string {
	r := strings.NewReplacer("&amp;", "&", "&lt;", "<", "&gt;", ">", "&quot;", `"`, "&#39;", "'")
	return r.Replace(s)
}

func stripTags(s string) string {
	return regexp.MustCompile(`<[^>]+>`).ReplaceAllString(s, "")
}

func (c *SearchAPIClient) Search(ctx context.Context, query string) (string, error) {
	timeout := 15 * time.Second
	if v := strings.TrimSpace(os.Getenv("MCP_SEARCH_TIMEOUT")); v != "" {
		if d, err := time.ParseDuration(v); err == nil && d > 0 {
			timeout = d
		}
	}
	searchCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var lastBaiduErr error
	results, err := c.baiduHTMLSearch(searchCtx, query, 5)
	if err == nil && len(results) > 0 {
		return formatSearchResults(results), nil
	}
	if err != nil {
		lastBaiduErr = err
	} else {
		lastBaiduErr = fmt.Errorf("empty results")
	}

	var lastBingRSSErr error
	results, err = c.bingRSSSearch(searchCtx, query, 5)
	if err == nil && len(results) > 0 {
		return formatSearchResults(results), nil
	}
	if err != nil {
		lastBingRSSErr = err
	} else {
		lastBingRSSErr = fmt.Errorf("empty results")
	}

	var lastBingHTMLErr error
	results, err = c.bingHTMLSearch(searchCtx, query, 5)
	if err == nil && len(results) > 0 {
		return formatSearchResults(results), nil
	}
	if err != nil {
		lastBingHTMLErr = err
	} else {
		lastBingHTMLErr = fmt.Errorf("empty results")
	}

	return "", fmt.Errorf("search failed: baidu=%v; bing_rss=%v; bing_html=%v", lastBaiduErr, lastBingRSSErr, lastBingHTMLErr)
}

func formatSearchResults(results []SearchResult) string {
	var b strings.Builder
	b.WriteString("搜索结果:\n")
	for _, r := range results {
		b.WriteString("- [")
		b.WriteString(r.Title)
		b.WriteString("](")
		b.WriteString(r.URL)
		b.WriteString(")\n")
		if s := strings.TrimSpace(r.Snippet); s != "" {
			b.WriteString("  摘要: ")
			b.WriteString(s)
			b.WriteString("\n")
		}
	}
	return b.String()
}

/*
	========================
	MCP Server
	========================
*/

func NewMCPServer() *server.MCPServer {
	weatherClient := NewWeatherAPIClient()
	searchClient := NewSearchAPIClient()

	mcpServer := server.NewMCPServer(
		"weather-query-server",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithLogging(),
	)

	mcpServer.AddTool(
		mcp.NewTool(
			"get_weather",
			mcp.WithDescription("获取指定城市的天气信息"),
			mcp.WithString(
				"city",
				mcp.Description("城市名称，如 Beijing、上海"),
				mcp.Required(),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := request.GetArguments()
			city, ok := args["city"].(string)
			if !ok || city == "" {
				return nil, fmt.Errorf("invalid city argument")
			}

			weather, err := weatherClient.GetWeather(ctx, city)
			if err != nil {
				return nil, err
			}

			resultText := fmt.Sprintf(
				"城市: %s\n温度: %.1f°C\n天气: %s\n湿度: %d%%\n风速: %.1f km/h",
				weather.Location,
				weather.Temperature,
				weather.Condition,
				weather.Humidity,
				weather.WindSpeed,
			)

			return &mcp.CallToolResult{
				Content: []mcp.Content{
					mcp.TextContent{
						Type: "text",
						Text: resultText,
					},
				},
			}, nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool(
			"web_search",
			mcp.WithDescription("通过搜索引擎检索互联网信息"),
			mcp.WithString(
				"query",
				mcp.Description("搜索关键词"),
				mcp.Required(),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := request.GetArguments()
			raw, ok := args["query"].(string)
			if !ok {
				return nil, fmt.Errorf("invalid query argument")
			}
			query := strings.TrimSpace(raw)
			if query == "" {
				return nil, fmt.Errorf("query cannot be empty")
			}

			text, err := searchClient.Search(ctx, query)
			if err != nil {
				return nil, err
			}

			return &mcp.CallToolResult{
				Content: []mcp.Content{
					mcp.TextContent{
						Type: "text",
						Text: text,
					},
				},
			}, nil
		},
	)

	return mcpServer
}

// StartServer 启动MCP服务器
// httpAddr: HTTP服务器监听的地址（例如":8080"）
func StartServer(httpAddr string) error {
	mcpServer := NewMCPServer()

	httpServer := server.NewStreamableHTTPServer(mcpServer)
	log.Printf("HTTP MCP server listening on %s/mcp", httpAddr)
	return httpServer.Start(httpAddr)
}
