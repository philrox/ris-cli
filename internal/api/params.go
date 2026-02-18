package api

import "net/url"

// Params wraps url.Values to provide a convenient builder for API query parameters.
type Params struct {
	values url.Values
}

// NewParams creates a new empty Params instance.
func NewParams() *Params {
	return &Params{values: url.Values{}}
}

// Set sets a single key-value pair. If the key already exists, it is replaced.
func (p *Params) Set(key, value string) {
	p.values.Set(key, value)
}

// Get returns the value for a key, or empty string if not set.
func (p *Params) Get(key string) string {
	return p.values.Get(key)
}

// Encode returns the URL-encoded query string.
func (p *Params) Encode() string {
	return p.values.Encode()
}

// Values returns the underlying url.Values.
func (p *Params) Values() url.Values {
	return p.values
}
