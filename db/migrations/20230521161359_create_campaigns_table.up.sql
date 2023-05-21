CREATE TABLE campaigns (
  id INT AUTO_INCREMENT PRIMARY KEY,
  user_id INT,
  name VARCHAR(255),
  short_description VARCHAR(255),
  description TEXT,
  goal_amount INT,
  current_amount INT,
  perks TEXT,
  backer_count INT,
  slug VARCHAR(255),
  created_at DATETIME,
  updated_at DATETIME,
  FOREIGN KEY (user_id) REFERENCES users(id)
);
