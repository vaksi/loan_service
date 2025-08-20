package domain

import "time"

// Disbursement represents the final state of a loan where funds
// are handed over to the borrower. It captures a link to the
// signed agreement letter, the employee responsible for the
// disbursement and the date it occurred. A loan may only have one
// disbursement record.
type Disbursement struct {
    ID               string    `gorm:"type:uuid;primaryKey" json:"id"`
    LoanID           string    `gorm:"type:uuid;not null;unique" json:"loan_id"`
    AgreementURL     string    `gorm:"not null" json:"agreement_url"`
    EmployeeID       string    `gorm:"size:50;not null" json:"employee_id"`
    DisbursementDate time.Time `gorm:"not null" json:"disbursement_date"`
    CreatedAt        time.Time `json:"created_at"`
}