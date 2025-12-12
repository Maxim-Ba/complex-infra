CREATE TABLE
  IF NOT EXISTS account (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    login TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW (),
    email TEXT UNIQUE,
    UNIQUE (login, password_hash),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW ()
  );

CREATE TABLE
  IF NOT EXISTS class (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    name TEXT NOT NULL UNIQUE,
    description TEXT
  );

CREATE TABLE
  IF NOT EXISTS character (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    account_id UUID REFERENCES account (id),
    class_id UUID REFERENCES class (id),
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW (),
    level INT NOT NULL DEFAULT 1,
    last_played_at TIMESTAMP DEFAULT NOW ()
  );



CREATE TABLE
  IF NOT EXISTS skill (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    name TEXT NOT NULL,
    description TEXT,
    required_level INT DEFAULT 1,
    mana_cost INT DEFAULT 0,
    cooldown INT DEFAULT 0,
    max_level INT DEFAULT 10
  );

CREATE TABLE
  IF NOT EXISTS item (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    item_type INT NOT NULL REFERENCES item_type (id),
    name TEXT NOT NULL,
    description TEXT,
    max_stack INT DEFAULT 1,
    slots_cost INT
  );

CREATE TABLE
  IF NOT EXISTS item_type (
    id INT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT
  );

-- Таблица связи персонаж-предметы
CREATE TABLE
  IF NOT EXISTS characters_item (
    character_id UUID NOT NULL REFERENCES character(id) ON DELETE CASCADE,
    item_id UUID NOT NULL REFERENCES item (id),
    PRIMARY KEY (character_id, item_id),
    is_equipped BOOLEAN DEFAULT FALSE,
    slot_position INT
  );



CREATE TABLE
  IF NOT EXISTS characters_equipment_slots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    character_id UUID NOT NULL REFERENCES character(id) ON DELETE CASCADE,
    slot_type_id UUID NOT NULL REFERENCES slot_type (id), 
    item_id UUID REFERENCES item(id),
    UNIQUE (character_id, slot_type_id)
  );

CREATE TABLE
  IF NOT EXISTS slot_type (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    name TEXT NOT NULL -- например: 'head', 'chest', 'weapon', 'ring1', 'ring2' 
    
  );

CREATE TABLE
  IF NOT EXISTS characteristic (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    agility INT DEFAULT 10,
    strength INT DEFAULT 10,
    intelligence INT DEFAULT 10,
    charisma INT DEFAULT 10,
    vitality INT DEFAULT 10,
    armor INT DEFAULT 0,
    magic_resist INT DEFAULT 0,
    health INT DEFAULT 100,
    mana INT DEFAULT 100,
    character_id UUID REFERENCES character(id)
  );

-- Таблица связи персонаж-навыки (изученные навыки)
CREATE TABLE
  IF NOT EXISTS characters_skill (
    character_id UUID NOT NULL REFERENCES character(id) ON DELETE CASCADE,
    skill_id UUID NOT NULL REFERENCES skill (id),
    learned_at TIMESTAMP DEFAULT NOW (),
    skill_level INT DEFAULT 1,
    slots_cost INT DEFAULT 1,
    PRIMARY KEY (character_id, skill_id)
  );

CREATE TABLE
  IF NOT EXISTS characters_skill_slots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    character_id UUID NOT NULL REFERENCES character(id) ON DELETE CASCADE,
    slot_number INT NOT NULL, -- номер слота (1, 2, 3, ...)
    skill_id UUID REFERENCES skill(id), 
    UNIQUE (character_id, slot_number)
  );


-- TRIGERS --

-- Для слотов экипировки
CREATE OR REPLACE FUNCTION check_equipment_item_exists()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.item_id IS NOT NULL AND NOT EXISTS (
        SELECT 1 FROM characters_item 
        WHERE character_id = NEW.character_id 
        AND item_id = NEW.item_id
    ) THEN
        RAISE EXCEPTION 'Item % does not belong to character %', NEW.item_id, NEW.character_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER validate_equipment_item
BEFORE INSERT OR UPDATE ON characters_equipment_slots
FOR EACH ROW EXECUTE FUNCTION check_equipment_item_exists();

-- Для слотов навыков
CREATE OR REPLACE FUNCTION check_skill_exists()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.skill_id IS NOT NULL AND NOT EXISTS (
        SELECT 1 FROM characters_skill 
        WHERE character_id = NEW.character_id 
        AND skill_id = NEW.skill_id
    ) THEN
        RAISE EXCEPTION 'Skill % is not learned by character %', NEW.skill_id, NEW.character_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER validate_skill_slot
BEFORE INSERT OR UPDATE ON characters_skill_slots
FOR EACH ROW EXECUTE FUNCTION check_skill_exists();
