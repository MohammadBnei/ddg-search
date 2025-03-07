package duckduckgogo

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"ddg-search/util"
)

type SearchClient interface {
	Search(ctx context.Context, query string) ([]Result, error)
	SearchLimited(ctx context.Context, query string, limit int) ([]Result, error)
}

type DuckDuckGoSearchClient struct {
	baseUrl string
}

func NewDuckDuckGoSearchClient() *DuckDuckGoSearchClient {
	return &DuckDuckGoSearchClient{
		baseUrl: "https://duckduckgo.com/html/",
	}
}
func (c *DuckDuckGoSearchClient) Search(ctx context.Context, query string) ([]Result, error) {
	return c.SearchLimited(ctx, query, 0)
}

func (c *DuckDuckGoSearchClient) SearchLimited(ctx context.Context, query string, limit int) ([]Result, error) {
	queryURLStr := c.baseUrl + "?q=" + url.QueryEscape(query)
	queryURL, err := url.Parse(queryURLStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", queryURLStr, err)
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, queryURL.String(), nil)

	req.Header.Add("User-Agent", util.GetRandomUserAgent())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("return status code %d", resp.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	results := make([]Result, 0)
	doc.Find(".results .web-result").Each(func(i int, s *goquery.Selection) {
		if i > limit-1 && limit > 0 {
			return
		}
		results = append(results, c.collectResult(s))
	})
	return results, nil
}

func (c *DuckDuckGoSearchClient) collectResult(s *goquery.Selection) Result {
	resURLHTML := html(s.Find(".result__url").Html())
	resURL := clean(s.Find(".result__url").Text())
	titleHTML := html(s.Find(".result__a").Html())
	title := clean(s.Find(".result__a").Text())
	snippetHTML := html(s.Find(".result__snippet").Html())
	snippet := clean(s.Find(".result__snippet").Text())
	icon := s.Find(".result__icon__img")
	src, _ := icon.Attr("src")
	width, _ := icon.Attr("width")
	height, _ := icon.Attr("height")
	return Result{
		HTMLFormattedURL: resURLHTML,
		HTMLTitle:        titleHTML,
		HTMLSnippet:      snippetHTML,
		FormattedURL:     resURL,
		Title:            title,
		Snippet:          snippet,
		Icon: Icon{
			Src:    src,
			Width:  toInt(width),
			Height: toInt(height),
		},
	}
}

func html(html string, err error) string {
	if err != nil {
		return ""
	}
	return clean(html)
}

func clean(text string) string {
	return strings.TrimSpace(strings.ReplaceAll(text, "\n", ""))
}

func toInt(n string) int {
	res, err := strconv.Atoi(n)
	if err != nil {
		return 0
	}
	return res
}
