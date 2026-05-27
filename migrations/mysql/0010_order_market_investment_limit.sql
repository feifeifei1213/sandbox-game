ALTER TABLE sg_order_market_config
    ADD COLUMN market_investment_limit DECIMAL(18,2) NULL COMMENT '单市场投入上限，NULL表示无上限' AFTER market_enabled;
