# TP2
Script de ejecución: _**`ejecutar.sh`**_


#### Para ejecutar la app, utilizar el siguiente comando:
> ```bash
> chmod +x ejecutar.sh 
>
> ./ejecutar.sh
> ```


## Dependencias:
- [Docker](https://www.docker.com/)


```mermaid

erDiagram

    users {
        integer id PK
        varchar_100 name
        timestamp_with_time_zone created_at
    }

    shopping_items {
        integer id PK
        varchar_255 title
        integer quantity
        numeric_10_2 price
        boolean is_purchased
        integer paid_by_user_id FK
        timestamp_with_time_zone created_at
    }

    item_splits {
        integer item_id FK
        integer user_id FK
    }

    item_splits }o--|| shopping_items : "item_id -> id"
    item_splits }o--|| users : "user_id -> id"
    shopping_items }o--|| users : "paid_by_user_id -> id"
```
