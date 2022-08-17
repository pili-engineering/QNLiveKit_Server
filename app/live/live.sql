CREATE TABLE `live_entities` (
   `id` int unsigned NOT NULL AUTO_INCREMENT,
   `live_id` varchar(64) DEFAULT NULL,
   `title` varchar(128) DEFAULT NULL,
   `notice` varchar(255) DEFAULT NULL,
   `cover_url` varchar(255) DEFAULT NULL,
   `extends` varchar(1024) DEFAULT NULL,
   `anchor_id` varchar(64) DEFAULT NULL,
   `status` int DEFAULT NULL,
   `pk_id` varchar(64) DEFAULT NULL,
   `online_count` int DEFAULT NULL,
   `start_at` datetime DEFAULT NULL,
   `end_at` datetime DEFAULT NULL,
   `chat_id` bigint DEFAULT NULL,
   `push_url` varchar(255) DEFAULT NULL,
   `rtmp_play_url` varchar(255) DEFAULT NULL,
   `flv_play_url` varchar(255) DEFAULT NULL,
   `hls_play_url` varchar(255) DEFAULT NULL,
   `last_heartbeat_at` datetime DEFAULT NULL,
   `created_at` datetime DEFAULT NULL,
   `updated_at` datetime DEFAULT NULL,
   `deleted_at` datetime DEFAULT NULL,
   PRIMARY KEY (`id`),
   UNIQUE KEY `live_id` (`live_id`),
   KEY `idx_status_deleted_at` (`status`,`deleted_at`),
   KEY `idx_anchor_deleted_at` (`anchor_id`,`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `live_mic_entities` (
   `id` int unsigned NOT NULL AUTO_INCREMENT,
   `live_id` varchar(64) DEFAULT NULL,
   `user_id` varchar(64) DEFAULT NULL,
   `mic` tinyint(1) DEFAULT NULL,
   `camera` tinyint(1) DEFAULT NULL,
   `status` int DEFAULT NULL,
   `extends` varchar(1024) DEFAULT NULL,
   `created_at` datetime DEFAULT NULL,
   `updated_at` datetime DEFAULT NULL,
   `deleted_at` datetime DEFAULT NULL,
   PRIMARY KEY (`id`),
   KEY `idx_live` (`live_id`,`deleted_at`),
   KEY `idx_user` (`user_id`,`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `live_users` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `user_id` varchar(64) DEFAULT NULL,
    `nick` varchar(64) DEFAULT NULL,
    `avatar` varchar(255) DEFAULT NULL,
    `extends` json DEFAULT NULL,
    `im_userid` bigint DEFAULT NULL,
    `im_username` varchar(64) DEFAULT NULL,
    `im_password` varchar(32) DEFAULT NULL,
    `created_at` datetime DEFAULT NULL,
    `updated_at` datetime DEFAULT NULL,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `user_id` (`user_id`),
    KEY `idx_im_userid` (`im_userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `live_room_user_entities` (
   `id` int unsigned NOT NULL AUTO_INCREMENT,
   `live_id` varchar(64) DEFAULT NULL,
   `user_id` varchar(64) DEFAULT NULL,
   `status` int DEFAULT NULL,
   `heart_beat_at` datetime DEFAULT NULL,
   `created_at` datetime DEFAULT NULL,
   `updated_at` datetime DEFAULT NULL,
   `deleted_at` datetime DEFAULT NULL,
   PRIMARY KEY (`id`),
   KEY `idx_live` (`live_id`),
   KEY `idx_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `relay_sessions` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `sid` varchar(64) DEFAULT NULL,
    `init_user_id` varchar(64) DEFAULT NULL,
    `init_room_id` varchar(64) DEFAULT NULL,
    `recv_user_id` varchar(64) DEFAULT NULL,
    `recv_room_id` varchar(64) DEFAULT NULL,
    `extends` varchar(1024) DEFAULT NULL,
    `status` int DEFAULT NULL,
    `start_at` datetime DEFAULT NULL,
    `stop_at` datetime DEFAULT NULL,
    `created_at` datetime DEFAULT NULL,
    `updated_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `sid` (`sid`),
    KEY `idx_init_room` (`init_room_id`),
    KEY `idx_recv_room` (`recv_room_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `items` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime DEFAULT NULL,
    `updated_at` datetime DEFAULT NULL,
    `deleted_at` datetime DEFAULT NULL,
    `live_id` varchar(64) DEFAULT NULL,
    `item_id` varchar(64) DEFAULT NULL,
    `order` int unsigned DEFAULT NULL,
    `title` varchar(64) DEFAULT NULL,
    `tags` varchar(128) DEFAULT NULL,
    `thumbnail` varchar(255) DEFAULT NULL,
    `link` varchar(255) DEFAULT NULL,
    `current_price` varchar(64) DEFAULT NULL,
    `origin_price` varchar(64) DEFAULT NULL,
    `status` int unsigned DEFAULT NULL,
    `record` varchar(64) DEFAULT NULL,
    `extends` varchar(1024) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_live_status` (`live_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `item_demonstrate` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `live_id` varchar(255) DEFAULT NULL,
    `item_id` varchar(255) DEFAULT NULL,
    `updated_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uix_live` (`live_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE  TABLE `item_demonstrate_log`(
        `id` int unsigned NOT NULL AUTO_INCREMENT,
        `live_id` varchar(255) DEFAULT NULL,
        `item_id` varchar(255) DEFAULT NULL,
        `start` datetime DEFAULT NULL,
        `end` datetime DEFAULT NULL,
        `fname` varchar(255) DEFAULT NULL,
        `status` int DEFAULT NULL,
        `format` int DEFAULT NULL,
        `expireDays` int DEFAULT NULL,
        `persistentID` varchar(255) DEFAULT NULL,
        PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;