package contract

import "fmt"

type SealType string

const (
	SealTypeOfficial SealType = "official"
	SealTypeFinance  SealType = "finance"
	SealTypeContract SealType = "contract"
)

type Urgency string

const (
	UrgencyNormal Urgency = "normal"
	UrgencyUrgent Urgency = "urgent"
)

type Submission struct {
	ContractName string   `json:"contractName"`
	Department   string   `json:"department"`
	SealType     SealType `json:"sealType"`
	Urgency      Urgency  `json:"urgency"`
}

type FileInput struct {
	Name string `json:"name"`
}

type CloseRecord struct {
	FileName   string `json:"fileName"`
	ResourceID string `json:"resourceId"`
}

type Application struct {
	ID           string        `json:"id"`
	ContractName string        `json:"contractName"`
	Department   string        `json:"department"`
	SealType     SealType      `json:"sealType"`
	Urgency      Urgency       `json:"urgency"`
	Status       string        `json:"status"`
	CloseRecords []CloseRecord `json:"closeRecords"`
}

func (s Submission) Validate() error {
	if s.ContractName == "" {
		return fmt.Errorf("contractName is required")
	}
	if s.Department == "" {
		return fmt.Errorf("department is required")
	}
	if s.SealType != SealTypeOfficial && s.SealType != SealTypeFinance && s.SealType != SealTypeContract {
		return fmt.Errorf("sealType is invalid")
	}
	if s.Urgency != UrgencyNormal && s.Urgency != UrgencyUrgent {
		return fmt.Errorf("urgency is invalid")
	}
	return nil
}
