CREATE DATABASE test
    WITH
    OWNER = postgres
    ENCODING = 'UTF8';

CREATE SCHEMA test;

CREATE TABLE test.wallets (
                              ID int8 primary key,
                              Balance float8 NOT NULL
);

INSERT INTO test.wallets(ID,Balance)
VALUES(1,0);

CREATE TABLE test.transactions (
                                   ID TEXT primary key,
                                   State TEXT NOT NULL,
                                   Amount float8 NOT NULL
);

CREATE OR REPLACE FUNCTION "test".transfer(trId TEXT,_state TEXT,amount float8)
    RETURNS void
    LANGUAGE plpgsql
AS $function$
DECLARE
    _tr_exist BOOLEAN;
    _balance float8;
BEGIN

    SELECT TRUE INTO _tr_exist
    FROM "test".transactions
    WHERE id = trId
    LIMIT 1;


    IF _tr_exist THEN
        ROLLBACK;
    END IF;

    INSERT INTO "test".transactions
    VALUES (trId,_state,amount);

    SELECT balance INTO _balance
    FROM "test".wallets
    WHERE id = 1
    LIMIT 1;

    IF amount>0 THEN
        UPDATE "test".wallets
        SET balance=balance+amount
        WHERE id=1;
    ELSE
        IF _balance>=amount THEN
            UPDATE "test".wallets
            SET balance=balance-amount
            WHERE id=1;
        ELSE
            ROLLBACK;
        END IF;
    END IF;

END;
$function$