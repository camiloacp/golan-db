# golang_db

Proyecto en Go para gestionar una base de datos PostgreSQL con tablas de productos y facturacion.

## Estructura del proyecto

```
golang_db/
├── main.go
├── pkg/
│   ├── product/          # Modelo y servicio de productos
│   ├── invoiceheader/    # Modelo y servicio de cabeceras de factura
│   └── invoiceitem/      # Modelo y servicio de items de factura
└── storage/
    ├── storage.go              # Conexion a PostgreSQL (singleton)
    ├── psql_product.go         # Storage de productos
    ├── psql_invoiceheader.go   # Storage de cabeceras de factura
    └── psql_invoiceitem.go     # Storage de items de factura
```

## Base de datos

**Conexion:** `postgres://golang_db_user:golang_db_password@localhost:7530/godb?sslmode=disable`

## Modelos

### Product

| Campo        | Tipo        | Restriccion          |
|--------------|-------------|----------------------|
| id           | SERIAL      | PRIMARY KEY, NOT NULL|
| name         | VARCHAR(25) | NOT NULL             |
| observations | VARCHAR(100)|                      |
| price        | INT         | NOT NULL             |
| created_at   | TIMESTAMP   | NOT NULL, DEFAULT now() |
| updated_at   | TIMESTAMP   |                      |

### InvoiceHeader

| Campo      | Tipo        | Restriccion             |
|------------|-------------|-------------------------|
| id         | SERIAL      | PRIMARY KEY, NOT NULL   |
| client     | VARCHAR(25) | NOT NULL                |
| created_at | TIMESTAMP   | NOT NULL, DEFAULT now() |
| updated_at | TIMESTAMP   |                         |

### InvoiceItem

| Campo             | Tipo      | Restriccion                    |
|-------------------|-----------|--------------------------------|
| id                | SERIAL    | PRIMARY KEY, NOT NULL          |
| invoice_header_id | INT       | NOT NULL, FK -> invoice_headers|
| product_id        | INT       | NOT NULL, FK -> products       |
| created_at        | TIMESTAMP | NOT NULL, DEFAULT now()        |
| updated_at        | TIMESTAMP |                                |

**Foreign keys en invoice_items:**
- `invoice_header_id` -> `invoice_headers(id)` ON UPDATE RESTRICT ON DELETE RESTRICT
- `product_id` -> `products(id)` ON UPDATE RESTRICT ON DELETE RESTRICT

## Migraciones (ya ejecutadas)

El siguiente codigo se uso para crear las tablas y fue removido de `main.go` ya que las tablas ya existen:

```go
storageProduct := storage.NewPsqlProduct(storage.Pool())
serviceProduct := product.NewService(storageProduct)

if err := serviceProduct.Migrate(); err != nil {
    log.Fatalf("product.Migrate: %v", err)
}

storageInvoiceHeader := storage.NewPsqlInvoiceHeader(storage.Pool())
serviceInvoiceHeader := invoiceheader.NewService(storageInvoiceHeader)

if err := serviceInvoiceHeader.Migrate(); err != nil {
    log.Fatalf("invoiceheader.Migrate: %v", err)
}

storageInvoiceItem := storage.NewPsqlInvoiceItem(storage.Pool())
serviceInvoiceItem := invoiceitem.NewService(storageInvoiceItem)

if err := serviceInvoiceItem.Migrate(); err != nil {
    log.Fatalf("invoiceitem.Migrate: %v", err)
}
```

Cada `Migrate()` ejecuta un `CREATE TABLE IF NOT EXISTS` con el esquema correspondiente.

## Insertar registros

Para crear un nuevo producto, se instancia el modelo y se llama al metodo `Create` del servicio:

```go
storageProduct := storage.NewPsqlProduct(storage.Pool())
serviceProduct := product.NewService(storageProduct)

m := &product.Model{
    Name:         "Curso de Go Poo",
    Price:        230,
    Observations: "On Fire",
}

if err := serviceProduct.Create(m); err != nil {
    log.Fatalf("product.Create: %v", err)
}

fmt.Printf("Producto creado con ID: %d\n", m.ID)
```

Los campos `id`, `created_at` y `updated_at` se manejan automaticamente (id es autoincrementable, created_at toma `now()`).

## Obtener todos los registros

Para obtener todos los productos, se llama al metodo `GetAll` del servicio:

```go
storageProduct := storage.NewPsqlProduct(storage.Pool())
serviceProduct := product.NewService(storageProduct)

ms, err := serviceProduct.GetAll()
if err != nil {
    log.Fatalf("product.GetAll: %v", err)
}
fmt.Printf("%+v\n", ms)
```

Retorna un slice de `product.Model` con todos los productos almacenados en la tabla.

## Obtener un registro por ID

Para obtener un producto por su ID, se llama al metodo `GetByID` del servicio. Se usa un `switch` para manejar el caso en que no exista el registro:

```go
m, err := serviceProduct.GetByID(4)
switch {
case errors.Is(err, sql.ErrNoRows):
    fmt.Println("Product not found with this ID")
case err != nil:
    log.Fatalf("product.GetByID: %v", err)
default:
    fmt.Printf("%+v\n", m)
}
```

Si el ID no existe en la tabla, `sql.ErrNoRows` permite detectarlo sin que el programa falle. Requiere importar `database/sql` y `errors`.

## Actualizar un registro

Para actualizar un producto, se instancia el modelo con el `ID` del registro a modificar y los campos actualizados, luego se llama al metodo `Update` del servicio:

```go
m := &product.Model{
    ID:    10,
    Name:  "Curso de Go",
    Price: 150,
    //Observations: "Wow",
}
err := serviceProduct.Update(m)
if err != nil {
    log.Fatalf("product.Update: %v", err)
}
```

Los campos que no se asignen se actualizaran con su zero value. El campo `updated_at` se actualiza automaticamente con el timestamp actual.

## Eliminar un registro

Para eliminar un producto por su ID, se llama al metodo `Delete` del servicio pasando el ID como argumento:

```go
storageProduct := storage.NewPsqlProduct(storage.Pool())
serviceProduct := product.NewService(storageProduct)

err := serviceProduct.Delete(3)
if err != nil {
    log.Fatalf("product.Delete: %v", err)
}
```

El metodo ejecuta un `DELETE FROM products WHERE id = $1`. Una vez eliminado el registro, el codigo fue removido de `main.go`.
