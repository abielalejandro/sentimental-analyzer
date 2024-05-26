import json
import pika
from config import Config
from event_intf import EventBus
from analyzer import Analyzer

class RabbitEventBus(EventBus):
    def __init__(self, config:Config, analyzer:Analyzer):
        self.config=config
        self.analyzer = analyzer

    def Listen(self):
        credentials = pika.PlainCredentials(self.config.rabbitmq.user, self.config.rabbitmq.pwd)
        self.connection = pika.BlockingConnection(
            pika.ConnectionParameters(self.config.rabbitmq.host,
                self.config.rabbitmq.port,"/",credentials))
        self.channel = self.connection.channel()
        self.ListenMaster()

    def ListenMaster(self): 
        self.channel.exchange_declare(
                durable=True,
                exchange=self.config.rabbitmq.exchange, 
                exchange_type=self.config.rabbitmq.exchange_type)
        result = self.channel.queue_declare(self.config.rabbitmq.queue, exclusive=False)
        queue_name = result.method.queue
        self.channel.queue_bind(
            exchange=self.config.rabbitmq.exchange, 
            queue=queue_name, 
            routing_key=self.config.rabbitmq.consumer_master_routing_key)

        self.channel.basic_consume(
                queue=queue_name, 
                on_message_callback=self.callback, 
                auto_ack=self.config.rabbitmq.auto_ack)
        self.channel.start_consuming()

    def callback(self,ch, method, properties, body):
        event = json.loads(body.decode('utf-8'))
        result = self.analyzer.Analyze(event["data"]["Msg"])

        data= {
                "Id":event["id"],
                "Label": result.label,
                "Score":result.score
                }

        event["data"] = data
        event["source"] = "sentimental/analyzer"
        event["type"] = self.config.rabbitmq.producer_master_routing_key

        print(json.dumps(event))
        self.channel.basic_publish(
            exchange=self.config.rabbitmq.exchange,
            routing_key=self.config.rabbitmq.producer_master_routing_key, 
            body=json.dumps(event))
