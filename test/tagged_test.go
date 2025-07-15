package test

import (
	"os"
	"strings"
	"testing"
)

// TaggedTest holds a test function and its tags
type TaggedTest struct {
	Name string
	Tags []string
	Func func(t *testing.T)
}

var taggedTests []TaggedTest

// RegisterTaggedTest registers a test with tags
func RegisterTaggedTest(name string, tags []string, fn func(t *testing.T)) {
	taggedTests = append(taggedTests, TaggedTest{Name: name, Tags: tags, Func: fn})
}

// RunTaggedTests runs only tests matching the specified tags
func RunTaggedTests(t *testing.T) {
	tagsEnv := os.Getenv("TAGS")
	if tagsEnv == "" {
		tagsEnv = os.Getenv("TEST_TAGS")
	}
	var filterTags []string
	if tagsEnv != "" {
		filterTags = strings.Split(tagsEnv, ",")
	}

	for _, test := range taggedTests {
		if len(filterTags) == 0 || hasAnyTag(test.Tags, filterTags) {
			t.Run(test.Name, test.Func)
		}
	}
}

func hasAnyTag(testTags, filterTags []string) bool {
	tagSet := make(map[string]struct{}, len(testTags))
	for _, tag := range testTags {
		tagSet[tag] = struct{}{}
	}
	for _, tag := range filterTags {
		if _, ok := tagSet[strings.TrimSpace(tag)]; ok {
			return true
		}
	}
	return false
}
