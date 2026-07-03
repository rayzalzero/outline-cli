package api

import (
	"encoding/json"
	"log"
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

	var doc Document
	if err := json.Unmarshal(resp.Data, &doc); err != nil {
		return nil, err
	}

	return &doc, nil
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
func (c *Client) UpdateDocument(id, text, title string, revision int) (*Document, error) {
	req := UpdateDocumentRequest{
		ID:       id,
		Text:     text,
		Revision: revision,
	}
	if title != "" {
		req.Title = title
	}
	
	resp, err := c.post("documents.update", req)
	if err != nil {
		return nil, err
	}

	var doc Document
	if err := json.Unmarshal(resp.Data, &doc); err != nil {
		return nil, err
	}

	return &doc, nil
}

// CreateDocumentRequest is the request payload for documents.create
type CreateDocumentRequest struct {
	Title            string `json:"title"`
	Text             string `json:"text,omitempty"`
	CollectionID     string `json:"collectionId,omitempty"`
	ParentDocumentID string `json:"parentDocumentId,omitempty"`
	Publish          bool   `json:"publish"`
}

// CreateDocumentResponse is the response from documents.create
type CreateDocumentResponse struct {
	Data Document `json:"data"`
}

// CreateDocument creates a new document
func (c *Client) CreateDocument(title, text, collectionID string) (*Document, error) {
	return c.CreateDocumentWithParent(title, text, collectionID, "")
}

// CreateDocumentWithParent creates a new document with optional parent
func (c *Client) CreateDocumentWithParent(title, text, collectionID, parentDocumentID string) (*Document, error) {
	req := CreateDocumentRequest{
		Title:   title,
		Text:    text,
		Publish: true,
	}
	
	if parentDocumentID != "" {
		req.ParentDocumentID = parentDocumentID
	} else if collectionID != "" {
		req.CollectionID = collectionID
	}
	
	log.Printf("[DEBUG] CreateDocument: title=%s, textLen=%d, collectionID=%s, parentID=%s", 
		title, len(text), collectionID, parentDocumentID)
	
	resp, err := c.post("documents.create", req)
	if err != nil {
		return nil, err
	}

	var doc Document
	if err := json.Unmarshal(resp.Data, &doc); err != nil {
		return nil, err
	}
	
	log.Printf("[DEBUG] CreateDocument response: id=%s, title=%s, textLen=%d", 
		doc.ID, doc.Title, len(doc.Text))

	return &doc, nil
}

// MoveDocumentRequest is the request payload for documents.move
type MoveDocumentRequest struct {
	ID               string `json:"id"`
	CollectionID     string `json:"collectionId,omitempty"`
	ParentDocumentID string `json:"parentDocumentId,omitempty"`
	Index            int    `json:"index,omitempty"`
}

// MoveDocument moves a document to a new parent (or root if parentID is empty)
func (c *Client) MoveDocument(id, parentID string) (*Document, error) {
	return c.MoveDocumentWithIndex(id, parentID, "", 0)
}

func (c *Client) MoveDocumentWithCollection(id, parentID, collectionID string) (*Document, error) {
	return c.MoveDocumentWithIndex(id, parentID, collectionID, 0)
}

func (c *Client) MoveDocumentWithIndex(id, parentID, collectionID string, index int) (*Document, error) {
	req := MoveDocumentRequest{
		ID:               id,
		ParentDocumentID: parentID,
	}
	if collectionID != "" {
		req.CollectionID = collectionID
	}
	if index > 0 {
		req.Index = index
	}
	
	resp, err := c.post("documents.move", req)
	if err != nil {
		return nil, err
	}

	var doc Document
	if err := json.Unmarshal(resp.Data, &doc); err != nil {
		return nil, err
	}

	return &doc, nil
}
