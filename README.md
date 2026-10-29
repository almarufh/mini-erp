# [BACKEND-FRONTEND] -  Mini ERP Application

<p>
  <a href="https://golang.org/"><img style="border-radius:4px" src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go"></a>
  <a href="https://github.com/gin-gonic/gin"><img style="border-radius:4px" src="https://img.shields.io/badge/Gin-008B8B?style=for-the-badge&logo=go&logoColor=white" alt="Gin"></a>
  <a href="https://github.com/joho/godotenv"><img style="border-radius:4px" src="https://img.shields.io/badge/Godotenv-ECD53F?style=for-the-badge&logo=dotenv&logoColor=black" alt="Godotenv"></a>
  <a href="https://www.postgresql.org/"><img style="border-radius:4px" src="https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL"></a>
  <a href="https://github.com/jackc/pgx"><img style="border-radius:4px" src="https://img.shields.io/badge/pgx-339933?style=for-the-badge&logo=postgresql&logoColor=white" alt="pgx"></a>
  <a href="https://redis.io/"><img style="border-radius:4px" src="https://img.shields.io/badge/Redis-DC382D?style=for-the-badge&logo=redis&logoColor=white" alt="Redis"></a>
</p>

A robust fullstack Mini ERP & Inventory Management System built with Go (Gin Framework) on the backend and React + Redux Toolkit + Tailwind CSS v4 on the frontend. This project is engineered to handle atomic stock adjustments, sales order processing, and stock movement auditing safely with database transactions.


## Tech Stack

- [![Go](https://img.shields.io/badge/Go-v1.27.1-blue?logo=go&logoColor=white)](https://go.dev/learn/)
- [![Gin-Gonic](https://img.shields.io/badge/Gin_Gonic-v1.12.0-green?logo=Gin&logoColor=white)](https://gin-gonic.com/en/)
- [![PostgreSQL](https://img.shields.io/badge/PostgreSQL-v18.6-blue?logo=PostgreSQL&logoColor=white)](https://www.postgresql.org/)
- [![pgx](https://img.shields.io/badge/Pgx-v5.11.0-blue?)](https://github.com/jackc/pgx)
- [![GodotEnv](https://img.shields.io/badge/GodotEnv-v1.5.1-blue?)](https://github.com/joho/godotenv)
- [![React](https://img.shields.io/badge/React-v19.2.8-blue?logo=react&logoColor=white)](https://react.dev/)
- [![Vite](https://img.shields.io/badge/Vite-v8.2.0-purple?logo=vite&logoColor=white)](https://vite.dev/)
- [![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-v4.3.3-blue?logo=tailwindcss&logoColor=white)](https://tailwindcss.com/) <
- [![React Router](https://img.shields.io/badge/React_Router-v8.3.0-red?logo=react-router&logoColor=white)](https://reactrouter.com/)

## Features

- **Product Catalog Management**: CRUD data produk, pencarian SKU/nama produk, dan paginasi interaktif.
- **Atomic Stock Adjustment**: Modul _Stock In_ & _Stock Out_ beserta pencatatan alasan/catatan perubahan.
- **Stock Movement Audit Trail**: Riwayat log pergerakan stok lengkap untuk setiap produk.
- **Sales Order Processing**: Pembuatan pesanan multi-item dengan kalkulasi subtotal otomatis secara real-time.
- **Concurrency & Stock Guard**: Validasi stok presisi yang mengembalikan status `409 INSUFFICIENT_STOCK` jika kuantitas order melebihi stok yang ada.
- **Order Cancellation**: Pembatalan sales order yang secara otomatis mengembalikan (_restore_) kuantitas stok ke gudang.

---

## Usage Instruction

### Clone this repository

```bash
$ git clone https://github.com/almarufh/mini-erp
```

### Run Entire Application via Make

```bash
$ make up
```

## Run Application

For run application, visit [**`localhost:8080`** ](http://localhost:8080/)

## Routes

| Endpoint                      | Method   | Description                                    |
| :---------------------------- | :------- | :--------------------------------------------- |
| `/api/products`               | `GET`    | Get paginated products list with search filter |
| `/api/products`               | `POST`   | Create a new product                           |
| `/api/products/:id`           | `PUT`    | Update product details                         |
| `/api/products/:id`           | `DELETE` | Delete a product (fails if movements exist)    |
| `/api/products/:id/stock`     | `POST`   | Adjust product stock (Stock In / Stock Out)    |
| `/api/products/:id/movements` | `GET`    | Get stock movement history log                 |
| `/api/orders`                 | `GET`    | Get all sales orders history                   |
| `/api/orders`                 | `POST`   | Create a new sales order                       |
| `/api/orders/:id`             | `GET`    | Get sales order details by ID                  |
| `/api/orders/:id/cancel`      | `POST`   | Cancel sales order and restore stock           |

### Documentation

For complete documentation, visit [docs swagger endpooint **`localhost:5000/api/docs/index.html`** ](http://localhost:5000/api/docs/index.html#/)

## How to Contribute

- Fork this repository
- Create your changes
- Pull Request

## License [![License: MIT](https://img.shields.io/badge/License-MIT-red)](https://opensource.org/license/mit)

This project is licensed under the MIT License

## CONTACTS

[![Email](https://img.shields.io/badge/Email-almarufhidayat99@gmail.com-blue?logo=maildotru&logoColor=red)](mailto:almarufhidayat99@gmail.com)<br>
[![WhatsApp](https://img.shields.io/badge/WhatsApp-6281973779380-blue?logo=whatsapp&logoColor=green)](https://wa.me/6281973779380)
