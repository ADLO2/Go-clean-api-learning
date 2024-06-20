CREATE TABLE news
(
    news_id    INT                     PRIMARY KEY AUTO_INCREMENT,
    author_id  VARCHAR(250)                     NOT NULL,
    title      VARCHAR(250)             NOT NULL,
    content    TEXT                     NOT NULL,
    image_url  VARCHAR(1024),
    category   VARCHAR(250),
    created_at TIMESTAMP                NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP                NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE product
(
    product_id  INT                     PRIMARY KEY AUTO_INCREMENT,
    author_id   VARCHAR(250)                     NOT NULL,
    name        VARCHAR(250)            NOT NULL,
    description TEXT                    NOT NULL,
    price       INT                     NOT NULL,
    image_url   VARCHAR(1024),
    category    VARCHAR(250),
    created_at  TIMESTAMP               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP               NOT NULL DEFAULT CURRENT_TIMESTAMP
);