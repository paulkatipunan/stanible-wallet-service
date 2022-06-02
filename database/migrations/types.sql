-- ALTER TYPE user_type ADD VALUE 'creator' AFTER 'admin';
-- ALTER TYPE tx_status ADD VALUE 'cancelled' AFTER 'success';
CREATE TYPE user_type AS ENUM (
    'gcash_deposit',
    'gcash_withdraw',
    'paymaya_deposit',
    'paymaya_withdraw',
    'grabpay_withdraw',
    'grabpay_deposit',
    'treasury',
    'admin',
    'creator',
    'regular_user'
);

CREATE TYPE tx_status AS ENUM ('pending', 'success', 'cancelled', 'failed');

CREATE TYPE tx_type AS ENUM ('deposit', 'withdraw', 'buy', 'refund', 'fee');