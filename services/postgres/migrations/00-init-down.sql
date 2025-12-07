-- Удаляем триггеры, если они были созданы
DROP TRIGGER IF EXISTS validate_equipment_item ON character_equipment_slots;
DROP TRIGGER IF EXISTS validate_skill_slot ON character_skill_slots;
DROP FUNCTION IF EXISTS check_equipment_item_exists();
DROP FUNCTION IF EXISTS check_skill_exists();

-- Удаляем таблицы в обратном порядке (учитывая зависимости внешних ключей)
DROP TABLE IF EXISTS character_skill_slots;
DROP TABLE IF EXISTS character_skill;
DROP TABLE IF EXISTS characteristic;
DROP TABLE IF EXISTS character_equipment_slots;
DROP TABLE IF EXISTS character_item;
DROP TABLE IF EXISTS item;
DROP TABLE IF EXISTS item_type;
DROP TABLE IF EXISTS skill;
DROP TABLE IF EXISTS class;
DROP TABLE IF EXISTS character;
DROP TABLE IF EXISTS account;

-- Если были созданы хранимые процедуры, их тоже нужно удалить
DROP FUNCTION IF EXISTS equip_item_to_slot(UUID, UUID, TEXT);
DROP PROCEDURE IF EXISTS equip_item_to_slot(UUID, UUID, TEXT);
DROP FUNCTION IF EXISTS equip_item(UUID, TEXT, UUID);
