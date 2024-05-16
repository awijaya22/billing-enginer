-- user_loan table
create index idx_user_id on user_loan(loan_id);

-- user_billing table
create index idx_loan_id_is_paid_due_at ON user_billing(loan_id, is_paid, due_at);