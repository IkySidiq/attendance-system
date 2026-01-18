CREATE TABLE branches (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    branch_code VARCHAR(50) NOT NULL UNIQUE,
    location TEXT,
    qr_code_attendance TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE INDEX idx_branches_branch_code ON branches(branch_code);
CREATE INDEX idx_branches_is_active ON branches(is_active);
