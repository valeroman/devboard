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
