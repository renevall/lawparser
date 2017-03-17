
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS book (
    book_id integer NOT NULL,
    name text,
    law_id integer,
    reviewed boolean
);

CREATE SEQUENCE IF NOT EXISTS book_book_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE book_book_id_seq OWNER TO "Penshiru";

ALTER SEQUENCE book_book_id_seq OWNED BY book.book_id;

ALTER TABLE ONLY book ALTER COLUMN book_id SET DEFAULT nextval('book_book_id_seq'::regclass);

ALTER TABLE ONLY book
    ADD CONSTRAINT pk_book PRIMARY KEY (book_id);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE book;
DROP SEQUENCE IF EXISTS book_book_id_seq;