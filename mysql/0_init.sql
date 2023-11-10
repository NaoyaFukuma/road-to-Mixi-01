CREATE TABLE `users` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL UNIQUE,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `friend_link` (
  `user1_id` bigint(20) NOT NULL,
  `user2_id` bigint(20) NOT NULL,
  PRIMARY KEY (`user1_id`, `user2_id`),
  CONSTRAINT `fk_friend_link_user1` FOREIGN KEY (`user1_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_friend_link_user2` FOREIGN KEY (`user2_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CHECK (`user1_id` != `user2_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `block_list` (
  `user1_id` bigint(20) NOT NULL,
  `user2_id` bigint(20) NOT NULL,
  PRIMARY KEY (`user1_id`, `user2_id`),
  CONSTRAINT `fk_block_list_user1` FOREIGN KEY (`user1_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_block_list_user2` FOREIGN KEY (`user2_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CHECK (`user1_id` != `user2_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `friend_requests` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `requester_id` bigint(20) NOT NULL,
  `requested_id` bigint(20) NOT NULL,
  `status` enum('pending','accepted','declined') NOT NULL DEFAULT 'pending',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_requester_requested_status` (`requester_id`, `requested_id`, `status`),
  CONSTRAINT `fk_friend_requests_requester` FOREIGN KEY (`requester_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_friend_requests_requested` FOREIGN KEY (`requested_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
