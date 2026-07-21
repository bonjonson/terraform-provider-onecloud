package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Domain struct {
	ID            int      `json:"ID"`
	Name          string   `json:"Name"`
	TechName      string   `json:"TechName"`
	State         string   `json:"State"`
	DateCreate    string   `json:"DateCreate"`
	IsDelegate    *bool    `json:"IsDelegate"`
	LinkedRecords []Record `json:"LinkedRecords"`
}

type createDomainRequest struct {
	DomainName string `json:"DomainName"`
	Migrate    bool   `json:"Migrate"`
}

func (c *Client) CreateDomain(name string, migrate bool) (*Domain, error) {
	req := createDomainRequest{DomainName: name, Migrate: migrate}
	resp, err := c.doRequest("POST", "/dns", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create domain, status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}
	var domain Domain
	if err := json.NewDecoder(resp.Body).Decode(&domain); err != nil {
		return nil, fmt.Errorf("decode domain: %w", err)
	}
	return &domain, nil
}

func (c *Client) GetDomain(id string) (*Domain, error) {
	resp, err := c.doRequest("GET", "/dns/"+id, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get domain, status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}
	var domain Domain
	if err := json.NewDecoder(resp.Body).Decode(&domain); err != nil {
		return nil, fmt.Errorf("decode domain: %w", err)
	}
	return &domain, nil
}

func (c *Client) DeleteDomain(id string) error {
	resp, err := c.doRequest("DELETE", "/dns/"+id, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete domain, status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}
	return nil
}
