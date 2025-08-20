package service

import (
	"context"
	"testing"
	"time"

	"loan_service/internal/domain"

	mock_loan_repo "loan_service/internal/service/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateLoan(t *testing.T) {
	repo := new(mock_loan_repo.MockLoanRepo)
	svc := NewLoanService(repo)
	input := domain.Loan{
		Principal: 1000,
	}
	repo.On("CreateLoan", mock.Anything, mock.AnythingOfType("*domain.Loan")).Return(nil)

	loan, err := svc.CreateLoan(context.Background(), input)
	assert.NoError(t, err)
	assert.NotEmpty(t, loan.ID)
	assert.Equal(t, domain.LoanStateProposed, loan.State)
	assert.WithinDuration(t, time.Now().UTC(), loan.CreatedAt, time.Second)
}

func TestApproveLoan_Success(t *testing.T) {
	repo := new(mock_loan_repo.MockLoanRepo)
	svc := NewLoanService(repo)
	loanID := uuid.New().String()
	loan := &domain.Loan{
		ID:    loanID,
		State: domain.LoanStateProposed,
	}
	repo.On("GetLoanByID", mock.Anything, loanID).Return(loan, nil)
	repo.On("CreateApproval", mock.Anything, mock.AnythingOfType("*domain.Approval")).Return(nil)
	repo.On("UpdateLoan", mock.Anything, loan).Return(nil)

	result, err := svc.ApproveLoan(context.Background(), loanID, "pic.jpg", "emp1", time.Now())
	assert.NoError(t, err)
	assert.Equal(t, domain.LoanStateApproved, result.State)
	assert.NotNil(t, result.Approval)
}

func TestApproveLoan_AlreadyApproved(t *testing.T) {
	repo := new(mock_loan_repo.MockLoanRepo)
	svc := NewLoanService(repo)
	loanID := uuid.New().String()
	loan := &domain.Loan{
		ID:       loanID,
		State:    domain.LoanStateProposed,
		Approval: &domain.Approval{ID: "appr1"},
	}
	repo.On("GetLoanByID", mock.Anything, loanID).Return(loan, nil)

	_, err := svc.ApproveLoan(context.Background(), loanID, "pic.jpg", "emp1", time.Now())
	assert.Error(t, err)
	assert.Equal(t, "loan already approved", err.Error())
}

func TestInvestInLoan_NewInvestor_Success(t *testing.T) {
	repo := new(mock_loan_repo.MockLoanRepo)
	svc := NewLoanService(repo)
	loanID := uuid.New().String()
	loan := &domain.Loan{
		ID:        loanID,
		State:     domain.LoanStateApproved,
		Principal: 1000,
	}
	repo.On("GetLoanByID", mock.Anything, loanID).Return(loan, nil)
	repo.On("FindInvestorByEmail", mock.Anything, "test@investor.com").Return(nil, nil)
	repo.On("CreateInvestor", mock.Anything, mock.AnythingOfType("*domain.Investor")).Return(nil)
	repo.On("GetTotalInvested", mock.Anything, loanID).Return(float64(0), nil)
	repo.On("CreateInvestment", mock.Anything, mock.AnythingOfType("*domain.Investment")).Return(nil)

	result, err := svc.InvestInLoan(context.Background(), loanID, "", "Test Investor", "test@investor.com", 500)
	assert.NoError(t, err)
	assert.Len(t, result.Investments, 1)
	assert.Equal(t, 500.0, result.Investments[0].Amount)
}

func TestDisburseLoan_Success(t *testing.T) {
	repo := new(mock_loan_repo.MockLoanRepo)
	svc := NewLoanService(repo)
	loanID := uuid.New().String()
	loan := &domain.Loan{
		ID:    loanID,
		State: domain.LoanStateInvested,
	}
	repo.On("GetLoanByID", mock.Anything, loanID).Return(loan, nil)
	repo.On("CreateDisbursement", mock.Anything, mock.AnythingOfType("*domain.Disbursement")).Return(nil)
	repo.On("UpdateLoan", mock.Anything, loan).Return(nil)

	result, err := svc.DisburseLoan(context.Background(), loanID, "agreement.pdf", "emp2", time.Now())
	assert.NoError(t, err)
	assert.Equal(t, domain.LoanStateDisbursed, result.State)
	assert.NotNil(t, result.Disbursement)
}

func TestDisburseLoan_AlreadyDisbursed(t *testing.T) {
	repo := new(mock_loan_repo.MockLoanRepo)
	svc := NewLoanService(repo)
	loanID := uuid.New().String()
	loan := &domain.Loan{
		ID:           loanID,
		State:        domain.LoanStateInvested,
		Disbursement: &domain.Disbursement{ID: "disb1"},
	}
	repo.On("GetLoanByID", mock.Anything, loanID).Return(loan, nil)

	_, err := svc.DisburseLoan(context.Background(), loanID, "agreement.pdf", "emp2", time.Now())
	assert.Error(t, err)
	assert.Equal(t, "loan already disbursed", err.Error())
}

func TestListLoans(t *testing.T) {
	repo := new(mock_loan_repo.MockLoanRepo)
	svc := NewLoanService(repo)
	loans := []domain.Loan{
		{ID: "loan1"},
		{ID: "loan2"},
	}
	repo.On("ListLoans", mock.Anything).Return(loans, nil)

	result, err := svc.ListLoans(context.Background())
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestGetLoanByID(t *testing.T) {
	repo := new(mock_loan_repo.MockLoanRepo)
	svc := NewLoanService(repo)
	loan := &domain.Loan{ID: "loan1"}
	repo.On("GetLoanByID", mock.Anything, "loan1").Return(loan, nil)

	result, err := svc.GetLoanByID(context.Background(), "loan1")
	assert.NoError(t, err)
	assert.Equal(t, "loan1", result.ID)
}

func TestInvestInLoan_InvalidAmount(t *testing.T) {
	repo := new(mock_loan_repo.MockLoanRepo)
	svc := NewLoanService(repo)
	_, err := svc.InvestInLoan(context.Background(), "loanid", "", "name", "email", 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "amount must be positive")
}

func TestApproveLoan_InvalidState(t *testing.T) {
	repo := new(mock_loan_repo.MockLoanRepo)
	svc := NewLoanService(repo)
	loanID := uuid.New().String()
	loan := &domain.Loan{
		ID:    loanID,
		State: domain.LoanStateApproved,
	}
	repo.On("GetLoanByID", mock.Anything, loanID).Return(loan, nil)

	_, err := svc.ApproveLoan(context.Background(), loanID, "pic.jpg", "emp1", time.Now())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "loan must be in proposed state to approve")
}

func TestDisburseLoan_InvalidState(t *testing.T) {
	repo := new(mock_loan_repo.MockLoanRepo)
	svc := NewLoanService(repo)
	loanID := uuid.New().String()
	loan := &domain.Loan{
		ID:    loanID,
		State: domain.LoanStateApproved,
	}
	repo.On("GetLoanByID", mock.Anything, loanID).Return(loan, nil)

	_, err := svc.DisburseLoan(context.Background(), loanID, "agreement.pdf", "emp2", time.Now())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "loan must be invested to disburse")
}
