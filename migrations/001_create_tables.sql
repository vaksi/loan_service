-- migration: create initial schema for Amartha loan service
-- This file contains SQL statements to create all tables required by
-- the loan service. It should be executed on a clean PostgreSQL
-- database. The schema uses UUID primary keys generated via the
-- uuid‑ossp extension. Relationships are enforced with foreign
-- keys.

-- enable uuid generation extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- loans table stores basic loan information and state
CREATE TABLE IF NOT EXISTS loans (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    borrower_id         VARCHAR(50) NOT NULL,
    principal           NUMERIC(12,2) NOT NULL,
    rate                NUMERIC(6,2) NOT NULL,
    roi                 NUMERIC(6,2) NOT NULL,
    agreement_letter_url TEXT,
    state               VARCHAR(20) NOT NULL,
    created_at          TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMP NOT NULL DEFAULT NOW()
);

-- approvals table stores a single approval record per loan
CREATE TABLE IF NOT EXISTS approvals (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    loan_id       UUID NOT NULL UNIQUE REFERENCES loans(id) ON DELETE CASCADE,
    picture_url   TEXT NOT NULL,
    employee_id   VARCHAR(50) NOT NULL,
    approval_date DATE NOT NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT NOW()
);

-- investors table stores investor identities
CREATE TABLE IF NOT EXISTS investors (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name       VARCHAR(100),
    email      VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- investments table stores contributions from investors to loans
CREATE TABLE IF NOT EXISTS investments (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    loan_id     UUID NOT NULL REFERENCES loans(id) ON DELETE CASCADE,
    investor_id UUID NOT NULL REFERENCES investors(id) ON DELETE CASCADE,
    amount      NUMERIC(12,2) NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_investments_loan_id ON investments (loan_id);
CREATE INDEX IF NOT EXISTS idx_investments_investor_id ON investments (investor_id);

-- disbursements table stores final hand‑over details
CREATE TABLE IF NOT EXISTS disbursements (
    id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    loan_id           UUID NOT NULL UNIQUE REFERENCES loans(id) ON DELETE CASCADE,
    agreement_url     TEXT NOT NULL,
    employee_id       VARCHAR(50) NOT NULL,
    disbursement_date DATE NOT NULL,
    created_at        TIMESTAMP NOT NULL DEFAULT NOW()
);