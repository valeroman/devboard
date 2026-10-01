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
