CREATE TABLE campaign_images (
  id INT AUTO_INCREMENT PRIMARY KEY,
  campaign_id INT,
  file_name VARCHAR(255),
  is_primary TINYINT(1),
  created_at DATETIME,
  updated_at DATETIME,
  FOREIGN KEY (campaign_id) REFERENCES campaigns(id)
);