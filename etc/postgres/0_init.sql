CREATE TABLE public.subscriptions (
extra_nonce_1 SERIAL,
extra_nonce_2 INT NOT NULL,
set_difficulty VARCHAR(255) NOT NULL,
notify VARCHAR(255) NOT NULL,
subscriber VARCHAR(255) NOT NULL,
created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now() NOT NULL,
active_session BOOLEAN NOT NULL DEFAULT TRUE
);