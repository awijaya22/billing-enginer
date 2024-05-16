package loan

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type repo struct {
	db *sql.DB
}

type Repo interface {
	// user-loan
	InsertLoan(ctx context.Context, loan Loan) (int, error)
	GetLoanByID(ctx context.Context, loanID int) (Loan, error)

	// user billing
	InsertBilling(ctx context.Context, billing Billing) (int, error)
	UpdateBilling(ctx context.Context, billingID int) error
	GetBillingByID(ctx context.Context, billingID int) (Billing, error)
	GetTotalPaidBillingByLoanID(ctx context.Context, loanID int) (int, error)
	CountUserUnpaidLoan(ctx context.Context, userID string) (int, error)
}

func NewRepo(db *sql.DB) Repo {
	return &repo{db: db}
}

func (r *repo) InsertLoan(ctx context.Context, loan Loan) (int, error) {
	query := `
	insert into user_loan (user_id, principal_loan_amt, interest_rate, start_at, end_at)
	values ($1, $2, $3, $4, $5)
	returning loan_id;
	`
	var loanID int
	err := r.db.QueryRow(query, loan.UserID, loan.PrincipalLoanAmt, loan.InterestRate, loan.StartAt, loan.EndAt).Scan(&loanID)
	if err != nil {
		fmt.Println("[InsertLoan] error = ", err)
		return 0, err
	}
	return loanID, nil
}

func (r *repo) GetLoanByID(ctx context.Context, loanID int) (Loan, error) {
	query := `
	select loan_id, user_id, principal_loan_amt, interest_rate, start_at, end_at, created_at, updated_at
	from user_loan
	where loan_id = $1;
	`
	var loan Loan
	row := r.db.QueryRow(query, loanID)
	err := row.Scan(&loan.ID, &loan.UserID, &loan.PrincipalLoanAmt, &loan.InterestRate, &loan.StartAt, &loan.EndAt, &loan.CreatedAt, &loan.UpdatedAt)
	if err != nil {
		return loan, err
	}
	return loan, nil
}

func (r *repo) InsertBilling(ctx context.Context, billing Billing) (int, error) {
	query := `
	insert into user_billing (loan_id, pay_amt, due_at)
	values ($1, $2, $3)
	returning billing_id;
	`
	var billingID int
	err := r.db.QueryRow(query, billing.LoanID, billing.PaymentAmt, billing.DueAt).Scan(&billingID)
	if err != nil {
		fmt.Println("[CreateBilling] error = ", err)
		return 0, err
	}
	return billingID, nil
}

func (r *repo) UpdateBilling(ctx context.Context, billingID int) error {
	query := `
	update user_billing
	set is_paid = $2, paid_at = current_timestamp
	where billing_id = $1;
	`
	_, err := r.db.Exec(query, billingID, true)
	if err != nil {
		fmt.Println("[UpdateBilling] error = ", err)
		return err
	}
	return nil
}

func (r *repo) GetBillingByID(ctx context.Context, billingID int) (Billing, error) {
	query := `
	select billing_id, loan_id, pay_amt, due_at, is_paid, paid_at, created_at
	from user_billing
	where billing_id = $1;
	`

	var billing Billing
	row := r.db.QueryRow(query, billingID)
	err := row.Scan(&billing.ID, &billing.LoanID, &billing.PaymentAmt, &billing.DueAt, &billing.IsPaid, &billing.PaidAt, &billing.CreatedAt)
	if err != nil {
		fmt.Println("[GetBillingByID] error = ", err)
		return billing, err
	}
	return billing, nil
}

func (r *repo) GetTotalPaidBillingByLoanID(ctx context.Context, loanID int) (int, error) {
	query := `
	select COALESCE(SUM(pay_amt), 0)
	from user_billing 
	where loan_id = $1 and is_paid is true;
	`
	var totalPayAmt int
	err := r.db.QueryRow(query, loanID).Scan(&totalPayAmt)
	if err != nil {
		return 0, err
	}
	return totalPayAmt, nil
}

func (r *repo) CountUserUnpaidLoan(ctx context.Context, userID string) (int, error) {
	query := `
	select count(*)
	from user_billing left join user_loan on user_billing.loan_id = user_loan.loan_id
	where user_loan.user_id = $1 and user_billing.is_paid is false and user_billing.due_at < $2
	`
	var count int
	err := r.db.QueryRow(query, userID, time.Now()).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
