package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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
	Title string
	URL   string
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

	// Extract Baidu search results: <h3 class="t"><a href="...">Title</a>
	re := regexp.MustCompile(`<h3[^>]*class="t"[^>]*>\s*<a[^>]*href="([^"]+)"[^>]*>(.*?)</a>`)
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

func htmlUnescape(s string) string {
	r := strings.NewReplacer("&amp;", "&", "&lt;", "<", "&gt;", ">", "&quot;", `"`, "&#39;", "'")
	return r.Replace(s)
}

func stripTags(s string) string {
	return regexp.MustCompile(`<[^>]+>`).ReplaceAllString(s, "")
}

func (c *SearchAPIClient) Search(ctx context.Context, query string) (string, error) {
	apiURL, err := url.Parse("https://api.duckduckgo.com/")
	if err != nil {
		return "", fmt.Errorf("parse base url failed: %w", err)
	}

	q := apiURL.Query()
	q.Set("q", query)
	q.Set("format", "json")
	q.Set("no_redirect", "1")
	q.Set("no_html", "1")
	apiURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL.String(), nil)
	if err != nil {
		return "", fmt.Errorf("create request failed: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		// Fallback to HTML providers when IA fails
		results, derr := c.ddgHTMLSearch(ctx, query, 5)
		if derr != nil || len(results) == 0 {
			results, berr := c.baiduHTMLSearch(ctx, query, 5)
			if berr != nil || len(results) == 0 {
				return "", fmt.Errorf("http request failed: %w", err)
			}
		}
		var b strings.Builder
		b.WriteString("搜索结果:\n")
		for _, r := range results {
			b.WriteString("- [")
			b.WriteString(r.Title)
			b.WriteString("](")
			b.WriteString(r.URL)
			b.WriteString(")\n")
		}
		return b.String(), nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response failed: %w", err)
	}

	var ia DDGInstantAnswer
	if err := json.Unmarshal(body, &ia); err != nil {
		trimmed := strings.TrimSpace(string(body))
		if len(trimmed) > 2000 {
			trimmed = trimmed[:2000] + "..."
		}
		// Fallback to HTML search parsing
		results, err := c.ddgHTMLSearch(ctx, query, 5)
		if err != nil || len(results) == 0 {
			results, err = c.baiduHTMLSearch(ctx, query, 5)
			if err != nil || len(results) == 0 {
				return trimmed, nil
			}
		}
		var b strings.Builder
		b.WriteString("搜索结果:\n")
		for _, r := range results {
			b.WriteString("- [")
			b.WriteString(r.Title)
			b.WriteString("](")
			b.WriteString(r.URL)
			b.WriteString(")\n")
		}
		return b.String(), nil
	}

	var b strings.Builder
	if ia.AbstractText != "" {
		b.WriteString("摘要:\n")
		b.WriteString(ia.AbstractText)
		b.WriteString("\n")
	}
	if ia.AbstractURL != "" {
		b.WriteString("参考链接: [")
		b.WriteString(ia.AbstractURL)
		b.WriteString("](")
		b.WriteString(ia.AbstractURL)
		b.WriteString(")\n")
	}

	count := 0
	for _, t := range ia.RelatedTopics {
		if t.Text == "" {
			continue
		}
		if count == 0 {
			if b.Len() > 0 {
				b.WriteString("\n")
			}
			b.WriteString("相关结果:\n")
		}
		if t.FirstURL != "" {
			b.WriteString("- [")
			b.WriteString(t.Text)
			b.WriteString("](")
			b.WriteString(t.FirstURL)
			b.WriteString(")")
		} else {
			b.WriteString("- ")
			b.WriteString(t.Text)
		}
		b.WriteString("\n")
		count++
		if count >= 5 {
			break
		}
	}

	result := strings.TrimSpace(b.String())
	if result == "" {
		trimmed := strings.TrimSpace(string(body))
		if len(trimmed) > 2000 {
			trimmed = trimmed[:2000] + "..."
		}
		// Fallback to HTML search parsing
		results, err := c.ddgHTMLSearch(ctx, query, 5)
		if err != nil || len(results) == 0 {
			results, err = c.baiduHTMLSearch(ctx, query, 5)
			if err != nil || len(results) == 0 {
				return trimmed, nil
			}
		}
		var bb strings.Builder
		bb.WriteString("搜索结果:\n")
		for _, r := range results {
			bb.WriteString("- [")
			bb.WriteString(r.Title)
			bb.WriteString("](")
			bb.WriteString(r.URL)
			bb.WriteString(")\n")
		}
		return bb.String(), nil
	}
	return result, nil
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
