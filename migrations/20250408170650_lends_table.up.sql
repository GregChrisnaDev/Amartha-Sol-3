CREATE TABLE lends
(
    id                      BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    loan_id                 BIGINT,
    user_id                 BIGINT,
    amount                  BIGINT,
    agreement_file_path     VARCHAR(255),
    user_sign_path          VARCHAR(255),
    created_at              TIMESTAMP WITH TIME ZONE,
    updated_at              TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_lends_user_id ON lends(user_id);
CREATE INDEX idx_lends_loan_id_user_id ON lends(loan_id, user_id);