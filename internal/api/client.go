package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// ClientOptions configures the API client.
type ClientOptions struct {
	Timeout time.Duration
	Verbose bool
	BaseURL string // Override API base URL (defaults to DefaultBaseURL)
}

// Client is the HTTP client for the RIS API.
type Client struct {
	baseURL    string
	httpClient *http.Client
	verbose    bool
}

// NewClient creates a new API client.
func NewClient(opts ClientOptions) *Client {
	baseURL := opts.BaseURL
	if baseURL == "" {
		baseURL = os.Getenv("RIS_BASE_URL")
	}
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &Client{
		baseURL: strings.TrimRight(baseURL, "/") + "/",
		httpClient: &http.Client{
			Timeout: timeout,
		},
		verbose: opts.Verbose,
	}
}

// Search performs a search query against the given API endpoint.
// Returns the raw JSON response body.
func (c *Client) Search(endpoint string, params *Params) ([]byte, error) {
	reqURL := c.baseURL + endpoint
	if params != nil && params.Encode() != "" {
		reqURL += "?" + params.Encode()
	}

	if c.verbose {
		fmt.Fprintf(os.Stderr, "GET %s\n", reqURL)
	}

	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		if isTimeout(err) {
			return nil, &TimeoutError{URL: reqURL, Err: err}
		}
		return nil, fmt.Errorf("HTTP-Anfrage fehlgeschlagen: %w", err)
	}
	defer resp.Body.Close()

	if c.verbose {
		fmt.Fprintf(os.Stderr, "HTTP %d %s\n", resp.StatusCode, resp.Status)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &HTTPError{StatusCode: resp.StatusCode, Status: resp.Status, URL: reqURL}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Antwort konnte nicht gelesen werden: %w", err)
	}

	return body, nil
}

// FetchDocument retrieves HTML content from a document URL.
// Validates the URL for SSRF protection before fetching.
func (c *Client) FetchDocument(docURL string) (string, error) {
	if err := validateDocURL(docURL); err != nil {
		return "", err
	}

	if c.verbose {
		fmt.Fprintf(os.Stderr, "GET %s\n", docURL)
	}

	resp, err := c.httpClient.Get(docURL)
	if err != nil {
		if isTimeout(err) {
			return "", &TimeoutError{URL: docURL, Err: err}
		}
		return "", fmt.Errorf("HTTP-Anfrage fehlgeschlagen: %w", err)
	}
	defer resp.Body.Close()

	if c.verbose {
		fmt.Fprintf(os.Stderr, "HTTP %d %s\n", resp.StatusCode, resp.Status)
	}

	if resp.StatusCode != http.StatusOK {
		return "", &HTTPError{StatusCode: resp.StatusCode, Status: resp.Status, URL: docURL}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Dokument konnte nicht gelesen werden: %w", err)
	}

	return string(body), nil
}

// validateDocURL checks that the URL is HTTPS and points to an allowed host.
func validateDocURL(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("Ungültige URL: %w", err)
	}
	if u.Scheme != "https" {
		return fmt.Errorf("Nur HTTPS-URLs erlaubt, erhalten: %q", u.Scheme)
	}
	host := strings.ToLower(u.Hostname())
	if !AllowedHosts[host] {
		return fmt.Errorf("Host %q nicht erlaubt", host)
	}
	return nil
}

// isTimeout checks if an error is a timeout error.
func isTimeout(err error) bool {
	type timeouter interface {
		Timeout() bool
	}
	if t, ok := err.(timeouter); ok {
		return t.Timeout()
	}
	return false
}

// TimeoutError indicates an HTTP request timed out.
type TimeoutError struct {
	URL string
	Err error
}

func (e *TimeoutError) Error() string {
	return fmt.Sprintf("Zeitüberschreitung: %s", e.URL)
}

func (e *TimeoutError) Unwrap() error {
	return e.Err
}

// HTTPError indicates an HTTP response with a non-200 status code.
type HTTPError struct {
	StatusCode int
	Status     string
	URL        string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s (%s)", e.StatusCode, e.Status, e.URL)
}
