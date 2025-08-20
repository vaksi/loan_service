package mocks

import (
	"context"
	"loan_service/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockLoanRepo struct {
	mock.Mock
}

func (m *MockLoanRepo) CreateLoan(ctx context.Context, loan *domain.Loan) error {
	args := m.Called(ctx, loan)
	return args.Error(0)
}

func (m *MockLoanRepo) GetLoanByID(ctx context.Context, id string) (*domain.Loan, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Loan), args.Error(1)
}

func (m *MockLoanRepo) UpdateLoan(ctx context.Context, loan *domain.Loan) error {
	args := m.Called(ctx, loan)
	return args.Error(0)
}

func (m *MockLoanRepo) CreateApproval(ctx context.Context, appr *domain.Approval) error {
	args := m.Called(ctx, appr)
	return args.Error(0)
}

func (m *MockLoanRepo) CreateInvestment(ctx context.Context, inv *domain.Investment) error {
	args := m.Called(ctx, inv)
	return args.Error(0)
}

func (m *MockLoanRepo) FindInvestorByEmail(ctx context.Context, email string) (*domain.Investor, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Investor), args.Error(1)
}

func (m *MockLoanRepo) CreateDisbursement(ctx context.Context, disb *domain.Disbursement) error {
	args := m.Called(ctx, disb)
	return args.Error(0)
}

func (m *MockLoanRepo) ListLoans(ctx context.Context) ([]domain.Loan, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Loan), args.Error(1)
}

func (m *MockLoanRepo) GetTotalInvested(ctx context.Context, loanID string) (float64, error) {
	args := m.Called(ctx, loanID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockLoanRepo) GetInvestorByID(ctx context.Context, id string) (*domain.Investor, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Investor), args.Error(1)
}

func (m *MockLoanRepo) CreateInvestor(ctx context.Context, inv *domain.Investor) error {
	args := m.Called(ctx, inv)
	return args.Error(0)
}
