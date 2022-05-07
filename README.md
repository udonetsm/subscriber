# subscriber

database contains next fields:


        users=# \d orders
                            Table "public.orders"
        Column   |       Type        | Collation | Nullable | Default 
        -----------+-------------------+-----------+----------+---------
        id        | character varying |           | not null | 
        orderjson | json              |           |          | 
        pubdate   | bigint            |           |          | 
        Indexes:
            "orders_pkey" PRIMARY KEY, btree (id)

For work with messages got by HTTP, you may unmarshall it. Mannualy it doesn't unmarshal and returns by HTTP from memory-cache of service as json.

