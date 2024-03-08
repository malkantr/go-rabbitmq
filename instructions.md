# RabbitMQ

```
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.11-management

http://localhost:15672
default guest guest

```

```
docker exec rabbitmq rabbitmqctl

docker exec rabbitmq rabbitmqctl add_user dave 1234                                 // create user
docker exec rabbitmq rabbitmqctl set_user_tags dave administrator                   // set permission

docker exec rabbitmq rabbitmqctl delete_user guest

docker exec rabbitmq rabbitmqctl add_vhost customers                                // add vhost
docker exec rabbitmq rabbitmqctl set_permissions -p customers dave ".*" ".*" ".*"   // configuration, write, read
```

```
docker exec rabbitmq rabbitmqadmin declare exchange --vhost=customers name=customer_events type=topic -u dave -p 1234 durable=true 
docker exec rabbitmq rabbitmqctl set_topic_permissions -p customers dave customer_events "^customers.*" "^customers.*"
```

```
// reset exchange for pub-sub/fanout
docker exec rabbitmq rabbitmqadmin delete exchange name=customer_events --vhost=customers -u dave -p 1234
docker exec rabbitmq rabbitmqadmin declare exchange --vhost=customers name=customer_events type=fanout -u dave -p 1234 durable=true
docker exec rabbitmq rabbitmqctl set_topic_permissions -p customers dave customer_events ".*" ".*"

```
```
docker exec rabbitmq rabbitmqadmin declare exchange --vhost=customers name=customer_callbacks type=direct -u dave -p 1234 durable=true
docker exec rabbitmq rabbitmqctl set_topic_permissions -p customers dave customer_callbacks ".*" ".*"
```

// tls connection - delete and re-create rabbitmq container
docker container rm -f rabbitmq
docker run -d --name rabbitmq -v "$(pwd)"/rabbitmq.conf:/etc/rabbitmq/rabbitmq.conf:ro -v "$(pwd)"/tls-gen/basic/result:/certs -p 5671:5671 -p 15671:15671 rabbitmq:3.11-management


// # final conf for rabbitmq container with definitions
// reset docker bind configuration 
docker run -d --name rabbitmq 
    -v "$(pwd)"/rabbitmq_definitions.json:/etc/rabbitmq/rabbitmq_definitions.json:ro 
    -v "$(pwd)"/rabbitmq.conf:/etc/rabbitmq/rabbitmq.conf:ro 
    -v "$(pwd)"/tls-gen/basic/result:/certs 
    -p 5671:5671 -p 15671:15671 rabbitmq:3.11-management




## Rabbit repository for tls
```
git clone https://github.com/rabbitmq/tls-gen tls-gen

cd tls-gen/basic
make PASSWORD=
make verify
sudo chmod 644 result/*

// "reset docker container"

//rabbitmq.conf
# Disable NON TCP
listeners.tcp = none
listeners.ssl.default = 5671
# ssl_certs
ssl_options.cacertfile = /certs/ca_certificate.pem
ssl_options.certfile = /certs/server_PCNAME_certificate.pem
ssl_options.keyfile = /certs/server_PCNAME_key.pem

# Peer vertification
ssl_options.verify = verify_peer
ssl_options.fail_if_no_peed_cert = true


docker restart rabbitmq
docker logs rabbitmq
```

//encodepassword.sh // https://stackoverflow.com/questions/41306350/how-to-generate-password-hash-for-rabbitmq-management-http-api
bash encodepassword.sh

// reset docker bind configuration 
docker run -d --name rabbitmq 
    -v "$(pwd)"/rabbitmq_definitions.json:/etc/rabbitmq/rabbitmq_definitions.json:ro 
    -v "$(pwd)"/rabbitmq.conf:/etc/rabbitmq/rabbitmq.conf:ro 
    -v "$(pwd)"/tls-gen/basic/result:/certs 
    -p 5671:5671 -p 15671:15671 rabbitmq:3.11-management

## Queue Parameters
- Durable       (persists)
- AutoDelete

## Exchange Types
- Direct
- Fanout
- Topic
- Header

## Routing Keys
binds exchanges with queues

## Delivery Modes
- Persistent 
- Non-Persistent (transient)


--- 

# Structure

```
git init
touch .gitignore

cd main/
mkdir -p cmd/producer
mkdir -p cmd/consumer
mkdir internal
touch cmd/producer/main.go
touch cmd/consumer/main.go
touch internal/rabbitmq.go

go mod init
go get github.com/rabbitmq/amqp091-go

go run cmd/producer/main.go
go run cmd/consumer/main.go
```
