CREATE TABLE books (
   id SERIAL PRIMARY KEY,
   title VARCHAR(255) NOT NULL,
   isbn VARCHAR(20) UNIQUE,
   author_id INT NOT NULL,
   category_id INT NOT NULL,
   stock INT NOT NULL DEFAULT 0,
   published_date DATE,
   description TEXT,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE borrowing_records (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    book_id INT NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    borrowed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    returned_at TIMESTAMP,
    due_date TIMESTAMP,
);

CREATE TABLE book_recommendation (
    id SERIAL PRIMARY KEY,
    book_id INT NOT NULL,
    recommended_book_id INT NOT NULL,
    score FLOAT DEFAULT 1.0,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
    FOREIGN KEY (recommended_book_id) REFERENCES books(id) ON DELETE CASCADE
)