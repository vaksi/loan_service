package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"loan_service/internal/domain"
	"loan_service/internal/handler"
	mock_loan_service "loan_service/internal/handler/mocks"
	"loan_service/internal/repository"
)

func TestCreateLoan_WithMockService(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ms := new(mock_loan_service.MockLoanService)
	created := &domain.Loan{ID: "L123", BorrowerID: "BRW", Principal: 1000, Rate: 0.1, ROI: 0.12, State: domain.LoanStateProposed}
	ms.On("CreateLoan", mock.Anything, mock.AnythingOfType("domain.Loan")).Return(created, nil).Once()

	h := handler.NewLoanHandler(ms)
	r := gin.Default()
	h.RegisterRoutes(r)

	body := map[string]any{"borrower_id": "BRW", "principal": 1000, "rate": 0.1, "roi": 0.12}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/loans", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	ms.AssertExpectations(t)
}
func TestInvestInLoan_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ms := new(mock_loan_service.MockLoanService)
	loanID := "L123"
	investorID := "INV1"
	investorName := "John Doe"
	investorEmail := "john@example.com"
	amount := 500.0
	expected := &domain.Loan{ID: loanID}

	ms.On("InvestInLoan", mock.Anything, loanID, investorID, investorName, investorEmail, amount).Return(expected, nil).Once()

	h := handler.NewLoanHandler(ms)
	r := gin.Default()
	h.RegisterRoutes(r)

	body := map[string]any{
		"investor_id":    investorID,
		"investor_name":  investorName,
		"investor_email": investorEmail,
		"amount":         amount,
	}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/loans/"+loanID+"/invest", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	ms.AssertExpectations(t)
}

func TestInvestInLoan_BadRequest_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ms := new(mock_loan_service.MockLoanService)
	h := handler.NewLoanHandler(ms)
	r := gin.Default()
	h.RegisterRoutes(r)

	req, _ := http.NewRequest("POST", "/loans/L123/invest", bytes.NewReader([]byte(`{invalid json`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestInvestInLoan_BadRequest_MissingAmount(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ms := new(mock_loan_service.MockLoanService)
	h := handler.NewLoanHandler(ms)
	r := gin.Default()
	h.RegisterRoutes(r)

	body := map[string]any{
		"investor_id":    "INV1",
		"investor_name":  "John Doe",
		"investor_email": "john@example.com",
		// "amount" is missing
	}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/loans/L123/invest", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestInvestInLoan_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ms := new(mock_loan_service.MockLoanService)
	loanID := "L123"
	ms.On("InvestInLoan", mock.Anything, loanID, "", "", "", 100.0).Return(nil, assert.AnError).Once()

	h := handler.NewLoanHandler(ms)
	r := gin.Default()
	h.RegisterRoutes(r)

	body := map[string]any{
		"amount": 100.0,
	}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/loans/"+loanID+"/invest", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	ms.AssertExpectations(t)
}
func TestListLoans_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ms := new(mock_loan_service.MockLoanService)
	expected := []domain.Loan{
		{ID: "L1", BorrowerID: "B1", Principal: 1000},
		{ID: "L2", BorrowerID: "B2", Principal: 2000},
	}
	ms.On("ListLoans", mock.Anything).Return(expected, nil).Once()

	h := handler.NewLoanHandler(ms)
	r := gin.Default()
	h.RegisterRoutes(r)

	req, _ := http.NewRequest("GET", "/loans", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp []domain.Loan
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, expected, resp)
	ms.AssertExpectations(t)
}

func TestListLoans_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ms := new(mock_loan_service.MockLoanService)
	ms.On("ListLoans", mock.Anything).Return(nil, assert.AnError).Once()

	h := handler.NewLoanHandler(ms)
	r := gin.Default()
	h.RegisterRoutes(r)

	req, _ := http.NewRequest("GET", "/loans", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	ms.AssertExpectations(t)
}

func TestGetLoan_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ms := new(mock_loan_service.MockLoanService)
	loanID := "L123"
	expected := &domain.Loan{ID: loanID, BorrowerID: "BRW"}
	ms.On("GetLoanByID", mock.Anything, loanID).Return(expected, nil).Once()

	h := handler.NewLoanHandler(ms)
	r := gin.Default()
	h.RegisterRoutes(r)

	req, _ := http.NewRequest("GET", "/loans/"+loanID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp domain.Loan
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, *expected, resp)
	ms.AssertExpectations(t)
}

func TestGetLoan_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ms := new(mock_loan_service.MockLoanService)
	loanID := "L404"
	ms.On("GetLoanByID", mock.Anything, loanID).Return(nil, repository.ErrNotFound).Once()

	h := handler.NewLoanHandler(ms)
	r := gin.Default()
	h.RegisterRoutes(r)

	req, _ := http.NewRequest("GET", "/loans/"+loanID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
	ms.AssertExpectations(t)
}

func TestGetLoan_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ms := new(mock_loan_service.MockLoanService)
	loanID := "L500"
	ms.On("GetLoanByID", mock.Anything, loanID).Return(nil, assert.AnError).Once()

	h := handler.NewLoanHandler(ms)
	r := gin.Default()
	h.RegisterRoutes(r)

	req, _ := http.NewRequest("GET", "/loans/"+loanID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	ms.AssertExpectations(t)
}

func TestApproveLoan_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ms := new(mock_loan_service.MockLoanService)
	loanID := "L123"
	pictureURL := "http://pic"
	employeeID := "EMP1"
	approvalDate := "2023-01-01T10:00:00Z"
	parsedDate, _ := time.Parse(time.RFC3339, approvalDate)
	expected := &domain.Loan{ID: loanID}

	ms.On("ApproveLoan", mock.Anything, loanID, pictureURL, employeeID, parsedDate).Return(expected, nil).Once()

	h := handler.NewLoanHandler(ms)
	r := gin.Default()
	h.RegisterRoutes(r)

	body := map[string]any{
		"picture_url":   pictureURL,
		"employee_id":   employeeID,
		"approval_date": approvalDate,
	}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/loans/"+loanID+"/approve", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	ms.AssertExpectations(t)
}

