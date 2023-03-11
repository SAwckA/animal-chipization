create table public.account (
    id serial primary key,
    firstName varchar(255),
    lastName varchar(255),
    email varchar(255),
    password varchar(255),

    constraint account_email_key unique(email)
);

create table public.location (
    id serial primary key,
    latitude double precision,
    longitude double precision,

    constraint location_latitude_longitude_key unique(latitude, longitude)
);

create table public.animal_type (
    id bigserial primary key,
    type varchar(255),

    constraint animal_type_type_key unique(type)
);

create table public.animal (
    id bigserial primary key,
    weight real,
    length real,
    height real,
    gender varchar(6),
    lifeStatus varchar(5),
    chippingDateTime timestamptz,
    chipperId int references account(id),
    chippingLocationId int references location(id),
    deathDateTime timestamptz
);

create table animal_types_list (
    animal_id int references animal(id) on delete cascade,
    type_id int references animal_type(id),

    unique(animal_id, type_id)
);

create table animal_locations_list (
   id serial primary key,
   animal_id int references animal(id) on delete cascade,
   location_id int references location(id),
   date_time_of_visited_location_point timestamptz
);
