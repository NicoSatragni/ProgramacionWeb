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
        integer item_id PK, FK
        integer user_id PK, FK
    }

    item_splits }o--|| shopping_items : "item_id -> id"
    item_splits }o--|| users : "user_id -> id"
    shopping_items }o--|| users : "paid_by_user_id -> id"
```

```mermaid

classDiagram
direction BT
class item_splits {
   integer item_id
   integer user_id
}
class settlements {
   integer payer_id
   integer receiver_id
   "numeric(10,2)" amount
   "timestamp with time zone" created_at
   integer id
}
class shopping_items {
   "varchar(255)" title
   integer quantity
   "numeric(10,2)" price
   boolean is_purchased
   integer paid_by_user_id
   "timestamp with time zone" created_at
   integer id
}
class users {
   "varchar(100)" name
   "timestamp with time zone" created_at
   integer id
}

item_splits  -->  shopping_items : item_id:id
item_splits  -->  users : user_id:id
settlements  -->  users : payer_id:id
settlements  -->  users : receiver_id:id
shopping_items  -->  users : paid_by_user_id:id


```
