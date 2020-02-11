CREATE DATABASE test
    WITH 
    OWNER = postgres
    ENCODING = 'UTF8';

CREATE SCHEMA test;

SET TRANSACTION ISOLATION LEVEL SERIALIZABLE;

CREATE TABLE test.wallets (
    ID int8 primary key,
    Balance float8 NOT NULL
);

INSERT INTO test.wallets(ID,Balance)
VALUES(1,0);

CREATE TABLE test.transactions (
    ID int8 primary key,
    Amount float8 NOT NULL
);

CREATE OR REPLACE FUNCTION "test".transfer(trId INT8,amount INT8)
    RETURNS void
    LANGUAGE plpgsql
AS $function$
DECLARE
    _tr_exist BOOLEAN;
BEGIN

    SELECT TRUE INTO _tr_exist
    FROM "test".transactions
    WHERE id = trId
    LIMIT 1;

    IF _tr_exist THEN
        ROLLBACK;
    END IF;

    INSERT INTO "test".transactions
    VALUES (trId,amount);

    UPDATE "test".wallets
    SET balance=balance+amount
    WHERE id=1;

END;
$function$