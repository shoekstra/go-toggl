package toggl

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// TagsService handles operations related to tags.
type TagsService struct {
	client *Client
}

// ListTagsOptions specifies the optional parameters to
// TagsService.ListTags.
type ListTagsOptions struct {
	// Search filters tags by name (case-insensitive substring match).
	Search *string
	// Page is the page number for pagination.
	Page *int
	// PerPage is the number of results per page.
	PerPage *int
}

// CreateTagOptions specifies the parameters to TagsService.CreateTag.
type CreateTagOptions struct {
	// Name is required. The tag name.
	Name string
}

// UpdateTagOptions specifies the parameters to TagsService.UpdateTag.
type UpdateTagOptions struct {
	// Name is required. The new tag name.
	Name string
}

// ListTags lists all tags in the given workspace.
//
// API: GET /api/v9/workspaces/{workspace_id}/tags
//
// See: https://engineering.toggl.com/docs/api/tags#get-tags
func (s *TagsService) ListTags(ctx context.Context, workspaceID int, opts *ListTagsOptions) ([]*Tag, *Response, error) {
	path := fmt.Sprintf("/api/v9/workspaces/%d/tags", workspaceID)

	if opts != nil {
		params := url.Values{}
		if opts.Search != nil {
			params.Set("search", *opts.Search)
		}
		if opts.Page != nil {
			params.Set("page", strconv.Itoa(*opts.Page))
		}
		if opts.PerPage != nil {
			params.Set("per_page", strconv.Itoa(*opts.PerPage))
		}
		if q := params.Encode(); q != "" {
			path += "?" + q
		}
	}

	var tags []*Tag
	resp, err := s.client.get(ctx, path, &tags)
	if err != nil {
		return nil, resp, err
	}

	return tags, resp, nil
}

// CreateTag creates a new tag in the given workspace.
//
// API: POST /api/v9/workspaces/{workspace_id}/tags
//
// See: https://engineering.toggl.com/docs/api/tags#post-create-tag
func (s *TagsService) CreateTag(ctx context.Context, workspaceID int, opts *CreateTagOptions) (*Tag, *Response, error) {
	if opts == nil {
		return nil, nil, fmt.Errorf("options required")
	}
	if opts.Name == "" {
		return nil, nil, fmt.Errorf("name is required")
	}

	path := fmt.Sprintf("/api/v9/workspaces/%d/tags", workspaceID)

	body := map[string]interface{}{
		"name": opts.Name,
	}

	tag := new(Tag)
	resp, err := s.client.post(ctx, path, body, tag)
	if err != nil {
		return nil, resp, err
	}

	return tag, resp, nil
}

// UpdateTag updates an existing tag. The name is required.
//
// API: PUT /api/v9/workspaces/{workspace_id}/tags/{tag_id}
//
// See: https://engineering.toggl.com/docs/api/tags#put-update-tag
func (s *TagsService) UpdateTag(ctx context.Context, workspaceID, tagID int, opts *UpdateTagOptions) (*Tag, *Response, error) {
	if opts == nil {
		return nil, nil, fmt.Errorf("options required")
	}
	if opts.Name == "" {
		return nil, nil, fmt.Errorf("name is required")
	}

	path := fmt.Sprintf("/api/v9/workspaces/%d/tags/%d", workspaceID, tagID)

	body := map[string]interface{}{
		"name": opts.Name,
	}

	tag := new(Tag)
	resp, err := s.client.put(ctx, path, body, tag)
	if err != nil {
		return nil, resp, err
	}

	return tag, resp, nil
}

// DeleteTag deletes a tag from the given workspace.
//
// API: DELETE /api/v9/workspaces/{workspace_id}/tags/{tag_id}
//
// See: https://engineering.toggl.com/docs/api/tags#delete-delete-tag
func (s *TagsService) DeleteTag(ctx context.Context, workspaceID, tagID int) (*Response, error) {
	path := fmt.Sprintf("/api/v9/workspaces/%d/tags/%d", workspaceID, tagID)
	return s.client.delete(ctx, path)
}
