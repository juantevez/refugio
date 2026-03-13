# Refugio Core API 🐾

Backend modular desarrollado en **Go** para la gestión integral de un refugio de animales. El sistema permite el registro de rescates, seguimiento de estados de salud, y la gestión de donaciones vinculadas a animales específicos.

## 🏗️ Arquitectura

El proyecto sigue los principios de **Arquitectura Hexagonal (Ports & Adapters)** y **Domain-Driven Design (DDD)**, asegurando que la lógica de negocio esté aislada de las dependencias externas (DB, HTTP Frameworks).

* **Domain**: Entidades de negocio y reglas principales.
* **Application/Services**: Casos de uso del sistema.
* **Infrastructure/Handlers**: Adaptadores de entrada (Gin Gonic).
* **Infrastructure/Repositories**: Adaptadores de salida (PostgreSQL con Sqlx).

## 🚀 Tecnologías

* **Lenguaje:** Go 1.26+
* **Framework Web:** Gin Gonic
* **Base de Datos:** PostgreSQL
* **Migraciones:** Flyway
* **Contenedores:** Docker & Docker Compose

## 🛠️ Configuración del Entorno

1. **Clonar el repositorio:**
   ```bash
   git clone [https://github.com/tu-usuario/refugio-core.git](https://github.com/tu-usuario/refugio-core.git)
   cd refugio-core

2. **Variables de Entorno:**

Crea un archivo .env en la raíz basado en .env.example:

```
Fragmento de código
DB_HOST=localhost
DB_PORT=5432
DB_USER=shelter_user
DB_PASSWORD=shelter_pass
DB_NAME=shelter_db
Levantar Infraestructura (Docker):
```

   ```
docker-compose up -d
   ```


### Correr la aplicación:

   ```bash
go run cmd/main.go
   ```bash

### 🔌 API Endpoints

   ``
Animales
POST /api/v1/animals - Registrar un nuevo rescate.

GET /api/v1/animals - Listar todos los animales en el refugio.

GET /api/v1/animals/:id - Obtener detalle de un animal.
   ``

Donaciones
POST /api/v1/donations - Registrar una donación (soporta Transferencia y Mercado Pago).

📊 ### Base de Datos
El esquema de la base de datos está versionado con Flyway.
Principales esquemas:

animal_management: Tablas de animales y su historial de estados.

adoptions_donations: Gestión de registros financieros y adopciones.

🧪 ### Tests
Para ejecutar la suite de tests:

   ```bash
go test ./...
   ```bash

### Desarrollado con ❤️ para ayudar a los animales.
