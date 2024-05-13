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

DROP DATABASE IF EXISTS crm;
--
-- Name: crm; Type: DATABASE; Schema: -; Owner: postgres
--

CREATE DATABASE crm WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'Russian_Russia.1251';


ALTER DATABASE crm OWNER TO postgres;

\connect crm

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
-- Name: sessions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sessions (
    id integer NOT NULL,
    id_user integer,
    access_token character varying,
    date_expiration_access_token timestamp without time zone,
    refresh_token character varying,
    date_expiration_refresh_token timestamp without time zone
);


ALTER TABLE public.sessions OWNER TO postgres;

--
-- Name: sessions_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.sessions ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.sessions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id integer NOT NULL,
    login character varying NOT NULL,
    password character varying NOT NULL,
    name character varying,
    surname character varying,
    patronymic character varying,
    role integer DEFAULT 0 NOT NULL
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
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: sessions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.sessions (id, id_user, access_token, date_expiration_access_token, refresh_token, date_expiration_refresh_token) FROM stdin;
10	35	6552ecca9c8742e9ad2e7594ca67d87ff940fd92c56c5bc6114754346054aa75	2024-05-12 17:23:00	b1f288e7f995761c3aa92ed3392793d80c1e42797b365d7c0d87d93ddc87a7db	2024-05-12 17:28:31.199679
11	36	fbe2e4c47850e22ef829a42e4ba4032223ea357c8d31f9f20cbb698ebc45927f	2024-05-12 17:25:00	eeab146e682673376cf86015ed07cd310fe7751b781e58585581a0d92083df04	2024-05-12 17:30:33.130748
50	18	e764f575c80e4b6b92c5c2f133b601b98bfab2f73cf4ee043dd18472f9e937ff	2024-05-12 23:17:51.66655	69c8aa59445518dfb19a805c2ff810fb3bd46368f03eccd485f99424073f5e67	2024-05-12 23:19:22.462235
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, login, password, name, surname, patronymic, role) FROM stdin;
18	123	d13e40a29eb822914e3dc8098d3d1b05bf8fe1c602e3d021cd9ae73b18216582	\N	\N	\N	0
21	1234	d13e40a29eb822914e3dc8098d3d1b05bf8fe1c602e3d021cd9ae73b18216582	\N	\N	\N	0
22	12345	d13e40a29eb822914e3dc8098d3d1b05bf8fe1c602e3d021cd9ae73b18216582	\N	\N	\N	0
23	123456	d13e40a29eb822914e3dc8098d3d1b05bf8fe1c602e3d021cd9ae73b18216582	\N	\N	\N	0
24	asdasd	10291d744512a7e669b80e18ae4a243a5e7b6c5ec4f728cdd519e4212d514732	\N	\N	\N	0
26		10291d744512a7e669b80e18ae4a243a5e7b6c5ec4f728cdd519e4212d514732	\N	\N	\N	0
27	1234567	d13e40a29eb822914e3dc8098d3d1b05bf8fe1c602e3d021cd9ae73b18216582	\N	\N	\N	0
28	12345678	d13e40a29eb822914e3dc8098d3d1b05bf8fe1c602e3d021cd9ae73b18216582	\N	\N	\N	0
29	123456789	d13e40a29eb822914e3dc8098d3d1b05bf8fe1c602e3d021cd9ae73b18216582	\N	\N	\N	0
30	1234567891	d13e40a29eb822914e3dc8098d3d1b05bf8fe1c602e3d021cd9ae73b18216582	\N	\N	\N	0
31	12345678912	d13e40a29eb822914e3dc8098d3d1b05bf8fe1c602e3d021cd9ae73b18216582	\N	\N	\N	0
32	123456789123	d13e40a29eb822914e3dc8098d3d1b05bf8fe1c602e3d021cd9ae73b18216582	\N	\N	\N	0
33	1234567891234	d13e40a29eb822914e3dc8098d3d1b05bf8fe1c602e3d021cd9ae73b18216582	\N	\N	\N	0
34	12345678912345	d13e40a29eb822914e3dc8098d3d1b05bf8fe1c602e3d021cd9ae73b18216582	\N	\N	\N	0
35	123456789123456	d13e40a29eb822914e3dc8098d3d1b05bf8fe1c602e3d021cd9ae73b18216582	\N	\N	\N	0
36	1234567891234567	d13e40a29eb822914e3dc8098d3d1b05bf8fe1c602e3d021cd9ae73b18216582	\N	\N	\N	0
37	12345678912345678	d13e40a29eb822914e3dc8098d3d1b05bf8fe1c602e3d021cd9ae73b18216582	\N	\N	\N	0
38	123456789123456789	d13e40a29eb822914e3dc8098d3d1b05bf8fe1c602e3d021cd9ae73b18216582	\N	\N	\N	0
40	12377	d13e40a29eb822914e3dc8098d3d1b05bf8fe1c602e3d021cd9ae73b18216582	\N	\N	\N	0
42	123778	d13e40a29eb822914e3dc8098d3d1b05bf8fe1c602e3d021cd9ae73b18216582	\N	\N	\N	0
43	333	d13e40a29eb822914e3dc8098d3d1b05bf8fe1c602e3d021cd9ae73b18216582	\N	\N	\N	0
\.


--
-- Name: sessions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.sessions_id_seq', 50, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 43, true);


--
-- Name: sessions sessions_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_pk PRIMARY KEY (id);


--
-- Name: users users_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pk UNIQUE (login);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: sessions sessions_users_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_users_id_fk FOREIGN KEY (id_user) REFERENCES public.users(id);


--
-- PostgreSQL database dump complete
--

