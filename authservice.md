```sql
CREATE TABLE `role`
(
    `id`         INT             PRIMARY KEY AUTO_INCREMENT,
    `name`       VARCHAR(20)     NOT NULL UNIQUE
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `perm`
(
    `id`         INT             PRIMARY KEY AUTO_INCREMENT,
    `path`       VARCHAR(20)     NOT NULL UNIQUE
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `user_role`
(
    `id`         INT             PRIMARY KEY AUTO_INCREMENT,
    `user_id`    INT             NOT NULL,
    `role_id`    INT             NOT NULL
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `role_perm`
(
    `id`         INT             PRIMARY KEY AUTO_INCREMENT,
    `role_id`    INT             NOT NULL,
    `perm_id`    INT             NOT NULL 
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `role` (`name`) VALUES ("ROOT");

INSERT INTO `role` (`name`) VALUES ("GESTURE");

INSERT INTO `perm` (`path`) VALUES ("/user/register");

INSERT INTO `perm` (`path`) VALUES ("/user/login");

INSERT INTO `perm` (`path`) VALUES ("order/create");

INSERT INTO `user_role` (`user_id`,`role_id`) VALUES (3,1);

INSERT INTO `user_role` (`user_id`,`role_id`) VALUES (3,2);

INSERT INTO `user_role` (`user_id`,`role_id`) VALUES (4,2);

INSERT INTO `user_role` (`user_id`,`role_id`) VALUES (5,1);

INSERT INTO `role_perm` (`role_id`,`perm_id`) VALUES (1,1);

INSERT INTO `role_perm` (`role_id`,`perm_id`) VALUES (1,2);

INSERT INTO `role_perm` (`role_id`,`perm_id`) VALUES (1,3);

INSERT INTO `role_perm` (`role_id`,`perm_id`) VALUES (2,1);

INSERT INTO `role_perm` (`role_id`,`perm_id`) VALUES (2,2);

select perm.path from user_role 
join role on user_role.role_id = role.id 
join role_perm on role_perm.role_id = role.id
join perm on perm.id = role_perm.perm_id
where user_id = 5;


```