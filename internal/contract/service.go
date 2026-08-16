package contract

import (
	"fmt"
	"sync"
)

type Service struct {
	mu      sync.Mutex
	nextID  int
	items   map[string]*Application
	factory *FixtureFactory
}

func NewService(factory *FixtureFactory) *Service {
	if factory == nil {
		factory = NewFixtureFactory()
	}
	return &Service{nextID: 1, items: make(map[string]*Application), factory: factory}
}

func (s *Service) Submit(input Submission) (Application, error) {
	if err := input.Validate(); err != nil {
		return Application{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	id := fmt.Sprintf("seal-%03d", s.nextID)
	s.nextID++
	item := &Application{
		ID:           id,
		ContractName: input.ContractName,
		Department:   input.Department,
		SealType:     input.SealType,
		Urgency:      input.Urgency,
		Status:       "pending",
		CloseRecords: []CloseRecord{},
	}
	s.items[id] = item
	return cloneApplication(item), nil
}

func (s *Service) Process(id string, files []FileInput) (Application, error) {
	if id == "" {
		return Application{}, fmt.Errorf("id is required")
	}
	if len(files) == 0 {
		return Application{}, fmt.Errorf("files are required")
	}
	for _, file := range files {
		if file.Name == "" {
			return Application{}, fmt.Errorf("file name is required")
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.items[id]
	if !ok {
		return Application{}, fmt.Errorf("application %s not found", id)
	}
	records := s.processFiles(files)

	item.Status = "processed"
	item.CloseRecords = records
	return cloneApplication(item), nil
}

func (s *Service) processFiles(files []FileInput) []CloseRecord {
	records := make([]CloseRecord, len(files))
	var current Resource
	for index, file := range files {
		current = s.factory.Open(file)
		defer func(holder *Resource, position int) {
			resource := *holder
			_ = resource.Close()
			records[position] = CloseRecord{FileName: resource.FileName(), ResourceID: resource.ID()}
		}(&current, index)
	}
	return records
}

func (s *Service) Get(id string) (Application, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.items[id]
	if !ok {
		return Application{}, false
	}
	return cloneApplication(item), true
}

func cloneApplication(item *Application) Application {
	result := *item
	result.CloseRecords = append([]CloseRecord{}, item.CloseRecords...)
	return result
}
