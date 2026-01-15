package resolver

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var ErrNotFound = errors.New("resource not found")

type IdentifierType int

const (
	TypeID IdentifierType = iota
	TypeName
	TypeURL
)

type SearchResult struct {
	ID   string
	Name string
}

type AmbiguousError struct {
	Query   string
	Matches []SearchResult
}

func (e *AmbiguousError) Error() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("ambiguous name %q matches multiple resources:\n", e.Query))
	for _, m := range e.Matches {
		sb.WriteString(fmt.Sprintf("  - %s (%s)\n", m.Name, m.ID))
	}
	return sb.String()
}

type Searcher interface {
	SearchTasks(query string) ([]SearchResult, error)
	SearchLists(query string) ([]SearchResult, error)
	SearchFolders(query string) ([]SearchResult, error)
	SearchUsers(query string) ([]SearchResult, error)
}

type Resolver struct {
	searcher     Searcher
	strictResolve bool
}

func New(searcher Searcher, strictResolve bool) *Resolver {
	return &Resolver{
		searcher:     searcher,
		strictResolve: strictResolve,
	}
}

var (
	taskURLPattern   = regexp.MustCompile(`^https://app\.clickup\.com/t/(?:\d+/)?([a-zA-Z0-9]+)$`)
	listURLPattern   = regexp.MustCompile(`^https://app\.clickup\.com/\d+/v/li/(\d+)`)
	folderURLPattern = regexp.MustCompile(`^https://app\.clickup\.com/\d+/v/f/(\d+)/`)
	idPattern        = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	hasDigitPattern  = regexp.MustCompile(`\d`)
)

func DetectIdentifierType(input string) IdentifierType {
	if strings.HasPrefix(input, "https://") {
		return TypeURL
	}
	if idPattern.MatchString(input) && hasDigitPattern.MatchString(input) {
		return TypeID
	}
	return TypeName
}

func ParseTaskURL(url string) (string, error) {
	matches := taskURLPattern.FindStringSubmatch(url)
	if matches == nil {
		return "", fmt.Errorf("invalid task URL: %s", url)
	}
	return matches[1], nil
}

func ParseListURL(url string) (string, error) {
	matches := listURLPattern.FindStringSubmatch(url)
	if matches == nil {
		return "", fmt.Errorf("invalid list URL: %s", url)
	}
	return matches[1], nil
}

func ParseFolderURL(url string) (string, error) {
	matches := folderURLPattern.FindStringSubmatch(url)
	if matches == nil {
		return "", fmt.Errorf("invalid folder URL: %s", url)
	}
	return matches[1], nil
}

func (r *Resolver) ResolveTask(input string) (string, error) {
	switch DetectIdentifierType(input) {
	case TypeID:
		return input, nil
	case TypeURL:
		return ParseTaskURL(input)
	case TypeName:
		return r.resolveByName(input, r.searcher.SearchTasks)
	}
	return "", fmt.Errorf("unknown identifier type")
}

func (r *Resolver) ResolveList(input string) (string, error) {
	switch DetectIdentifierType(input) {
	case TypeID:
		return input, nil
	case TypeURL:
		return ParseListURL(input)
	case TypeName:
		return r.resolveByName(input, r.searcher.SearchLists)
	}
	return "", fmt.Errorf("unknown identifier type")
}

func (r *Resolver) ResolveFolder(input string) (string, error) {
	switch DetectIdentifierType(input) {
	case TypeID:
		return input, nil
	case TypeURL:
		return ParseFolderURL(input)
	case TypeName:
		return r.resolveByName(input, r.searcher.SearchFolders)
	}
	return "", fmt.Errorf("unknown identifier type")
}

func (r *Resolver) ResolveUser(input string) (string, error) {
	switch DetectIdentifierType(input) {
	case TypeID:
		return input, nil
	case TypeName:
		return r.resolveByName(input, r.searcher.SearchUsers)
	}
	return "", fmt.Errorf("unknown identifier type")
}

func (r *Resolver) resolveByName(query string, searchFn func(string) ([]SearchResult, error)) (string, error) {
	results, err := searchFn(query)
	if err != nil {
		return "", err
	}
	if len(results) == 0 {
		return "", ErrNotFound
	}
	if len(results) > 1 && r.strictResolve {
		return "", &AmbiguousError{Query: query, Matches: results}
	}
	return results[0].ID, nil
}
