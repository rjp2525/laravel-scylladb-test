CREATE KEYSPACE IF NOT EXISTS test_keyspace
WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 3}
AND durable_writes = true;

USE test_keyspace;
