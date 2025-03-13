# API Documentos - Estructura y Uso

## Introducción

Esta API actúa como intermediaria entre un cliente y un servicio de gestión documental (aplicado actualmente al servicio Alfresco), gestionando subida, listado y descarga de documentos y el proceso de autenticación para ingresar a solicitudes propias del servicio. Está diseñada para ser escalable y fácil de mantener.

## Estructura del Proyecto

```
APIDOCUMENTOS/
│
├── 📁 application/                             # Capa de Aplicación (Casos de Uso)
│   └── 📁 documento/                           # Orquesta lógica de negocio
│       ├── 📄 servicio_documento.go            # Valida reglas de negocio y delega a la capa de dominio
│       └── 📄 servicio_dto.go                  # DTOs específicos de Alfresco (mapeo JSON → Dominio)
│
├── 📁 domain/                                  # Capa de Dominio (Lógica Central)
│   ├── 📁 auth/
│   │   └── 📄 servicio_auth.go                 # Define lógica de autenticación: tickets, roles
│   │
│   └── 📁 documento/
│       ├── 📄 entidad.go                        # Modelo de dominio: Define estructura Documento, validaciones
│       └── 📄 repositorio.go                     # Interfaz abstracta (AlmacenamientoDocumentos)
│
├── 📁 infrastructure/                            # Capa de Infraestructura
│   ├── 📁 api/
│   │   └── 📁 handlers/                           # Controladores HTTP
│   │       ├── 📄 auth_handler.go                 # Maneja endpoints de autenticación
│   │       ├── 📄 documento_handler.go            # Maneja endpoints de documentos
│   │       ├── 📄 busqueda_descarga_handler.go    # Maneja búsqueda y descarga de archivos
│   │       └── 📄 lote_documentos_handler.go      # Maneja carga de múltiples documentos
│   │
│   ├── 📁 persistence/                        # Implementación de almacenamiento externo
│   │   ├── 📁 servicio/                       # Adaptador para Alfresco
│   │   │   ├── 📄 cliente.go                  # Cliente HTTP para Alfresco
│   │   │   ├── 📄 mock_cliente.go             # Simulación de respuestas para tests
│   │   │   ├── 📄 repositorio.go              # Implementa AlmacenamientoDocumentos
│   │   │   └── 📄 auth_cliente.go             # Gestión de autenticación con Alfresco
│   │   │
│   │   └── 📁 db/
│   │       ├── 📄 bd_mongo.go                 # Gestión de base de datos MongoDB
│   │       └── 📄 mongo_dto.go                # DTOs para persistencia en MongoDB
│   │
│   └── 📁 routes/
│       └── 📄 API_router.go                   # Define rutas, middlewares y enlaza handlers
|
|__📁settings/
|  └── 📄 appsettings.go                       # Configura las variables de entorno para despliegue de API en Azure Devops
|  └── 📄 appsettings.json                     # Informa las variables de entorno globales de la aplicación (En formato JSON)
│
├── 📁 shared/                                 # Utilidades transversales
│   ├── 📁 middleware/
│   │   └── 📄 logging.go                      # Registra métricas, tiempos y errores
│   │
│   ├── 📁 config/
│   │   └── 📄 config.go                       # Centraliza configuración (Variables de Entorno se obtienen de archivo appsettings.json)
│   │
│   └── 📁 logger/
│       └── 📄 logger.go                       # Servicio de logging unificado
│
├── 📄 appsettings-example.json                # Plantilla de variables de entorno
├── 📄 go.mod                                  # Define dependencias del proyecto
├── 📄 go.sum                                  # Checksum de dependencias
└── 📄 main.go                                 # Punto de entrada: Inicializa componentes
```

## Descripción de los Componentes Principales

### 1. Capa de Aplicación (`application/`)

Esta capa implementa los casos de uso de la aplicación, orquestando la lógica de negocio:

- **`servicio_documento.go`**: Valida las reglas de negocio y delega operaciones a la capa de dominio.
- **`servicio_dto.go`**: Define los objetos de transferencia de datos (DTOs) específicos para la integración con Alfresco.

### 2. Capa de Dominio (`domain/`)

Contiene la lógica central y los modelos de negocio de la aplicación:

- **`auth/servicio_auth.go`**: Define toda la lógica relacionada con autenticación, incluyendo generación y validación de tickets y gestión de roles.
- **`documento/entidad.go`**: Establece el modelo de dominio, definiendo la estructura de Documento y sus validaciones.
- **`documento/repositorio.go`**: Define interfaces abstractas que establecen contratos para las operaciones CRUD.

### 3. Capa de Infraestructura (`infrastructure/`)

Implementa los detalles técnicos y la comunicación con sistemas externos:

- **Handlers HTTP**: Gestionan los endpoints para autenticación, documentos, búsqueda/descarga y manejo de lotes.
- **Persistencia**:
  - **Servicio Alfresco**: Implementa la comunicación con el servicio de gestión documental.
  - **Base de datos MongoDB**: Gestiona la persistencia local de información.
- **Router**: Configura las rutas de la API, middleware global y enlaza los controladores.

### 4. Utilidades Compartidas (`shared/`)

Componentes transversales utilizados por las distintas capas:

- **Middleware**: Interceptores de solicitudes HTTP para logging, métricas y gestión de errores.
- **Configuración**: Centraliza parámetros como URLs, claves de API y límites operativos.
- **Logging**: Proporciona un servicio unificado para el registro de eventos y errores.

### 5. Archivos Raíz

- **`.env.example`**: Plantilla para configurar variables de entorno.
- **`go.mod` y `go.sum`**: Gestión de dependencias de Go.
- **`main.go`**: Punto de entrada que inicializa la configuración, el logger, el router y el servidor HTTP.