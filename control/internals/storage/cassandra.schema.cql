CREATE KEYSPACE IF NOT EXISTS sentimental WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };

DROP TYPE sentimental.sentimental_result;
CREATE TYPE IF NOT EXISTS sentimental.sentimental_result (
    label text,
    score double
);

DROP TABLE sentimental.messages;
CREATE TABLE IF NOT EXISTS sentimental.messages
    (
    id uuid, 
    msg text, 
    msg_analized FROZEN<sentimental_result>, 
    created_at timestamp, 
    updated_at timestamp, 
    expires_at timestamp,
    PRIMARY KEY(id));

SELECT * FROM sentimental.messages;