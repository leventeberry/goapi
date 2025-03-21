DROP TABLE IF EXISTS users;
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20) DEFAULT '9999999999',
    role ENUM('customer', 'admin', 'bartender') DEFAULT 'customer',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users (first_name, last_name, email, password_hash, role)
VALUES ('John', 'Doe', 'johndoe@example.com', 'password', 'admin'),
       ('Jane', 'Doe', 'janedoe@example.com', 'password', 'customer')

DROP TABLE IF EXISTS drinks;
CREATE TABLE drinks (
    drink_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(5, 2) NOT NULL,
    image_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO drinks (name, description, price, image_url)
VALUES ('Margarita', 'A classic cocktail made with tequila, lime juice, and triple sec.', 8.00, 'https://www.thecocktaildb.com/images/media/drink/5noda61589575158.jpg'),
       ('Old Fashioned', 'A cocktail made by muddling sugar with bitters, then adding alcohol, such as whiskey or brandy, and a twist of citrus rind.', 10.00, 'https://www.thecocktaildb.com/images/media/drink/vrwquq1478252802.jpg'),
       ('Martini', 'A cocktail made with gin and vermouth, and garnished with an olive or a lemon twist.', 12.00, 'https://www.thecocktaildb.com/images/media/drink/71t8581504353095.jpg'),
       ('Mojito', 'A cocktail made with rum, sugar, lime juice, soda water, and mint.', 9.00, 'https://www.thecocktaildb.com/images/media/drink/3z6xdi1589574603.jpg'),
       ('Cosmopolitan', 'A cocktail made with vodka, triple sec, cranberry juice, and freshly squeezed lime juice.', 10.00, 'https://www.thecocktaildb.com/images/media/drink/upxxpq1439907580.jpg');