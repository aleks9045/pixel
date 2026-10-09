--
-- PostgreSQL database dump
--

\restrict mEmh2vKwpX2EV8M22DfgBHa5oPf2dNC7WRTCLhjCBNCMeIbFWAvxXfioecbFfCk

-- Dumped from database version 18.6
-- Dumped by pg_dump version 18.6

SET statement_timeout = 0;

SET lock_timeout = 0;

SET idle_in_transaction_session_timeout = 0;

SET transaction_timeout = 0;

SET client_encoding = 'UTF8';

SET standard_conforming_strings = on;

SELECT pg_catalog.set_config ('search_path', '', false);

SET check_function_bodies = false;

SET xmloption = content;

SET client_min_messages = warning;

SET row_security = off;

--
-- Name: public; Type: SCHEMA; Schema: -; Owner: postgres
--

-- *not* creating schema, since initdb creates it

ALTER SCHEMA public OWNER TO postgres;

--
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: postgres
--

COMMENT ON SCHEMA public IS '';

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: pixel; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pixel (
    x smallint NOT NULL,
    y smallint NOT NULL,
    color smallint NOT NULL,
    placed_by integer
);

ALTER TABLE public.pixel OWNER TO postgres;

CREATE TABLE public.pixel_history (
    x smallint NOT NULL,
    y smallint NOT NULL,
    placed_on timestamp with time zone DEFAULT now() NOT NULL,
    placed_by integer
);

ALTER TABLE public.pixel_history OWNER TO postgres;

CREATE TABLE public.users (
    id integer NOT NULL,
    user_id character varying NOT NULL,
    username character varying(255) NOT NULL,
    name character varying(255) NOT NULL,
    avatar_url character varying(255) NOT NULL,
    student_group character varying(16),
    is_admin boolean DEFAULT false NOT NULL,
    time_limit smallint DEFAULT 5000 NOT NULL,
    last_placed_pixel timestamp with time zone
);

ALTER TABLE public.users OWNER TO postgres;

CREATE SEQUENCE public.users_id_seq AS integer START
WITH
    1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;

ALTER SEQUENCE public.users_id_seq OWNER TO postgres;

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;

ALTER TABLE ONLY public.users
ALTER COLUMN id
SET DEFAULT nextval(
    'public.users_id_seq'::regclass
);

ALTER TABLE ONLY public.pixel
ADD CONSTRAINT pixel_pkey PRIMARY KEY (x, y);

ALTER TABLE ONLY public.users
ADD CONSTRAINT users_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.users
ADD CONSTRAINT users_username_key UNIQUE (username);

CREATE INDEX idx_pixel_placed_by ON public.pixel USING btree (placed_by);

CREATE INDEX idx_pixel_y_x ON public.pixel USING btree (y, x);

CREATE INDEX idx_pixel_history_placed_by ON public.pixel_history (placed_by);

CREATE INDEX idx_pixel_history_user_time ON public.pixel_history (placed_by, placed_on DESC);

ALTER TABLE ONLY public.pixel_history
ADD CONSTRAINT pixel_history_placed_by_fkey FOREIGN KEY (placed_by) REFERENCES public.users (id);

ALTER TABLE ONLY public.pixel
ADD CONSTRAINT pixel_placed_by_fkey FOREIGN KEY (placed_by) REFERENCES public.users (id);

REVOKE USAGE ON SCHEMA public FROM PUBLIC;

INSERT INTO
    pixel (x, y, color, placed_by)
SELECT x, y, 0, NULL
FROM generate_series(0, 499) AS x
    CROSS JOIN generate_series(0, 499) AS y;
--
-- PostgreSQL database dump complete
--

\unrestrict mEmh2vKwpX2EV8M22DfgBHa5oPf2dNC7WRTCLhjCBNCMeIbFWAvxXfioecbFfCk