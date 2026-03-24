package toggl

import "context"

// MeService handles operations related to the authenticated user.
type MeService struct {
	client *Client
}

// GetMe returns the profile of the authenticated user.
//
// API: GET /api/v9/me
//
// See: https://engineering.toggl.com/docs/api/me#get-me
func (s *MeService) GetMe(ctx context.Context) (*Me, *Response, error) {
	me := new(Me)
	resp, err := s.client.get(ctx, "/api/v9/me", me)
	if err != nil {
		return nil, resp, err
	}

	return me, resp, nil
}
