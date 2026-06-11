BEGIN TRANSACTION;
CREATE TABLE "schedule"(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date TEXT CHECK(date IS strftime('%Y%m%d', date)) NOT NULL DEFAULT (strftime('%Y%m%d', 'now')),
	title TEXT NOT NULL DEFAULT "",
	comment TEXT DEFAULT "",
	repeat VARCHAR(128)	
);
CREATE INDEX idx_schedule_date ON schedule(date);
COMMIT;