CREATE TABLE kvstore (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL
);
CREATE TABLE dictionary (
    word text PRIMARY KEY,
    exp text NOT NULL,
    translator text NOT NULL
);
CREATE INDEX idx_dictionary_word ON dictionary(word);
