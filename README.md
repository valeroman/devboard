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