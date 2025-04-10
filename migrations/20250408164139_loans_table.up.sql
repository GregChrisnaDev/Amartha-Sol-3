CREATE TABLE loans
(
    id                      BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id                 BIGINT,
    principal_amount        BIGINT,
    rate                    INTEGER,
    loan_duration           INTEGER,
    status                  INTEGER,
    proposed_date           TIMESTAMP WITH TIME ZONE,
    picture_proof_filepath  VARCHAR(255),
    approver_uid            BIGINT,
    approval_date           TIMESTAMP WITH TIME ZONE,
    disburser_uid           BIGINT,
    user_sign_path          VARCHAR(255),
    disbursement_date       TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_loans_user_id ON loans(user_id);