CREATE TYPE status_type AS ENUM ('wait', 'error', 'success', 'new'); 

CREATE TABLE users (
 	id         serial        PRIMARY KEY,
	name       VARCHAR (355) NOT NULL,
	password   VARCHAR (50)  NOT NULL,
	created_at TIMESTAMP     DEFAULT now() NOT NULL,
	last_login TIMESTAMP,
	
	UNIQUE (name)
);


CREATE TABLE channels(
	id          serial      PRIMARY KEY,
	rss_url     TEXT        NOT NULL,
	description TEXT, 
	pub_date    TIMESTAMP   NOT NULL,
	parsed_at   TIMESTAMP,
	status      status_type,

	UNIQUE (rss_url)
);
CREATE TABLE user_channels(
	user_id    INTEGER NOT NULL,
	channel_id INTEGER NOT NULL,

	PRIMARY KEY (user_id, channel_id),
	
	FOREIGN KEY (user_id)          REFERENCES users(id)       ON DELETE CASCADE,
	FOREIGN KEY (channel_id)       REFERENCES channels(id)    ON DELETE CASCADE
);

CREATE TABLE channel_content(
	channel_id  INTEGER   NOT NULL,
	link        TEXT,
	title       TEXT,
	description TEXT      NOT NULL,
	pub_date    TIMESTAMP NOT NULL,

	PRIMARY KEY (channel_id),

	FOREIGN KEY (channel_id)       REFERENCES channels(id)    ON DELETE CASCADE
);

