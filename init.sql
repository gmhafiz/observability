CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create table uuid
(
    id integer primary key ,
    uuid UUID NOT NULL DEFAULT uuid_generate_v4()
);



do $$
    begin
        for i in 1..1000000 loop
                INSERT INTO uuid (id)  VALUES (i);
        end loop;
end; $$