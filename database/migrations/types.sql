CREATE TYPE user_type AS ENUM ('gcash', 'paymaya', 'treasury', 'admin', 'regular_user');
ALTER TYPE user_type ADD VALUE 'dragonpay' AFTER 'paymaya';
ALTER TYPE user_type ADD VALUE 'creator' AFTER 'admin';

CREATE TYPE tx_type AS ENUM ('deposit', 'withdraw', 'buy', 'refund');

CREATE TYPE tx_status AS ENUM ('pending', 'success', 'failed');