package loan

import (
	"context"
	"fmt"
	"time"
)

type service struct {
	repo Repo
}

type Service interface {
	CreateLoan(ctx context.Context, userID string, principalAmt int, startAt *time.Time) error
	MakePayment(ctx context.Context, userID string, billingID int, paymentAmt int) error
	GetOutstanding(ctx context.Context, loanID int) (int, error)
	IsDelinquent(ctx context.Context, userID string) (bool, error)
}

func NewService(repo Repo) Service {
	return &service{repo: repo}
}

func (s *service) CreateLoan(ctx context.Context, userID string, principalAmt int, startAt *time.Time) error {
	endAt := startAt.Add(TotalRepayment * 7 * 24 * time.Hour)
	loanID, err := s.repo.InsertLoan(ctx, Loan{
		UserID:           userID,
		PrincipalLoanAmt: principalAmt,
		InterestRate:     InterestRate,
		StartAt:          startAt,
		EndAt:            &endAt,
	})
	if err != nil {
		fmt.Println("[CreateLoan][InsertLoan] error = ", err)
		return err
	}
	interestVal := (principalAmt * InterestRate / 100)

	payAmt := (principalAmt + interestVal) / TotalRepayment

	// create billing schedule
	for i := 0; i < TotalRepayment; i++ {
		t := time.Now().Local().Add(7 * 24 * time.Hour)
		// change time to 23:59
		dueAt := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 59, t.Location())
		billing := Billing{
			LoanID:     loanID,
			PaymentAmt: payAmt,
			DueAt:      &dueAt,
		}
		_, err := s.repo.InsertBilling(ctx, billing)
		if err != nil {
			fmt.Println("[CreateLoan][InsertBilling] error = ", err)
			return err
		}
	}
	return nil
}

func (s *service) MakePayment(ctx context.Context, userID string, billingID int, paymentAmt int) error {
	// get billing
	billing, err := s.repo.GetBillingByID(ctx, billingID)
	if err != nil {
		fmt.Println("[MakePayment][GetBillingByID] error = ", err)
		return err
	}

	if billing.IsPaid {
		return fmt.Errorf("billing ID is paid")
	}

	if billing.PaymentAmt != paymentAmt {
		// payment not exact
		return fmt.Errorf("payment amount must be exact")
	}

	// assuming payment process is done, update billingID is_paid
	err = s.repo.UpdateBilling(ctx, billingID)
	if err != nil {
		fmt.Println("[MakePayment][UpdateBilling] error = ", err)
		return err
	}

	fmt.Println("[info] successful payment of billingID =", billingID)
	return nil
}

func (s *service) GetOutstanding(ctx context.Context, loanID int) (int, error) {
	// get loan
	loan, err := s.repo.GetLoanByID(ctx, loanID)
	if err != nil {
		return 0, err
	}
	totalPayment := loan.PrincipalLoanAmt + (loan.PrincipalLoanAmt * loan.InterestRate / 100)

	// get total payment
	paidAmt, err := s.repo.GetTotalPaidBillingByLoanID(ctx, loanID)
	if err != nil {
		return 0, err
	}
	return totalPayment - paidAmt, nil
}

func (s *service) IsDelinquent(ctx context.Context, userID string) (bool, error) {
	unpaidLoanCount, err := s.repo.CountUserUnpaidLoan(ctx, userID)
	if err != nil {
		return false, err
	}
	if unpaidLoanCount >= 2 {
		return true, nil
	}
	return false, nil
}
