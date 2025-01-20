CREATE TABLE "pets"(
    "id" BIGSERIAL NOT NULL,
    "name" TEXT NOT NULL,
    "breed" TEXT NOT NULL,
    "host_id" BIGINT NOT NULL,
    "animal_type" TEXT NOT NULL
);
ALTER TABLE
    "pets" ADD PRIMARY KEY("id");
CREATE TABLE "host"(
    "id" BIGSERIAL NOT NULL,
    "name" TEXT NOT NULL,
    "surname" TEXT NOT NULL,
    "midname" TEXT NULL,
    "phone_number" TEXT NOT NULL,
    "email" TEXT NULL
);
ALTER TABLE
    "host" ADD PRIMARY KEY("id");
CREATE TABLE "services"(
    "id" BIGSERIAL NOT NULL,
    "type" TEXT NOT NULL,
    "price" DECIMAL(8, 2) NOT NULL,
    "duration" BIGINT NOT NULL
);
ALTER TABLE
    "services" ADD PRIMARY KEY("id");
CREATE TABLE "appointments"(
    "id" BIGSERIAL NOT NULL,
    "pet_id" BIGINT NOT NULL,
    "service_id" BIGINT NOT NULL,
    "time" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "salon_master_id" BIGINT NOT NULL,
    "review_id" BIGINT NULL,
    "status" VARCHAR(255) NOT NULL,
    "last_update" TIMESTAMP
);
ALTER TABLE
    "appointments" ADD PRIMARY KEY("id");
CREATE TABLE "groomers"(
    "id" BIGSERIAL NOT NULL,
    "name" TEXT NOT NULL,
    "surname" TEXT NOT NULL,
    "description" TEXT NULL
);
ALTER TABLE
    "groomers" ADD PRIMARY KEY("id");
CREATE TABLE "salons"(
    "id" BIGSERIAL NOT NULL,
    "address" TEXT NOT NULL,
    "phone_number" TEXT NOT NULL
);
ALTER TABLE
    "salons" ADD PRIMARY KEY("id");
CREATE TABLE "salon_masters"(
    "id" BIGSERIAL NOT NULL,
    "salon_id" BIGINT NOT NULL,
    "groomer_id" BIGINT NOT NULL
);
ALTER TABLE
    "salon_masters" ADD PRIMARY KEY("id");
CREATE TABLE "review"(
    "id" BIGSERIAL NOT NULL,
    "content" TEXT NOT NULL,
    "score" BIGINT NOT NULL
);
ALTER TABLE
    "review" ADD PRIMARY KEY("id");
CREATE TABLE "auth"(
    "user_id" BIGINT NOT NULL,
    "login" TEXT UNIQUE NOT NULL,
    "password_hash" TEXT NOT NULL,
    "role" TEXT NOT NULL
);
ALTER TABLE
    "appointments" ADD CONSTRAINT "appointments_service_id_foreign" FOREIGN KEY("service_id") REFERENCES "services"("id") ON DELETE CASCADE;
ALTER TABLE
    "salon_masters" ADD CONSTRAINT "salon_masters_groomer_id_foreign" FOREIGN KEY("groomer_id") REFERENCES "groomers"("id") ON DELETE CASCADE;
ALTER TABLE
    "appointments" ADD CONSTRAINT "appointments_review_id_foreign" FOREIGN KEY("review_id") REFERENCES "review"("id") ON DELETE CASCADE;
ALTER TABLE
    "appointments" ADD CONSTRAINT "appointments_salon_master_id_foreign" FOREIGN KEY("salon_master_id") REFERENCES "salon_masters"("id") ON DELETE CASCADE;
ALTER TABLE
    "appointments" ADD CONSTRAINT "appointments_pet_id_foreign" FOREIGN KEY("pet_id") REFERENCES "pets"("id") ON DELETE CASCADE;
ALTER TABLE
    "salon_masters" ADD CONSTRAINT "salon_masters_salon_id_foreign" FOREIGN KEY("salon_id") REFERENCES "salons"("id") ON DELETE CASCADE;
ALTER TABLE
    "pets" ADD CONSTRAINT "pets_host_id_foreign" FOREIGN KEY("host_id") REFERENCES "host"("id") ON DELETE CASCADE;



INSERT INTO services (type, price, duration) VALUES ('стрижка', 1000, 1), ('мытьё', 500, 1), ('и то, и другое', 1499.9, 2);

INSERT INTO salons (address, phone_number) VALUES
('г. Москва', '+70987654321'),
('г. Санкт-Петербург', '+71234567890');

CREATE OR REPLACE FUNCTION update_last_update_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.last_update := NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_last_update
BEFORE UPDATE ON appointments
FOR EACH ROW
EXECUTE FUNCTION update_last_update_column();
