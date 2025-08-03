-- Add transaction fee column to transfers table
ALTER TABLE "transfers" ADD COLUMN "fee" DECIMAL(20,2) DEFAULT 0;

-- Add comment for documentation
COMMENT ON COLUMN "transfers"."fee" IS 'Transaction fee amount in source currency'; 