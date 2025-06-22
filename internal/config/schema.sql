CREATE TABLE kvstore (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL
);
CREATE TABLE history (
    word text NOT NULL,
    exp text NOT NULL,
    lang text NOT NULL,
    translator text NOT NULL,
    context text NOT NULL
);
CREATE TABLE notes (
    note text PRIMARY KEY,
    comment text NOT NULL,
    occurrence int NOT NULL,
    context text NOT NULL,
    CHECK(note <> '')
);
CREATE INDEX idx_history_word ON history(word);
