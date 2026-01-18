CREATE TYPE employee_status AS ENUM (
    'active',
    'inactive',
    'suspended'
);

CREATE TABLE employees (
    id UUID PRIMARY KEY,

    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    username VARCHAR(100) NOT NULL,
    password TEXT NOT NULL,
    employee_code VARCHAR(50) NOT NULL UNIQUE,
    branch_id UUID NOT NULL,
    status employee_status NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT employees_email_unique UNIQUE (email),
    CONSTRAINT employees_username_unique UNIQUE (username),

    CONSTRAINT fk_employees_branch
        FOREIGN KEY (branch_id)
        REFERENCES branches(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);

CREATE INDEX idx_employees_email ON employees(email);
CREATE INDEX idx_employees_username ON employees(username);
CREATE INDEX idx_employees_employee_code ON employees(employee_code);
CREATE INDEX idx_employees_branch_id ON employees(branch_id);
CREATE INDEX idx_employees_status ON employees(status);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_employees_updated_at
BEFORE UPDATE ON employees
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();