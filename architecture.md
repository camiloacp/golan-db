# Architecture Diagrams

---

## 1. Arquitectura en Capas

Este diagrama muestra cómo está organizado el proyecto en **4 capas** que fluyen de arriba hacia abajo:

| Capa | Carpeta | Responsabilidad |
|------|---------|-----------------|
| **Entry Point** | `main.go` | Inicializa la DB e inyecta dependencias |
| **Domain Layer** | `pkg/` | Contiene la lógica de negocio (servicios) y los contratos (interfaces) |
| **Infrastructure Layer** | `storage/` | Implementa los contratos conectándose a PostgreSQL |
| **Base de Datos** | PostgreSQL | Almacena los datos en 3 tablas relacionadas |

**Cómo leer las flechas:**
- `→` flecha sólida: una capa llama a otra directamente
- `⇒` flecha gruesa: escritura de datos en la base de datos
- `··→` flecha punteada: relación de clave foránea (FK) entre tablas

**Colores:**
- 🟣 Violeta oscuro: `main.go` (punto de entrada)
- 🔵 Azul: Servicios del dominio (`pkg/`)
- 🟢 Verde: Implementaciones de storage CRUD simples
- 🟠 Naranja: `PsqlInvoice` — orquesta la transacción atómica
- 💜 Lila: Singleton de conexión a la DB (`sync.Once`)
- ⬛ Gris oscuro: Tablas de PostgreSQL

```mermaid
flowchart TB
    classDef entry    fill:#312e81,stroke:#818cf8,stroke-width:2px,color:#e0e7ff,font-weight:bold
    classDef service  fill:#1e3a5f,stroke:#60a5fa,stroke-width:2px,color:#dbeafe
    classDef infra    fill:#14532d,stroke:#4ade80,stroke-width:2px,color:#dcfce7
    classDef tx       fill:#7c2d12,stroke:#fb923c,stroke-width:2px,color:#ffedd5,font-weight:bold
    classDef pool     fill:#4a1d96,stroke:#c084fc,stroke-width:2px,color:#f3e8ff
    classDef db       fill:#1c1917,stroke:#a78bfa,stroke-width:2px,color:#ede9fe

    %% ── Layer 1: Entry ───────────────────────────────────────────
    main["main.go\nInit DB · Dependency Injection"]

    %% ── Layer 2: Domain Services ─────────────────────────────────
    subgraph PKG["  pkg/  ·  Domain Layer  "]
        direction LR
        svcProduct["product.Service\nMigrate · CRUD"]
        svcInvoice["invoice.Service\nCreate"]
    end

    %% ── Layer 3: Infrastructure ──────────────────────────────────
    subgraph STORE["  storage/  ·  Infrastructure Layer  "]
        direction LR
        pool["Pool · sync.Once\nNewPostgresDB()"]
        psqlProduct["PsqlProduct\nCRUD directo"]
        psqlInvoice["PsqlInvoice\nOrquesta TX"]
        psqlHeader["PsqlInvoiceHeader\nCreateTx()"]
        psqlItem["PsqlInvoiceItem\nCreateTx()"]
    end

    %% ── Layer 4: Database ────────────────────────────────────────
    subgraph DB["  PostgreSQL · godb  "]
        direction LR
        tProd[("products")]
        tHead[("invoice_headers")]
        tItem[("invoice_items")]
    end

    %% ── Connections ──────────────────────────────────────────────
    main --> svcProduct & svcInvoice

    svcProduct --> psqlProduct
    svcInvoice --> psqlInvoice

    pool --> psqlProduct & psqlInvoice & psqlHeader & psqlItem

    psqlInvoice --> psqlHeader & psqlItem

    psqlProduct ==> tProd
    psqlHeader  ==> tHead
    psqlItem    ==> tItem

    tItem -. "FK: invoice_header_id" .-> tHead
    tItem -. "FK: product_id" .-> tProd

    %% ── Styles ───────────────────────────────────────────────────
    class main entry
    class svcProduct,svcInvoice service
    class psqlProduct,psqlHeader,psqlItem infra
    class psqlInvoice tx
    class pool pool
    class tProd,tHead,tItem db
```

---

## 2. Flujo de Transaccion Atomica

Cuando se crea una factura, el sistema debe insertar datos en **dos tablas distintas** dentro de una sola transacción. Esto garantiza que nunca quede una cabecera sin sus ítems (o viceversa) ante un error.

El diagrama muestra cada paso del proceso y qué ocurre si algo falla:
- **Camino verde** (`→`): todo salió bien → `COMMIT`, los datos quedan guardados
- **Camino rojo** (`→`): algo falló → `ROLLBACK`, se deshacen **todos** los cambios

```mermaid
flowchart TD
    classDef step   fill:#1e3a5f,stroke:#60a5fa,stroke-width:2px,color:#dbeafe
    classDef ok     fill:#14532d,stroke:#4ade80,stroke-width:2px,color:#dcfce7,font-weight:bold
    classDef fail   fill:#7f1d1d,stroke:#f87171,stroke-width:2px,color:#fee2e2,font-weight:bold
    classDef decide fill:#451a03,stroke:#fb923c,stroke-width:2px,color:#ffedd5
    classDef tx     fill:#312e81,stroke:#818cf8,stroke-width:2px,color:#e0e7ff

    START(["Inicio: serviceInvoice.Create(model)"])

    BEGIN["PsqlInvoice abre la transaccion\ndb.Begin()"]

    INSERT_HEADER["PsqlInvoiceHeader.CreateTx\nINSERT INTO invoice_headers\nRetorna header.ID"]

    D1{{"¿Header\ninsertado?"}}

    INSERT_ITEMS["PsqlInvoiceItem.CreateTx\nINSERT INTO invoice_items\npor cada item"]

    D2{{"¿Items\ninsertados?"}}

    COMMIT["tx.Commit()\nCambios guardados en PostgreSQL"]
    SUCCESS(["Exito: factura creada"])

    ROLLBACK_1["tx.Rollback()\nSe deshace el header"]
    ERR_1(["Error: 'failed to create header'"])

    ROLLBACK_2["tx.Rollback()\nSe deshacen items y header"]
    ERR_2(["Error: 'failed to create items'"])

    START         --> BEGIN
    BEGIN         --> INSERT_HEADER
    INSERT_HEADER --> D1

    D1 -- "SI" --> INSERT_ITEMS
    D1 -- "NO" --> ROLLBACK_1
    ROLLBACK_1   --> ERR_1

    INSERT_ITEMS --> D2
    D2 -- "SI" --> COMMIT
    D2 -- "NO" --> ROLLBACK_2
    ROLLBACK_2   --> ERR_2

    COMMIT       --> SUCCESS

    class START,BEGIN tx
    class INSERT_HEADER,INSERT_ITEMS step
    class D1,D2 decide
    class COMMIT,SUCCESS ok
    class ROLLBACK_1,ROLLBACK_2,ERR_1,ERR_2 fail
```
