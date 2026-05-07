CREATE TABLE IF NOT EXISTS wb_scraper.products (
    id SERIAL PRIMARY KEY,
    wb_id INTEGER UNIQUE NOT NULL,
    category_id INTEGER,
    name VARCHAR(255) NOT NULL, 
    brand VARCHAR(90) NOT NULL,
    "supplierId" INTEGER NOT NULL,
    "reviewRating" DOUBLE PRECISION NOT NULL,
    feedbacks INTEGER NOT NULL,
    pics INTEGER NOT NULL,

    CONSTRAINT fk_category
        FOREIGN KEY(category_id)
        REFERENCES wb_scraper.categories(id)
        ON DELETE SET NULL 
);

CREATE TABLE IF NOT EXISTS wb_scraper.sizes (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL,
    name VARCHAR(64),
    "priceBasic" DOUBLE PRECISION NOT NULL,
    "priceProduct" DOUBLE PRECISION NOT NULL,

    CONSTRAINT fk_product 
        FOREIGN KEY(product_id)
        REFERENCES wb_scraper.products(id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_products_category_id ON wb_scraper.products (category_id);
CREATE INDEX IF NOT EXISTS idx_sizes_product_id ON wb_scraper.sizes (product_id);
