package repository

import (
	"context"
	"errors"

	"loan_service/internal/domain"

	"gorm.io/gorm"
)

// LoanRepository provides persistence methods for loans and their
// associated entities. It encapsulates all database interactions and
// exposes high level CRUD functions used by the service layer.
type LoanRepository struct {
	db *gorm.DB
}

// NewLoanRepository instantiates a new repository bound to the given
// GORM database handle. It is typically created once during
// application startup and injected into services.
func NewLoanRepository(db *gorm.DB) *LoanRepository {
	return &LoanRepository{db: db}
}

// CreateLoan inserts a new loan record into the database. The caller
// should set all required fields on the loan before invoking this
// method. The ID will be generated automatically via a database
// function in the migration.
func (r *LoanRepository) CreateLoan(ctx context.Context, loan *domain.Loan) error {
	return r.db.WithContext(ctx).Create(loan).Error
}

// GetLoanByID retrieves a loan by its ID. It preloads related
// Approval, Investments and Disbursement records. If the loan is not
// found a gorm.ErrRecordNotFound is returned.
func (r *LoanRepository) GetLoanByID(ctx context.Context, id string) (*domain.Loan, error) {
	var loan domain.Loan
	if err := r.db.WithContext(ctx).
		Preload("Approval").
		Preload("Investments").
		Preload("Disbursement").
		First(&loan, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &loan, nil
}

// UpdateLoan updates the given loan record in the database. GORM will
// generate an UPDATE statement based on the dirty fields. Use this
// method when modifying the state or other top level fields of the
// loan. It returns an error if the update fails.
func (r *LoanRepository) UpdateLoan(ctx context.Context, loan *domain.Loan) error {
	return r.db.WithContext(ctx).Save(loan).Error
}

// ListLoans returns all loans in the database. It preloads
// relationships to provide a complete snapshot of each loan. In a
// production system this method should support pagination.
func (r *LoanRepository) ListLoans(ctx context.Context) ([]domain.Loan, error) {
	var loans []domain.Loan
	if err := r.db.WithContext(ctx).
		Preload("Approval").Preload("Investments").Preload("Disbursement").
		Find(&loans).Error; err != nil {
		return nil, err
	}
	return loans, nil
}

// CreateApproval inserts a new approval record into the database.
// Enforces that each loan may only have one approval record by
// delegating uniqueness constraints to the database schema. If the
// insert fails due to a uniqueness violation, an error is returned.
func (r *LoanRepository) CreateApproval(ctx context.Context, approval *domain.Approval) error {
	return r.db.WithContext(ctx).Create(approval).Error
}

// CreateInvestment inserts a new investment record into the
// database. Multiple investments by the same investor for the same
// loan are allowed and aggregated at query time. It returns any
// resulting error.
func (r *LoanRepository) CreateInvestment(ctx context.Context, investment *domain.Investment) error {
	return r.db.WithContext(ctx).Create(investment).Error
}

// CreateDisbursement inserts a new disbursement record into the
// database. Each loan may have only one disbursement record, which
// should be enforced by the database schema. An error is returned if
// the insert fails.
func (r *LoanRepository) CreateDisbursement(ctx context.Context, d *domain.Disbursement) error {
	return r.db.WithContext(ctx).Create(d).Error
}

// GetTotalInvested returns the sum of all investments for the given
// loan ID. If no investments exist the returned total will be zero.
func (r *LoanRepository) GetTotalInvested(ctx context.Context, loanID string) (float64, error) {
	var total float64
	if err := r.db.WithContext(ctx).
		Model(&domain.Investment{}).
		Where("loan_id = ?", loanID).
		Select("COALESCE(SUM(amount),0)").
		Scan(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

// ErrNotFound wraps gorm.ErrRecordNotFound to decouple the service
// layer from the underlying ORM implementation. It can be used to
// differentiate between not found and other errors in handlers.
var ErrNotFound = gorm.ErrRecordNotFound

// CreateInvestor inserts a new investor record into the database. If
// the insert fails (for example due to a unique constraint
// violation) the error is returned. Investors can be created
// separately or on the fly when investing in a loan.
func (r *LoanRepository) CreateInvestor(ctx context.Context, inv *domain.Investor) error {
	return r.db.WithContext(ctx).Create(inv).Error
}

// GetInvestorByID fetches an investor by primary key. Returns
// ErrNotFound if the investor does not exist.
func (r *LoanRepository) GetInvestorByID(ctx context.Context, id string) (*domain.Investor, error) {
	var inv domain.Investor
	if err := r.db.WithContext(ctx).First(&inv, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &inv, nil
}

// FindInvestorByEmail returns the investor with the given email
// address if one exists. It returns nil and nil error if no investor
// matches the email. This can be used to look up an investor when
// performing an investment based on email rather than ID.
func (r *LoanRepository) FindInvestorByEmail(ctx context.Context, email string) (*domain.Investor, error) {
	var inv domain.Investor
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&inv).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &inv, nil
}
