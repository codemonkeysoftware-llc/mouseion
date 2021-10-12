CREATE TABLE entry_tag (
    entry_id INTEGER,
    tag TEXT,
    FOREIGN KEY (entry_id) REFERENCES log_entry(id) 
            ON DELETE CASCADE ON UPDATE CASCADE,
    PRIMARY KEY(entry_id,tag)
);