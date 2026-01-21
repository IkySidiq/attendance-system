CREATE TABLE sessions (
    id UUID PRIMARY KEY,
    employee_id UUID NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMPTZ NOT NULL,

    ip INET,
    user_agent TEXT,

    revoked BOOLEAN NOT NULL DEFAULT false,

    refresh_token TEXT,
    refresh_expires_at TIMESTAMPTZ,

    device_type VARCHAR(50),
    device_name VARCHAR(100),
    platform VARCHAR(50),

    last_activity TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    remember_me BOOLEAN NOT NULL DEFAULT false,

    session_data JSONB,

    CONSTRAINT fk_sessions_employee
        FOREIGN KEY (employee_id)
        REFERENCES employees(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT check_expires_after_created
        CHECK (expires_at > created_at),

    CONSTRAINT check_refresh_expires_after_expires
        CHECK (
            refresh_expires_at IS NULL
            OR refresh_expires_at > expires_at
        ),

    CONSTRAINT check_device_type
        CHECK (
            device_type IS NULL
            OR device_type IN (
                'mobile',
                'desktop',
                'tablet',
                'tv',
                'watch',
                'other'
            )
        ),

    CONSTRAINT check_platform
        CHECK (
            platform IS NULL
            OR platform IN (
                'android',
                'ios',
                'web',
                'windows',
                'macos',
                'linux',
                'other'
            )
        )
);

CREATE INDEX idx_sessions_employee_id ON sessions(employee_id);
CREATE INDEX idx_sessions_employee_active ON sessions(employee_id, revoked);

CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
CREATE INDEX idx_sessions_refresh_expires_at ON sessions(refresh_expires_at);
CREATE INDEX idx_sessions_refresh_token ON sessions(refresh_token);

CREATE INDEX idx_sessions_last_activity ON sessions(last_activity);
CREATE INDEX idx_sessions_revoked ON sessions(revoked);

CREATE INDEX idx_sessions_device_type ON sessions(device_type);
CREATE INDEX idx_sessions_platform ON sessions(platform);
