package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// DomainSpec represents the parsed domain specification from AI
type DomainSpec struct {
	DomainName          string              `json:"domain_name"`
	Entities            []EntitySpec        `json:"entities"`
	ValueObjects        []ValueObjectSpec   `json:"value_objects"`
	RepositoryInterface RepositorySpec      `json:"repository_interface"`
	ServiceInterface    ServiceSpec         `json:"service_interface"`
}

// EntitySpec represents an entity specification
type EntitySpec struct {
	Name            string       `json:"name"`
	IsAggregateRoot bool         `json:"is_aggregate_root"`
	Fields          []FieldSpec  `json:"fields"`
	Methods         []MethodSpec `json:"methods"`
}

// FieldSpec represents a field specification
type FieldSpec struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Validation  string `json:"validation"`
}

// MethodSpec represents a method specification
type MethodSpec struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	Signature      string `json:"signature"`
	Implementation string `json:"implementation"`
}

// ValueObjectSpec represents a value object specification
type ValueObjectSpec struct {
	Name       string      `json:"name"`
	Fields     []FieldSpec `json:"fields"`
	Validation string      `json:"validation"`
}

// RepositorySpec represents a repository interface specification
type RepositorySpec struct {
	Name    string              `json:"name"`
	Methods []InterfaceMethod   `json:"methods"`
}

// ServiceSpec represents a service interface specification
type ServiceSpec struct {
	Name    string            `json:"name"`
	Methods []InterfaceMethod `json:"methods"`
}

// InterfaceMethod represents a method in an interface
type InterfaceMethod struct {
	Name        string `json:"name"`
	Signature   string `json:"signature"`
	Description string `json:"description"`
}

// ParseDomainSpec parses the AI response into a DomainSpec
func ParseDomainSpec(content string) (*DomainSpec, error) {
	// Clean up the content - remove markdown code blocks if present
	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	// Parse JSON
	var spec DomainSpec
	if err := json.Unmarshal([]byte(content), &spec); err != nil {
		return nil, fmt.Errorf("parse JSON: %w\nContent:\n%s", err, content)
	}

	// Validate required fields
	if spec.DomainName == "" {
		return nil, fmt.Errorf("domain_name is required")
	}

	if len(spec.Entities) == 0 {
		return nil, fmt.Errorf("at least one entity is required")
	}

	return &spec, nil
}

// GenerateDomain generates domain code using AI
func GenerateDomain(ctx context.Context, orchestrator *Orchestrator, description string) (*DomainSpec, error) {
	// Create request
	req := &GenerateRequest{
		SystemPrompt: SystemPromptDDD,
		UserPrompt:   UserPromptTemplate(description),
		Temperature:  0.3, // Lower temperature for more consistent output
		MaxTokens:    8000, // Increased for complex domain specs
		TopP:         0.9,
	}

	// Generate
	resp, err := orchestrator.Generate(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("generate: %w", err)
	}

	// Parse response
	spec, err := ParseDomainSpec(resp.Content)
	if err != nil {
		return nil, fmt.Errorf("parse spec: %w", err)
	}

	return spec, nil
}
