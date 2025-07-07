Create Table spending_records (
    id SERIAL PRIMARY KEY,
    uuid UUID NOT NULL DEFAULT gen_random_uuid(),
    amount NUMERIC(10, 2) NOT NULL,
    remark TEXT,
    spending_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    category TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
)