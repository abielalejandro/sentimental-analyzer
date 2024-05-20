import json
import pika
from config import Config
from event_intf import EventBus
from transformers import pipeline
from pysentimiento import create_analyzer

class RabbitEventBus(EventBus):
    def __init__(self, config:Config):
        self.config=config

        if config.analyzer=="pysentimiento":
            self.analyzer = create_analyzer(task="sentiment", lang="es")

        if config.analyzer=="pysentimiento":            
            #self.sentiment_pipeline = pipeline("sentiment-analysis")
            self.sentiment_pipeline = pipeline(model="finiteautomata/bertweet-base-sentiment-analysis")

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

        call=self.callbackPy
        if self.config.analyzer=="pysentimiento":
            call = self.callbackPy
        if self.config.analyzer=="transformers":   
            call = self.callback

        self.channel.basic_consume(
                queue=queue_name, 
                on_message_callback=call, 
                auto_ack=self.config.rabbitmq.auto_ack)
        self.channel.start_consuming()

    def callbackPy(self,ch, method, properties, body):
        event = json.loads(body.decode('utf-8'))
        result = self.analyzer.predict(event["data"]["Msg"])
        print(result)

        label="POSITIVE"
        score=result.probas["POS"]
        if result.probas["POS"]<result.probas["NEG"]:
            label="NEGATIVE"
            score=result.probas["NEG"]

        if result.probas["NEU"]>result.probas["POS"] and result.probas["NEU"]>result.probas["NEG"]:
            label="POSITIVE"
            score=0.5

        data= {"Id":event["id"],"Label": label, "Score":score}
        event["data"] = data
        event["source"] = "sentimental/analyzer"
        event["type"] = self.config.rabbitmq.producer_master_routing_key
        print(event)
        self.channel.basic_publish(
            exchange=self.config.rabbitmq.exchange,
            routing_key=self.config.rabbitmq.producer_master_routing_key, 
            body=json.dumps(event))
        
    def callback(self,ch, method, properties, body):
        event = json.loads(body.decode('utf-8'))
        result = self.sentiment_pipeline(event["data"]["Msg"])
        print(result)

        if result[0]["label"]=="NEU":
            result[0]["label"] = "POSITIVE"

        if result[0]["label"]=="NEG":
            result[0]["label"] = "NEGATIVE"

        if result[0]["label"]=="POS":
            result[0]["label"] = "POSITIVE"

        data= {"Id":event["id"],"Label": result[0]["label"], "Score":result[0]["score"]}
        event["data"] = data
        event["source"] = "sentimental/analyzer"
        event["type"] = self.config.rabbitmq.producer_master_routing_key
        print(event)
        self.channel.basic_publish(
            exchange=self.config.rabbitmq.exchange,
            routing_key=self.config.rabbitmq.producer_master_routing_key, 
            body=json.dumps(event))