package api

import (
	"encoding/json"
	"time"
)

// Collection represents an Outline collection
type Collection struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
	Icon        string    `json:"icon"`
	URL         string    `json:"url"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// ListCollectionsRequest is the request payload for collections.list
type ListCollectionsRequest struct {
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

// ListCollectionsResponse is the response from collections.list
type ListCollectionsResponse struct {
	Data       []Collection `json:"data"`
	Pagination Pagination   `json:"pagination"`
}

// Pagination contains pagination info
type Pagination struct {
	Offset   int    `json:"offset"`
	Limit    int    `json:"limit"`
	NextPath string `json:"nextPath,omitempty"`
}

// ListCollections fetches all accessible collections
func (c *Client) ListCollections() ([]Collection, error) {
	resp, err := c.post("collections.list", ListCollectionsRequest{
		Limit: 100,
	})
	if err != nil {
		return nil, err
	}

	var collections []Collection
	if err := json.Unmarshal(resp.Data, &collections); err != nil {
		return nil, err
	}

	return collections, nil
}

// GetCollectionRequest is the request payload for collections.info
type GetCollectionRequest struct {
	ID string `json:"id"`
}

// GetCollectionResponse is the response from collections.info
type GetCollectionResponse struct {
	Data Collection `json:"data"`
}

// GetCollection fetches a single collection by ID
func (c *Client) GetCollection(id string) (*Collection, error) {
	resp, err := c.post("collections.info", GetCollectionRequest{
		ID: id,
	})
	if err != nil {
		return nil, err
	}

	var coll Collection
	if err := json.Unmarshal(resp.Data, &coll); err != nil {
		return nil, err
	}

	return &coll, nil
}

// DocumentNode represents a document in the collection tree
type DocumentNode struct {
	ID       string         `json:"id"`
	Title    string         `json:"title"`
	URL      string         `json:"url"`
	Children []DocumentNode `json:"children"`
}

// GetCollectionDocumentsRequest is the request payload for collections.documents
type GetCollectionDocumentsRequest struct {
	ID string `json:"id"`
}

// GetCollectionDocumentsResponse is the response from collections.documents
type GetCollectionDocumentsResponse struct {
	Data []DocumentNode `json:"data"`
}

// GetCollectionDocuments fetches the document tree for a collection
func (c *Client) GetCollectionDocuments(collectionID string) ([]DocumentNode, error) {
	resp, err := c.post("collections.documents", GetCollectionDocumentsRequest{
		ID: collectionID,
	})
	if err != nil {
		return nil, err
	}

	var nodes []DocumentNode
	if err := json.Unmarshal(resp.Data, &nodes); err != nil {
		return nil, err
	}

	return nodes, nil
}
