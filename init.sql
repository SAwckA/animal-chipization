create table public.account (
    id serial primary key,
    firstName varchar(50),
    lastName varchar(50),
    email varchar(50) unique,
    password varchar(255)
);

create table public.location (
    id serial primary key,
    latitude double precision,
    longitude double precision

    -- FIXME: Unique lat + long
);

create table public.animal_type (
    id bigserial primary key,
    type varchar(255)

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
    chippingDateTime time,      -- Дата и время чипирования в формате
                                    -- ISO-8601 (устанавливается автоматически на момент добавления животного)

    chipperId int,              -- Идентификатор аккаунта чиппера
    chippingLocationId bigint,  -- // Идентификатор точки локации животных
    
    -- visitedLocations "[long]", // Массив идентификаторов объектов с
    -- информацией о посещенных точках локаций

    deathDateTime time          --  Дата и время смерти животного в
                                    -- формате ISO-8601 (устанавливается автоматически при смене lifeStatus на “DEAD”).
                                    -- Равняется null, пока lifeStatus = “ALIVE”.
);

-- TODO: 2 промежуточные таблицы many:many