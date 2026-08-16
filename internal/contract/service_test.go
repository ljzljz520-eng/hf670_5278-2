package contract

import "testing"

func TestSubmitProducesPendingApplication(t *testing.T) {
	service := NewService(NewFixtureFactory())
	item, err := service.Submit(Submission{
		ContractName: "办公楼租赁合同",
		Department:   "采购部",
		SealType:     SealTypeContract,
		Urgency:      UrgencyNormal,
	})
	if err != nil {
		t.Fatalf("submit failed: %v", err)
	}
	if item.ID != "seal-001" || item.Status != "pending" {
		t.Fatalf("unexpected application: %#v", item)
	}
}

func TestSubmitRejectsIncompleteApplication(t *testing.T) {
	service := NewService(NewFixtureFactory())
	if _, err := service.Submit(Submission{Department: "采购部", SealType: SealTypeContract, Urgency: UrgencyNormal}); err == nil {
		t.Fatal("expected validation error")
	}
}
