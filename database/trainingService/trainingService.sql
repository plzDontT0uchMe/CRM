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

DROP DATABASE IF EXISTS "crm-trainingService";
--
-- Name: crm-trainingService; Type: DATABASE; Schema: -; Owner: postgres
--

CREATE DATABASE "crm-trainingService" WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'Russian_Russia.1251';


ALTER DATABASE "crm-trainingService" OWNER TO postgres;

\connect -reuse-previous=on "dbname='crm-trainingService'"

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
    id_program integer,
    id_user integer
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
-- Name: exercise; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.exercise (
    id integer NOT NULL,
    name character varying,
    description character varying,
    image character varying,
    video character varying
);


ALTER TABLE public.exercise OWNER TO postgres;

--
-- Name: exercise_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.exercise_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.exercise_id_seq OWNER TO postgres;

--
-- Name: exercise_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.exercise_id_seq OWNED BY public.exercise.id;


--
-- Name: muscles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.muscles (
    id integer NOT NULL,
    id_exercise integer,
    muscle character varying
);


ALTER TABLE public.muscles OWNER TO postgres;

--
-- Name: muscles_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.muscles_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.muscles_id_seq OWNER TO postgres;

--
-- Name: muscles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.muscles_id_seq OWNED BY public.muscles.id;


--
-- Name: programs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.programs (
    id integer NOT NULL,
    id_creator integer,
    name character varying,
    description character varying
);


ALTER TABLE public.programs OWNER TO postgres;

--
-- Name: programs_exercises; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.programs_exercises (
    id integer NOT NULL,
    id_program integer,
    id_exercise integer
);


ALTER TABLE public.programs_exercises OWNER TO postgres;

--
-- Name: programs_exercises_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.programs_exercises_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.programs_exercises_id_seq OWNER TO postgres;

--
-- Name: programs_exercises_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.programs_exercises_id_seq OWNED BY public.programs_exercises.id;


--
-- Name: programs_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.programs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.programs_id_seq OWNER TO postgres;

--
-- Name: programs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.programs_id_seq OWNED BY public.programs.id;


--
-- Name: active id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.active ALTER COLUMN id SET DEFAULT nextval('public.active_id_seq'::regclass);


--
-- Name: exercise id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.exercise ALTER COLUMN id SET DEFAULT nextval('public.exercise_id_seq'::regclass);


--
-- Name: muscles id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.muscles ALTER COLUMN id SET DEFAULT nextval('public.muscles_id_seq'::regclass);


--
-- Name: programs id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.programs ALTER COLUMN id SET DEFAULT nextval('public.programs_id_seq'::regclass);


--
-- Name: programs_exercises id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.programs_exercises ALTER COLUMN id SET DEFAULT nextval('public.programs_exercises_id_seq'::regclass);


--
-- Data for Name: active; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.active (id, id_program, id_user) FROM stdin;
1	1	44
3	3	54
2	2	54
4	4	54
5	5	54
6	6	54
\.


