CREATE TABLE lends
(
    id                      BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    loan_id                 BIGINT,
    user_id                 BIGINT,
    amount                  BIGINT,
    AgreementFilePath       VARCHAR(255),
    created_at              TIMESTAMP WITH TIME ZONE,
    updated_at              TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_lends_user_id ON lends(user_id);
CREATE INDEX idx_lends_loan_id ON lends(loan_id);