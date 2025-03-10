CREATE TABLE IF NOT EXISTS files (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  file_path TEXT NOT NULL,
  description TEXT NOT NULL,
  key TEXT NOT NULL,
  salt TEXT NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER IF NOT EXISTS update_files_updated_at
AFTER UPDATE ON files
FOR EACH ROW
BEGIN
  UPDATE files SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;
