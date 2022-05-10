#action_flow table
DROP TABLE
IF
	EXISTS `action_flow`;
CREATE TABLE
IF
	NOT EXISTS `action_flow` (
		`id` VARCHAR ( 36 ) NOT NULL COMMENT 'id',
		`appId` VARCHAR ( 36 ) NOT NULL COMMENT 'appId',
		`displayName` VARCHAR ( 128 ) NOT NULL DEFAULT '' COMMENT 'display name',
		`name` VARCHAR ( 128 ) NOT NULL DEFAULT '' COMMENT 'flow name',
		`description` VARCHAR ( 256 ) NOT NULL DEFAULT '' COMMENT 'description',
		`inputParameter` VARCHAR ( 1000 ) NOT NULL DEFAULT '' COMMENT 'input parameter',
		`outputParameter` VARCHAR ( 1000 ) NOT NULL DEFAULT '' COMMENT 'output parameter',
		`timeout` BIGINT NOT NULL COMMENT 'timeout',
		`retryOptionsId` VARCHAR ( 36 ) NOT NULL DEFAULT '' COMMENT 'retry options',
		`createdBy` VARCHAR ( 36 ) NOT NULL DEFAULT '' COMMENT 'created user id',
		`createdAt` BIGINT NOT NULL COMMENT 'created time',
		`ModifiedAt` BIGINT NOT NULL DEFAULT 0 COMMENT 'modified user id',
		`modifiedBy` VARCHAR ( 36 ) NOT NULL DEFAULT '' COMMENT 'modified time',
		PRIMARY KEY ( `id` ) 
	) ENGINE = INNODB DEFAULT CHARSET = utf8mb4;

#action table
DROP TABLE
IF
	EXISTS `action`;
CREATE TABLE
IF
	NOT EXISTS `action` (
		`id` VARCHAR ( 36 ) NOT NULL COMMENT 'id',
        `displayName` VARCHAR ( 128 ) NOT NULL  COMMENT 'display name',
        `name` VARCHAR ( 128 ) NOT NULL  COMMENT 'action name',
		`flowId` VARCHAR ( 36 ) NOT NULL COMMENT 'action flow id',
		`type` VARCHAR ( 36 ) NOT NULL COMMENT 'action type [conditional | action | variable]',
		`parameter` VARCHAR ( 1000 ) NOT NULL DEFAULT '' COMMENT 'parameter',
		`command` VARCHAR ( 5000 ) NOT NULL DEFAULT '' COMMENT 'handle object',
		`retryOptionsId` VARCHAR ( 36 ) NOT NULL DEFAULT '' COMMENT 'retry options',
		`createdBy` VARCHAR ( 36 ) NOT NULL COMMENT 'created user id',
		`createdAt` BIGINT NOT NULL COMMENT 'created time',
		`ModifiedAt` BIGINT NOT NULL DEFAULT 0 COMMENT 'modified user id',
		`modifiedBy` VARCHAR ( 36 ) NOT NULL DEFAULT '' COMMENT 'modified time',
		PRIMARY KEY ( `id` ) 
	) ENGINE = INNODB DEFAULT CHARSET = utf8mb4;

#retry options table
DROP TABLE
    IF
        EXISTS `retry_options`;
