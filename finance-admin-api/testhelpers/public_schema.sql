CREATE ROLE api;

CREATE TABLE public.persons
(
    id            INTEGER NOT NULL
        PRIMARY KEY,
    firstname     VARCHAR(255) DEFAULT NULL,
    surname       VARCHAR(255) DEFAULT NULL,
    caserecnumber VARCHAR(255) DEFAULT NULL,
    feepayer_id   INTEGER      DEFAULT NULL
        CONSTRAINT fk_a25cc7d3aff282de
            REFERENCES public.persons,
    deputytype    VARCHAR(255) DEFAULT NULL
);

ALTER TABLE public.persons
    OWNER TO api;

CREATE TABLE public.cases
(
    id          INTEGER NOT NULL
        PRIMARY KEY,
    client_id   INTEGER
        CONSTRAINT fk_1c1b038b19eb6921
            REFERENCES public.persons,
    orderstatus VARCHAR(255) DEFAULT NULL
);

CREATE INDEX cases_orderstatus_index ON public.cases (orderstatus);

CREATE INDEX idx_1c1b038b19eb6921 ON public.cases (client_id);

CREATE SEQUENCE persons_id_seq;

ALTER SEQUENCE persons_id_seq OWNER TO api;

CREATE SEQUENCE cases_id_seq;

ALTER SEQUENCE cases_id_seq OWNER TO api;