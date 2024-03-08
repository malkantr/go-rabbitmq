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
