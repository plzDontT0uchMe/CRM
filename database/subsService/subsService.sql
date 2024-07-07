--
-- PostgreSQL database dump
--

-- Dumped from database version 16.2
-- Dumped by pg_dump version 16.2

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

DROP DATABASE IF EXISTS "crm-subsService";
--
-- Name: crm-subsService; Type: DATABASE; Schema: -; Owner: postgres
--

CREATE DATABASE "crm-subsService" WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'Russian_Russia.1251';


ALTER DATABASE "crm-subsService" OWNER TO postgres;

\connect -reuse-previous=on "dbname='crm-subsService'"

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: active; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.active (
    id integer NOT NULL,
    id_client integer,
    id_subscription integer,
    id_trainer integer,
    date_expires timestamp without time zone
);


ALTER TABLE public.active OWNER TO postgres;

--
-- Name: active_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.active_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.active_id_seq OWNER TO postgres;

--
-- Name: active_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.active_id_seq OWNED BY public.active.id;


--
-- Name: applications; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.applications (
    id integer NOT NULL,
    id_client integer NOT NULL,
    id_subcription integer NOT NULL,
    id_trainer integer,
    duration integer
);


ALTER TABLE public.applications OWNER TO postgres;

--
-- Name: applications_id_client_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.applications_id_client_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.applications_id_client_seq OWNER TO postgres;

--
-- Name: applications_id_client_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.applications_id_client_seq OWNED BY public.applications.id_client;


--
-- Name: applications_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.applications_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.applications_id_seq OWNER TO postgres;

--
-- Name: applications_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.applications_id_seq OWNED BY public.applications.id;


--
-- Name: applications_id_subcription_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.applications_id_subcription_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.applications_id_subcription_seq OWNER TO postgres;

--
-- Name: applications_id_subcription_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.applications_id_subcription_seq OWNED BY public.applications.id_subcription;


--
-- Name: possibilities; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.possibilities (
    id integer NOT NULL,
    id_subscription integer,
    possibility character varying
);


ALTER TABLE public.possibilities OWNER TO postgres;

--
-- Name: possibilities_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.possibilities_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.possibilities_id_seq OWNER TO postgres;

--
-- Name: possibilities_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.possibilities_id_seq OWNED BY public.possibilities.id;


--
-- Name: subscription; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.subscription (
    id integer NOT NULL,
    name character varying,
    price double precision,
    description character varying
);


ALTER TABLE public.subscription OWNER TO postgres;

--
-- Name: subs_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.subs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.subs_id_seq OWNER TO postgres;

--
-- Name: subs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.subs_id_seq OWNED BY public.subscription.id;


--
-- Name: active id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.active ALTER COLUMN id SET DEFAULT nextval('public.active_id_seq'::regclass);


--
-- Name: applications id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.applications ALTER COLUMN id SET DEFAULT nextval('public.applications_id_seq'::regclass);


--
-- Name: applications id_client; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.applications ALTER COLUMN id_client SET DEFAULT nextval('public.applications_id_client_seq'::regclass);


--
-- Name: applications id_subcription; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.applications ALTER COLUMN id_subcription SET DEFAULT nextval('public.applications_id_subcription_seq'::regclass);


--
-- Name: possibilities id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.possibilities ALTER COLUMN id SET DEFAULT nextval('public.possibilities_id_seq'::regclass);


--
-- Name: subscription id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.subscription ALTER COLUMN id SET DEFAULT nextval('public.subs_id_seq'::regclass);


--
-- Data for Name: active; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.active (id, id_client, id_subscription, id_trainer, date_expires) FROM stdin;
7	54	1	\N	\N
4	61	2	\N	\N
3	67	3	61	2024-06-09 20:43:52
1	63	3	61	2025-06-11 14:55:40.088055
2	64	3	61	2024-12-13 17:07:02.949549
8	0	1	\N	\N
9	75	1	\N	\N
10	76	1	\N	\N
\.


--
-- Data for Name: applications; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.applications (id, id_client, id_subcription, id_trainer, duration) FROM stdin;
41	64	1	\N	\N
\.


--
-- Data for Name: possibilities; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.possibilities (id, id_subscription, possibility) FROM stdin;
2	1	Make training plans
1	1	View exercises
3	2	View exercises
4	2	Make training plans
5	2	Sign up for a gym without a trainer
6	3	View exercises
7	3	Make training plans
8	3	Sign up for a gym without a trainer
9	3	Sign up for a gym with a trainer
\.


--
-- Data for Name: subscription; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.subscription (id, name, price, description) FROM stdin;
1	Free	0	For all users
3	Premium	29.99	For a playboy, philanthropist, billionaire and just a genius
2	Standard	14.99	For half users
\.


--
-- Name: active_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.active_id_seq', 10, true);


--
-- Name: applications_id_client_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.applications_id_client_seq', 1, false);


--
-- Name: applications_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.applications_id_seq', 41, true);


--
-- Name: applications_id_subcription_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.applications_id_subcription_seq', 1, false);


--
-- Name: possibilities_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.possibilities_id_seq', 10, true);


--
-- Name: subs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.subs_id_seq', 3, true);


--
-- Name: active active_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.active
    ADD CONSTRAINT active_pk PRIMARY KEY (id);


--
-- Name: applications applications_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.applications
    ADD CONSTRAINT applications_pk PRIMARY KEY (id);


--
-- Name: applications applications_pk_2; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.applications
    ADD CONSTRAINT applications_pk_2 UNIQUE (id_client);


--
-- Name: possibilities possibilities_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.possibilities
    ADD CONSTRAINT possibilities_pk PRIMARY KEY (id);


--
-- Name: subscription subs_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.subscription
    ADD CONSTRAINT subs_pk PRIMARY KEY (id);


--
-- Name: active active_subscription_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.active
    ADD CONSTRAINT active_subscription_id_fk FOREIGN KEY (id_subscription) REFERENCES public.subscription(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: applications applications_subscription_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.applications
    ADD CONSTRAINT applications_subscription_id_fk FOREIGN KEY (id_subcription) REFERENCES public.subscription(id);


--
-- Name: possibilities possibilities_subscription_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.possibilities
    ADD CONSTRAINT possibilities_subscription_id_fk FOREIGN KEY (id_subscription) REFERENCES public.subscription(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

