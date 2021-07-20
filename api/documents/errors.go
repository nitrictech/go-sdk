package documents

import "fmt"

// CollectionDepthExceededError - Error indicating that requested collection depth exceeds the maximum allowed
type CollectionDepthExceededError struct {
}

func (c *CollectionDepthExceededError) Error() string {
	return fmt.Sprintf("Maximum collection depth %d exceeded", MaxCollectionDepth)
}

func newCollectionDepthExceededError() error {
	return &CollectionDepthExceededError{}
}
