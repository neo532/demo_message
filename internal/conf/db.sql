
CREATE DATABASE IF NOT EXISTS `message` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;

use message;

CREATE TABLE IF NOT EXISTS `mg_campaign` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
	`origin_type` tinyint NOT NULL DEFAULT '0' COMMENT 'origin type:1=file,2=url',
	`origin_content` varchar(200) NOT NULL DEFAULT '' COMMENT 'origin content',
    `message_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT 'message type:1=fixed,2=custom',
    `message` varchar(255) NOT NULL DEFAULT '' COMMENT 'message',
    `time_send` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'send time',
    `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT 'campaign status:1=on,2=off',
    `time_create` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
    `time_update` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB COMMENT='campaign-neo';

CREATE TABLE IF NOT EXISTS `mg_recipient` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `mobile` varchar(20) NOT NULL DEFAULT '' COMMENT 'mobile',
    `name` varchar(120) NOT NULL DEFAULT '' COMMENT 'name',
    `time_create` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
    `time_update` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_mobile` (`mobile`)
) ENGINE=InnoDB COMMENT='recipient-neo';

CREATE TABLE IF NOT EXISTS `mg_message` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `campaign_id` int(11) NOT NULL DEFAULT '0' COMMENT 'campaign_id',
    `recipient_id` int(11) NOT NULL DEFAULT '0' COMMENT 'recipient_id',
    `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT 'message status:1=to send,2=has sended',
    `time_send` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'send time',
    `log_id` varchar(50) NOT NULL DEFAULT '' COMMENT 'log id',
    `time_create` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
    `time_update` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB COMMENT='message-neo';


-- csv-campaign 
-- message_type,message,time_send

-- csv-message
-- mobile,name,campaign_id
