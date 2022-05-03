CREATE OR REPLACE FUNCTION on_update_trigger() 
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW; 
END;
$$ language 'plpgsql';

CREATE TRIGGER update_table_modtime BEFORE UPDATE ON fiat_transactions FOR EACH ROW EXECUTE PROCEDURE on_update_trigger();
