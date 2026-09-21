CREATE TABLE merchants (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE accounts (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    merchant_id BIGINT NOT NULL REFERENCES merchants(id),

    balance BIGINT NOT NULL DEFAULT 0 
        CHECK (balance >= 0),

    currency TEXT NOT NULL DEFAULT 'TWD'
        CHECK (currency IN ('TWD')),

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE payments (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    account_id BIGINT NOT NULL REFERENCES accounts(id),

    amount BIGINT NOT NULL 
        CHECK (amount > 0),

    status TEXT NOT NULL
        CHECK (status IN ('processing', 'completed', 'failed')),

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE transfers (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    from_account_id BIGINT NOT NULL REFERENCES accounts(id),
    to_account_id BIGINT NOT NULL REFERENCES accounts(id),
    
    amount BIGINT NOT NULL
        CHECK (amount > 0),

    status TEXT NOT NULL
        CHECK (status IN ('processing', 'completed', 'failed')),

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CHECK (from_account_id <> to_account_id)
);