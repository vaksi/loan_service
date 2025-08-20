package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"loan_service/internal/domain"

	"github.com/google/uuid"
)

// LoanRepo abstracts the persistence layer for easier testing (can be mocked in unit tests).
// The concrete implementation is repository.LoanRepository.
//
//go:generate mockery --name=LoanRepo --output=./mocks --outpkg=mocks --case=underscore
type LoanRepo interface {
	CreateLoan(ctx context.Context, loan *domain.Loan) error
	GetLoanByID(ctx context.Context, id string) (*domain.Loan, error)
	UpdateLoan(ctx context.Context, loan *domain.Loan) error
	CreateApproval(ctx context.Context, appr *domain.Approval) error
	CreateInvestment(ctx context.Context, inv *domain.Investment) error
	FindInvestorByEmail(ctx context.Context, email string) (*domain.Investor, error)
	CreateDisbursement(ctx context.Context, disb *domain.Disbursement) error
	ListLoans(ctx context.Context) ([]domain.Loan, error)
	GetTotalInvested(ctx context.Context, loanID string) (float64, error)
	GetInvestorByID(ctx context.Context, id string) (*domain.Investor, error)
	CreateInvestor(ctx context.Context, inv *domain.Investor) error
}

// LoanService orchestrates business logic for loans. It sits
// between handlers and repositories, enforcing state transitions and
// computing derived data such as the total invested amount. Errors
// returned from this service are suitable for consumption by HTTP
// handlers.
type LoanService struct{ repo LoanRepo }

// NewLoanService constructs a new LoanService using the given
// repository. Typically there is a single instance of the service
// created during application startup.
func NewLoanService(repo LoanRepo) *LoanService { return &LoanService{repo: repo} }

// Repo returns the underlying repository. It is exposed to allow
// handlers to perform read‑only operations not encapsulated by the
// service. Exposing the repository in this way keeps the service
// focused on business logic and avoids cluttering it with simple
// pass‑through methods.
func (s *LoanService) Repo() LoanRepo { return s.repo }

// CreateLoan creates a new loan with initial state `proposed`. It
// populates the ID with a new UUID. The loan is persisted via the
// repository and returned with default timestamps.
func (s *LoanService) CreateLoan(ctx context.Context, input domain.Loan) (*domain.Loan, error) {
	// Generate a new UUID for the loan.
	input.ID = uuid.New().String()
	input.State = domain.LoanStateProposed
	now := time.Now().UTC()
	input.CreatedAt = now
	input.UpdatedAt = now
	if err := s.repo.CreateLoan(ctx, &input); err != nil {
		return nil, err
	}
	return &input, nil
}

// ApproveLoan approves the loan with the given ID. It requires a
// picture proof URL, the employee ID of the validator and the
// approval date. The loan must currently be in the `proposed` state
// and must not already have an approval record. On success the loan
// state transitions to `approved` and the Approval record is
// persisted.
func (s *LoanService) ApproveLoan(ctx context.Context, loanID, pictureURL, employeeID string, approvalDate time.Time) (*domain.Loan, error) {
	loan, err := s.repo.GetLoanByID(ctx, loanID)
	if err != nil {
		return nil, err
	}
	// Validate current state
	if loan.State != domain.LoanStateProposed {
		return nil, fmt.Errorf("loan must be in proposed state to approve, current state: %s", loan.State)
	}
	// Check if already approved
	if loan.Approval != nil {
		return nil, errors.New("loan already approved")
	}
	// Create approval record
	approval := &domain.Approval{
		ID:           uuid.New().String(),
		LoanID:       loan.ID,
		PictureURL:   pictureURL,
		EmployeeID:   employeeID,
		ApprovalDate: approvalDate,
		CreatedAt:    time.Now().UTC(),
	}
	// Update loan state
	loan.State = domain.LoanStateApproved
	loan.UpdatedAt = time.Now().UTC()
	if err := s.repo.CreateApproval(ctx, approval); err != nil {
		return nil, err
	}
	if err := s.repo.UpdateLoan(ctx, loan); err != nil {
		return nil, err
	}
	// Reload loan with approval for return
	loan.Approval = approval
	return loan, nil
}

