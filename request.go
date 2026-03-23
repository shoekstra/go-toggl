package toggl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// Pagination contains pagination metadata parsed from API response headers.
//
// Not all endpoints populate every field. The page-based fields (CurrentPage,
// TotalPages, PageSize) are set when the response includes X-Page, X-Pages, and
// X-Page-Size headers. The cursor-based fields (NextID, NextRowNumber) are set
// for the detailed reports endpoint via X-Next-ID and X-Next-Row-Number headers.
//
// To iterate through all pages using page-based pagination:
//
//	perPage := 50
//	page := 1
//	for {
//	    items, resp, err := client.Projects.ListProjects(ctx, wsID, &toggl.ListProjectsOptions{
//	        Page:    toggl.Int(page),
//	        PerPage: toggl.Int(perPage),
//	    })
//	    if err != nil {
//	        return err
//	    }
//	    // process items...
//	    if resp.Pagination.TotalPages > 0 && resp.Pagination.CurrentPage >= resp.Pagination.TotalPages {
//	        break
//	    }
//	    if resp.Pagination.TotalPages == 0 && len(items) < perPage {
//	        break
//	    }
//	    page++
//	}
//
// To iterate through detailed report pages using cursor-based pagination:
//
//	var firstID, firstRowNumber *int
//	for {
//	    opts := &toggl.DetailedReportOptions{FirstID: firstID, FirstRowNumber: firstRowNumber}
//	    entries, resp, err := client.Reports.DetailedReport(ctx, wsID, opts)
//	    if err != nil {
//	        return err
//	    }
//	    // process entries...
//	    if resp.Pagination.NextID == 0 {
//	        break
//	    }
//	    firstID = toggl.Int(resp.Pagination.NextID)
//	    firstRowNumber = toggl.Int(resp.Pagination.NextRowNumber)
//	}
type Pagination struct {
	// CurrentPage is the current page number (X-Page header).
	CurrentPage int
	// TotalPages is the total number of pages (X-Pages header).
	TotalPages int
	// PageSize is the number of items per page (X-Page-Size header).
	PageSize int
	// NextID is the cursor for the next page of detailed report results (X-Next-ID header).
	NextID int
	// NextRowNumber is the row number cursor for the next page of detailed report results (X-Next-Row-Number header).
	NextRowNumber int
}

// Response wraps the standard http.Response with additional context.
type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
	Response   *http.Response
	Pagination Pagination
}

// newRequest creates a new HTTP request with common headers.
func (c *Client) newRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	url := c.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	// Set headers
	req.SetBasicAuth(c.token, basicAuthPassword)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)

	return req, nil
}

// do executes an HTTP request and returns the response.
func (c *Client) do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	pagination, err := parsePagination(resp.Header)
	if err != nil {
		return nil, err
	}

	response := &Response{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       respBody,
		Response:   resp,
		Pagination: pagination,
	}

	// Check for HTTP errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return response, &ErrorResponse{
			StatusCode: resp.StatusCode,
			Message:    string(respBody),
			Response:   resp,
		}
	}

	// Decode response if provided
	if v != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, v); err != nil {
			return response, fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return response, nil
}

// get performs a GET request.
func (c *Client) get(ctx context.Context, path string, v interface{}) (*Response, error) {
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	return c.do(ctx, req, v)
}

// post performs a POST request.
func (c *Client) post(ctx context.Context, path string, body, v interface{}) (*Response, error) {
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	return c.do(ctx, req, v)
}

// put performs a PUT request.
func (c *Client) put(ctx context.Context, path string, body, v interface{}) (*Response, error) {
	req, err := c.newRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return nil, err
	}
	return c.do(ctx, req, v)
}

// patch performs a PATCH request.
func (c *Client) patch(ctx context.Context, path string, body, v interface{}) (*Response, error) {
	req, err := c.newRequest(ctx, http.MethodPatch, path, body)
	if err != nil {
		return nil, err
	}
	return c.do(ctx, req, v)
}

// delete performs a DELETE request.
func (c *Client) delete(ctx context.Context, path string) (*Response, error) {
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	return c.do(ctx, req, nil)
}

// parsePagination extracts pagination metadata from response headers.
// It returns an error if any present header contains a non-integer value.
func parsePagination(h http.Header) (Pagination, error) {
	var (
		p   Pagination
		err error
	)
	if p.CurrentPage, err = headerInt(h, "X-Page"); err != nil {
		return Pagination{}, err
	}
	if p.TotalPages, err = headerInt(h, "X-Pages"); err != nil {
		return Pagination{}, err
	}
	if p.PageSize, err = headerInt(h, "X-Page-Size"); err != nil {
		return Pagination{}, err
	}
	if p.NextID, err = headerInt(h, "X-Next-ID"); err != nil {
		return Pagination{}, err
	}
	if p.NextRowNumber, err = headerInt(h, "X-Next-Row-Number"); err != nil {
		return Pagination{}, err
	}
	return p, nil
}

// headerInt reads an HTTP response header as an integer.
// It returns (0, nil) when the header is absent, and (0, err) when the
// value is present but cannot be parsed as an integer.
func headerInt(h http.Header, key string) (int, error) {
	v := h.Get(key)
	if v == "" {
		return 0, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("invalid pagination header %s: %w", key, err)
	}
	return n, nil
}
