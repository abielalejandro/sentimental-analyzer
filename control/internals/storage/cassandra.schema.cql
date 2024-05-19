CREATE KEYSPACE IF NOT EXISTS sentimental WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };

CREATE TYPE IF NOT EXISTS sentimental.sentimental_result (
    label text,
    score decimal
);

CREATE TABLE IF NOT EXISTS sentimental.messages
    (
    id uuid, 
    msg text, 
    msg_analized FROZEN<sentimental_result>, 
    created_at timestamp, 
    expires_at timestamp,
    PRIMARY KEY(id))
    WITH CLUSTERING ORDER BY (created_at DESC);

