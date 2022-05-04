CREATE TYPE user_type AS ENUM ('gcash', 'paymaya', 'treasury', 'admin', 'regular_user');
ALTER TYPE user_type ADD VALUE 'dragonpay' AFTER 'paymaya';

CREATE TYPE tx_type AS ENUM ('deposit', 'withdraw', 'buy', 'refund');