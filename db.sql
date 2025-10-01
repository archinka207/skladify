
SET search_path TO warehouse_db;
-- Скрипт для создания таблиц базы данных "Складской Учет" для PostgreSQL

-- Удаляем таблицы, если они уже существуют, для чистого запуска скрипта.
-- Порядок удаления обратный порядку создания из-за зависимостей.
DROP TABLE IF EXISTS WarehouseReceipts;
DROP TABLE IF EXISTS Materials;
DROP TABLE IF EXISTS MaterialGroups;
DROP TABLE IF EXISTS MaterialClasses;
DROP TABLE IF EXISTS Suppliers;
DROP TABLE IF EXISTS UnitsOfMeasure;
DROP TABLE IF EXISTS DocumentTypes;

-- =================================================================
-- 1. Справочник: Классы материалов (MaterialClasses)
-- Самая верхняя категория для классификации материалов.
-- =================================================================
CREATE TABLE MaterialClasses (
    class_id SERIAL PRIMARY KEY,
    class_name VARCHAR(100) NOT NULL UNIQUE
);

COMMENT ON TABLE MaterialClasses IS 'Справочник: Классы материалов (например, "Строительные материалы")';
COMMENT ON COLUMN MaterialClasses.class_id IS 'Код класса (PK)';
COMMENT ON COLUMN MaterialClasses.class_name IS 'Наименование класса';

-- =================================================================
-- 2. Справочник: Группы материалов (MaterialGroups)
-- Подкатегория, связанная с классом.
-- =================================================================
CREATE TABLE MaterialGroups (
    group_id SERIAL PRIMARY KEY,
    group_name VARCHAR(100) NOT NULL,
    class_id INT NOT NULL,

    CONSTRAINT fk_materialgroups_class
        FOREIGN KEY (class_id)
        REFERENCES MaterialClasses(class_id)
        ON DELETE RESTRICT -- Запретить удаление класса, если у него есть группы
);

COMMENT ON TABLE MaterialGroups IS 'Справочник: Группы материалов (например, "Сухие смеси")';
COMMENT ON COLUMN MaterialGroups.group_id IS 'Код группы (PK)';
COMMENT ON COLUMN MaterialGroups.group_name IS 'Наименование группы';
COMMENT ON COLUMN MaterialGroups.class_id IS 'Внешний ключ к MaterialClasses';

-- =================================================================
-- 3. Справочник: Материалы (Materials)
-- Основной справочник номенклатуры.
-- =================================================================
CREATE TABLE Materials (
    material_id BIGSERIAL PRIMARY KEY,
    material_name VARCHAR(255) NOT NULL,
    group_id INT NOT NULL,

    CONSTRAINT fk_materials_group
        FOREIGN KEY (group_id)
        REFERENCES MaterialGroups(group_id)
        ON DELETE RESTRICT -- Запретить удаление группы, если в ней есть материалы
);

COMMENT ON TABLE Materials IS 'Справочник: Конкретные материалы или товары';
COMMENT ON COLUMN Materials.material_id IS 'Код материала (PK)';
COMMENT ON COLUMN Materials.material_name IS 'Наименование материала';
COMMENT ON COLUMN Materials.group_id IS 'Внешний ключ к MaterialGroups';

-- =================================================================
-- 4. Справочник: Поставщики (Suppliers)
-- Информация о компаниях-поставщиках.
-- =================================================================
CREATE TABLE Suppliers (
    supplier_id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    inn VARCHAR(12) NOT NULL UNIQUE,
    legal_zip_code VARCHAR(10),
    legal_city VARCHAR(100),
    legal_street_address VARCHAR(255),
    bank_zip_code VARCHAR(10),
    bank_city VARCHAR(100),
    bank_street_address VARCHAR(255),
    bank_account VARCHAR(20)
);

COMMENT ON TABLE Suppliers IS 'Справочник: Поставщики материалов';
COMMENT ON COLUMN Suppliers.supplier_id IS 'Код поставщика (PK)';
COMMENT ON COLUMN Suppliers.name IS 'Наименование организации';
COMMENT ON COLUMN Suppliers.inn IS 'ИНН (уникальный)';

