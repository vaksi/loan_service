package domain

import "time"

// Investment represents an investment made by an investor into a
// specific loan. Each call to the invest endpoint will create a new
// record, capturing the invested amount and linking it to both the
// loan and the investor. Multiple investments by the same investor
// toward the same loan are allowed and aggregated by the service
// layer.
type Investment struct {
    ID         string    `gorm:"type:uuid;primaryKey" json:"id"`
    LoanID     string    `gorm:"type:uuid;not null" json:"loan_id"`
    InvestorID string    `gorm:"type:uuid;not null" json:"investor_id"`
    Amount     float64   `gorm:"not null" json:"amount"`
    CreatedAt  time.Time `json:"created_at"`
}