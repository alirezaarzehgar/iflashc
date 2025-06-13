CREATE TABLE kvstore (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL
);
CREATE TABLE dictionary (
    word text NOT NULL,
    exp text NOT NULL,
    lang text NOT NULL,
    translator text NOT NULL,
    context text NOT NULL
);
CREATE INDEX idx_dictionary_word ON dictionary(word);
