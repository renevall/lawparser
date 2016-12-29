--
-- PostgreSQL database dump
--

-- Dumped from database version 9.5.4
-- Dumped by pg_dump version 9.5.1

-- Started on 2016-09-06 16:25:12 CST

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 1 (class 3079 OID 12361)
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- TOC entry 2150 (class 0 OID 0)
-- Dependencies: 1
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- TOC entry 186 (class 1259 OID 16424)
-- Name: article; Type: TABLE; Schema: public; Owner: Penshiru
--

CREATE TABLE article (
    article_id integer NOT NULL,
    name text,
    text text,
    chapter_id integer,
    law_id integer,
    reviewed boolean
);


ALTER TABLE article OWNER TO "Penshiru";

--
-- TOC entry 189 (class 1259 OID 16441)
-- Name: article_article_id_seq; Type: SEQUENCE; Schema: public; Owner: Penshiru
--

CREATE SEQUENCE article_article_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE article_article_id_seq OWNER TO "Penshiru";

--
-- TOC entry 2151 (class 0 OID 0)
-- Dependencies: 189
-- Name: article_article_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: Penshiru
--

ALTER SEQUENCE article_article_id_seq OWNED BY article.article_id;


--
-- TOC entry 190 (class 1259 OID 16452)
-- Name: chapter_chapter_id_seq; Type: SEQUENCE; Schema: public; Owner: Penshiru
--

CREATE SEQUENCE chapter_chapter_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE chapter_chapter_id_seq OWNER TO "Penshiru";

--
-- TOC entry 185 (class 1259 OID 16416)
-- Name: chapter; Type: TABLE; Schema: public; Owner: Penshiru
--

CREATE TABLE chapter (
    chapter_id integer DEFAULT nextval('chapter_chapter_id_seq'::regclass) NOT NULL,
    name text,
    title_id integer,
    law_id integer,
    reviewed boolean
);


ALTER TABLE chapter OWNER TO "Penshiru";

--
-- TOC entry 182 (class 1259 OID 16396)
-- Name: law; Type: TABLE; Schema: public; Owner: Penshiru
--

CREATE TABLE law (
    law_id integer NOT NULL,
    name text,
    approval_date date,
    publish_date date,
    journal text,
    intro text,
    reviewed boolean,
    revision integer
);


ALTER TABLE law OWNER TO "Penshiru";

--
-- TOC entry 181 (class 1259 OID 16394)
-- Name: law_law_id_seq; Type: SEQUENCE; Schema: public; Owner: Penshiru
--

CREATE SEQUENCE law_law_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE law_law_id_seq OWNER TO "Penshiru";

--
-- TOC entry 2152 (class 0 OID 0)
-- Dependencies: 181
-- Name: law_law_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: Penshiru
--

ALTER SEQUENCE law_law_id_seq OWNED BY law.law_id;


--
-- TOC entry 184 (class 1259 OID 16407)
-- Name: title; Type: TABLE; Schema: public; Owner: Penshiru
--

CREATE TABLE title (
    title_id integer NOT NULL,
    name text,
    law_id integer,
    reviewed boolean
);


ALTER TABLE title OWNER TO "Penshiru";

--
-- TOC entry 183 (class 1259 OID 16405)
-- Name: title_title_id_seq; Type: SEQUENCE; Schema: public; Owner: Penshiru
--

CREATE SEQUENCE title_title_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE title_title_id_seq OWNER TO "Penshiru";

--
-- TOC entry 2153 (class 0 OID 0)
-- Dependencies: 183
-- Name: title_title_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: Penshiru
--

ALTER SEQUENCE title_title_id_seq OWNED BY title.title_id;


--
-- TOC entry 187 (class 1259 OID 16427)
-- Name: user; Type: TABLE; Schema: public; Owner: Penshiru
--

CREATE TABLE "user" (
    user_id integer NOT NULL,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone,
    first_name text,
    last_name text,
    email text,
    address text,
    contact_number text,
    status_id integer,
    user_level integer,
    password text,
    gender_id integer,
    pic_url text
);


ALTER TABLE "user" OWNER TO "Penshiru";

--
-- TOC entry 188 (class 1259 OID 16430)
-- Name: user_user_id_seq; Type: SEQUENCE; Schema: public; Owner: Penshiru
--

CREATE SEQUENCE user_user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE user_user_id_seq OWNER TO "Penshiru";

--
-- TOC entry 2154 (class 0 OID 0)
-- Dependencies: 188
-- Name: user_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: Penshiru
--

ALTER SEQUENCE user_user_id_seq OWNED BY "user".user_id;


--
-- TOC entry 2017 (class 2604 OID 16443)
-- Name: article_id; Type: DEFAULT; Schema: public; Owner: Penshiru
--

ALTER TABLE ONLY article ALTER COLUMN article_id SET DEFAULT nextval('article_article_id_seq'::regclass);


--
-- TOC entry 2014 (class 2604 OID 16399)
-- Name: law_id; Type: DEFAULT; Schema: public; Owner: Penshiru
--

ALTER TABLE ONLY law ALTER COLUMN law_id SET DEFAULT nextval('law_law_id_seq'::regclass);


--
-- TOC entry 2015 (class 2604 OID 16410)
-- Name: title_id; Type: DEFAULT; Schema: public; Owner: Penshiru
--

ALTER TABLE ONLY title ALTER COLUMN title_id SET DEFAULT nextval('title_title_id_seq'::regclass);


--
-- TOC entry 2018 (class 2604 OID 16432)
-- Name: user_id; Type: DEFAULT; Schema: public; Owner: Penshiru
--

ALTER TABLE ONLY "user" ALTER COLUMN user_id SET DEFAULT nextval('user_user_id_seq'::regclass);


--
-- TOC entry 2026 (class 2606 OID 16451)
-- Name: pk_article; Type: CONSTRAINT; Schema: public; Owner: Penshiru
--

ALTER TABLE ONLY article
    ADD CONSTRAINT pk_article PRIMARY KEY (article_id);


--
-- TOC entry 2024 (class 2606 OID 16423)
-- Name: pk_chapter; Type: CONSTRAINT; Schema: public; Owner: Penshiru
--

ALTER TABLE ONLY chapter
    ADD CONSTRAINT pk_chapter PRIMARY KEY (chapter_id);


--
-- TOC entry 2020 (class 2606 OID 16404)
-- Name: pk_law; Type: CONSTRAINT; Schema: public; Owner: Penshiru
--

ALTER TABLE ONLY law
    ADD CONSTRAINT pk_law PRIMARY KEY (law_id);


--
-- TOC entry 2022 (class 2606 OID 16415)
-- Name: pk_title; Type: CONSTRAINT; Schema: public; Owner: Penshiru
--

ALTER TABLE ONLY title
    ADD CONSTRAINT pk_title PRIMARY KEY (title_id);


--
-- TOC entry 2028 (class 2606 OID 16440)
-- Name: pk_user; Type: CONSTRAINT; Schema: public; Owner: Penshiru
--

ALTER TABLE ONLY "user"
    ADD CONSTRAINT pk_user PRIMARY KEY (user_id);


--
-- TOC entry 2149 (class 0 OID 0)
-- Dependencies: 6
-- Name: public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM postgres;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;


-- Completed on 2016-09-06 16:25:13 CST

--
-- PostgreSQL database dump complete
--

