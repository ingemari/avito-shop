CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    from_user INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    to_user INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount INT NOT NULL CHECK (amount > 0)
);
