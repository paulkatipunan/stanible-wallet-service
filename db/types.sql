CREATE TYPE ramp_type AS ENUM ('gcash', 'paymaya');
ALTER TYPE ramp_type ADD VALUE 'dragonpay' AFTER 'paymaya';

CREATE TYPE tx_type AS ENUM ('deposit', 'withdraw', 'buy', 'refund');