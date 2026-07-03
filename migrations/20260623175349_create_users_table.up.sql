CREATE TABLE users (

    id uuid PRIMARY KEY,

    full_name VARCHAR(100) NOT NULL,

    email VARCHAR(100) UNIQUE NOT NULL,

    phone_number VARCHAR(30) NOT NULL,

    date_birthday DATE NOT NULL,

    passport UNIQUE VARCHAR(12) NOT NULL,

    password VARCHAR(256) NOT NULL,

    role UNIQUE VARCHAR(15) NOT NULL,

    is_banned BOOLEAN NOT NULL DEFAULT FALSE


);



