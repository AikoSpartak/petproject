CREATE TABLE users (

    id BIGSERIAL PRIMARY KEY,

    full_name VARCHAR NOT NULL,

    email VARCHAR UNIQUE NOT NULL,

    phone_number VARCHAR(30) NOT NULL,

    date_birthday DATE NOT NULL,

    serial_passport VARCHAR NOT NULL,

    number_passport VARCHAR NOT NULL,

    history_operation TEXT,

    history_refund TEXT,

    history_moving TEXT

);



