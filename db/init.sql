-- =========================
-- TABLES
-- =========================

CREATE TABLE IF NOT EXISTS authors (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS customers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    author_id INT,
    price DECIMAL(6,2),
    stock INT DEFAULT 0,
    FOREIGN KEY (author_id) REFERENCES authors(id)
);

CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    customer_id INT,
    created_at TIMESTAMP DEFAULT now(),
    FOREIGN KEY (customer_id) REFERENCES customers(id)
);

CREATE TABLE IF NOT EXISTS order_items (
    id SERIAL PRIMARY KEY,
    order_id INT,
    book_id INT,
    quantity INT NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders(id),
    FOREIGN KEY (book_id) REFERENCES books(id)
);

-- =========================
-- DATA (SEEDING)
-- =========================

-- AUTHORS
INSERT INTO authors (name) VALUES
('J. K. Rowling'),
('George R. R. Martin'),
('J. R. R. Tolkien'),
('Frank Herbert'),
('Isaac Asimov'),
('Brandon Sanderson');

-- CUSTOMERS
INSERT INTO customers (name, email) VALUES
('Alice', 'alice@example.com'),
('Bob', 'bob@example.com'),
('Charlie', 'charlie@example.com'),
('Diana', 'diana@example.com'),
('Eve', 'eve@example.com');

-- BOOKS
INSERT INTO books (title, author_id, price, stock) VALUES
('Harry Potter and the Philosopher''s Stone', 1, 19.99, 12),
('Harry Potter and the Chamber of Secrets', 1, 20.99, 8),
('Game of Thrones', 2, 25.50, 5),
('A Clash of Kings', 2, 26.50, 4),
('The Hobbit', 3, 18.99, 10),
('The Lord of the Rings', 3, 29.99, 7),
('Dune', 4, 24.99, 6),
('Foundation', 5, 22.50, 9),
('Mistborn', 6, 21.99, 11);

-- ORDERS
INSERT INTO orders (customer_id) VALUES
(1),
(2),
(3),
(4),
(5),
(1),
(3);

-- ORDER ITEMS
INSERT INTO order_items (order_id, book_id, quantity) VALUES
(1, 1, 2),
(1, 3, 1),
(2, 2, 1),
(2, 5, 2),
(3, 6, 1),
(3, 7, 1),
(4, 8, 3),
(5, 9, 1),
(6, 4, 2),
(7, 1, 1),
(7, 6, 1);