func TestApproveLoan_BadRequest_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ms := new(mock_loan_service.MockLoanService)
	h := handler.NewLoanHandler(ms)
	r := gin.Default()
	h.RegisterRoutes(r)

	req, _ := http.NewRequest("POST", "/loans/L123/approve", bytes.NewReader([]byte(`{invalid json`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestApproveLoan_BadRequest_InvalidDate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ms := new(mock_loan_service.MockLoanService)
	h := handler.NewLoanHandler(ms)
	r := gin.Default()
	h.RegisterRoutes(r)

	body := map[string]any{
		"picture_url":   "http://pic",
		"employee_id":   "EMP1",
		"approval_date": "not-a-date",
	}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/loans/L123/approve", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestApproveLoan_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ms := new(mock_loan_service.MockLoanService)
	loanID := "L123"
	pictureURL := "http://pic"
	employeeID := "EMP1"
	approvalDate := "2023-01-01T10:00:00Z"
	parsedDate, _ := time.Parse(time.RFC3339, approvalDate)

	ms.On("ApproveLoan", mock.Anything, loanID, pictureURL, employeeID, parsedDate).Return(nil, assert.AnError).Once()

	h := handler.NewLoanHandler(ms)
	r := gin.Default()
	h.RegisterRoutes(r)

	body := map[string]any{
		"picture_url":   pictureURL,
		"employee_id":   employeeID,
		"approval_date": approvalDate,
	}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/loans/"+loanID+"/approve", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	ms.AssertExpectations(t)
}

func TestDisburseLoan_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ms := new(mock_loan_service.MockLoanService)
	loanID := "L123"
	agreementURL := "http://agreement"
	employeeID := "EMP1"
	disbursementDate := "2023-01-01T10:00:00Z"
	parsedDate, _ := time.Parse(time.RFC3339, disbursementDate)
	expected := &domain.Loan{ID: loanID}

	ms.On("DisburseLoan", mock.Anything, loanID, agreementURL, employeeID, parsedDate).Return(expected, nil).Once()

	h := handler.NewLoanHandler(ms)
	r := gin.Default()
	h.RegisterRoutes(r)

	body := map[string]any{
		"agreement_url":     agreementURL,
		"employee_id":       employeeID,
		"disbursement_date": disbursementDate,
	}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/loans/"+loanID+"/disburse", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	ms.AssertExpectations(t)
}

func TestDisburseLoan_BadRequest_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ms := new(mock_loan_service.MockLoanService)
	h := handler.NewLoanHandler(ms)
	r := gin.Default()
	h.RegisterRoutes(r)

	req, _ := http.NewRequest("POST", "/loans/L123/disburse", bytes.NewReader([]byte(`{invalid json`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDisburseLoan_BadRequest_InvalidDate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ms := new(mock_loan_service.MockLoanService)
	h := handler.NewLoanHandler(ms)
	r := gin.Default()
	h.RegisterRoutes(r)

	body := map[string]any{
		"agreement_url":     "http://agreement",
		"employee_id":       "EMP1",
		"disbursement_date": "not-a-date",
	}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/loans/L123/disburse", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDisburseLoan_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ms := new(mock_loan_service.MockLoanService)
	loanID := "L123"
	agreementURL := "http://agreement"
	employeeID := "EMP1"
	disbursementDate := "2023-01-01T10:00:00Z"
	parsedDate, _ := time.Parse(time.RFC3339, disbursementDate)

	ms.On("DisburseLoan", mock.Anything, loanID, agreementURL, employeeID, parsedDate).Return(nil, assert.AnError).Once()

	h := handler.NewLoanHandler(ms)
	r := gin.Default()
	h.RegisterRoutes(r)

	body := map[string]any{
		"agreement_url":     agreementURL,
		"employee_id":       employeeID,
		"disbursement_date": disbursementDate,
	}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/loans/"+loanID+"/disburse", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	ms.AssertExpectations(t)
}
