create table public.account (
    id serial primary key,
    firstName varchar(255),
    lastName varchar(255),
    email varchar(255) unique,
    password varchar(255)
);

create table public.location (
    id serial primary key,
    latitude double precision,
    longitude double precision,

    constraint unique_location_points unique(latitude, longitude)

    -- FIXME: Unique lat + long
);

create table public.animal_type (
    id bigserial primary key,
    type varchar(255),

    constraint unique_type unique(type)
    -- FIXME: Unique type
);

create table public.animal (
    id bigserial primary key,
    
    -- animalTypes "[long]", Массив идентификаторов типов животного
    
    weight real,                -- Масса животного, кг
    length real,                -- Длина животного, м
    height real,                -- Высота животного, м
    gender varchar(6),          -- Гендерный признак животного, доступные значения “MALE”, “FEMALE”, “OTHER”
    lifeStatus varchar(5),      -- Жизненный статус животного, доступные значения
                                    -- “ALIVE”(устанавливается автоматически при добавлении нового животного), 
                                    -- “DEAD”(можно установить при обновлении информации о животном)
    chippingDateTime timestamptz,      -- Дата и время чипирования в формате
                                    -- ISO-8601 (устанавливается автоматически на момент добавления животного)

    chipperId int references account(id),              -- Идентификатор аккаунта чиппера
    chippingLocationId int references location(id),  -- // Идентификатор точки локации животных
    
    -- visitedLocations "[long]", // Массив идентификаторов объектов с
    -- информацией о посещенных точках локаций

    deathDateTime timestamptz          --  Дата и время смерти животного в
                                    -- формате ISO-8601 (устанавливается автоматически при смене lifeStatus на “DEAD”).
                                    -- Равняется null, пока lifeStatus = “ALIVE”.
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


-- TODO: 2 промежуточные таблицы many:many