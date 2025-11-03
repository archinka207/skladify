ALTER TABLE WarehouseReceipts 
ADD COLUMN notes TEXT;

COMMENT ON COLUMN WarehouseReceipts.notes IS 'Дополнительные примечания к приходному ордеру';