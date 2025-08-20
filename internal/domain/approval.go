package domain

import "time"

// Approval represents information captured when a loan is approved by
// field staff. It stores a link to a photographic proof that the
// borrower has been visited, the employee identifier of the field
// validator and the date of the approval. A loan may only have one
// approval record.
type Approval struct {
    ID            string    `gorm:"type:uuid;primaryKey" json:"id"`
    LoanID        string    `gorm:"type:uuid;not null;unique" json:"loan_id"`
    PictureURL    string    `gorm:"not null" json:"picture_url"`
    EmployeeID    string    `gorm:"size:50;not null" json:"employee_id"`
    ApprovalDate  time.Time `gorm:"not null" json:"approval_date"`
    CreatedAt     time.Time `json:"created_at"`
}