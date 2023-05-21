CREATE TABLE transactions (
  id INT AUTO_INCREMENT PRIMARY KEY,
  campaign_id INT,
  user_id INT,
  amount INT,
  status VARCHAR(255),
  code VARCHAR(255),
  created_at DATETIME,
  updated_at DATETIME,
  FOREIGN KEY (campaign_id) REFERENCES campaigns(id),
  FOREIGN KEY (user_id) REFERENCES users(id)
);
