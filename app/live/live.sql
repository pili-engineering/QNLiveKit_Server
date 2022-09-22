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

   `stop_reason` varchar(64) DEFAULT '',
   `stop_user_id` varchar(64) DEFAULT '',
   `stop_at` datetime DEFAULT NULL,

   `unaudit_censor_count` int DEFAULT  NULL,
   `last_censor_time` int DEFAULT NULL,
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
    `record_id` int unsigned DEFAULT NULL,
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

CREATE TABLE `stats_single_live`(
      `id` int unsigned NOT NULL AUTO_INCREMENT,
      `live_id` varchar(255) DEFAULT NULL,
      `biz_id` varchar(255) DEFAULT NULL,
      `user_id` varchar(64) DEFAULT NULL,
      `type` int DEFAULT NULL,
      `count` int DEFAULT 0,
      `updated_at` datetime DEFAULT NULL,
      PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `admin_user` (
     `id` int unsigned NOT NULL AUTO_INCREMENT,
     `user_name` varchar(255) DEFAULT NULL,
     `user_id` varchar(255) DEFAULT NULL,
     `password` varchar(255) DEFAULT NULL,
     `description` varchar(255) DEFAULT NULL,
      PRIMARY KEY (`id`),
      UNIQUE KEY `uix_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `censor_config` (
      `id` int unsigned NOT NULL,
      `enable` BOOLEAN DEFAULT FALSE,
      `interval` int DEFAULT NULL,
      `pulp` BOOLEAN DEFAULT FALSE,
      `terror` BOOLEAN DEFAULT FALSE,
      `politician` BOOLEAN DEFAULT FALSE,
      `ads` BOOLEAN DEFAULT FALSE,
       PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `operation_log` (
     `id` int unsigned NOT NULL AUTO_INCREMENT,
     `user_id` varchar(30) NOT NULL DEFAULT '',
     `ip` varchar(30) DEFAULT NULL,
     `method` varchar(10) DEFAULT NULL,
     `url` varchar(1024) DEFAULT NULL,
     `args` varchar(1024) DEFAULT NULL,
     `status_code` int NOT NULL DEFAULT 0,
     `created_at` timestamp NULL DEFAULT NULL,
     PRIMARY KEY (`id`),
     KEY `idx_user`(`user_id`, `created_at`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `live_censor` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `live_id` varchar(255) DEFAULT NULL,
    `job_id` varchar(255) DEFAULT NULL,
    `pulp` BOOLEAN DEFAULT FALSE,
    `terror` BOOLEAN DEFAULT FALSE,
    `politician` BOOLEAN DEFAULT FALSE,
    `ads` BOOLEAN DEFAULT FALSE,
    `interval` int DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_live_job` (`live_id`,`job_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `censor_image`(
     `id` int unsigned NOT NULL AUTO_INCREMENT,
     `url` varchar(255)  NOT NULL,
     `job_id` varchar(255)  NOT NULL,
     `created_at` int DEFAULT NULL,
     `suggestion` varchar(255) DEFAULT NULL,
     `pulp`  varchar(255) DEFAULT FALSE,
     `terror` varchar(255)  DEFAULT FALSE,
     `politician` varchar(255) DEFAULT FALSE,
     `ads` varchar(255) DEFAULT FALSE,
     `live_id` varchar(255) DEFAULT NULL,
     `is_review` int DEFAULT 0,
     `review_answer` int DEFAULT NULL,
     `review_user_id` varchar(255) DEFAULT NULL,
     `review_time` datetime DEFAULT NULL,
     PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `live_like` (
     `id` int unsigned NOT NULL AUTO_INCREMENT,
     `live_id` varchar(64) NOT NULL  COMMENT '直播间ID',
     `user_id` varchar(64) NOT NULL  COMMENT '用户ID，空表示直播间内总点赞数',
     `count` int NOT NULL  COMMENT '点赞数量',
     `created_at` datetime DEFAULT NULL,
     `updated_at` datetime DEFAULT NULL,
     PRIMARY KEY (`id`),
     UNIQUE KEY `live_user` (`live_id`, `user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `live_gift` (
     `id` int unsigned NOT NULL AUTO_INCREMENT,
     `biz_id` varchar(64) NOT NULL COMMENT '业务ID，用于接口幂等处理',
     `live_id` varchar(64) NOT NULL COMMENT '直播间ID',
     `user_id` varchar(64) NOT NULL COMMENT '用户ID',
     `anchor_id` varchar(64)  NOT NULL COMMENT '主播ID',
     `gift_id` int NOT NULL  COMMENT '礼物id',
     `amount` int NOT NULL DEFAULT 0 COMMENT '礼物金额',
     `status` int NOT NULL DEFAULT 0 COMMENT '状态',
     `created_at` datetime DEFAULT NULL,
     `updated_at` datetime DEFAULT NULL,
     PRIMARY KEY (`id`),
     UNIQUE KEY `uix_biz` (`biz_id`),
     KEY `idx_live_user` (`live_id`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `gift_config` (
     `id` int unsigned NOT NULL AUTO_INCREMENT,
     `gift_id` int NOT NULL COMMENT '礼物id',
     `type` int NOT NULL DEFAULT 0 COMMENT '礼物类型',
     `name` varchar(64) NOT NULL COMMENT '礼物名称',
     `amount` int NOT NULL COMMENT '礼物金额，0 表示自定义金额',
     `img` varchar(512) NOT NULL DEFAULT '' COMMENT '礼物图片',
     `animation_type` int NOT NULL DEFAULT 0 COMMENT '动态效果类型',
     `animation_img` varchar(512) NOT NULL DEFAULT '' COMMENT '动态效果图片',
     `order` int NOT NULL DEFAULT 0 COMMENT '排序，从小到大排序，相同order 根据创建时间排序',
     `created_at` datetime DEFAULT NULL,
     `updated_at` datetime DEFAULT NULL,
     `deleted_at` datetime DEFAULT NULL,
     `extends` varchar(1024) DEFAULT NULL,
     PRIMARY KEY (`id`),
     UNIQUE KEY uid_gift_id (`gift_id`),
     KEY idx_type_order (`type`,`order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
