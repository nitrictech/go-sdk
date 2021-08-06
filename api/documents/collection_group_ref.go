package documents

import v1 "github.com/nitrictech/go-sdk/interfaces/nitric/v1"

// CollectionGroupRef - A reference to a chain of collections not tied to document keys for building sub collection queries
type CollectionGroupRef interface {
	// Query - Create a query for this collection group reference
	Query() Query
	// TODO: Add when deeper collection group references are supported
	// Collection() CollectionGroupRef

	// Parent - Get the parent collection of this collection group reference
	Parent() CollectionGroupRef

	// Name - Get the name of this collection group reference
	Name() string
}

type collectionGroupRefImpl struct {
	dc     v1.DocumentServiceClient
	parent *collectionGroupRefImpl
	name   string
}

func (c *collectionGroupRefImpl) Name() string {
	return c.name
}

func (c *collectionGroupRefImpl) Query() Query {
	return newQuery(c.toColRef(), c.dc)
}

func (c *collectionGroupRefImpl) Parent() CollectionGroupRef {
	return c.parent
}

// converts to a collection reference with nil document keys to chain
// collections together ready for query
func (c *collectionGroupRefImpl) toColRef() CollectionRef {
	if c.parent != nil {
		return &collectionRefImpl{
			name: c.name,
			dc:   c.dc,
			parentDocument: &documentRefImpl{
				dc:  c.dc,
				col: c.parent.toColRef(),
			},
		}
	}

	return &collectionRefImpl{
		name: c.name,
		dc:   c.dc,
	}
}

// Create a collection group reference from a collection
func fromColRef(col CollectionRef, dc v1.DocumentServiceClient) *collectionGroupRefImpl {
	if col.Parent() != nil {
		pDoc := col.Parent()
		return &collectionGroupRefImpl{
			dc:     dc,
			parent: fromColRef(pDoc.Parent(), dc),
			name:   col.Name(),
		}
	}

	return &collectionGroupRefImpl{
		dc:   dc,
		name: col.Name(),
	}
}
