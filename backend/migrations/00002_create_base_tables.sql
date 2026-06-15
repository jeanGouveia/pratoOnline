-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS users (
    id            INTEGER  PRIMARY KEY AUTOINCREMENT,
    name          TEXT     NOT NULL,
    email         TEXT     NOT NULL UNIQUE,
    password_hash TEXT     NOT NULL,
    created_at    DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at    DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

CREATE TABLE IF NOT EXISTS products (
    id             INTEGER  PRIMARY KEY AUTOINCREMENT,
    name           TEXT     NOT NULL,
    price          REAL     NOT NULL DEFAULT 0,
    is_composto    INTEGER  NOT NULL DEFAULT 0 CHECK(is_composto IN (0,1)),
    stock_quantity REAL     NOT NULL DEFAULT 0,
    created_at     DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at     DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS ingredients (
    id             INTEGER  PRIMARY KEY AUTOINCREMENT,
    name           TEXT     NOT NULL,
    unit           TEXT     NOT NULL DEFAULT 'un',
    stock_quantity REAL     NOT NULL DEFAULT 0,
    created_at     DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at     DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS product_ingredients (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    product_id    INTEGER NOT NULL REFERENCES products(id)    ON DELETE CASCADE,
    ingredient_id INTEGER NOT NULL REFERENCES ingredients(id) ON DELETE RESTRICT,
    quantity      REAL    NOT NULL CHECK(quantity > 0),
    UNIQUE(product_id, ingredient_id)
);

CREATE INDEX IF NOT EXISTS idx_product_ingredients_product ON product_ingredients(product_id);

CREATE TABLE IF NOT EXISTS product_compositions (
    id                   INTEGER PRIMARY KEY AUTOINCREMENT,
    parent_product_id    INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    component_product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    quantity             REAL    NOT NULL CHECK(quantity > 0),
    UNIQUE(parent_product_id, component_product_id),
    CHECK(parent_product_id != component_product_id)
);

CREATE INDEX IF NOT EXISTS idx_compositions_parent ON product_compositions(parent_product_id);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_compositions_parent;
DROP TABLE IF EXISTS product_compositions;
DROP INDEX IF EXISTS idx_product_ingredients_product;
DROP TABLE IF EXISTS product_ingredients;
DROP TABLE IF EXISTS ingredients;
DROP TABLE IF EXISTS products;
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;

-- +goose StatementEnd
