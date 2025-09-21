package gosolar

import (
	"context"
	"fmt"
)

// BulkSetCustomProperty sets a custom property on multiple entities
func (c *Client) BulkSetCustomProperty(uris []string, name string, value interface{}) error {
	return c.BulkSetCustomPropertyContext(context.Background(), uris, name, value)
}

// BulkSetCustomPropertyContext sets a custom property on multiple entities with context
func (c *Client) BulkSetCustomPropertyContext(ctx context.Context, uris []string, name string, value interface{}) error {
	if len(uris) == 0 {
		return NewError(ErrorTypeValidation, "bulk_set_custom_property", "no URIs provided")
	}
	if name == "" {
		return NewError(ErrorTypeValidation, "bulk_set_custom_property", "property name cannot be empty")
	}

	// Prepare URIs for custom properties endpoint
	cpuris := make([]string, 0, len(uris))
	for _, uri := range uris {
		cpuris = append(cpuris, uri+"/CustomProperties")
	}

	bulkRequest := struct {
		URIs       []string               `json:"uris"`
		Properties map[string]interface{} `json:"properties"`
	}{
		URIs: cpuris,
		Properties: map[string]interface{}{
			name: value,
		},
	}

	_, err := c.PostContext(ctx, "BulkUpdate", &bulkRequest)
	if err != nil {
		return WrapError(err, ErrorTypeInternal, "bulk_set_custom_property", "failed to update custom properties")
	}

	return nil
}

// SetCustomProperty sets a single custom property on an entity
func (c *Client) SetCustomProperty(uri, name string, value interface{}) error {
	return c.SetCustomPropertyContext(context.Background(), uri, name, value)
}

// SetCustomPropertyContext sets a single custom property with context
func (c *Client) SetCustomPropertyContext(ctx context.Context, uri, name string, value interface{}) error {
	if uri == "" {
		return NewError(ErrorTypeValidation, "set_custom_property", "URI cannot be empty")
	}
	if name == "" {
		return NewError(ErrorTypeValidation, "set_custom_property", "property name cannot be empty")
	}

	property := map[string]interface{}{
		name: value,
	}

	_, err := c.PostContext(ctx, uri+"/CustomProperties", &property)
	if err != nil {
		return WrapError(err, ErrorTypeInternal, "set_custom_property", "failed to update custom property")
	}

	return nil
}

// SetCustomProperties sets multiple custom properties on an entity
func (c *Client) SetCustomProperties(uri string, properties map[string]interface{}) error {
	return c.SetCustomPropertiesContext(context.Background(), uri, properties)
}

// SetCustomPropertiesContext sets multiple custom properties with context
func (c *Client) SetCustomPropertiesContext(ctx context.Context, uri string, properties map[string]interface{}) error {
	if uri == "" {
		return NewError(ErrorTypeValidation, "set_custom_properties", "URI cannot be empty")
	}
	if len(properties) == 0 {
		return NewError(ErrorTypeValidation, "set_custom_properties", "no properties provided")
	}

	_, err := c.PostContext(ctx, uri+"/CustomProperties", &properties)
	if err != nil {
		return WrapError(err, ErrorTypeInternal, "set_custom_properties", "failed to update custom properties")
	}

	return nil
}

// CustomPropertyType represents the type of a custom property
type CustomPropertyType string

const (
	CustomPropertyTypeString   CustomPropertyType = "string"
	CustomPropertyTypeInteger  CustomPropertyType = "integer"
	CustomPropertyTypeFloat    CustomPropertyType = "float"
	CustomPropertyTypeBoolean  CustomPropertyType = "boolean"
	CustomPropertyTypeDateTime CustomPropertyType = "datetime"
)

// CreateCustomPropertyRequest represents the parameters for creating a custom property
type CreateCustomPropertyRequest struct {
	Entity      string             `json:"entity"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Type        CustomPropertyType `json:"type"`
	Length      int                `json:"length,omitempty"`
}

// CreateCustomProperty creates a new custom property definition
func (c *Client) CreateCustomProperty(req CreateCustomPropertyRequest) error {
	return c.CreateCustomPropertyContext(context.Background(), req)
}

// CreateCustomPropertyContext creates a new custom property with context
func (c *Client) CreateCustomPropertyContext(ctx context.Context, req CreateCustomPropertyRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	// Determine length based on type
	length := req.Length
	if req.Type == CustomPropertyTypeString && length == 0 {
		length = 400 // Default string length
	}

	// Build the properties array as expected by SolarWinds
	props := []interface{}{
		req.Name,
		req.Description,
		string(req.Type),
		fmt.Sprintf("%d", length),
		"", // default value
		"", // validation expression
		"", // units
		"", // format
		"", // tooltip
		"", // category
		"", // subcategory
		"false", // mandatory
		"", // display name
	}

	endpoint := fmt.Sprintf("Invoke/%s/CreateCustomProperty", req.Entity)

	_, err := c.PostContext(ctx, endpoint, &props)
	if err != nil {
		return WrapError(err, ErrorTypeInternal, "create_custom_property", "failed to create custom property")
	}

	return nil
}

// Validate checks if the CreateCustomPropertyRequest is valid
func (req *CreateCustomPropertyRequest) Validate() error {
	if req.Entity == "" {
		return NewError(ErrorTypeValidation, "create_custom_property", "entity cannot be empty")
	}
	if req.Name == "" {
		return NewError(ErrorTypeValidation, "create_custom_property", "name cannot be empty")
	}
	if req.Type == "" {
		return NewError(ErrorTypeValidation, "create_custom_property", "type cannot be empty")
	}

	// Validate type
	validTypes := []CustomPropertyType{
		CustomPropertyTypeString,
		CustomPropertyTypeInteger,
		CustomPropertyTypeFloat,
		CustomPropertyTypeBoolean,
		CustomPropertyTypeDateTime,
	}

	validType := false
	for _, validT := range validTypes {
		if req.Type == validT {
			validType = true
			break
		}
	}

	if !validType {
		return NewError(ErrorTypeValidation, "create_custom_property", fmt.Sprintf("invalid type: %s", req.Type))
	}

	return nil
}
