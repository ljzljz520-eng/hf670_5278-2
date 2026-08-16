package contract

import "fmt"

type Resource interface {
	Close() error
	FileName() string
	ID() string
}

type FixtureResource struct {
	name       string
	id         string
	closeCount int
}

func (r *FixtureResource) Close() error {
	r.closeCount++
	return nil
}

func (r *FixtureResource) FileName() string {
	return r.name
}

func (r *FixtureResource) ID() string {
	return r.id
}

func (r *FixtureResource) CloseCount() int {
	return r.closeCount
}

type FixtureFactory struct {
	opened []*FixtureResource
}

func NewFixtureFactory() *FixtureFactory {
	return &FixtureFactory{}
}

func (f *FixtureFactory) Open(file FileInput) Resource {
	id := fmt.Sprintf("resource-%d", len(f.opened)+1)
	resource := &FixtureResource{name: file.Name, id: id}
	f.opened = append(f.opened, resource)
	return resource
}

func (f *FixtureFactory) Opened() []*FixtureResource {
	resources := make([]*FixtureResource, len(f.opened))
	copy(resources, f.opened)
	return resources
}
