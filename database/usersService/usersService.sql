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

DROP DATABASE IF EXISTS "crm-usersService";
--
-- Name: crm-usersService; Type: DATABASE; Schema: -; Owner: postgres
--

CREATE DATABASE "crm-usersService" WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'Russian_Russia.1251';


ALTER DATABASE "crm-usersService" OWNER TO postgres;

\connect -reuse-previous=on "dbname='crm-usersService'"

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
-- Name: trainers; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.trainers (
    id integer NOT NULL,
    id_account integer,
    exp integer,
    sport character varying,
    achievement character varying
);


ALTER TABLE public.trainers OWNER TO postgres;

--
-- Name: trainers_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.trainers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.trainers_id_seq OWNER TO postgres;

--
-- Name: trainers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.trainers_id_seq OWNED BY public.trainers.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id integer NOT NULL,
    id_account integer NOT NULL,
    name character varying,
    surname character varying,
    patronymic character varying,
    gender integer DEFAULT 0,
    date_born date,
    "position" character varying
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: trainers id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.trainers ALTER COLUMN id SET DEFAULT nextval('public.trainers_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: trainers; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.trainers (id, id_account, exp, sport, achievement) FROM stdin;
1	61	4	Body-building	KMS of Body-building
3	65	5	Powerlifter	Master of Sports in Powerlifter
2	61	6	Hand to hand combat	\nMaster of Sports in Hand-to-Hand Combat
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, id_account, name, surname, patronymic, gender, date_born, "position") FROM stdin;
3	61	Kostya	The	Beast	1	1999-02-10	trainer
4	63	3211	2222	2223	2	2005-05-10	client
1	54	Dmitry1	Homyakov	Sergeevich	1	2002-03-20	manager
6	65	Minion	Champion		0	\N	trainer
10	67	test123			0	\N	client
5	64	233	34	23552	2	1999-03-20	client
\.


--
-- Name: trainers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.trainers_id_seq', 3, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 10, true);


--
-- Name: trainers trainers_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.trainers
    ADD CONSTRAINT trainers_pk PRIMARY KEY (id);


--
-- Name: users users_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pk PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

