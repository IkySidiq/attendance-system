ALTER TABLE employees
ADD COLUMN last_login_at timestamp NOT NULL DEFAULT now();