-- =================================================================
-- 5. Справочник: Единицы измерения (UnitsOfMeasure)
-- =================================================================
CREATE TABLE UnitsOfMeasure (
    unit_id SERIAL PRIMARY KEY,
    unit_name VARCHAR(50) NOT NULL UNIQUE,
    abbreviation VARCHAR(10) NOT NULL UNIQUE
);

COMMENT ON TABLE UnitsOfMeasure IS 'Справочник: Единицы измерения (кг, шт, м.п.)';
COMMENT ON COLUMN UnitsOfMeasure.unit_id IS 'Код единицы измерения (PK)';
COMMENT ON COLUMN UnitsOfMeasure.unit_name IS 'Полное наименование';
COMMENT ON COLUMN UnitsOfMeasure.abbreviation IS 'Краткое обозначение';

-- =================================================================
-- 6. Справочник: Типы сопроводительных документов (DocumentTypes)
-- =================================================================
CREATE TABLE DocumentTypes (
    doc_type_id SERIAL PRIMARY KEY,
    doc_type_name VARCHAR(100) NOT NULL UNIQUE
);

COMMENT ON TABLE DocumentTypes IS 'Справочник: Типы сопроводительных документов (УПД, ТОРГ-12)';
COMMENT ON COLUMN DocumentTypes.doc_type_id IS 'Код типа документа (PK)';
COMMENT ON COLUMN DocumentTypes.doc_type_name IS 'Наименование типа документа';

-- =================================================================
-- 7. Основная таблица: Приходные ордера (WarehouseReceipts)
-- Транзакционная таблица, фиксирующая каждую операцию поступления.
-- =================================================================
CREATE TABLE WarehouseReceipts (
    receipt_id BIGSERIAL PRIMARY KEY,
    order_number VARCHAR(50) NOT NULL,
    receipt_date DATE NOT NULL,
    supplier_id BIGINT NOT NULL,
    balance_account VARCHAR(20) NOT NULL,
    doc_type_id INT NOT NULL,
    document_number VARCHAR(50) NOT NULL,
    material_id BIGINT NOT NULL,
    material_account VARCHAR(20) NOT NULL,
    unit_id INT NOT NULL,
    quantity DECIMAL(18, 3) NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(18, 2) NOT NULL CHECK (unit_price >= 0),

    -- Внешние ключи
    CONSTRAINT fk_receipts_supplier
        FOREIGN KEY (supplier_id) REFERENCES Suppliers(supplier_id) ON DELETE RESTRICT,
    CONSTRAINT fk_receipts_doctype
        FOREIGN KEY (doc_type_id) REFERENCES DocumentTypes(doc_type_id) ON DELETE RESTRICT,
    CONSTRAINT fk_receipts_material
        FOREIGN KEY (material_id) REFERENCES Materials(material_id) ON DELETE RESTRICT,
    CONSTRAINT fk_receipts_unit
        FOREIGN KEY (unit_id) REFERENCES UnitsOfMeasure(unit_id) ON DELETE RESTRICT
);

COMMENT ON TABLE WarehouseReceipts IS 'Основная таблица: Приходные ордера (поступления на склад)';
COMMENT ON COLUMN WarehouseReceipts.receipt_id IS 'ID записи о поступлении (PK)';
COMMENT ON COLUMN WarehouseReceipts.order_number IS 'Внутренний номер приходного ордера';
COMMENT ON COLUMN WarehouseReceipts.receipt_date IS 'Дата поступления материала на склад';
COMMENT ON COLUMN WarehouseReceipts.supplier_id IS 'Код поставщика (FK)';
COMMENT ON COLUMN WarehouseReceipts.quantity IS 'Количество пришедшего материала (не может быть <= 0)';
COMMENT ON COLUMN WarehouseReceipts.unit_price IS 'Цена за единицу (не может быть отрицательной)';

-- =================================================================
-- Создание индексов для внешних ключей
-- Это ускоряет операции JOIN и выборки данных.
-- =================================================================
CREATE INDEX idx_materialgroups_class_id ON MaterialGroups(class_id);
CREATE INDEX idx_materials_group_id ON Materials(group_id);
CREATE INDEX idx_warehousereceipts_supplier_id ON WarehouseReceipts(supplier_id);
CREATE INDEX idx_warehousereceipts_material_id ON WarehouseReceipts(material_id);
CREATE INDEX idx_warehousereceipts_receipt_date ON WarehouseReceipts(receipt_date);