--
-- Data for Name: exercise; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.exercise (id, name, description, image, video) FROM stdin;
1	Приседания со штангой	Если вы тренируетесь в зале, то это упражнение со штангой лучше выполнять в силовой раме. Положите гриф на область трапеций. Упритесь пятками в пол и начинайте приседание. Подбородок приподнят, плечи развернуты, спина прямая. Отводите таз максимально назад, не сгибаясь в пояснице. Приседайте до параллели с полом, следите, чтобы носки находились на уровне коленей. Не подворачивайте таз при подъеме и держите колени немного согнутыми, чтобы снизить на них нагрузку.	prisedanie_so_shtangoj.gif	\N
2	Становая тяга	Штангу в этом упражнении берут с пола прямым или смешанным хватом. Подойдя к снаряду на максимальное расстояние, согните ноги в коленях и наклонитесь с прямой спиной, чтобы взять штангу с пола. Плечи должны быть развернуты, подбородок приподнят. Поднимите штангу, выпрямляя ноги, но не отклоняйтесь назад. Тяните снаряд до уровня бедер, не выпрямляя колени полностью, а затем опускайте на пол, сгибая колени и наклоняясь вперед.	stanovaya_tyaga.gif	\N
3	Мертвая тяга со штангой	Для выполнения возьмите штангу прямым хватом, стопы поставьте на ширине плеч, ноги немного согнуты в коленях. Теперь наклонитесь вперед с прямой спиной, можно немного прогнуться в пояснице. Опускайте штангу до уровня голеней, не следует наклоняться слишком низко. Во время упражнения вы должны чувствовать, как растягиваются мышцы задней поверхности бедер.	rumynskaya_tyga_shtanga.gif	\N
4	Выпады на месте со штангой	Выполнять выпады на месте можно со свободной штангой в силовой раме или в тренажере Смита с зафиксированным грифом. Для выполнения положите штангу на плечи и сделайте широкий шаг назад. Затем согните ноги в коленях, опускаясь в выпад. Приседайте до параллели переднего бедра с полом и следите, чтобы колено не выходило за носок. На протяжении упражнения спину держите прямо, а плечи развернутыми. Не забудьте повторить выпады для другой ноги.	vypad_na_meste_so_shtangoj.gif	\N
5	Жим штанги лежа	Жим лежа выполняется на прямой или наклонной скамье со стойками для грифа. Для выполнения лягте на скамью и возьмите гриф прямым закрытым хватом. Расстояние между ладонями должно быть чуть шире плеч. Опустите штангу к груди, разводя локти в стороны до параллели с полом. Затем выжмите гриф вверх, выпрямляя руки. Во время жима локти должны быть согнуты под прямым углом и смотреть в стороны, а не вниз.	zhim_lezha_so_shtangoj.gif	\N
6	Тяга к поясу	Перед выполнением упражнения положите гриф на стойку перед собой, а затем возьмите штангу прямым хватом. Согните ноги в коленях и подайтесь вперед, наклоняя спину. Приподнимите подбородок и слегка прогнитесь в спине. Теперь согните руки в локтях, притягивая штангу к поясу. Локти должны двигаться строго назад и сгибаться под прямым углом. Затем выпрямите руки, возвращаясь в начальное положение и снова повторите движение.	tyaga_shtangi_v_naklone.gif	\N
7	Армейский жим	Для выполнения возьмите штангу прямым хватом и согните руки в локтях, укладывая ее на плечи перед собой. Стойте прямо, плечи должны быть разведены, подбородок приподнят. Выжмите штангу вверх над головой, сохраняя локти немного согнутыми, затем вернитесь обратно. Во время жима не помогайте себе корпусом, он должен быть стабильным, работают только руки и плечи. В нижней точке сводите лопатки для максимальной нагрузки.	zhim_shtangi_na_plechi.gif	\N
8	Тяга к подбородку	Чтобы правильно выполнить тягу к подбородку, возьмите штангу прямым хватом и опустите руки вниз. Теперь тяните гриф вверх, сгибая руки в локтях и разводя их в стороны. В верхней точки локти образуют острый угол, так как поднимаются вверх. Тяните штангу до уровня плеч, а затем опускайте вниз. Снаряд должен «скользить» вдоль тела, двигаясь в одной плоскости.	tyaga_shtangi_k_podborodku.gif	\N
9	Сгибание рук со штангой	Для выполнения возьмите штангу обратным хватом и опустите прямые руки вниз, согнув их немного в локтях. Теперь сгибайте руки с полной амплитудой, приводя штангу к груди. В нижней точке оставляйте локти немного согнутыми, что поможет предотвратить травмы. Выполняйте упражнение медленно без резких движений, сосредоточившись на напряжении бицепсов рук.	sgibanie_ruk_so_shtangoj.gif	\N
10	Французский жим лежа	Чтобы выполнить упражнение, возьмите гриф прямым хватом и лягте на прямую скамью без стоек. Поднимите штангу над собой и согните руки в локтях, опуская гриф за голову. Не сгибайте локти более чем на 90 градусов, чтобы не травмировать плечи. Выполняйте упражнение осторожно в медленном темпе, акцентируя внимание на работе задней части рук. Для упражнения подойдет не только прямой, но и Z-гриф.	francuzsckij_zhim_lezha_so_shtangoj.gif	\N
11	Сумо-приседания	В зале сумо-приседания со штангой выполняют в силовой раме, как и другие виды приседов. Положите гриф на область трапеций. Затем отступите несколько шагов назад и поставьте ноги шире плеч, носки разведите в стороны. Теперь приседайте, отводя таз назад до параллели бедер с полом. Не забывайте держать спину ровно и сводить лопатки.	sumo_prisedanie_so_shtangoj.gif	\N
12	Фронтальные приседания	Для выполнения возьмите гриф прямым хватом и положите на плечи перед собой, сгибая руки в запястьях. Из этого положения согните ноги в коленях, выполняя приседания. Отводите таз назад до параллели с полом, при этом сохраняя легкий прогиб в пояснице. Не сводите плечи, подбородок держите приподнятым. Поднимаясь из приседа, не разгибайте колени полностью, что поможет уберечь их от чрезмерной нагрузки.	prisedanie_frontalnoe_so_shtangoj.gif	\N
13	Становая тяга сумо	Перед выполнением встаньте вплотную к штанге, которая лежит на полу перед вами. Поставьте ноги как можно шире и согните их в коленях, опускаясь в присед сумо. Наклонитесь с прямой спиной и возьмите гриф прямым хватом, затем выпрямите ноги, поднимая штангу до уровня бедер. Затем снова согните ноги в коленях, опуская снаряд вниз. При этом сохраняйте легкий прогиб в спине и не сводите плечи. В верхней точке не отклоняйтесь назад, чтобы не нагружать поясницу.	stanovaya_sumo_tyaga_shtangi.gif	\N
14	Выпады вперед	Для выполнения положите штангу на плечи и встаньте прямо, ноги на ширине плеч. Сделайте широкий шаг вперед и согните ноги в коленях, чтобы опуститься в выпад. Усилием ягодичных и бедренных мышц вернитесь обратно и повторите другой ногой. Выполняя упражнение, следите, чтобы ноги сгибались под прямым углом, что обеспечит стабильность коленным суставам.	vypad_vpered_so_shtangoj.gif	\N
15	Выпады назад	Чтобы выполнить выпады назад, положите штангу на плечи и встаньте прямо, ноги на ширине плеч. Теперь сделайте широкий шаг назад и согните ноги в коленях. Следите, чтобы переднее колено не выходило за носок. Усилием ягодичных вернитесь в начальное положение и повторите выпад другой ногой.	vypad_nazad_so_shtangoj.gif	\N
16	Болгарские сплит-приседания	Для упражнения понадобится прямая лавка высотой до колена и штанга. Положите гриф на плечи и встаньте спиной к лавке на расстоянии одного шага. Положите одну ногу на лавку сводом стопы вниз. Теперь согните в колене переднюю ногу до параллели с полом. Следите, чтобы колено не сгибалось под острым углом. Выполнив все повторения для одной ноги, сделайте сплит-приседания для другой.	vypad_bolgarskij_so_shtangoj.gif	\N
17	Тяга к поясу обратным хватом	Для выполнения возьмите гриф обратным хватом и поставьте ноги на ширине плеч. Немного согните ноги в коленях и наклоните спину вперед, прогибаясь в пояснице. Сведите лопатки и согните руки в локтях, притягивая штангу к поясу. Затем выпрямите руки, возвращаясь в исходное положение.	obratnaya_tyaga_shtangi_v_naklone.gif	\N
18	Жим к груди узким хватом	Выполняют жим лежа узким хватом на наклонной или прямой скамье с меньшим весом, чем в классической технике. Для выполнения лягте на лавку и возьмите гриф со стоек прямым хватом, поставив ладони на грифе уже ширины плеч. Постановка не должна быть слишком узкой, между ладонями должно сохраняться расстояние не менее 20-30 см. Далее выполняйте жим в классической технике, опуская штангу как можно ниже к груди. При этом локти не разводятся в стороны, а уходят немного вперед, формируя острый угол между плечом и предплечьем.	zhim_uzkij_so_shtangoj.gif	\N
19	Жим штанги из-за головы	Для выполнения возьмите штангу прямым хватом и поднимите ее вверх, как в армейском жиме, а затем опустите за голову. Из этого положения выжимайте гриф над головой, опуская его за затылок. В верхней точке не выпрямляйте локти полностью, а в нижней – следите за движением плечевых суставов и не выжимайте гриф до конца, пытаясь положить его на трапеции. В нижней точке локти должны быть согнуты под прямым углом, что обеспечит вам безопасность во время выполнения упражнения.	zhim_shtangi_za_golovu.gif	\N
20	Французский жим сидя	Упражнение можно выполнять сидя или стоя. Возьмите штангу прямым хватом и поднимите ее над головой на прямых руках. Локти должны быть немного согнуты, чтобы снизить нагрузку на суставы. Теперь медленно опустите гриф за голову, сгибая руки в локтях. Опускайте до прямого угла между плечом и предплечьем, так как слишком большая амплитуда движения может быть небезопасна. Акцентируйте внимание на работе трицепса и выполняйте упражнение очень медленно.	razgibanie_iz_za_golovy_so_shtangoj_sidya.gif	\N
21	Наклоны со штангой вперед	Для выполнения возьмите штангу со стойки и положите на трапеции. Поставьте ноги на ширине плеч и немного согните их в коленях. Теперь наклонитесь вперед, слегка прогибаясь в пояснице. Сводите лопатки, чтобы держать спину прямо. Наклоняйтесь до параллели корпуса с полом, ориентируясь на ощущения в задней части бедер. При правильном выполнении вы почувствуете, как растягиваются подколенные сухожилия.	good_morning_sto_shtangoj.gif	\N
22	Подъем штанги перед собой	Чтобы выполнить упражнение, возьмите гриф прямым хватом и опустите руки вниз. Не забывайте немного сгибать их в локтях, чтобы не нагружать чрезмерно суставы. Теперь поднимите руки вверх до параллели с полом, выполняя движение усилием плечевых мышц, а не корпуса или рук. Контролируемым движением опустите руки в исходное положение, сохраняя мышцы в напряжении, а затем снова повторите подъем.	podem_shtangi_pered_soboj.gif	\N
23	Шраги со штангой	Для выполнения возьмите штангу прямым хватом и опустите руки вниз. Теперь поднимите плечи, не выполняя никаких движений руками и корпусом. Фиксируйтесь в верхней точке на 1-2 секунды, чтобы усилить эффект. Затем опустите плечи вниз. Выполняйте упражнение с полной амплитудой, чтобы увеличить нагрузку на целевую мышцу.	shragi_so_shtangoj.gif	\N
24	Сгибание рук прямым хватом	Для упражнения возьмите штангу прямым хватом и опустите руки вниз, немного согнув их в локтях. Теперь согните руки с полной амплитудой, приводя штангу к груди. Во время движения не сгибайте руки в запястьях, что может быть травматично. Держите гриф прямым закрытым хватом и выполняйте движения медленно и размеренно. Для этой разновидности сгибаний возьмите меньший вес, чем для классических.	sgibanie_ruk_obratnyj_hvat_so_shtangoj.gif	\N
25	Ягодичный мостик со штангой	Для выполнения понадобится горизонтальная скамья. Лягте на скамью лопатками, согните ноги в коленях. Гриф положите на область сгиба бедер. В нижней точке вы не должны касаться ягодицами пола. Теперь поднимите таз вверх до параллели корпуса с полом. Держите гриф обеими руками, сильно не прогибайтесь в пояснице. Выполняйте движение за счет ягодичных, а не мышц спины. Сосредоточьтесь на работе целевых мышц, чтобы увеличить эффективность упражнения.	mostik_shtanga.gif	\N
26	Зашагивания на платформу	Для упражнения понадобится горизонтальная скамейка или платформа для прыжков высотой до колена или ниже. Положите штангу на плечи и встаньте лицом к скамейке. Теперь шагните на нее правой ногой, следом за ней приставляя и левую. Спускайтесь с противоположной ноги. Можно выполнить все подходы для одной ноги, а затем для другой, или чередовать ноги при подъеме. Во время зашагивания ступайте на платформу полной стопой. Не забывайте держать спину прямо, расправив плечи и приподняв подбородок.	zashagivanie_na_skamiu_so_shtangoj.gif	\N
27	Подъем на носки	Чтобы выполнить подъем на носки, вам понадобится степ-платформа или небольшое стабильное возвышение, например, порог. Положите штангу на плечи и поставьте носки на платформу. Ноги должны находиться на ширине плеч, спина прямая, плечи развернуты. Поднимите пятки, поднимаясь на носочки на платформе, а затем вернитесь обратно. Задерживайтесь в верхней точке, чтобы усилить нагрузку. Не выполняйте упражнение на нестабильной опоре, так как это может быть опасно.	podem_na_cypochki_so_shtangoj.gif	\N
\.


