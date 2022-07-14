CREATE TABLE products( id INTEGER PRIMARY KEY AUTOINCREMENT,
                                              guid VARCHAR(255) UNIQUE NOT NULL,
                                                                       name VARCHAR(255) UNIQUE NOT NULL,
                                                                                                price REAL NOT NULL,
                                                                                                           description TEXT, createdAt TEXT NOT NULL);