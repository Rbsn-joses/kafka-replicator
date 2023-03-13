Repositório feito em golang para replicar mensagens de tópicos de um kafka para o outro

Configurável por meio de variáveis de ambiente. Caso preciso incluir configuração de certificado ssl, apenas incluir as respectivas envs:

Variáveis para o kafka-replicator mandar as mensagens para o kafka replica
KAFKA_PRODUCER_TLS_ENABLED=true
KAFKA_PRODUCER_SSL_CERTIFICATE_LOCATION=/etc/ssl/ca.crt
KAFKA_PRODUCER_SSL_PRIVATE_KEY_LOCATION=/etc/ssl/private.key
KAFKA_PRODUCER_SSL_PRIVATE_KEY_PASSWORD= topSecretRandom

Caso não precise de conexão ssl, apenas incluir as envs para o producer, ou seja o kafka que receberá as mensagens replicadas
KAFKA_PRODUCER_BOOSTRAP_SERVERS
KAFKA_PRODUCER_TOPIC

Caso essa env não tenha valor true, as mensagens não serão mandadas pro kafka replica, apenas será mostrado nos logs do container.
KAFKA_REPLICATION_ENABLED=true

Variáveis para o que o kafka-replicator leia o tópico e puxe as mensagens do kafka original
KAFKA_CONSUMER_TLS_ENABLED

KAFKA_CONSUMER_SSL_CERTIFICATE_LOCATION

KAFKA_CONSUMER_SSL_PRIVATE_KEY_LOCATION

KAFKA_CONSUMER_SSL_PRIVATE_KEY_PASSWORD

Caso não precise de conexão ssl, apenas incluir as envs para o kafka-replicator
KAFKA_CONSUMER_BOOSTRAP_SERVERS
KAFKA_CONSUMER_TOPIC
KAFKA_CONSUMER_GROUP_ID
