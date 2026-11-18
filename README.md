# Comandos de Swagger para Go (devboard)

Guía rápida de comandos para configurar, generar y mantener la documentación Swagger del proyecto.

## Instalación

Instalar la herramienta de línea de comandos `swag` (solo se hace una vez):

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

## Dependencias

Agregar las librerías de Swagger al proyecto:

```bash
go get github.com/swaggo/swag@latest
go get github.com/swaggo/http-swagger@latest
```

## Generar la documentación

Leer las anotaciones del código y generar los archivos de documentación en la carpeta `docs`:

```bash
swag init -g cmd/api/main.go -o docs
```

El flag `-g` indica el archivo principal con la info general de la API, y `-o docs` define la carpeta de salida.

## Mantenimiento

Actualizar y ordenar las dependencias del proyecto:

```bash
go mod tidy
```

Borrar la carpeta de documentación generada (por si quieres regenerarla desde cero):

```bash
rm -rf docs
```

Formatear todo el código Go al estilo estándar:

```bash
gofmt -w .
```

## Flujo típico

Cuando cambias las anotaciones y quieres regenerar la doc limpia:

```bash
rm -rf docs
swag init -g cmd/api/main.go -o docs
go mod tidy
gofmt -w .
```

---

# Patrones Útiles en Go

**Patrón Repository**
Nos ayuda a separar la lógica del negocio del acceso a datos.

**Patrón Factory**
Centraliza la creación de objetos cuando construirlos empieza a ser muy repetitivo o muy complejo.

**Patrón Strategy**
Permite cambiar un comportamiento sin modificar el código principal.

**Patrón Adapter**
Permite conectar la aplicación con una librería, API externa o sistema que tiene una interfaz distinta a la que necesitamos.

**Patrón Middleware**
Permite ejecutar lógica antes y después de una petición HTTP.

**Inyección de Dependencias**
Sirve para pasar dependencias desde afuera en lugar de crearlas dentro del componente.

## Las tres preguntas que la arquitectura responde

1. ¿Puedo testear la lógica del negocio sin levantar una base de datos real?
2. ¿Puedo cambiar de PostgreSQL a otra base de datos sin reescribir toda la lógica del negocio?
3. ¿Puede un desarrollador nuevo entender el flujo sin leer 200 líneas de un solo archivo?

## Clean Architecture

### 1. Entities (Domain)

- Contiene las reglas de negocio más importantes y estables.
- El núcleo de la aplicación.
- No depende de nada externo.
- Responde a: qué es, cómo se representa y qué validaciones tiene. Ejemplos: Usuario, Tarea.

### 2. Use Cases (use case)

- Contiene la lógica específica de la aplicación.
- Coordina las entidades y define las interfaces que necesitamos.

### 3. Interface Adapters / Delivery (handler)

- Traduce datos entre el mundo externo y los casos de uso.
- Handlers HTTP, Controllers, DTOs, Presenters, implementaciones de repositorios, Mappers, etc.

### 4. Frameworks & Drivers / Infrastructure (repository)

- Contiene herramientas y detalles técnicos.
- Bases de datos, frameworks, servidores, APIs externas, sistemas de archivos, etc.
- Puede cambiar sin afectar el negocio.

### Regla principal: dependencia hacia adentro

Las capas externas pueden depender de las internas, pero las internas no deben conocer las externas.

**Dirección de las dependencias:**

```
Frameworks & Drivers --> Interface Adapters --> Use Cases --> Entities
```

## Concurrencia: maps y Mutex

- Los maps **no son seguros** para lecturas y escrituras concurrentes.
  Error típico: `fatal error: concurrent map read and map write`

**¿Qué es un Mutex?**
Mutex significa _mutual exclusion_, es decir, exclusión mutua. Sirve para que solo una goroutine acceda a un recurso compartido a la vez, evitando ese error de acceso concurrente.

---

# Base de datos y migraciones

## Dependencias

Agregar las librerías de Postgres (pgx) y migraciones al proyecto:

```bash
go get github.com/jackc/pgx/v5
go get -tags 'postgres' github.com/golang-migrate/migrate/v4
go get github.com/golang-migrate/migrate/v4/database/pgx/v5
go get github.com/golang-migrate/migrate/v4/source/file
```

## Herramientas de línea de comandos (CLI)

Instalar `migrate` (con soporte para Postgres) y `sqlc` como comandos globales:

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

## Verificar instalaciones

```bash
migrate --version
sqlc version
```

## Error en migraciones

```bash
$ make migrate-up
migrate -path ./migrations -database "postgresql://postgres:password@localhost:5433/devboard?sslmode=disable" up
error: Dirty database version 2. Fix and force version.
make: *** [migrate-up] Error 1

$ migrate -path ./migrations -database "postgresql://postgres:password@localhost:5433/devboard?sslmode=disable" force 1
```

```bash
migrate -path ./migrations -database "postgresql://postgres:password@localhost:5433/devboard?sslmode=disable" force 2
```

# Comando para migrag una tabla de una bd a otra bd

```
docker compose exec -T postgres \
> pg_dump \
> -U postgres \
> -d devboard \
> --data-only \
> --table=public.users \
> --column-inserts \
> --no-owner \
> --no-privileges \
> | docker compose exec -T identity-postgres \
pipe> psql \
pipe> -U identity \
pipe> -d identity

```

# - Pasos para crear un microservice partiendo de un monolito

1 - Crear la carpeta **task** dentro de service

