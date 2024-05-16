package loan

import "time"

const (
	// for simplicity hard coded
	InterestRate   = 10 // in percentage
	TotalRepayment = 10 // weeks
)

type Loan struct {
	ID               int
	UserID           string
	PrincipalLoanAmt int
	InterestRate     int
	StartAt          *time.Time
	EndAt            *time.Time
	CreatedAt        *time.Time
	UpdatedAt        *time.Time
}

type Billing struct {
	ID         int
	LoanID     int
	PaymentAmt int
	IsPaid     bool
	PaidAt     *time.Time
	DueAt      *time.Time
	CreatedAt  *time.Time
}
