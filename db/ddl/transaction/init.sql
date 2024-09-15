SET FOREIGN_KEY_CHECKS = 0;

CREATE TABLE `user` (
  `id` VARCHAR(255) NOT NULL COMMENT 'ID',
  `user_id` VARCHAR(255) NOT NULL COMMENT 'ユーザー作成ID',
  `email` VARCHAR(255) NOT NULL COMMENT 'メールアドレス',
  `password` VARCHAR(255) NOT NULL COMMENT 'パスワード',
  `name` VARCHAR(255) NOT NULL COMMENT 'ユーザー名',
  `profile` VARCHAR(255) NOT NULL COMMENT 'プロフィール',
  `icon_img` VARCHAR(255) NOT NULL COMMENT 'アイコン画像パス',
  `header_img` VARCHAR(255) NOT NULL COMMENT 'ヘッダー画像パス',
  `created_at` DATETIME NOT NULL COMMENT '作成日時',
  `updated_at` DATETIME NOT NULL COMMENT '更新日時',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET=utf8mb4
COMMENT='ユーザー';

SET FOREIGN_KEY_CHECKS = 1;
