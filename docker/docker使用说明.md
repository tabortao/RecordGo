**本地存储 docker-compose**

- 文件名建议：`docker-compose.local.yml`
- 用法：`docker compose -f docker-compose.local.yml --env-file .env.local up -d`

**本地存储 .env**

- 文件名建议：`.env.local`

**S3 存储 docker-compose（MinIO / path-style）**

- 文件名建议：`docker-compose.s3.yml`
- 用法：`docker compose -f docker-compose.s3.yml --env-file .env.s3 up -d`

**S3 存储 .env（MinIO 示例）**

- 文件名建议：`.env.s3`

**关键区别与说明**

- 本地存储
  - 使用 `STORAGE_DRIVER=local`，文件写入容器内 `/app/storage/uploads/...`；通过 `volumes` 将 `./backend/storage` 挂载到容器，持久化 SQLite 与上传文件。
- S3 存储（MinIO）
  - 使用 `STORAGE_DRIVER=s3`，上传通过预签名直传到 `S3_ENDPOINT`；读取展示使用 `S3_PUBLIC_BASE_URL`。
  - MinIO 必须 `S3_FORCE_PATH_STYLE=true`；展示地址为 `S3_PUBLIC_BASE_URL + '/' + S3_BUCKET + '/' + object_key`。
  - 当桶名为 `blog` 时，`S3_KEY_PREFIX` 留空，避免路径重复成 `blog/blog/...`。
- 运行命令
  - 本地存储：`docker compose -f docker-compose.local.yml --env-file .env.local up -d`
  - S3 存储：`docker compose -f docker-compose.s3.yml --env-file .env.s3 up -d`
