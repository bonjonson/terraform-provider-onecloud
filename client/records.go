package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Record struct {
	ID           int    `json:"ID,omitempty"`
	DomainID     string `json:"-"`
	TypeRecord   string `json:"TypeRecord,omitempty"`
	HostName     string `json:"HostName,omitempty"`
	IP           string `json:"IP,omitempty"`
	Priority     string `json:"Priority,omitempty"` // строка!
	Text         string `json:"Text,omitempty"`
	MnemonicName string `json:"MnemonicName,omitempty"`
	ExtHostName  string `json:"ExtHostName,omitempty"`
	Service      string `json:"Service,omitempty"`
	Proto        string `json:"Proto,omitempty"`
	Weight       string `json:"Weight,omitempty"` // строка!
	Port         string `json:"Port,omitempty"`   // строка!
	Target       string `json:"Target,omitempty"`
	TTL          int    `json:"TTL,omitempty"`
	State        string `json:"State,omitempty"`
}

type recordRequest struct {
	DomainID     string `json:"DomainId"`
	Name         string `json:"Name,omitempty"`
	IP           string `json:"IP,omitempty"`
	HostName     string `json:"HostName,omitempty"`
	Priority     string `json:"Priority,omitempty"`
	Text         string `json:"Text,omitempty"`
	MnemonicName string `json:"MnemonicName,omitempty"`
	Service      string `json:"Service,omitempty"`
	Proto        string `json:"Proto,omitempty"`
	Weight       string `json:"Weight,omitempty"`
	Port         string `json:"Port,omitempty"`
	Target       string `json:"Target,omitempty"`
	TTL          string `json:"TTL,omitempty"`
}

func getRecordEndpoint(recordType string) string {
	switch strings.ToUpper(recordType) {
	case "A":
		return "recorda"
	case "AAAA":
		return "recordaaaa"
	case "CNAME":
		return "recordcname"
	case "MX":
		return "recordmx"
	case "NS":
		return "recordns"
	case "SRV":
		return "recordsrv"
	case "TXT":
		return "recordtxt"
	default:
		return "recorda"
	}
}

func (c *Client) CreateRecord(domainID string, rec *Record) (*Record, error) {
	endpoint := getRecordEndpoint(rec.TypeRecord)
	req := recordRequest{
		DomainID: domainID,
		TTL:      strconv.Itoa(rec.TTL),
	}
	switch strings.ToUpper(rec.TypeRecord) {
	case "A", "AAAA":
		req.Name = rec.HostName
		req.IP = rec.IP
	case "CNAME":
		req.Name = rec.HostName
		req.MnemonicName = rec.MnemonicName
	case "MX":
		req.HostName = rec.ExtHostName
		req.Priority = rec.Priority
	case "NS":
		req.HostName = rec.HostName
		req.Name = rec.ExtHostName
	case "SRV":
		req.Service = rec.Service
		req.Proto = rec.Proto
		req.Name = rec.HostName
		req.Priority = rec.Priority
		req.Weight = rec.Weight
		req.Port = rec.Port
		req.Target = rec.Target
	case "TXT":
		req.Name = rec.HostName
		req.Text = rec.Text
	}
	resp, err := c.doRequest("POST", "/dns/"+endpoint, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create record, status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}
	var created Record
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		return nil, fmt.Errorf("decode record response: %w", err)
	}
	created.DomainID = domainID
	return &created, nil
}

// GetRecordByDomainID получает запись через домен (надёжно)
func (c *Client) GetRecordByDomainID(domainID, recordID string) (*Record, error) {
	domain, err := c.GetDomain(domainID)
	if err != nil {
		return nil, err
	}
	if domain == nil {
		return nil, nil
	}
	for _, r := range domain.LinkedRecords {
		if strconv.Itoa(r.ID) == recordID {
			r.DomainID = domainID
			return &r, nil
		}
	}
	return nil, nil
}

// GetRecord получает запись через универсальный эндпоинт (может не работать, не используется для ресурса)
func (c *Client) GetRecord(recordID string) (*Record, error) {
	resp, err := c.doRequest("GET", "/dns/record/"+recordID, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get record, status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}
	var rec Record
	if err := json.NewDecoder(resp.Body).Decode(&rec); err != nil {
		return nil, fmt.Errorf("decode record: %w", err)
	}
	return &rec, nil
}

func (c *Client) UpdateRecord(domainID, recordID string, rec *Record) (*Record, error) {
	endpoint := getRecordEndpoint(rec.TypeRecord)
	req := recordRequest{
		DomainID: domainID,
		TTL:      strconv.Itoa(rec.TTL),
	}
	switch strings.ToUpper(rec.TypeRecord) {
	case "A", "AAAA":
		req.Name = rec.HostName
		req.IP = rec.IP
	case "CNAME":
		req.Name = rec.HostName
		req.MnemonicName = rec.MnemonicName
	case "MX":
		req.HostName = rec.ExtHostName
		req.Priority = rec.Priority
	case "NS":
		req.HostName = rec.HostName
		req.Name = rec.ExtHostName
	case "SRV":
		req.Service = rec.Service
		req.Proto = rec.Proto
		req.Name = rec.HostName
		req.Priority = rec.Priority
		req.Weight = rec.Weight
		req.Port = rec.Port
		req.Target = rec.Target
	case "TXT":
		req.Name = rec.HostName
		req.Text = rec.Text
	}
	resp, err := c.doRequest("PUT", "/dns/"+endpoint+"/"+recordID, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to update record, status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}
	var updated Record
	if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
		return nil, fmt.Errorf("decode record response: %w", err)
	}
	updated.DomainID = domainID
	return &updated, nil
}

func (c *Client) DeleteRecord(domainID, recordID string) error {
	resp, err := c.doRequest("DELETE", "/dns/"+domainID+"/"+recordID, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete record, status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}
	return nil
}