--
-- Data for Name: muscles; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.muscles (id, id_exercise, muscle) FROM stdin;
1	1	квадрицепсы
2	1	ягодицы
3	1	бицепсы бедер
4	1	разгибатели спины
5	1	икроножные мышцы
6	2	широчайшие
7	2	ромбовидные
8	2	трапеции
9	2	разгибатели позвоночника
10	2	бицепсы и квадрицепсы бедер
11	2	ягодичные мышцы
12	3	стабилизаторы 
13	3	разгибатели позвоночника
14	3	мышцы кора
15	3	бицепсы бедер
16	3	ягодицы
17	4	ягодичные мышцы
18	4	бицепсы и квадрицепсы бедер
19	4	икроножные
20	4	мышцы кора
21	4	пресс
\.


--
-- Data for Name: programs; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.programs (id, id_creator, name, description) FROM stdin;
1	33	value2	value3
2	54	123	321
3	54	123	321
4	54	test program	await getPrograms()
5	54	test program v2	test program description v2
6	54	eahresres	werahewarhwaerh
\.


--
-- Data for Name: programs_exercises; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.programs_exercises (id, id_program, id_exercise) FROM stdin;
1	1	1
2	1	2
3	1	3
6	3	1
7	3	2
4	2	1
5	2	2
9	4	4
8	4	3
10	5	10
11	5	1
12	5	16
13	6	6
14	6	7
\.


--
-- Name: active_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.active_id_seq', 6, true);


--
-- Name: exercise_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.exercise_id_seq', 27, true);


--
-- Name: muscles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.muscles_id_seq', 21, true);


--
-- Name: programs_exercises_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.programs_exercises_id_seq', 14, true);


--
-- Name: programs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.programs_id_seq', 6, true);


--
-- Name: active active_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.active
    ADD CONSTRAINT active_pk PRIMARY KEY (id);


--
-- Name: exercise exercise_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.exercise
    ADD CONSTRAINT exercise_pk PRIMARY KEY (id);


--
-- Name: muscles muscles_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.muscles
    ADD CONSTRAINT muscles_pk PRIMARY KEY (id);


--
-- Name: programs_exercises programs_exercises_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.programs_exercises
    ADD CONSTRAINT programs_exercises_pk PRIMARY KEY (id);


--
-- Name: programs programs_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.programs
    ADD CONSTRAINT programs_pk PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

