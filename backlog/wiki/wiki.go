package wiki

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/nekrassov01/backlog-utils/backlog"
)

// Wiki represents a Backlog wiki client
type Wiki struct {
	*backlog.Backlog
}

// Page represents a wiki page
type Page struct {
	ID        int64  `json:"id"`
	ProjectID int64  `json:"projectId"`
	Name      string `json:"name"`
	Content   string `json:"content,omitempty"`
}

// New creates a new Backlog wiki client
func New(w io.Writer, url, apiKey string) (*Wiki, error) {
	if url == "" {
		return nil, errors.New("empty URL")
	}
	if apiKey == "" {
		return nil, errors.New("empty api key")
	}
	o := &Wiki{
		Backlog: &backlog.Backlog{
			Writer:           w,
			BaseURL:          url,
			APIKey:           apiKey,
			MaxRetryAttempts: 5,
			MaxJitterMilli:   1000,
		},
	}
	return o, nil
}

// List returns a list of wiki pages for the specified project key
func (o *Wiki) List(projectKey, pattern string) ([]*Page, error) {
	if projectKey == "" {
		return nil, errors.New("empty project key")
	}

	uri := fmt.Sprintf("%s/api/v2/wikis?projectIdOrKey=%s&apiKey=%s", o.BaseURL, projectKey, o.APIKey)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	resp, err := o.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg := backlog.GetErrorMessage(resp)
		return nil, fmt.Errorf("failed to list wikis: %d: %s", resp.StatusCode, msg)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var pages []*Page
	if err := json.Unmarshal(body, &pages); err != nil {
		return nil, err
	}

	if pattern != "" {
		matched := make([]*Page, 0, len(pages))
		for _, page := range pages {
			r, err := regexp.Compile(pattern)
			if err != nil {
				return nil, err
			}
			if r.MatchString(page.Name) {
				matched = append(matched, page)
			}
		}
		pages = matched
	}

	return pages, nil
}

// Get returns a wiki page
func (o *Wiki) Get(id int64) (*Page, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid wikiId: %d", id)
	}

	uri := fmt.Sprintf("%s/api/v2/wikis/%d?apiKey=%s", o.BaseURL, id, o.APIKey)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	resp, err := o.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg := backlog.GetErrorMessage(resp)
		return nil, fmt.Errorf("failed to get wiki page: %d: %s", resp.StatusCode, msg)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var page *Page
	if err := json.Unmarshal(body, &page); err != nil {
		return nil, err
	}

	return page, nil
}

// Rename renames a wiki page
func (o *Wiki) Rename(page *Page, before, after string) error {
	if page == nil {
		return errors.New("empty wiki page")
	}
	if before == "" {
		return errors.New("old strings must not be empty")
	}

	oldName := page.Name
	newName := strings.ReplaceAll(page.Name, before, after)

	values := url.Values{
		"name": {newName},
	}

	uri := fmt.Sprintf("%s/api/v2/wikis/%d?apiKey=%s", o.BaseURL, page.ID, o.APIKey)
	req, err := http.NewRequest(http.MethodPatch, uri, strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := o.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg := backlog.GetErrorMessage(resp)
		return fmt.Errorf("failed to update wiki page: %d: %s", resp.StatusCode, msg)
	}

	fmt.Fprintf(o.Writer, "updated: %s => %s\n", oldName, newName)
	return nil
}

// Replace replaces strings in the wiki page content
func (o *Wiki) Replace(page *Page, pairs ...string) error {
	if page == nil {
		return errors.New("empty wiki page")
	}
	if len(pairs) == 0 || len(pairs)%2 != 0 {
		return fmt.Errorf("number of old/new strings to replace does not match: %d", len(pairs))
	}

	replacer := strings.NewReplacer(pairs...)
	newContent := replacer.Replace(page.Content)

	values := url.Values{
		"content": {newContent},
	}

	uri := fmt.Sprintf("%s/api/v2/wikis/%d?apiKey=%s", o.BaseURL, page.ID, o.APIKey)
	req, err := http.NewRequest(http.MethodPatch, uri, strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := o.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg := backlog.GetErrorMessage(resp)
		return fmt.Errorf("failed to update wiki page content: %d: %s", resp.StatusCode, msg)
	}

	fmt.Fprintf(o.Writer, "updated: %d: %s\n", page.ID, page.Name)
	return nil
}
