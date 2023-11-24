<!-- Drop table -->

```sql
DROP TABLE IF EXISTS users;
```

<!-- Create table -->

```sql
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  age INT,
  first_name TEXT,
  last_name TEXT,
  email TEXT UNIQUE NOT NULL
);
```

<!-- Insert values -->

```sql
INSERT INTO users (age, email, first_name, last_name)
VALUES (30, 'jon@calhoun.io', 'Jonathan', 'Calhoun');
```

<!-- Querying data -->

```sql
SELECT * FROM users
WHERE age < 24 OR last_name = 'Calhoun';
```

<!-- Update data -->

```sql
UPDATE users
SET first_name = 'Johnny', last_name = 'Appleseed'
WHERE id = 1;
```

<!-- Delete data -->

```sql
DELETE FROM users
WHERE id = 1;
```
