DROP TYPE IF EXISTS status_type; 

CREATE TABLE IF NOT EXISTS users (
 	id         serial        PRIMARY KEY,
	name       VARCHAR (355) NOT NULL,
	password   VARCHAR (50)  NOT NULL,
	created_at TIMESTAMP     DEFAULT now() NOT NULL,
	last_login TIMESTAMP,
	
	UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS channels(
	id          serial      PRIMARY KEY,
	rss_url     TEXT        NOT NULL,
	description TEXT,
        title       TEXT,
	link        TEXT,	
	pub_date    TIMESTAMP,
	created_at  TIMESTAMP   DEFAULT now() NOT NULL,

	UNIQUE (rss_url)
);

CREATE TABLE IF NOT EXISTS user_channels(
	user_id    INTEGER NOT NULL,
	channel_id INTEGER NOT NULL,

	PRIMARY KEY (user_id, channel_id),
	
	FOREIGN KEY (user_id)          REFERENCES users(id)       ON DELETE CASCADE,
	FOREIGN KEY (channel_id)       REFERENCES channels(id)    ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS channel_content(
	channel_id  INTEGER   NOT NULL,
	link        TEXT,
	title       TEXT,
	author      TEXT,
	category    TEXT,
	description TEXT      NOT NULL,
	pub_date    TIMESTAMP NOT NULL,

	FOREIGN KEY (channel_id)       REFERENCES channels(id)    ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS scheduler(
	channel_id  INTEGER     NOT NULL,
	rss_url     TEXT        NOT NULL,
	finish      TIMESTAMP,
	start       TIMESTAMP,
	plan_start  TIMESTAMP,
	status      VARCHAR(50),
	message     TEXT,
	
	UNIQUE (rss_url),

	PRIMARY KEY (channel_id),

	FOREIGN KEY (channel_id)       REFERENCES channels(id)    ON DELETE CASCADE
);


