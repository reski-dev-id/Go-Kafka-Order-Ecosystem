SELECT 'CREATE DATABASE order_db'
WHERE NOT EXISTS (
    SELECT FROM pg_database
    WHERE datname = 'order_db'
)\gexec

SELECT 'CREATE DATABASE payment_db'
WHERE NOT EXISTS (
    SELECT FROM pg_database
    WHERE datname = 'payment_db'
)\gexec

SELECT 'CREATE DATABASE notification_db'
WHERE NOT EXISTS (
    SELECT FROM pg_database
    WHERE datname = 'notification_db'
)\gexec