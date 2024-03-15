CREATE DATABASE pack_calculator;
GRANT ALL PRIVILEGES ON DATABASE pack_calculator to postgres;
\connect pack_calculator

CREATE TABLE public.pack (
    pack_size int NOT NULL,
    PRIMARY KEY(pack_size)
);

CREATE TABLE public.order (
    order_id uuid NOT NULL,
    order_quantity int NOT NULL,
    PRIMARY KEY(order_id)
);

CREATE TABLE public.order_packs (
    order_packs_id uuid NOT NULL,
    order_id uuid REFERENCES public.order(order_id),
    pack_size int NOT NULL,
    pack_quantity int NOT NULL,
    PRIMARY KEY(order_packs_id, order_id, pack_size)
);
