-- Table users
CREATE TABLE users (
  user_id  UUID NOT NULL,
  username VARCHAR(50) NOT NULL,
  password VARCHAR(100) NOT NULL
);

ALTER TABLE users
ADD CONSTRAINT pk_users PRIMARY KEY (user_id);

CREATE UNIQUE INDEX users_idx_username
ON users (username);


-- Create admin user
INSERT INTO users (user_id, username, password)
  values (
        'eeeeeb80-33f3-4d3f-befe-58694d2ac841',
        'admin',
        '$2a$10$NuzYobDOSTCx/EKBClGwGe0A9c8/yC7D4IP75hwz1jn.RCBfdEtb2' -- adm1n
  );
