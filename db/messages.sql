-- ----------------------------
-- Table structure for messages
-- ----------------------------
DROP TABLE IF EXISTS "public"."messages";
CREATE TABLE "public"."messages" (
  "id" int8 NOT NULL DEFAULT nextval('messages_id_seq'::regclass),
  "title" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "content" text COLLATE "pg_catalog"."default",
  "createAt" timestamp(0)
)
;

-- ----------------------------
-- Indexes structure for table messages
-- ----------------------------
CREATE UNIQUE INDEX "id" ON "public"."messages" USING btree (
  "id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table messages
-- ----------------------------
ALTER TABLE "public"."messages" ADD CONSTRAINT "messages_pkey" PRIMARY KEY ("id");