CREATE TABLE
    IF
    NOT EXISTS `retry_options` (
        `id` VARCHAR ( 36 ) NOT NULL COMMENT 'id',
        `resourceType` VARCHAR ( 16 ) NOT NULL COMMENT 'type [actionflow | action]',
        `resourceId` VARCHAR ( 36 ) NOT NULL COMMENT 'resource id',
        `Attempts` VARCHAR ( 128 ) NOT NULL DEFAULT '' COMMENT 'display name',
        `initialInterval` VARCHAR ( 128 ) NOT NULL DEFAULT '' COMMENT 'action name',
        `coefficient` VARCHAR ( 36 ) NOT NULL COMMENT 'action flow id',
        `maximumInterval` VARCHAR ( 36 ) NOT NULL COMMENT 'action type [conditional | action | variable]',
        `maximumAttempts` VARCHAR ( 1000 ) NOT NULL DEFAULT '' COMMENT 'parameter',
        `NonRetryableErrorReasons` VARCHAR ( 36 ) NOT NULL COMMENT 'handle type [rest]',
        `createdBy` VARCHAR ( 36 ) NOT NULL COMMENT 'created user id',
        `createdAt` BIGINT NOT NULL COMMENT 'created time',
        `ModifiedAt` BIGINT NOT NULL DEFAULT 0 COMMENT 'modified user id',
        `modifiedBy` VARCHAR ( 36 ) NOT NULL DEFAULT '' COMMENT 'modified time',
        PRIMARY KEY ( `id` )
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4;

#timer table
DROP TABLE
IF
	EXISTS `timer`;
CREATE TABLE
IF
	NOT EXISTS `timer` (
		`id` VARCHAR ( 36 ) NOT NULL COMMENT 'id',
		`name` VARCHAR ( 36 ) NOT NULL COMMENT 'name',
		`description` VARCHAR ( 36 ) NOT NULL DEFAULT '' COMMENT 'description',
		`actionFlowId` VARCHAR ( 1000 ) NOT NULL COMMENT 'action flow id',
		`actionFlowVariable` VARCHAR ( 36 ) NOT NULL DEFAULT '' COMMENT 'action flow variable',
		`cron` VARCHAR ( 24 ) NOT NULL COMMENT 'cron',
		`timeout` BIGINT NOT NULL COMMENT 'timeout',
		`status` VARCHAR ( 16 ) NOT NULL COMMENT 'status [open | closed ]',
		`endTime` BIGINT NOT NULL DEFAULT 0 COMMENT 'end time',
		`createdBy` VARCHAR ( 36 ) NOT NULL COMMENT 'created user id',
		`createdAt` BIGINT NOT NULL COMMENT 'created time',
		`ModifiedAt` BIGINT NOT NULL DEFAULT 0 COMMENT 'modified user id',
		`modifiedBy` VARCHAR ( 36 ) NOT NULL DEFAULT '' COMMENT 'modified time',
		PRIMARY KEY ( `id` ) 
	) ENGINE = INNODB DEFAULT CHARSET = utf8mb4;

#action_flow_job table
DROP TABLE
IF
	EXISTS `action_flow_job`;
CREATE TABLE
IF
	NOT EXISTS `action_flow_job` (
		`id` VARCHAR ( 36 ) NOT NULL COMMENT 'id',
		`flowId` VARCHAR ( 36 ) NOT NULL COMMENT 'action flow id',
		`flowVariable` VARCHAR ( 1000 ) NOT NULL DEFAULT '' COMMENT 'action flow variable',
		`timeout` BIGINT NOT NULL COMMENT 'timeout',
		`status` VARCHAR ( 16 ) NOT NULL COMMENT 'status [deleted | pause | active]',
        `message` VARCHAR ( 128 ) NOT NULL DEFAULT '' COMMENT 'result message',
		`retryOptionsId` VARCHAR ( 36 ) NOT NULL DEFAULT '' COMMENT 'retry options',
		`createdBy` VARCHAR ( 36 ) NOT NULL COMMENT 'created user id',
		`createdAt` BIGINT NOT NULL COMMENT 'created time',
		`ModifiedAt` BIGINT NOT NULL DEFAULT 0 COMMENT 'modified user id',
		`modifiedBy` VARCHAR ( 36 ) NOT NULL DEFAULT '' COMMENT 'modified time',
	PRIMARY KEY ( `id` )
	) ENGINE = INNODB DEFAULT CHARSET = utf8mb4;


#action_flow_job_run_instance table
DROP TABLE
    IF
        EXISTS `action_flow_job_run_instance`;
CREATE TABLE
    IF
    NOT EXISTS `action_flow_job_run_instance` (
        `id` VARCHAR ( 36 ) NOT NULL COMMENT 'id',
        `jobId` VARCHAR ( 36 ) NOT NULL COMMENT 'action flow job id',
        `flowId` VARCHAR ( 36 ) NOT NULL COMMENT 'action flow id',
        `flowVariable` VARCHAR ( 1000 ) NOT NULL DEFAULT '' COMMENT 'action flow variable',
        `timeout` BIGINT NOT NULL COMMENT 'timeout',
        `status` VARCHAR ( 16 ) NOT NULL COMMENT 'status [pending | executing | success | execute | cancelled]',
        `message` VARCHAR ( 256 ) NOT NULL DEFAULT '' COMMENT 'result message',
        `retryOptionsId` VARCHAR ( 36 ) NOT NULL DEFAULT '' COMMENT 'retry options',
        `executedAt` BIGINT NOT NULL COMMENT 'executed time',
        `finishedAt` BIGINT NOT NULL COMMENT 'finished time',
        `createdBy` VARCHAR ( 36 ) NOT NULL COMMENT 'created user id',
        `createdAt` BIGINT NOT NULL COMMENT 'created time',
        `ModifiedAt` BIGINT NOT NULL DEFAULT 0 COMMENT 'modified user id',
        `modifiedBy` VARCHAR ( 36 ) NOT NULL DEFAULT '' COMMENT 'modified time',
        PRIMARY KEY ( `id` )
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4;


DROP TABLE
    IF
        EXISTS flow_action_edge;
CREATE TABLE
    IF
    NOT EXISTS flow_action_edge (
        `id` bigint UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        `entryEdgeId` int NOT NULL,
        `directEdgeId` int NOT NULL,
        `exitEdgeId` int NOT NULL,
        `startActionId` varchar(36) NOT NULL,
        `endActionId` varchar(36) NOT NULL,
        `hops` int NOT NULL
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4;
