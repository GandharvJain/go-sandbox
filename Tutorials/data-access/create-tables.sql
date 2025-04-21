DROP TABLE IF EXISTS album;

CREATE TABLE album (
	id	serial NOT NULL PRIMARY KEY,
	title	VARCHAR(128) NOT NULL,
	artist	VARCHAR(255) NOT NULL,
	price	DECIMAL(5, 2) NOT NULL
);

INSERT INTO album
	(title, artist,  price)
VALUES
	('Long Live', 'Taylor Swift', 10),
	('Work From Home', 'Fifth Harmony', 5),
	('Fireflies', 'Owl City', 8);
