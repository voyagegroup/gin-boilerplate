-- +migrate Up
CREATE TABLE `todos` (
  `todo_id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `title` varchar(256) NOT NULL COMMENT 'タスクのタイトル',
  `completed` BOOL NOT NULL DEFAULT FALSE COMMENT 'タスクが完了したか否か',
  `created` datetime NOT NULL DEFAULT NOW() COMMENT '登録日',
  `updated` datetime DEFAULT NULL COMMENT '更新日',
  PRIMARY KEY (`todo_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='タスクリスト';

-- +migrate Down
DROP TABLE todos;
