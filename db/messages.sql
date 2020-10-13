CREATE TABLE public.messages
(
    content text COLLATE pg_catalog."default",
    id bigserial NOT NULL,
    title character varying COLLATE pg_catalog."default" NOT NULL,
    "createAt" timestamp without time zone,
    CONSTRAINT m_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE public.messages
    OWNER to postgres;