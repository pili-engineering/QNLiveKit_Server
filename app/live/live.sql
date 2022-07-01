CREATE TABLE `live_entities` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `live_id` varchar(255) DEFAULT NULL,
  `title` varchar(255) DEFAULT NULL,
  `notice` varchar(255) DEFAULT NULL,
  `cover_url` varchar(255) DEFAULT NULL,
  `extends` json DEFAULT NULL,
  `anchor_id` varchar(255) DEFAULT NULL,
  `status` int(11) DEFAULT NULL,
  `pk_id` varchar(255) DEFAULT NULL,
  `online_count` int(11) DEFAULT NULL,
  `start_at` datetime DEFAULT NULL,
  `end_at` datetime DEFAULT NULL,
  `relay_session_id` varchar(255) DEFAULT NULL,
  `chat_id` bigint(20) DEFAULT NULL,
  `push_url` varchar(255) DEFAULT NULL,
  `rtmp_play_url` varchar(255) DEFAULT NULL,
  `flv_play_url` varchar(255) DEFAULT NULL,
  `hls_play_url` varchar(255) DEFAULT NULL,
  `last_heartbeat_at` datetime DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;



CREATE TABLE `live_mic_entities` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `live_id` varchar(255) DEFAULT NULL,
  `user_id` varchar(255) DEFAULT NULL,
  `mic` tinyint(1) DEFAULT NULL,
  `camera` tinyint(1) DEFAULT NULL,
  `status` int(11) DEFAULT NULL,
  `extends` json DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;



CREATE TABLE `live_room_user_entities` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `live_id` varchar(255) DEFAULT NULL,
  `user_id` varchar(255) DEFAULT NULL,
  `status` int(11) DEFAULT NULL,
  `heart_beat_at` datetime DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;



CREATE TABLE `live_users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` varchar(255) DEFAULT NULL,
  `nick` varchar(255) DEFAULT NULL,
  `avatar` varchar(255) DEFAULT NULL,
  `extends` json DEFAULT NULL,
  `im_userid` bigint(20) DEFAULT NULL,
  `im_username` varchar(255) DEFAULT NULL,
  `im_password` varchar(255) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=utf8mb4;


CREATE TABLE `relay_sessions` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `sid` varchar(255) DEFAULT NULL,
  `init_user_id` varchar(255) DEFAULT NULL,
  `init_room_id` varchar(255) DEFAULT NULL,
  `recv_user_id` varchar(255) DEFAULT NULL,
  `recv_room_id` varchar(255) DEFAULT NULL,
  `extends` varchar(512) DEFAULT NULL,
  `status` int(11) DEFAULT NULL,
  `start_at` datetime DEFAULT NULL,
  `stop_at` datetime DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;