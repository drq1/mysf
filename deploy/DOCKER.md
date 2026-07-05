# mysf Bundled Redis Image

This image packages Sub2API and Redis in the same container.
You only need to provide PostgreSQL connection settings from the outside.

## Image

```bash
ghcr.io/<owner>/mysf:latest
```

## Quick Start

```bash
docker run -d \
  --name sub2api \
  -p 8080:8080 \
  -e DATABASE_HOST=host.docker.internal \
  -e DATABASE_PORT=5432 \
  -e DATABASE_USER=sub2api \
  -e DATABASE_PASSWORD=change_this_secure_password \
  -e DATABASE_DBNAME=sub2api \
  -e DATABASE_SSLMODE=disable \
  ghcr.io/<owner>/mysf:latest
```

If your PostgreSQL server is in another container on the same Docker network, set
`DATABASE_HOST` to that container name instead.

## Minimal Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `DATABASE_HOST` | PostgreSQL host | Yes | - |
| `DATABASE_PORT` | PostgreSQL port | No | `5432` |
| `DATABASE_USER` | PostgreSQL user | No | `postgres` |
| `DATABASE_PASSWORD` | PostgreSQL password | Yes | - |
| `DATABASE_DBNAME` | PostgreSQL database name | No | `sub2api` |
| `DATABASE_SSLMODE` | PostgreSQL SSL mode | No | `disable` |

Redis runs locally inside the container and is already wired to `127.0.0.1:6379`.
You do not need to set any Redis environment variables for the bundled image.

## Notes

- The container keeps application data under `/app/data`.
- The bundled Redis instance stores its data in `/app/data/redis`.
- If you want the previous deployment model with an external Redis service, use
  the standard Sub2API image instead of this bundled one.
