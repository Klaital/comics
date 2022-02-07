CREATE TABLE IF NOT EXISTS `users` (
    `user_id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    `email` varchar(128) NOT NULL,
    `password_digest` varchar(256) NOT NULL,
    PRIMARY KEY (`user_id`),
    UNIQUE INDEX `email` (`email`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
