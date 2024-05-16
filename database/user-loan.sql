CREATE TABLE user_loan (
    loan_id serial primary key, 
    user_id varchar(50) not null, 
    principal_loan_amt int not null, -- principal amount 
    interest_rate int not null, 
    start_at TIMESTAMP not null, 
    end_at TIMESTAMP not null,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);