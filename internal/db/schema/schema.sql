CREATE TABLE IF NOT EXISTS Applications (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS Packages (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS Versions (
    id SERIAL PRIMARY KEY,
    package_id INT NOT NULL,
    version VARCHAR(50) NOT NULL,
    FOREIGN KEY (package_id) REFERENCES Packages(id),
    UNIQUE (package_id, version)
);

CREATE TABLE IF NOT EXISTS Licenses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS ApplicationPackages (
    application_id INT NOT NULL,
    package_id INT NOT NULL,
    license_id INT NOT NULL,
    FOREIGN KEY (application_id) REFERENCES Applications(id),
    FOREIGN KEY (package_id) REFERENCES Packages(id),
    FOREIGN KEY (license_id) REFERENCES Licenses(id),
    PRIMARY KEY (application_id, package_id, license_id)
);

-- New table to store application versions and their SBOM URLs
CREATE TABLE IF NOT EXISTS ApplicationVersions (
    id SERIAL PRIMARY KEY,
    application_id INT NOT NULL,
    version VARCHAR(50) NOT NULL,
    sbom_url VARCHAR(255) NOT NULL,
    FOREIGN KEY (application_id) REFERENCES Applications(id),
    UNIQUE (application_id, version)
);