2 - Crear la estructura de carpetas **internal** y **cmd**/**api**

3 - En la terminal ir a la ruta del nuevo servicio

```bash
cd services/task
```

4 - ejecutar el comando:
Para crear el archivo **go.mod** del nuevo servicio

```go
go mod init github.com/valeroman/devboard/services/task
```

5 - Regresar a la ruta principal y ejecutamos el siguiente comando:
Lo que hace es incluir task en el workspace del monolito

```go
go work use ./services/task
```

6 - Crea el server de task copiando la carpeta **server** del servicio de workspace o identity

7 - Creamos las carpeta **domain** dentro del servicio task

8 - Copiamos los archivos del **domian** del monolito a nuestra carpeta **domain**, los siguiestes archivos:

- comment.go
- erros.go
- task_status.go
- task.go

9 - Creamos la carpeta **usecase** y copiamos los usecases que le pertenece a task del monolito y revisamos sus dependencias actualizando la rutas del servicio
ejemplo : "github.com/valeroman/devboard/services/task/internal/domain"

10 - Crear la carpeta dentro de interna **repository** y dentro de repository **postgres** y dentro de postgres **queries**

11 - Copiamos los archivos de **queries** del monolito a nuestra carpeta **quieries**

- comments.sql
- tasks.sql

12 - Crear un schema para darle la estuctura base de nuestra tabla

- creamos una carpeta dentro de postgres **schema** y el archivo **task.sql**, que tendra en sql la creacion de las tablas **tasks** y **comments**

13 - Creamos el archivo **sqlc.yml**, en la raiz del servicio **task** y copiamos el contenido del mismo archivo pero del servicio **workspace**

14 - y cambiamos

```yml
schema: "migrations/"
```

por

```yml
queries: "internal/repository/postgres/schema/"
```

15 - Ejecutamos el siguiente comando, pero en la ruta de

```bash
cd services/task
```

comando:

```go
sqlc generate
```

Y eso me genera la carpeta **generated** dentro de repository/postgres con los siguientes archivos:

- comments.sql.go
- db.go
- models.go
- tasks.sql.go

16 - Copiar el archivo **db.go** del servicio workspace en la ruta **repository/postgres**, a nuestro servicio task en la misma ruta

17 - Copiamos los archivos **comments_rapository.go** y _task_repository.go_ del monilito **repository/postgres**, a nuestro servicio task en la misma ruta **repository/postgres** y verificamos sus dependencias y las actualizamos.

18 - Creamos las carpetas **handler**, **validator** y **middlewares** dentro de **internal**

19 - Copiamos los archivos **comment_handler.go**, **task_handler.go**, **response.go** y **pagination.go**, a nuestro **handler** del servicio task y revisamos sus dependencias que hay que actualizar.

20 - Copiamos el archivo de **validator.go** de la carpeta **validator** dentro de **internal** del servicio workspace, a nuestra carpeta **validator**

21 - Copiamos toda la carpeta de **middleware** dentro de **internal** del servicio de workspace, a nuestro servicio de task, a nuestra carpeta **middleware** y verificamos las dependencias.

22 - Copiamos toda la carpeta de **auth** dentro de **internal** del servicio de workspace, a nuestro servicio de task, a nuestra carpeta **auth** y verificamos las dependencias.

23 - Copiamos toda la carpeta de **config** dentro de **internal** del servicio de workspace, a nuestro servicio de task, a nuestra carpeta **config** y verificamos las variables.

24 - Copiamos el archivo **main.go** dentro de **cmd/api** del servicio de workspace, a nuestro servicio de task, a nuestra carpeta **cmd/api** y verificamos las dependencias. Tambien realizamos los respectivos cambios en Repositorios, Use Cases y Handler

25 - Realizamos cambios en el archivo **Makefile**

```makefile
TASK_PATH=./services/task
TASK_MAIN_PATH=$(TASK_PATH)/cmd/api
TASK_MIGRATE_PATH=$(TASK_PATH)/migrations
TASK_DB_URL=

## run-task: Inicia el servidor de task service
run-task:
	TASK_HTTP_ADDR=:8093 \
	TASK_DATABASE_URL=${TASK_DB_URL} \
	JWT_SECRET=$(JWT) \
	go run $(TASK_MAIN_PATH)

## task-generate: Genera codigo Go con sqlc
task-generate:
	cd $(TASK_PATH) && sqlc generate
```

26 - Instalamos dependencias si es que falta, ejecutando el siguiente comando **(services/task)**:
crea los archivos:

- go.mod
- go.sum

```go
go mod tidy
```

27 - Regresamos a la ruta **devboard** y ejecutamos el comando de:
pero no genera nada

```makefile
make task-generate
```

28 - Hacer pruebas de los endpoints en postman:

    - http://localhost:8093/api/v1/tasks    POST
    - http://localhost:8093/api/v1/tasks/taskID    GET
    - http://localhost:8093/api/v1/tasks/taskID/comments   POST
    - http://localhost:8093/api/v1/tasks/taskID/comments   GET

29 - Cutover en el monilito, quitar lo que ya pasamos al servicio task
en el archivo **main.go** del monolito, eliminar, repositorios, casos de usos, hendler, rutas

30 - En el monolito ir a la carpeta **domain** y borrar los archivos: - comment.go - task_status_test.go - task_status.go - task.go

31 - En el monolito ir a la carpeta **usecase** y borrar los archivos: - comment_usecase.go - task_usecase.go

32 - En el monolito ir a la carpeta **handler** y borrar los archivos: - comment_handler.go - task_handler.go

33 - En el monolito ir a la carpeta **handler\postgres** y borrar los archivos: - comment_repository.go - task_repository.go

34 - En el monolito ir a la carpeta **handler\postgres\queries** y borrar los archivos: - comments.sql - tasks.sql

35 - En el monolito ir a la carpeta **handler\postgres\generated** y borrar los archivos: - comments.sql.go - tasks.sql.go

## Crear la BD de task

36 - Ir al archivo de **docker-compose.yml** y agregamos esto:

```yml
	task-postgres:
		image: postgres:16-alpine
		container_name: devboard_task_postgres
		environment:
		POSTGRES_USER: task-algo
		POSTGRES_PASSWORD: task-algo
		POSTGRES_DB: task-algo
		ports:
		- "5437:5432"
		volumes:
		- task_postgres_data:/var/lib/postgresql/data
		healthcheck:
		test: [ "CMD-SHELL", "pg_isready -U task-algo -d task-algo" ]
		interval: 5s
		timeout: 5s
		retries: 5

volumes:
  postgres_data:
  identity_postgres_data:
  workspace_postgres_data:
  task_postgres_data:
  redis_data
```

## Crear las migraciones en Task Service

37 - Creamos la carpeta **migrations** dentro de la carpeta **internal** del servicio task

38 - En el archivo Makefile crear los siguientes comandos:

```makefile
## task-generate: Genera codigo Go con sqlc
task-generate:
	cd $(TASK_PATH) && sqlc generate

## task-migrate-create: crea una nueva migración de task (OJO para usarlo su uso es: make migrate-create name=create_users)
task-migrate-create:
	migrate create -ext sql -dir $(TASK_MIGRATE_PATH) -seq $(name)

## task-migrate-up: aplicar todas las migraciones pendientes para task
task-migrate-up:
	migrate -path $(TASK_MIGRATE_PATH) -database "$(TASK_DB_URL)" up

## task-migrate-down: revertir la última migración de task
task-migrate-down:
	migrate -path $(TASK_MIGRATE_PATH) -database "$(TASK_DB_URL)" down 1
```

39 - Creamos la primera migracion con el comando:
Eso crea los archivos en la carpeta **migrations**:

- 000001_create_tasks.down.sql
- 000001_create_tasks.up.sql

```go
make task-migrate-create name=create_tasks
```

40 - Copiamos esto en el archivo de **000001_create_tasks.up.sql**

```sql

CREATE TABLE tasks (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id  UUID NOT NULL,
    title       TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status      TEXT NOT NULL DEFAULT 'todo',
                CHECK (status IN ('todo', 'in_progress','done')),
    assignee_id UUID,
    created_by  UUID NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tasks_project_id ON tasks(project_id);
CREATE INDEX idx_tasks_assignee_id ON tasks(assignee_id);
CREATE INDEX idx_tasks_created_by_id ON tasks(created_by);
```

41 - Copiamos esto en el archivo de **000001_create_tasks.down.sql**

```sql

DROP INDEX IF EXISTS idx_tasks_created_by;
DROP INDEX IF EXISTS idx_tasks_assignee_id;
DROP INDEX IF EXISTS idx_tasks_project_id;
DROP TABLE IF EXISTS tasks;

```

42 - Creamos la segunda migracion con el comando:
Eso crea los archivos en la carpeta **migrations**:

- 000001_create_comments.down.sql
- 000001_create_comments.up.sql

```go
make task-migrate-create name=create_comments
```

43 - Copiamos esto en el archivo de **000001_create_comments.up.sql**

```sql

CREATE TABLE comments (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id    UUID NOT NULL REFERENCES tasks(id)
               ON DELETE CASCADE,
    author_id  UUID NOT NULL,
    content    TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_comments_task_id ON comments(task_id);
CREATE INDEX idx_comments_author_id ON comments(author_id);
```

44 - Copiamos esto en el archivo de **000001_create_comments.down.sql**

```sql
DROP TABLE IF EXISTS comments;
```

45 - Eliminamos nuestro schema que era provisional **repository\postgres**

## Actualizamos el archivo sqlc.yml

46 - En el service task actualizamos el archivo **sqlc.yml** el schema:

```yml
schema: "migrations/"
```

## Actualizamos el archivo Makefie

47 - modificamos la ruta de la bd de task

```makefile
TASK_DB_URL=postgresql://task-algo:task-algo@localhost:5437/task-algo?sslmode=disable
```

## Levantamos la BD de task y aplicamos las migraciones

48 - Aplicamos el siguiente comando para levantar ls bds

```bash
docker compose up -d
```

49 - Aplicamos las migraciones, que crea las tablas en task-postgres con el siguiente comando:

```bash
make task-migrate-up
```

50 - Para comprbar que se crearon esas tablas en el contenedor de la bd de task

```bash
$ docker compose exec \
> task-postgres \
> psql \
> -U task \
> -d task \
> -c "\dt"
```

```sql
List of relations
 Schema |       Name        | Type  | Owner
--------+-------------------+-------+-------
 public | comments          | table | task
 public | schema_migrations | table | task
 public | tasks             | table | task
(3 rows)
```

## Migrar la data que tenemos en la bd devboar a la bd de task service db

51 - Nota: Cuando migramos datos relacionados debemos respetar sus dependencias para no violar los constraints (restricciones) del destino.
En pocas palabras tenemos que respetar el orden de migracion para que no genere problemas.

1. tasks
2. comments

## Copiar Task

```
docker compose exec -T postgres \
pg_dump \
-U postgres \
-d devboard \
--data-only \
--table=public.tasks \
--column-inserts \
--no-owner \
--no-privileges \
| docker compose exec -T task-postgres \
psql \
-U task \
-d task
```

## Copiar Comments

```
docker compose exec -T postgres \
pg_dump \
-U postgres \
-d devboard \
--data-only \
--table=public.comments \
--column-inserts \
--no-owner \
--no-privileges \
| docker compose exec -T task-postgres \
psql \
-U task \
-d task
```

## Verificar integridad

Origen:

```
docker compose exec postgres \
psql \
-U postgres \
-d devboard \
-c "
SELECT
    (SELECT COUNT(*) FROM tasks)
        AS tasks,
    (SELECT COUNT(*) FROM comments)
        AS comments;
"
```

Destino:

```
docker compose exec task-postgres \
psql \
-U task \
-d task \
-c "
SELECT
    (SELECT COUNT(*) FROM tasks)
        AS tasks,
    (SELECT COUNT(*) FROM comments)
        AS comments;
"
```

## Comprobar relaciones

Dentro de Task DB

```
docker compose exec task-postgres \
psql \
-U task \
-d task \
-c "
SELECT c.id
FROM comments c
LEFT JOIN tasks t
    ON t.id = c.task_id
WHERE t.id IS NULL;
"
```

## Crear migracion en monolito para detach

52 - Ejecutar el siguiente comando:

ese comando crea los siguientes archivos en la carpta **migrations** en el monolito:

- 000009_detach_task_data.down.sql
- 000009_detach_task_data.up.sql

```
make migrate-create name=detach_task_data
```

## Agregar query de detach up al archivo 000009_detach_task_data.up.sql

```
ALTER TABLE comments
    RENAME TO legacy_comments;

ALTER TABLE tasks
    RENAME TO legacy_tasks;
```

## Agregar query de detach down al archivo 000009_detach_task_data.down.sql

```
ALTER TABLE legacy_tasks
    RENAME TO tasks;

ALTER TABLE legacy_comments
    RENAME TO comments;
```

## Comprobamos tablas

```
docker compose exec postgres \
psql -U postgres -d devboard \
-c "\dt"
```

```sql
List of relations
 Schema |           Name           | Type  |  Owner
--------+--------------------------+-------+----------
 public | comments                 | table | postgres
 public | legacy_projects          | table | postgres
 public | legacy_users             | table | postgres
 public | legacy_workspace_members | table | postgres
 public | legacy_workspaces        | table | postgres
 public | schema_migrations        | table | postgres
 public | tasks                    | table | postgres
(7 rows)
```

53 - Ejecutamos el comando para que se corra las migraciones que modificamos (monolito):

```bash
make migrate-up
```

```sql
List of relations
 Schema |           Name           | Type  |  Owner
--------+--------------------------+-------+----------
 public | legacy_comments          | table | postgres
 public | legacy_projects          | table | postgres
 public | legacy_tasks             | table | postgres
 public | legacy_users             | table | postgres
 public | legacy_workspace_members | table | postgres
 public | legacy_workspaces        | table | postgres
 public | schema_migrations        | table | postgres
(7 rows)
```

# - Crear el Dockerfile para Identity, Workspace y Task

1 - Creamos el archivo **Dockerfile** entro de la ruta **services/identity**, **services/workspace** y **services/task**

y agregamos lo siguiente:

```dockerfile
FROM golang:1.27.3-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download \
    && go mod verify

COPY . .

RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/service \
    ./cmd/api


FROM alpine:3.24

RUN apk add --no-cache \
    ca-certificates \
    curl \
    && addgroup -S app \
    && adduser -S app -G app

WORKDIR /app

COPY --from=builder \
    /out/service \
    /app/service

USER app

ENTRYPOINT ["/app/service"]
```

## Dockerignore

2 - creamos los archivos **.dockerignore** en de la ruta **services/identity**, **services/workspace** y **services/task**

```
.git
.gitignore
coverage.out
*.log
tmp/
```

## Actualizamos el archivo docker-compose

```yml
services:
  # ====
  # Identity
  identity-postgres:
    image: postgres:16-alpine
    container_name: devboard_identity_postgres
    environment:
      POSTGRES_USER: identity
      POSTGRES_PASSWORD: identity
      POSTGRES_DB: identity
    ports:
      - "5435:5432"
    volumes:
      - identity_postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U identity -d identity"]
      interval: 5s
      timeout: 5s
      retries: 5
      start_period: 5s
    restart: unless-stopped
    networks:
      - identity-data

  redis:
    image: redis:7-alpine
    container_name: devboard_redis
    ports:
      - "6379:6379"
    command: redis-server --appendonly yes
    volumes:
      - redis_data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 5s
      retries: 5
      start_period: 5s
    restart: unless-stopped
    networks:
      - identity-data

  identity-migrate:
    image: migrate/migrate:v4.19.1
    container_name: devboard_identity_migrate
    volumes:
      - ./services/identity/migrations:/migrations:ro
    command:
      - "-path=/migrations"
      - "-database=postgresql://identity:identity@identity-postgres:5432/identity?sslmode=disable"
      - "up"
    depends_on:
      identity-postgres:
        condition: service_healthy
    restart: "no"
    networks:
      - identity-data

  identity:
    build:
      context: ./services/identity
      dockerfile: Dockerfile
    container_name: devboard_identity
    environment:
      IDENTITY_HTTP_ADDR: ":8091"

      IDENTITY_DATABASE_URL: >-
        postgresql://identity:identity@identity-postgres:5432/identity?sslmode=disable
      REDIS_ADDR: "redis:6379"
      JWT_SECRET: "${JWT_SECRET}"
    ports:
      - "8091:8091"
    depends_on:
      identity-postgres:
        condition: service_healthy

      redis:
        condition: service_healthy

      identity-migrate:
        condition: service_completed_successfully
    healthcheck:
      test: ["CMD", "curl", "-fsS", "http://localhost:8091/health"]
      interval: 10s
      timeout: 3s
      retries: 5
      start_period: 5s
    restart: unless-stopped
    networks:
      - app
      - identity-data

  # ===
  # Workspace
  workspace-postgres:
    image: postgres:16-alpine
    container_name: devboard_workspace_postgres
    environment:
      POSTGRES_USER: workspace
      POSTGRES_PASSWORD: workspace
      POSTGRES_DB: workspace
    ports:
      - "5436:5432"
    volumes:
      - workspace_postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U workspace -d workspace"]
      interval: 5s
      timeout: 5s
      retries: 5
      start_period: 5s
    restart: unless-stopped
    networks:
      - workspace-data

  workspace-migrate:
    image: migrate/migrate:v4.19.1
    container_name: devboard_workspace_migrate
    volumes:
      - ./services/workspace/migrations:/migrations:ro
    command:
      - "-path=/migrations"
      - "-database=postgresql://workspace:workspace@workspace-postgres:5432/workspace?sslmode=disable"
      - "up"
    depends_on:
      workspace-postgres:
        condition: service_healthy
    restart: "no"
    networks:
      - workspace-data

  workspace:
    build:
      context: ./services/workspace
      dockerfile: Dockerfile
    container_name: devboard_workspace
    environment:
      WORKSPACE_HTTP_ADDR: ":8092"

      WORKSPACE_DATABASE_URL: >-
        postgresql://workspace:workspace@workspace-postgres:5432/workspace?sslmode=disable
      JWT_SECRET: "${JWT_SECRET}"
    ports:
      - "8092:8092"
    depends_on:
      workspace-postgres:
        condition: service_healthy

      workspace-migrate:
        condition: service_completed_successfully
    healthcheck:
      test: ["CMD", "curl", "-fsS", "http://localhost:8092/health"]
      interval: 10s
      timeout: 3s
      retries: 5
      start_period: 5s
    restart: unless-stopped
    networks:
      - app
      - workspace-data

  # ====
  # Task
  task-postgres:
    image: postgres:16-alpine
    container_name: devboard_task_postgres
    environment:
      POSTGRES_USER: task
      POSTGRES_PASSWORD: task
      POSTGRES_DB: task
    ports:
      - "5437:5432"
    volumes:
      - task_postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U task -d task"]
      interval: 5s
      timeout: 5s
      retries: 5
      start_period: 5s
    restart: unless-stopped
    networks:
      - task-data

  task-migrate:
    image: migrate/migrate:v4.19.1
    container_name: devboard_task_migrate
    volumes:
      - ./services/task/migrations:/migrations:ro
    command:
      - "-path=/migrations"
      - "-database=postgresql://task:task@task-postgres:5432/task?sslmode=disable"
      - "up"
    depends_on:
      task-postgres:
        condition: service_healthy
    restart: "no"
    networks:
      - task-data

  task:
    build:
      context: ./services/task
      dockerfile: Dockerfile
    container_name: devboard_task
    environment:
      TASK_HTTP_ADDR: ":8093"

      TASK_DATABASE_URL: >-
        postgresql://task:task@task-postgres:5432/task?sslmode=disable
      JWT_SECRET: "${JWT_SECRET}"
    ports:
      - "8093:8093"
    depends_on:
      task-postgres:
        condition: service_healthy

      task-migrate:
        condition: service_completed_successfully
    healthcheck:
      test: ["CMD", "curl", "-fsS", "http://localhost:8093/health"]
      interval: 10s
      timeout: 3s
      retries: 5
      start_period: 5s
    restart: unless-stopped
    networks:
      - app
      - task-data

  # Legacy Database Monolito
  postgres:
    image: postgres:16-alpine
    container_name: devboard_postgres
    profiles:
      - legacy
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
      POSTGRES_DB: devboard
    ports:
      - "5434:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres -d devboard"]
      interval: 5s
      timeout: 5s
      retries: 5
    networks:
      - legacy-data

volumes:
  postgres_data:
  identity_postgres_data:
  workspace_postgres_data:
  task_postgres_data:
  redis_data:

networks:
  app:

  identity-data:
    internal: true

  workspace-data:
    internal: true

  task-data:
    internal: true

  legacy-data:
    internal: true
```

## Creamos nuevos comandos en el MakeFiles

```makefile
## compose-build: Construir contenedores
compose-build:
	docker compose build

## compose-up: Levanta nuestros contenedores
compose-up:
	docker compose up --build -d

## compose-down: Baja nuestros contenedores
compose-down:
	docker compose down

## compose-ps: Lista contenedores y su estado
compose-ps:
	docker compose ps

## compose-logs: Listado de logs para servicios
compose-logs:
	docker compose logs -f identity workspace task

## compose-config: Compose config
compose-config:
	docker compose config
```

# - Comunicación Sincrona gRPC

## crear estructura de carpetas para la creación de contratos del protocol buffers

1 - Creamos las siguientes carpetas:

**contracts -> proto -> devboard -> workspace -> v1**

2 - Dentro de la carpeta **v1**, creamos el archivo **workspace.proto**, donde vamos a definir el contrato.

3 - Agregamos esto al archivo **workspace.proto**, para definir el **contrato**

```proto
syntax = "proto3";
package devboard.workspace.v1;
option go_package = "github.com/valeroman/contracts/gen/go/devboard/workspace/v1;workspacev1";
```

4 - Ahora en el archivo **workspace.proto**, vamos a definir el **request**

- Primero vamos a solicitar la validación

```proto
  message ValidateProjectAccessRequest {
    string project_id = 1;
    string user_id = 2;
}
```

- Respondemos con esto

```proto
  message ValidateProjectAccessResponse {
    string workspace_id = 1;
}
```

5 - Definir el servicio

```proto
service WorkspaceService {
    rpc ValidateProjectAccess(ValidateProjectAccessRequest)
    returns (ValidateProjectAccessResponse);
}
```

## Instalamos Buf CLI

1 - Revisar la documentacion https://buf.build/docs/cli/installation/, para realizar la instalación

2 - Generamos el archivo **go.mod** en la ruta de **devboard/contracts/**, usando el siguiente comando:

**nota**: verificar que todos los go.mod tengan la misma version en todos los servicios y el de contracts.

```go
  go mod init github.com/valeroman/devboard/contracts
```

3 - Agregar al archivo **go.work**, los contratos de **protocol buffers**, ir a la ruta de **devboard** y ejecutar el siguiente comando:

```go
  go work use ./contracts
```

## Configuración Buf

1 - Creamos el archivo **buf.yaml**, en la ruta de **contracts** y agregamos esto:

```yaml
version: v2

modules:
  - path: proto

lint:
  use:
    - STANDARD

breaking:
  use:
    - FILE
```

2 - Ejecutamos estos comandos:

```bash
buf lint
buf build
```

3 - Creamos el archivo **buf.gen.yaml** para la configuración del archivo de generación de codigo y copiamos esto:

```yaml
version: v2

clean: true

plugins:
  - remote: buf.build/protocolbuffers/go
    out: gen/go
    opt: paths=source_relative
  - remote: buf.build/grpc/go
    out: gen/go
    opt: paths=source_relative
```

4 - Ejecutamos el siguiente comando para generar el codigo generado or Buf

estando en la ruta de **constracts**

```bash
buf generate
```

me genera los siguientes archivos:

- workspace_grpc.pb.go
- workspace.pb.go

5 - Ahora ejecutamos el siguiente comando para instalar o actualizar paquetes:

```go
go mod tidy
go mod verify
```

## Modificar el archivo Makefile

1 - Agregamos la ruta de la carpeta contracts:

```makefile
CONTRACTS_PATH=./contracts

...
## proto-format: Formatear proto
proto-format:
	cd $(CONTRACTS_PATH) && buf format -w

## proto-lint: APlicar linter al proto
proto-lint:
	cd $(CONTRACTS_PATH) && buf lint

## proto-build: Construir contractos
proto-build:
	cd $(CONTRACTS_PATH) && buf build

## proto-generate: Genera codigo con Buf e instala dependencias
proto-generate:
	cd  $(CONTRACTS_PATH) && buf generate
	cd  $(CONTRACTS_PATH) && go mod tidy

## proto-check: Revisa linter, construir buf y probar generación con test
proto-check:
	cd  $(CONTRACTS_PATH) && buf lint
	cd  $(CONTRACTS_PATH) && buf build
	cd  $(CONTRACTS_PATH) && go test ./...

```

## Crear Server (gRPC)

1 - Creamos el archivo **project_access_usecose.go** en la siguiente ruta **services/workspace/internal/usecase** y colocamos el siguuiente codigo

```go
package usecase

import (
	"context"
	"errors"

	"github.com/valeroman/devboard/services/workspace/internal/domain"
)

type ProjectAccessResult struct {
	WorkspaceID string
}

type ProjectAccessValidator interface {
	Validate(
		ctx context.Context,
		projectID string,
		userID string,
	) (ProjectAccessResult, error)
}

type ProjectAccessProjectRepository interface {
	GetByID(
		ctx context.Context,
		id string,
	) (*domain.Project, error)
}

type ProjectAccessWorkspaceRepository interface {
	GetMember(
		ctx context.Context,
		workspaceID string,
		userID string,
	) (*domain.WorkspaceMember, error)
}

type ProjectAccessUseCase struct {
	projectRepository   ProjectAccessProjectRepository
	workspaceRepository ProjectAccessWorkspaceRepository
}

func NewProjectAccessUseCase(
	projectRepository ProjectAccessProjectRepository,
	workspaceRepository ProjectAccessWorkspaceRepository,
) *ProjectAccessUseCase {
	return &ProjectAccessUseCase{
		projectRepository: projectRepository,

		workspaceRepository: workspaceRepository,
	}
}

func (
	uc *ProjectAccessUseCase,
) Validate(
	ctx context.Context,
	projectID string,
	userID string,
) (ProjectAccessResult, error) {
	if projectID == "" ||
		userID == "" {
		return ProjectAccessResult{},
			domain.ErrInvalidInput
	}

	project, err :=
		uc.projectRepository.GetByID(
			ctx,
			projectID,
		)

	if err != nil {
		return ProjectAccessResult{},
			err
	}

	member, err :=
		uc.workspaceRepository.GetMember(
			ctx,
			project.WorkspaceID,
			userID,
		)

	if err != nil {
		if errors.Is(
			err,
			domain.ErrNotFound,
		) {
			return ProjectAccessResult{},
				domain.ErrForbidden
		}

		return ProjectAccessResult{},
			err
	}

	if member == nil {
		return ProjectAccessResult{},
			domain.ErrForbidden
	}

	return ProjectAccessResult{
		WorkspaceID: project.WorkspaceID,
	}, nil
}

```

## Crear el Adapter del gRPC

1 - En la ruta de **services/workspace/internal** creamos la carpeta **transport** y dentro de transport creamos la carpeta **grpc** y dentro de la carpeta grpc creamos el archivo **handler.go**

y colocamos el siguiente codigo:

```go
package grpc

import (
	"context"
	"errors"
	"log/slog"

	workspacev1 "github.com/valeroman/devboard/contracts/gen/go/devboard/workspace/v1"
	"github.com/valeroman/devboard/services/workspace/internal/domain"
	"github.com/valeroman/devboard/services/workspace/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	workspacev1.UnimplementedWorkspaceServiceServer

	projectAccess usecase.ProjectAccessValidator
	logger        *slog.Logger
}

func NewHandler(
	projectAccess usecase.ProjectAccessValidator,
	logger *slog.Logger,
) *Handler {
	return &Handler{
		projectAccess: projectAccess,

		logger: logger,
	}
}

func (
	handler *Handler,
) ValidateProjectAccess(
	ctx context.Context,
	request *workspacev1.ValidateProjectAccessRequest,
) (
	*workspacev1.ValidateProjectAccessResponse,
	error,
) {
	if request == nil {
		return nil,
			status.Error(
				codes.InvalidArgument,
				"request requerida",
			)
	}

	result, err :=
		handler.projectAccess.Validate(
			ctx,
			request.GetProjectId(),
			request.GetUserId(),
		)

	if err != nil {
		return nil,
			handler.mapError(
				ctx,
				err,
			)
	}

	return &workspacev1.ValidateProjectAccessResponse{
		WorkspaceId: result.WorkspaceID,
	}, nil
}

func (
	handler *Handler,
) mapError(
	ctx context.Context,
	err error,
) error {
	switch {
	case errors.Is(
		err,
		domain.ErrInvalidInput,
	):
		return status.Error(
			codes.InvalidArgument,
			"project_id y user_id son requeridos",
		)

	case errors.Is(
		err,
		domain.ErrNotFound,
	):
		return status.Error(
			codes.NotFound,
			"project no encontrado",
		)

	case errors.Is(
		err,
		domain.ErrForbidden,
	):
		return status.Error(
			codes.PermissionDenied,
			"el usuario no tiene acceso al project",
		)

	default:
		handler.logger.ErrorContext(
			ctx,
			"workspace grpc request failed",
			slog.Any(
				"error",
				err,
			),
		)

		return status.Error(
			codes.Internal,
			"error interno",
		)
	}
}
```

## Configuración del gRPC de workspace

1 - Ir al archivo **config.go**, en la ruta **services/workspace/internal/config** y agregar el **GRPCAddr**

2 - Crear el archivo **server.go** en la ruta de **services/workspace/internal/transport/grpc** y agregamos esto

```go
package grpc

import (
	"fmt"
	"log/slog"
	"net"
	"time"

	workspacev1 "github.com/valeroman/devboard/contracts/gen/go/devboard/workspace/v1"

	googlegrpc "google.golang.org/grpc"
)

type Server struct {
	server   *googlegrpc.Server
	listener net.Listener
	logger   *slog.Logger
}

func NewServer(
	addr string,
	handler *Handler,
	logger *slog.Logger,
) (*Server, error) {
	listener, err := net.Listen("tcp", addr)

	if err != nil {
		return nil,
			fmt.Errorf(
				"listen grpc %s: %w",
				addr,
				err,
			)
	}

	grpcServer := googlegrpc.NewServer()

	workspacev1.RegisterWorkspaceServiceServer(
		grpcServer,
		handler,
	)

	return &Server{
		server:   grpcServer,
		listener: listener,
		logger:   logger,
	}, nil
}

func (server *Server) Serve() error {
	server.logger.Info(
		"starting workspace grpc server",
		slog.String(
			"addr",
			server.listener.Addr().String(),
		),
	)

	if err := server.server.Serve(server.listener); err != nil {
		return fmt.Errorf(
			"serve grpc: %w",
			err,
		)
	}

	return nil
}

func (server *Server) Shutdown(timeout time.Duration) {
	done := make(chan struct{})

	go func() {
		server.server.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		server.logger.Info(
			"workspace grpc server stopped gracefully",
		)

	case <-time.After(timeout):
		server.logger.Warn(
			"workspace grpc graceful shutdown timed out",
		)

		server.server.Stop()
	}
}

```

3 - Modificar o evolucionar el archivo **server.go** en la ruta **services/workspace/internal/server**, con lo del grpc, donde se modifica la funcion **Start** y se agrega **Shutdown**

```go
// Start server -> Metodo para arrancar el servidor http
func (server *Server) Start() error {

	server.logger.Info("iniciado servidor...", slog.String("addr", server.httpServer.Addr))

	err := server.httpServer.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server http: %w", err)
	}

	return nil
}

// Shutdown method
func (server *Server) Shutdown(ctx context.Context) error {

	if err := server.httpServer.Shutdown(ctx); err != nil {
		server.logger.Error("error en shutdown", slog.Any("error", err))
		return err
	}

	server.logger.Info("Servidor detenido correctamente")

	return nil
}
```

4 - modificamos el archivo **main.go** de la siguiente ruta **services/worspace/cmd/api**

```go
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/valeroman/devboard/services/workspace/internal/auth"
	"github.com/valeroman/devboard/services/workspace/internal/config"
	"github.com/valeroman/devboard/services/workspace/internal/handler"
	"github.com/valeroman/devboard/services/workspace/internal/middlewares"
	"github.com/valeroman/devboard/services/workspace/internal/repository/postgres"
	"github.com/valeroman/devboard/services/workspace/internal/server"
	"github.com/valeroman/devboard/services/workspace/internal/usecase"
	"github.com/valeroman/devboard/services/workspace/internal/validator"

	grpcTransport "github.com/valeroman/devboard/services/workspace/internal/transport/grpc"
)

func main() {
	if err := run(); err != nil {
		slog.Error("Workspace service se ha detenido", slog.Any("error", err))
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx := context.Background()

	pool, err := postgres.NewPool(
		ctx,
		postgres.Config{
			URL:             cfg.DatabaseURL,
			MaxConns:        10,
			MinConns:        2,
			MaxConnLifeTime: time.Hour,
			MaxConnIdleTime: 30 * time.Minute,
		},
	)

	if err != nil {
		return err
	}

	defer pool.Close()

	logger.Info("La conexion a la base de datos de workspace fue establecida")

	tokenVerifier, err := auth.NewTokenVerifier([]byte(cfg.JWTSecret))

	if err != nil {
		return err
	}

	// Repositories
	workspaceRepository := postgres.NewWorkspaceRepository(pool)
	projectRepository := postgres.NewProjectRepository(pool)

	// Use cases
	workspaceUseCase := usecase.NewWorkspaceUseCase(workspaceRepository)
	projectUseCase := usecase.NewProjectUseCase(projectRepository, workspaceRepository)
	projectAccessUseCase := usecase.NewProjectAccessUseCase(projectRepository, workspaceRepository)

	validate := validator.New()

	// Handlers
	workspaceHandler := handler.NewWorkspaceHandler(workspaceUseCase, validate, logger)
	projectHandler := handler.NewProjectHandler(projectUseCase, validate, logger)
	grpcHandler := grpcTransport.NewHandler(projectAccessUseCase, logger)

	// Server
	srv := server.New(cfg.HTTPAddr,
		server.Withlogger(logger),
		server.WithReadHeaderTimeout(15*time.Second),
		server.WithWriteTimeout(30*time.Second))

	grpcServer, err := grpcTransport.NewServer(cfg.GRPCAddr, grpcHandler, logger)

	if err != nil {
		return err
	}

	authMiddleware := middlewares.RequireAuth(tokenVerifier, logger)

	// Rutas Publicas
	srv.RegisterRoutes("GET /health", http.HandlerFunc(
		func(write http.ResponseWriter, _ *http.Request) {
			write.Header().Set(
				"Content-Type", "application/json",
			)

			write.WriteHeader(
				http.StatusOK,
			)

			_, _ = write.Write([]byte(
				`{"status": "ok", "service": "workspace"}`,
			))
		},
	))

	// Workspace routes
	srv.RegisterRoutes("POST /api/v1/workspaces", middlewares.Chain(
		http.HandlerFunc(workspaceHandler.Create),
		authMiddleware,
	))

	srv.RegisterRoutes("GET /api/v1/workspaces", middlewares.Chain(
		http.HandlerFunc(workspaceHandler.List),
		authMiddleware,
	))

	srv.RegisterRoutes("GET /api/v1/workspaces/{id}", middlewares.Chain(
		http.HandlerFunc(workspaceHandler.Get),
		authMiddleware,
	))

	srv.RegisterRoutes("POST /api/v1/workspaces/{id}/members", middlewares.Chain(
		http.HandlerFunc(workspaceHandler.AddMember),
		authMiddleware,
	))

	srv.RegisterRoutes("GET /api/v1/workspaces/{id}/members", middlewares.Chain(
		http.HandlerFunc(workspaceHandler.ListMembers),
		authMiddleware,
	))

	srv.RegisterRoutes("DELETE /api/v1/workspaces/{id}/members/{userID}", middlewares.Chain(
		http.HandlerFunc(workspaceHandler.RemoveMember),
		authMiddleware,
	))

	// Project Routes
	srv.RegisterRoutes("POST /api/v1/projects", middlewares.Chain(
		http.HandlerFunc(projectHandler.Create),
		authMiddleware,
	))

	srv.RegisterRoutes("GET /api/v1/projects/{id}", middlewares.Chain(
		http.HandlerFunc(projectHandler.Get),
		authMiddleware,
	))

	srv.RegisterRoutes("GET /api/v1/workspaces/{id}/projects", middlewares.Chain(
		http.HandlerFunc(projectHandler.ListByWorkspace),
		authMiddleware,
	))

	// HTTP middleware
	srv.Use(
		middlewares.SecurityHeaders(middlewares.SecurityConfig{
			AllowedOrigins: []string{
				"http://localhost:3000",
				"http://localhost:5173",
			},
		}),
	)

	srv.Use(
		middlewares.Logger(logger),
	)

	srv.Use(
		middlewares.Recovery(logger),
	)

	logger.Info("Se inicio WORKSPACE SERVICE", slog.String("addr", cfg.HTTPAddr))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	defer stop()

	errCh := make(chan error, 2)

	// Inicializando el servidor http
	go func() {
		errCh <- srv.Start()
	}()

	// Inicializando el servidor grpc
	go func() {
		errCh <- grpcServer.Serve()
	}()

	var runErr error

	select {
	case <-ctx.Done():
		logger.Info("Señal de apagado workspace recibida")

	case err := <-errCh:
		if err != nil {
			runErr = err
			logger.Error("El servidor workspace falló", slog.Any("error", err))
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	grpcServer.Shutdown(5 * time.Second)

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("http shutdown falló", slog.Any("error", err))

		if runErr == nil {
			runErr = err
		}
	}

	return runErr

}

```

5 - En la ruta de **services/worspace**, ejecutar el comando

```go
  go mod tidy
```

si da error como estos

```go
go mod tidy
go: finding module for package google.golang.org/grpc/status
go: finding module for package google.golang.org/grpc
go: finding module for package github.com/valeroman/devboard/contracts/gen/go/devboard/workspace/v1
go: finding module for package google.golang.org/grpc/codes
go: downloading github.com/valeroman/devboard v1.3.0
go: found google.golang.org/grpc in google.golang.org/grpc v1.85.0
go: found google.golang.org/grpc/codes in google.golang.org/grpc v1.85.0
go: found google.golang.org/grpc/status in google.golang.org/grpc v1.85.0
go: finding module for package github.com/valeroman/devboard/contracts/gen/go/devboard/workspace/v1
go: github.com/valeroman/devboard/services/workspace/internal/transport/grpc imports
        github.com/valeroman/devboard/contracts/gen/go/devboard/workspace/v1: github.com/valeroman/devboard@v1.3.0: verifying module: github.com/valeroman/devboard@v1.3.0: reading https://sum.golang.org/lookup/github.com/valeroman/devboard@v1.3.0: 404 Not Found
        server response:
        not found: github.com/valeroman/devboard@v1.3.0: invalid version: git ls-remote -q --end-of-options https://github.com/valeroman/devboard in /tmp/gopath/pkg/mod/cache/vcs/b0a7c1919f3f3cdd71dd557853fdc052548e3b1485fed4377159da6dfdb7c054: exit status 128:
                fatal: unable to access 'https://github.com/valeroman/devboard/': Could not resolve host: github.com
```

ejecutar lo siguiente:

```go
go mod edit -require=github.com/valeroman/devboard/contracts@v0.0.0
go mod edit -replace=github.com/valeroman/devboard/contracts=../../contracts
```

y despues un **go mod tidy**

6 - modificamos el archivo **docker-compose** agregando estos en workspace:

```yaml
workspace:
  build:
    context: .
    dockerfile: services/workspace/Dockerfile
  container_name: devboard_workspace
  environment:
    WORKSPACE_HTTP_ADDR: ":8092"
    WORKSPACE_GRPC_ADDR: ":9072"

    WORKSPACE_DATABASE_URL: >-
      postgresql://workspace:workspace@workspace-postgres:5432/workspace?sslmode=disable
    JWT_SECRET: "${JWT_SECRET}"
  ports:
    - "8092:8092"
    - "9072:9072"
  depends_on:
    workspace-postgres:
      condition: service_healthy

    workspace-migrate:
      condition: service_completed_successfully
  healthcheck:
    test: ["CMD", "curl", "-fsS", "http://localhost:8092/health"]
    interval: 10s
    timeout: 3s
    retries: 5
    start_period: 5s
  restart: unless-stopped
  networks:
    - app
    - workspace-data
```

Y tambien modificamos el archivo **Dockerfile** del **services/workspace**

```yaml
FROM golang:1.27.3-alpine AS builder

WORKDIR /src

COPY contracts/go.mod contracts/go.sum ./contracts/
COPY services/workspace/go.mod services/workspace/go.sum ./services/workspace/

WORKDIR /src/services/workspace

RUN go mod download \
&& go mod verify

WORKDIR /src

COPY contracts ./contracts
COPY services/workspace ./services/workspace

WORKDIR /src/services/workspace

RUN CGO_ENABLED=0 \
GOOS=linux \
go build \
-trimpath \
-ldflags="-s -w" \
-o /out/service \
./cmd/api


FROM alpine:3.24

RUN apk add --no-cache \
ca-certificates \
curl \
&& addgroup -S app \
&& adduser -S app -G app

WORKDIR /app

COPY --from=builder \
/out/service \
/app/service

USER app

ENTRYPOINT ["/app/service"]
```

Y ahora ejecutamos el comando **make dockercompose-up**

## Instalar gRPCurl para pruebas

1 - Instalar el siguiente paquete con este comando:

```bash
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

## Registrar Reflection

1 - Ir a la ruta de **services/workspace/transport/grpc** y en el archivo **server.go**, agregamos esto

```go
reflection.Register(
	grpcServer,
)
```

Y ahora ejecutamos el comando **make dockercompose-up**

## Probando que el servicio este registrado

```
grpcurl \
-plaintext \
localhost:9072 \
list


grpcurl \
-plaintext \
localhost:9072 \
list \
devboard.workspace.v1.WorkspaceService
```

## IDs existentes

```
docker compose exec workspace-postgres \
psql \
-U workspace \
-d workspace \
-c "
SELECT
    p.id AS project_id,
    wm.user_id,
    p.workspace_id
FROM projects p
JOIN workspace_members wm
    ON wm.workspace_id = p.workspace_id
LIMIT 1;
"
```

## Probando ValidateProjectAccess

```
grpcurl \
-plaintext \
-d '{
  "project_id": "$PROJECT_ID",
  "user_id": "$USER_ID"
}' \
localhost:9092 \
devboard.workspace.v1.WorkspaceService/ValidateProjectAccess
```

## Probando un Project inexistente

```
grpcurl \
-plaintext \
-d '{
  "project_id": "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa",
  "user_id": "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb"
}' \
localhost:9092 \
devboard.workspace.v1.WorkspaceService/ValidateProjectAccess
```

## Probando usuario sin acceso

```
grpcurl \
-plaintext \
-d "{
  "project_id": "$PROJECT_ID",
  "user_id": "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb"
}" \
localhost:9092 \
devboard.workspace.v1.WorkspaceService/ValidateProjectAccess
```

## Probando input vacio

```
grpcurl \
-plaintext \
-d '{}' \
localhost:9092 \
devboard.workspace.v1.WorkspaceService/ValidateProjectAccess
```

# - Implementando Task y Workspace (Services) gRPC Client, para que se comuniquen entre si.

## Creando port en Task service

1 - En la ruta de **services/task/internal/usecase/**, creamos el archivo **project_access.go**

```go
package usecase

import "context"

type ProjectAccess struct {
	WorkspaceID string
}

type ProjectAccessChecker interface {
	ValidateProjectAccess(ctx context.Context, projectID string, userID string) (ProjectAccess, error)
}
```

2 - Ahora modificamos el archivo **task_usecase.go**, agregando lo siguiente:

```go
type TaskUseCase struct {
	commands      TaskCommandRepository
	queries       TaskQueryRepository
	projectAccess ProjectAccessChecker  // nuevo
}

// NewTaskUseCase (constructor)
func NewTaskUseCase(commands TaskCommandRepository, queries TaskQueryRepository, projectAccess ProjectAccessChecker) *TaskUseCase {
	return &TaskUseCase{
		commands:      commands,
		queries:       queries,
		projectAccess: projectAccess, // nuevo
	}
}

// CreateTask método para crear tarea
func (uc *TaskUseCase) CreateTask(ctx context.Context, projectID, title, description, createdBy string) (*domain.Task, error) {

	_, err := uc.projectAccess.ValidateProjectAccess(ctx, projectID, createdBy) // nuevo

	if err != nil {                     // nuevo
		return nil, err
	}

	task := domain.NewTask(projectID, title, description, createdBy)

	if err := uc.commands.Create(ctx, task); err != nil {
		return nil, fmt.Errorf("creando tarea: %w", err)
	}

	return task, nil
}

```

## Creando el adapter de gRPC en Task service

1 - Creamos las carpetas **client/workspace**, dentro de la ruta **services/task/internal**

2 - Creamos el archivo **client.go**

```go
package workspace

import (
	"fmt"

	workspacev1 "github.com/valeroman/devboard/contracts/gen/go/devboard/workspace/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	connection *grpc.ClientConn
	client     workspacev1.WorkspaceServiceClient
}

func NewClient(target string) (*Client, error) {
	connection, err :=
		grpc.NewClient(
			target,
			grpc.WithTransportCredentials(
				insecure.NewCredentials(),
			),
		)

	if err != nil {
		return nil,
			fmt.Errorf(
				"create workspace grpc client: %w",
				err,
			)
	}

	return &Client{
		connection: connection,
		client:     workspacev1.NewWorkspaceServiceClient(connection),
	}, nil
}

```
