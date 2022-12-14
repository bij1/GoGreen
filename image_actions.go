package godo

import (
	"context"
	"fmt"
	"net/http"
)

// ImageActionsService is an interface for interfacing with the image actions
// endpoints of the DigitalOcean API
// See: https://docs.digitalocean.com/reference/api/api-reference/#tag/Image-Actions
type ImageActionsService interface {
	Get(context.Context, string, int) (*Action, *Response, error)
	Transfer(context.Context, string, *ActionRequest) (*Action, *Response, error)
	Convert(context.Context, string) (*Action, *Response, error)
}

// ImageActionsServiceOp handles communication with the image action related methods of the
// DigitalOcean API.
type ImageActionsServiceOp struct {
	client *Client
}

var _ ImageActionsService = &ImageActionsServiceOp{}

// Transfer an image
func (i *ImageActionsServiceOp) Transfer(ctx context.Context, imageID string, transferRequest *ActionRequest) (*Action, *Response, error) {
	if transferRequest == nil {
		return nil, nil, NewArgError("transferRequest", "cannot be nil")
	}

	path := fmt.Sprintf("v2/images/%s/actions", imageID)

	req, err := i.client.NewRequest(ctx, http.MethodPost, path, transferRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(actionRoot)
	resp, err := i.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Event, resp, err
}

// Convert an image to a snapshot
func (i *ImageActionsServiceOp) Convert(ctx context.Context, imageID string) (*Action, *Response, error) {
	path := fmt.Sprintf("v2/images/%s/actions", imageID)

	convertRequest := &ActionRequest{
		"type": "convert",
	}

	req, err := i.client.NewRequest(ctx, http.MethodPost, path, convertRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(actionRoot)
	resp, err := i.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Event, resp, err
}

// Get an action for a particular image by id.
func (i *ImageActionsServiceOp) Get(ctx context.Context, imageID string, actionID int) (*Action, *Response, error) {
	if actionID < 1 {
		return nil, nil, NewArgError("actionID", "cannot be less than 1")
	}

	path := fmt.Sprintf("v2/images/%s/actions/%d", imageID, actionID)

	req, err := i.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(actionRoot)
	resp, err := i.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Event, resp, err
}
