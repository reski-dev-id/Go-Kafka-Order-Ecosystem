INSERT INTO orders (
    id,
    customer_name,
    product_name,
    quantity,
    amount,
    status
)
VALUES
(
    gen_random_uuid(),
    'Reski',
    'Macbook Pro',
    1,
    2500,
    'PENDING'
);