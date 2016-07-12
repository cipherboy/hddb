CREATE TABLE IF NOT EXISTS files (
	size INTEGER,
	md5 STRING,
	sha1 STRING,
	sha256 STRING,
	tiger STRING,
	whirlpool STRING,
	path STRING,
	scandate DATE,
	ignore BOOLEAN
);
