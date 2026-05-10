import { defineConfig, env } from '@prisma/config';

export default defineConfig({
  schema: 'prisma/schema.prisma',
  datasource: {
    url: env('USER_SERVICE_DB_URL'),
  },
  migrations: {
    path: 'prisma/migrations',
  },
});
