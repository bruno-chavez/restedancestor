CREATE TABLE quotes (
                        id_quote INTEGER PRIMARY KEY,
                        content TEXT NOT NULL,
                        score INTEGER NOT NULL,
                        uuid TEXT NOT NULL UNIQUE
);
CREATE INDEX uuid on quotes (uuid);
CREATE TABLE indexes (
                         id_index INTEGER PRIMARY KEY,
                         word TEXT NOT NULL UNIQUE
);
CREATE TABLE indexes_quotes (
                                id_index INTEGER,
                                id_quote INTEGER
);
CREATE INDEX index_quote on indexes_quotes (id_index, id_quote);