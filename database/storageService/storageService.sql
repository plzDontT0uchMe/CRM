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
18	54	9edd9b2d4a07774e6664f7f9f3fb3839d8ebfde349597404ce48c21063ad725d.jpg	img\\9edd9b2d4a07774e6664f7f9f3fb3839d8ebfde349597404ce48c21063ad725d.jpg
19	61	2c99c04146587a7a8215f4d0361987f6b7f7efc3da716a22bd9603d2d5fea133.gif	img\\2c99c04146587a7a8215f4d0361987f6b7f7efc3da716a22bd9603d2d5fea133.gif
22	65	28460b43a6f9551869d2e427e0fddffca884ed70337e36c41165c359490a48a3.png	img\\28460b43a6f9551869d2e427e0fddffca884ed70337e36c41165c359490a48a3.png
23	67	5dba52a42fcfdd93f38a36da98dfdeaf1c819e265d207dc0371ca7a4aa1d8d99.jpg	img\\5dba52a42fcfdd93f38a36da98dfdeaf1c819e265d207dc0371ca7a4aa1d8d99.jpg
24	\N	francuzsckij_zhim_lezha_so_shtangoj.gif	exercise\\francuzsckij_zhim_lezha_so_shtangoj.gif
25	\N	good_morning_sto_shtangoj.gif	exercise\\good_morning_sto_shtangoj.gif
26	\N	mostik_shtanga.gif	exercise\\mostik_shtanga.gif
27	\N	obratnaya_tyaga_shtangi_v_naklone.gif	exercise\\obratnaya_tyaga_shtangi_v_naklone.gif
28	\N	podem_na_cypochki_so_shtangoj.gif	exercise\\podem_na_cypochki_so_shtangoj.gif
29	\N	podem_shtangi_pered_soboj.gif	exercise\\podem_shtangi_pered_soboj.gif
30	\N	prisedanie_frontalnoe_so_shtangoj.gif	exercise\\prisedanie_frontalnoe_so_shtangoj.gif
31	\N	prisedanie_so_shtangoj.gif	exercise\\prisedanie_so_shtangoj.gif
32	\N	razgibanie_iz_za_golovy_so_shtangoj_sidya.gif	exercise\\razgibanie_iz_za_golovy_so_shtangoj_sidya.gif
33	\N	rumynskaya_tyga_shtanga.gif	exercise\\rumynskaya_tyga_shtanga.gif
34	\N	sgibanie_ruk_obratnyj_hvat_so_shtangoj.gif	exercise\\sgibanie_ruk_obratnyj_hvat_so_shtangoj.gif
35	\N	sgibanie_ruk_so_shtangoj.gif	exercise\\sgibanie_ruk_so_shtangoj.gif
36	\N	shragi_so_shtangoj.gif	exercise\\shragi_so_shtangoj.gif
37	\N	stanovaya_sumo_tyaga_shtangi.gif	exercise\\stanovaya_sumo_tyaga_shtangi.gif
38	\N	stanovaya_tyaga.gif	exercise\\stanovaya_tyaga.gif
39	\N	sumo_prisedanie_so_shtangoj.gif	exercise\\sumo_prisedanie_so_shtangoj.gif
40	\N	tyaga_shtangi_k_podborodku.gif	exercise\\tyaga_shtangi_k_podborodku.gif
41	\N	tyaga_shtangi_v_naklone.gif	exercise\\tyaga_shtangi_v_naklone.gif
42	\N	vypad_bolgarskij_so_shtangoj.gif	exercise\\vypad_bolgarskij_so_shtangoj.gif
43	\N	vypad_na_meste_so_shtangoj.gif	exercise\\vypad_na_meste_so_shtangoj.gif
44	\N	vypad_nazad_so_shtangoj.gif	exercise\\vypad_nazad_so_shtangoj.gif
45	\N	vypad_vpered_so_shtangoj.gif	exercise\\vypad_vpered_so_shtangoj.gif
46	\N	zashagivanie_na_skamiu_so_shtangoj.gif	exercise\\zashagivanie_na_skamiu_so_shtangoj.gif
47	\N	zhim_lezha_so_shtangoj.gif	exercise\\zhim_lezha_so_shtangoj.gif
48	\N	zhim_shtangi_na_plechi.gif	exercise\\zhim_shtangi_na_plechi.gif
49	\N	zhim_shtangi_za_golovu.gif	exercise\\zhim_shtangi_za_golovu.gif
50	\N	zhim_uzkij_so_shtangoj.gif	exercise\\zhim_uzkij_so_shtangoj.gif
21	64	e7a55824b46fbe167c62d46953afc9ff8827c294f23dd1c27c35f388d716ad8d.jpg	img\\e7a55824b46fbe167c62d46953afc9ff8827c294f23dd1c27c35f388d716ad8d.jpg
\.


--
-- Name: images_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.images_id_seq', 50, true);


--
-- Name: files files_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.files
    ADD CONSTRAINT files_pk PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

