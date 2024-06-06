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

DROP DATABASE IF EXISTS "crm-storageService";
--
-- Name: crm-storageService; Type: DATABASE; Schema: -; Owner: postgres
--

CREATE DATABASE "crm-storageService" WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'Russian_Russia.1251';


ALTER DATABASE "crm-storageService" OWNER TO postgres;

\connect -reuse-previous=on "dbname='crm-storageService'"

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
-- Name: files; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.files (
    id integer NOT NULL,
    id_account integer,
    link character varying,
    path character varying
);


ALTER TABLE public.files OWNER TO postgres;

--
-- Name: images_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.images_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.images_id_seq OWNER TO postgres;

--
-- Name: images_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.images_id_seq OWNED BY public.files.id;


--
-- Name: files id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.files ALTER COLUMN id SET DEFAULT nextval('public.images_id_seq'::regclass);


--
-- Data for Name: files; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.files (id, id_account, link, path) FROM stdin;
20	63	\N	\N
21	64	\N	\N
18	54	9edd9b2d4a07774e6664f7f9f3fb3839d8ebfde349597404ce48c21063ad725d.jpg	img\\9edd9b2d4a07774e6664f7f9f3fb3839d8ebfde349597404ce48c21063ad725d.jpg
19	61	2c99c04146587a7a8215f4d0361987f6b7f7efc3da716a22bd9603d2d5fea133.gif	img\\2c99c04146587a7a8215f4d0361987f6b7f7efc3da716a22bd9603d2d5fea133.gif
22	65	28460b43a6f9551869d2e427e0fddffca884ed70337e36c41165c359490a48a3.png	img\\28460b43a6f9551869d2e427e0fddffca884ed70337e36c41165c359490a48a3.png
23	67	\N	\N
\.


--
-- Name: images_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.images_id_seq', 23, true);


--
-- Name: files files_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.files
    ADD CONSTRAINT files_pk PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

