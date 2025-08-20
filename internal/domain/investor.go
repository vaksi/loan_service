package domain

import "time"

// Investor represents an individual or entity that invests funds into a
// loan. Each investor may contribute to multiple loans and each loan
// may have multiple investors. The name and email fields are optional
// but can be used to send notifications such as agreement letters.
type Investor struct {
    ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
    Name      string    `gorm:"size:100" json:"name"`
    Email     string    `gorm:"size:100" json:"email"`
    CreatedAt time.Time `json:"created_at"`
}