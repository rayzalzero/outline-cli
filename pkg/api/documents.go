package api

import (
	"encoding/json"
	"time"
)

// Document represents an Outline document
type Document struct {
	ID             string    `json:"id"`
	CollectionID   string    `json:"collectionId"`
	Title          string    `json:"title"`
	Text           string    `json:"text"`
	URL            string    `json:"url"`
	URLId          string    `json:"urlId"`
	Revision       int       `json:"revision"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	PublishedAt    time.Time `json:"publishedAt"`
	ParentDocumentID string  `json:"parentDocumentId,omitempty"`
}

// GetDocumentRequest is the request payload for documents.info
type GetDocumentRequest struct {
	ID      string `json:"id,omitempty"`
	ShareID string `json:"shareId,omitempty"`
}

// GetDocumentResponse is the response from documents.info
type GetDocumentResponse struct {
	Data Document `json:"data"`
}

// GetDocument fetches a single document by ID
func (c *Client) GetDocument(id string) (*Document, error) {
	resp, err := c.post("documents.info", GetDocumentRequest{
		ID: id,
	})
	if err != nil {
		return nil, err
	}

	var docResp GetDocumentResponse
	if err := json.Unmarshal(resp.Data, &docResp); err != nil {
		return nil, err
	}

	return &docResp.Data, nil
}

// UpdateDocumentRequest is the request payload for documents.update
type UpdateDocumentRequest struct {
	ID       string `json:"id"`
	Text     string `json:"text"`
	Title    string `json:"title,omitempty"`
	Revision int    `json:"revision,omitempty"`
}

// UpdateDocumentResponse is the response from documents.update
type UpdateDocumentResponse struct {
	Data Document `json:"data"`
}

// UpdateDocument updates document content
func (c *Client) UpdateDocument(id, text string, revision int) (*Document, error) {
	resp, err := c.post("documents.update", UpdateDocumentRequest{
		ID:       id,
		Text:     text,
		Revision: revision,
	})
	if err != nil {
		return nil, err
	}

	var updateResp UpdateDocumentResponse
	if err := json.Unmarshal(resp.Data, &updateResp); err != nil {
		return nil, err
	}

	return &updateResp.Data, nil
}

// CreateDocumentRequest is the request payload for documents.create
type CreateDocumentRequest struct {
	Title        string `json:"title"`
	Text         string `json:"text"`
	CollectionID string `json:"collectionId"`
	Publish      bool   `json:"publish"`
}

// CreateDocumentResponse is the response from documents.create
type CreateDocumentResponse struct {
	Data Document `json:"data"`
}

// CreateDocument creates a new document
func (c *Client) CreateDocument(title, text, collectionID string) (*Document, error) {
	resp, err := c.post("documents.create", CreateDocumentRequest{
		Title:        title,
		Text:         text,
		CollectionID: collectionID,
		Publish:      true,
	})
	if err != nil {
		return nil, err
	}

	var createResp CreateDocumentResponse
	if err := json.Unmarshal(resp.Data, &createResp); err != nil {
		return nil, err
	}

	return &createResp.Data, nil
}
