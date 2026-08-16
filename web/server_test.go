package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"contractseal/internal/contract"
)

func TestBatchSealFlow(t *testing.T) {
	service := contract.NewService(contract.NewFixtureFactory())
	handler := NewHandler(service)

	createBody := []byte(`{"contractName":"年度服务合同","department":"市场部","sealType":"official","urgency":"urgent"}`)
	create := httptest.NewRequest(http.MethodPost, "/api/contract-seals", bytes.NewReader(createBody))
	create.Header.Set("Content-Type", "application/json")
	created := httptest.NewRecorder()
	handler.ServeHTTP(created, create)
	if created.Code != http.StatusCreated {
		t.Fatalf("create status: %d", created.Code)
	}
	var item contract.Application
	if err := json.NewDecoder(created.Body).Decode(&item); err != nil {
		t.Fatalf("decode create response: %v", err)
	}

	processBody := []byte(`{"files":[{"name":"contract.docx"},{"name":"appendix.pdf"},{"name":"seal-list.xlsx"}]}`)
	process := httptest.NewRequest(http.MethodPost, "/api/contract-seals/"+item.ID+"/process", bytes.NewReader(processBody))
	process.Header.Set("Content-Type", "application/json")
	processed := httptest.NewRecorder()
	handler.ServeHTTP(processed, process)
	if processed.Code != http.StatusOK {
		t.Fatalf("process status: %d", processed.Code)
	}
	var result contract.Application
	if err := json.NewDecoder(processed.Body).Decode(&result); err != nil {
		t.Fatalf("decode process response: %v", err)
	}
	want := []contract.CloseRecord{
		{FileName: "contract.docx", ResourceID: "resource-1"},
		{FileName: "appendix.pdf", ResourceID: "resource-2"},
		{FileName: "seal-list.xlsx", ResourceID: "resource-3"},
	}
	if len(result.CloseRecords) != len(want) {
		t.Fatalf("close count: got %d want %d", len(result.CloseRecords), len(want))
	}
	for i := range want {
		if result.CloseRecords[i] != want[i] {
			t.Fatalf("close record %d: got %#v want %#v", i, result.CloseRecords[i], want[i])
		}
	}
}
