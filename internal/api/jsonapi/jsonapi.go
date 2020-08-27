package jsonapi

import "encoding/json"

// RequestObject ...
type RequestObject struct {
	Data *ResourceObject `json:"data,omitempty" validate:"required"`
	Meta interface{}     `json:"meta,omitempty"`
}

// ResourceObject ...
type ResourceObject struct {
	ID            string      `json:"id,omitempty"`
	Type          string      `json:"type,omitempty"`
	Attributes    interface{} `json:"attributes,omitempty" validate:"required"`
	Meta          interface{} `json:"meta,omitempty"`
	Relationships interface{} `json:"relationships,omitempty"`
}

// ErrorObject ...
type ErrorObject struct {
	Status int    `json:"status,omitempty"`
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
}

// ResponseObject ...
type ResponseObject struct {
	Data     interface{}   `json:"data,omitempty"`
	Errors   []ErrorObject `json:"errors,omitempty"`
	Meta     interface{}   `json:"meta,omitempty"`
	Included interface{}   `json:"included,omitempty"`
}

// B ...
type B struct {
	*RequestObject
}

// UnmarshalJSON ...
func (i *B) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, i.RequestObject)
}