// InvestInLoan records a new investment in the specified loan. It
// accepts optional investor details. If an investor ID is provided it
// must exist; otherwise a new investor will be created using the
// provided name and email. The loan must be in the `approved` state
// and the total investment after this call must not exceed the
// principal amount. When the total invested equals the principal the
// loan state transitions to `invested`. A slice of investments is
// returned for convenience.
func (s *LoanService) InvestInLoan(ctx context.Context, loanID, investorID, investorName, investorEmail string, amount float64) (*domain.Loan, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}
	loan, err := s.repo.GetLoanByID(ctx, loanID)
	if err != nil {
		return nil, err
	}

	if loan.State == domain.LoanStateInvested {
		return nil, fmt.Errorf("loan already fully funded")
	}

	if loan.State != domain.LoanStateApproved {
		return nil, fmt.Errorf("loan must be approved to invest, current state: %s", loan.State)
	}

	// Retrieve or create investor
	var investor *domain.Investor
	if investorID != "" {
		investor, err = s.repo.GetInvestorByID(ctx, investorID)
		if err != nil {
			return nil, err
		}
	} else {
		// Try to find by email if provided
		if investorEmail != "" {
			existing, err := s.repo.FindInvestorByEmail(ctx, investorEmail)
			if err != nil {
				return nil, err
			}
			if existing != nil {
				investor = existing
			}
		}
		if investor == nil {
			investor = &domain.Investor{
				ID:        uuid.New().String(),
				Name:      investorName,
				Email:     investorEmail,
				CreatedAt: time.Now().UTC(),
			}
			if err := s.repo.CreateInvestor(ctx, investor); err != nil {
				return nil, err
			}
		}
	}
	// Check that investment will not exceed principal
	currentTotal, err := s.repo.GetTotalInvested(ctx, loan.ID)
	if err != nil {
		return nil, err
	}
	if currentTotal+amount > loan.Principal {
		return nil, fmt.Errorf("investment would exceed principal; current invested %.2f + new %.2f > principal %.2f", currentTotal, amount, loan.Principal)
	}
	// Create investment record
	invRec := &domain.Investment{
		ID:         uuid.New().String(),
		LoanID:     loan.ID,
		InvestorID: investor.ID,
		Amount:     amount,
		CreatedAt:  time.Now().UTC(),
	}
	if err := s.repo.CreateInvestment(ctx, invRec); err != nil {
		return nil, err
	}
	// Update state if fully funded
	newTotal := currentTotal + amount
	if newTotal == loan.Principal {
		loan.State = domain.LoanStateInvested
		loan.UpdatedAt = time.Now().UTC()
		if err := s.repo.UpdateLoan(ctx, loan); err != nil {
			return nil, err
		}
		// In a real system we would asynchronously send emails to
		// investors here. To preserve simplicity and avoid external
		// dependencies this implementation just logs the event.
		fmt.Printf("Loan %s fully funded. Total invested: %.2f. Sending agreement link to investors...\n", loan.ID, newTotal)
	}
	// Reload investments
	// Instead of reloading from database, append to loan's slice for return
	loan.Investments = append(loan.Investments, *invRec)
	return loan, nil
}

// DisburseLoan finalises the loan by marking it as disbursed. The
// caller must supply a URL pointing to the signed agreement letter,
// the employee responsible for the disbursement and the date. The
// loan must be in the `invested` state and must not already have a
// disbursement record. On success the state is set to `disbursed`.
func (s *LoanService) DisburseLoan(ctx context.Context, loanID, agreementURL, employeeID string, disbursementDate time.Time) (*domain.Loan, error) {
	loan, err := s.repo.GetLoanByID(ctx, loanID)
	if err != nil {
		return nil, err
	}
	if loan.State != domain.LoanStateInvested {
		return nil, fmt.Errorf("loan must be invested to disburse, current state: %s", loan.State)
	}
	if loan.Disbursement != nil {
		return nil, errors.New("loan already disbursed")
	}
	disb := &domain.Disbursement{
		ID:               uuid.New().String(),
		LoanID:           loan.ID,
		AgreementURL:     agreementURL,
		EmployeeID:       employeeID,
		DisbursementDate: disbursementDate,
		CreatedAt:        time.Now().UTC(),
	}
	loan.State = domain.LoanStateDisbursed
	loan.UpdatedAt = time.Now().UTC()
	if err := s.repo.CreateDisbursement(ctx, disb); err != nil {
		return nil, err
	}
	if err := s.repo.UpdateLoan(ctx, loan); err != nil {
		return nil, err
	}
	loan.Disbursement = disb
	return loan, nil
}

// ListLoans retrieves all loans from the repository. It returns
// loans with their nested Approval, Investments and Disbursement
// records. In a production system this method should support
// pagination and filtering.
func (s *LoanService) ListLoans(ctx context.Context) ([]domain.Loan, error) {
	loans, err := s.repo.ListLoans(ctx)
	if err != nil {
		return nil, err
	}
	return loans, nil
}

// GetLoanByID retrieves a single loan by its ID. It returns the loan
// with its nested Approval, Investments and Disbursement records.
func (s *LoanService) GetLoanByID(ctx context.Context, id string) (*domain.Loan, error) {
	loan, err := s.repo.GetLoanByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return loan, nil
}
