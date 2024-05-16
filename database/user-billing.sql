CREATE TABLE user_billing (
    billing_id serial primary key, 
    loan_id int not null,
    pay_amt int not null, --payment need to pay by user each billing cycle
    due_at TIMESTAMP NOT NULL, 
    is_paid boolean default false, 
    paid_at timestamp default null, 
    created_at TIMESTAMP default current_timestamp,
    foreign key (loan_id) references user_loan(loan_id)
);

