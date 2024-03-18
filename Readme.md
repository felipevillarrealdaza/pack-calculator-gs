## Running the Pack Calculator service
You must run `docker-compose up`, which will spin up the following services:
- API
- DB

That's all. Once the services are up and running, you may use the API targeting 0.0.0.0 as host and 8000 as port.

## Adding pack sizes

To add a pack size, you must make a request similar to this:

```bash
curl --location '0.0.0.0:8000/api/v1/pack' \
--header 'Content-Type: application/json' \
--data '{
    "size": 21
}'
```

You can modify the *size* field in the request body as needed.

## Removing pack sizes

To remove a pack size, you must make a request similar to this:

```bash
curl --location --request DELETE '0.0.0.0:8000/api/v1/pack' \
--header 'Content-Type: application/json' \
--data '{
    "size": 10
}'
```

You can modify the *size* field in the request body as needed.

## Creating and order to calculate packs

To create an order and calculate the packs needed for a specific amount of items, you must make a request similar to this:

```bash
curl --location '0.0.0.0:8000/api/v1/order' \
--header 'Content-Type: application/json' \
--data '{
    "quantity": 53
}'
```

You can modify the *quantity* field in the request body as needed.

## Pack algorithm used

The high level algorithm used to calculate the packs is the following:

1. Loop through all the quantities specified. If the order has 500 items, then the algorithm will loop from 1 to 500.
2. In each of the quantities, loop through possible pack sizes. This will allow to know all possible packing options.
3. Use previous iterations results to calculate the current quantity. This is a dynamic programming approach
4. Compare possible packing options and identify which is the best one, according to these rules:
- The least number of items must be packed to fulfill the order
- The least number of packs must be used to fulfill the order
5. Save the best option for that order quantity. This result will be used for next iterations.

