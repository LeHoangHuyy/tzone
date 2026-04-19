# Docker stack

This folder contains local Docker support files for TZone:

- `docker-compose.yml` starts the API, frontend, PostgreSQL, MongoDB, and MinIO services.
- `docker/postgres/init/001-init.sql` enables the `pgcrypto` extension required by GORM UUID defaults.
- `docker/mongo-seed/` seeds the MongoDB `Cluster0.brands` collection from `phoneExample.json`.
- `media/` is mirrored to MinIO bucket for seeded device images.

