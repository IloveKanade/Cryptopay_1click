-- 为payment_links表添加callback_url字段
ALTER TABLE payment_links ADD COLUMN callback_url VARCHAR(500) DEFAULT '' COMMENT '回调地址';

-- 更新现有记录的callback_url字段为空字符串
UPDATE payment_links SET callback_url = '' WHERE callback_url IS NULL;
