SET NAMES utf8;
SET time_zone = '+00:00';
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

CREATE TABLE IF NOT EXISTS `users`(
    `id` int NOT NULL AUTO_INCREMENT,
    `username` varchar(255) NOT NULL UNIQUE,
    `password` varchar(255) NOT NULL,
    PRIMARY KEY (`id`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE IF NOT EXISTS `films`
(
    `id` int NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(255) NOT NULL,
    `description` TEXT NOT NULL,
    `duration` int NOT NULL,
    `min_age` int NOT NULL,
    `country` varchar(255) NOT NULL,
    `producer_name` varchar(255) NOT NULL,
    `date_of_release` DATE NOT NULL,
    `sum_mark`      int NOT NULL,
    `num_of_marks`    int NOT NULL,
    `rating`        DECIMAL(3, 1) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `actors`
(
    `id` int NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(255) NOT NULL,
    `surname` VARCHAR(255) NOT NULL,
    `nationality` VARCHAR(255) NOT NULL,
    `birthday` DATE NOT NULL,
    PRIMARY KEY (`id`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `actor_films`
(
    `id` int NOT NULL AUTO_INCREMENT,
    `film_id` int NOT NULL,
    `actor_id` int NOT NULL,
    FOREIGN KEY (`film_id`)  REFERENCES `films`(`id`),
    FOREIGN KEY (`actor_id`)  REFERENCES `actors`(`id`),
    PRIMARY KEY (`id`)
 ) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE IF NOT EXISTS `genres`
(
    `id` int NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(255) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE IF NOT EXISTS `film_genres`
(
    `id` int NOT NULL AUTO_INCREMENT,
    `film_id` int NOT NULL,
    `genre_id` INTEGER NOT NULL,
    FOREIGN KEY (`film_id`)  REFERENCES `films`(`id`),
    FOREIGN KEY (`genre_id`)  REFERENCES `genres`(`id`),
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE IF NOT EXISTS `reviews`
(
    `id` int NOT NULL AUTO_INCREMENT,
    `film_id` int NOT NULL REFERENCES films (id),
    `user_id` int NOT NULL REFERENCES users(id),
    `mark` int NOT NULL,
    `comment` TEXT,
    FOREIGN KEY (`film_id`)  REFERENCES `films`(`id`),
    FOREIGN KEY (`user_id`)  REFERENCES `users`(`id`),
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE IF NOT EXISTS `favourite_films`
(
    `id` int NOT NULL AUTO_INCREMENT,
    `film_id` int NOT NULL REFERENCES films (id),
    `user_id` int NOT NULL REFERENCES users(id),
    FOREIGN KEY (`film_id`)  REFERENCES `films`(`id`),
    FOREIGN KEY (`user_id`)  REFERENCES `users`(`id`),
    PRIMARY KEY (`id`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

