CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    registered BOOLEAN DEFAULT false NOT NULL
);

CREATE TABLE groups (
    group_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_name TEXT NOT NULL,
    description TEXT,
    created_by UUID REFERENCES users (user_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL
);

CREATE TABLE group_members (
    user_id UUID REFERENCES users (user_id) ON DELETE CASCADE,
    group_id UUID REFERENCES groups (group_id) ON DELETE CASCADE,
    joined_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    PRIMARY KEY (user_id, group_id)
);

CREATE TABLE expenses (
    expense_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id UUID REFERENCES groups (group_id) ON DELETE CASCADE,
    added_by UUID REFERENCES users (user_id) ON DELETE SET NULL,
    title TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    amount DOUBLE PRECISION NOT NULL,
    is_incomplete_amount BOOLEAN DEFAULT FALSE,
    is_incomplete_split BOOLEAN DEFAULT FALSE,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION
);

CREATE TABLE expense_splits (
    expense_id UUID REFERENCES expenses (expense_id) ON DELETE CASCADE,
    user_id UUID REFERENCES users (user_id) ON DELETE CASCADE,
    amount DOUBLE PRECISION NOT NULL,
    user_role TEXT CHECK (user_role IN ('paid', 'owes')),
    PRIMARY KEY (expense_id, user_id)
);
