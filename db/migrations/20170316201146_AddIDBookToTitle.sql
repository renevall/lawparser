
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE title
ADD COLUMN book_id integer;

ALTER TABLE title
ADD CONSTRAINT fk_book_id
FOREIGN KEY(book_id) REFERENCES book(book_id);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE title
DROP CONSTRAINT fk_book_id;

ALTER TABLE title
DROP COLUMN book_id;
