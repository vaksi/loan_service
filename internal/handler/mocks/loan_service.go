package mocks

import (
	"context"
	"loan_service/internal/domain"
	"time"

	"github.com/stretchr/testify/mock"
)

// --- Mock service implementing handler.LoanUsecase ---
type MockLoanService struct{ mock.Mock }

func (m *MockLoanService) CreateLoan(ctx context.Context, input domain.Loan) (*domain.Loan, error) {
	args := m.Called(ctx, input)
	loan, _ := args.Get(0).(*domain.Loan)
	return loan, args.Error(1)
}
func (m *MockLoanService) ApproveLoan(ctx context.Context, loanID, pictureURL, employeeID string, approvalDate time.Time) (*domain.Loan, error) {
	args := m.Called(ctx, loanID, pictureURL, employeeID, approvalDate)
	loan, _ := args.Get(0).(*domain.Loan)
	return loan, args.Error(1)
}
func (m *MockLoanService) InvestInLoan(ctx context.Context, loanID, investorID, investorName, investorEmail string, amount float64) (*domain.Loan, error) {
	args := m.Called(ctx, loanID, investorID, investorName, investorEmail, amount)
	loan, _ := args.Get(0).(*domain.Loan)
	return loan, args.Error(1)
}
func (m *MockLoanService) DisburseLoan(ctx context.Context, loanID, agreementURL, employeeID string, disbursementDate time.Time) (*domain.Loan, error) {
	args := m.Called(ctx, loanID, agreementURL, employeeID, disbursementDate)
	loan, _ := args.Get(0).(*domain.Loan)
	return loan, args.Error(1)
}
func (m *MockLoanService) GetLoanByID(ctx context.Context, id string) (*domain.Loan, error) {
	args := m.Called(ctx, id)
	loan, _ := args.Get(0).(*domain.Loan)
	return loan, args.Error(1)
}
func (m *MockLoanService) ListLoans(ctx context.Context) ([]domain.Loan, error) {
	args := m.Called(ctx)
	loans, _ := args.Get(0).([]domain.Loan)
	return loans, args.Error(1)
}
