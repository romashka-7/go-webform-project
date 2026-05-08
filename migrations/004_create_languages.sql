CREATE TABLE IF NOT EXISTS languages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS application_languages (
    application_id INT NOT NULL,
    language_id INT NOT NULL,

    PRIMARY KEY (application_id, language_id),

    FOREIGN KEY (application_id) REFERENCES applications(id) ON DELETE CASCADE,
    FOREIGN KEY (language_id) REFERENCES languages(id) ON DELETE CASCADE
);

INSERT IGNORE INTO languages (name) VALUES
('Pascal'),
('C'),
('C++'),
('JavaScript'),
('PHP'),
('Python'),
('Java'),
('Haskell'),
('Clojure'),
('Prolog'),
('Scala'),
('Go');