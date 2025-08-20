package domain

import (
    "time"
)

// LoanState represents the valid states of a loan. Loans move
// forward through these states according to business rules and may
// never move backwards. See the service layer for state transition
// validation.
type LoanState string

const (
    // LoanStateProposed is the initial state when a loan is created.
    LoanStateProposed LoanState = "proposed"
    // LoanStateApproved indicates that the loan has been approved by
    // Amartha staff and is ready for investment.
    LoanStateApproved LoanState = "approved"
    // LoanStateInvested indicates that the loan has been fully funded
    // by investors.
    LoanStateInvested LoanState = "invested"
    // LoanStateDisbursed indicates that the loan principal has been
    // handed over to the borrower.
    LoanStateDisbursed LoanState = "disbursed"
)

// Loan represents a loan offered by Amartha. It contains basic
// information such as the borrower identifier, principal amount,
// interest rate, return on investment, a link to the generated
// agreement letter and the current state of the loan.
//
// The schema uses UUIDs as primary keys to ensure scalability when
// operating in distributed systems where auto‑incremented integers
// might collide. GORM handles UUID generation via default functions
// defined in migrations.
type Loan struct {
    ID                 string    `gorm:"type:uuid;primaryKey" json:"id"`
    BorrowerID         string    `gorm:"size:50;not null" json:"borrower_id"`
    Principal          float64   `gorm:"not null" json:"principal"`
    Rate               float64   `gorm:"not null" json:"rate"`
    ROI                float64   `gorm:"not null" json:"roi"`
    AgreementLetterURL string    `gorm:"column:agreement_letter_url" json:"agreement_letter_url"`
    State              LoanState `gorm:"size:20;not null" json:"state"`
    CreatedAt          time.Time `json:"created_at"`
    UpdatedAt          time.Time `json:"updated_at"`
    Approval           *Approval     `json:"approval,omitempty"`
    Investments        []Investment  `json:"investments,omitempty"`
    Disbursement       *Disbursement `json:"disbursement,omitempty"`
}