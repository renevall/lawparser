Testing User:

```
docker run -it --rm --link docker_db_1:postgres --net docker_default postgres psql -h postgres -U Penshiru -d Penshiru_Dev
```
Insert User

```
INSERT INTO "user"
 (first_name,
    last_name,
    email,
    address,
    contact_number,
    status_id,
    user_level,
    password,
    gender_id,
    pic_url,
    created_at,
    updated_at)
VALUEs(
    'test',
    'user',
    'email@test.com',
    '123 Fake Street',
    '1-800-fake',
    1,
    3,
    'JDJhJDEwJHZNTWc3cWFycDBwM1d6T0syRjllSGVqR013T0M2b1pZUWUzVlpKUWtqSGpEckdNaUVwYzBT',
    1,
    'https://dummyimage.com/300x200/000/fff.png&text=Test+Image',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
);
```