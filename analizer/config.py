from decouple import config

class EventBus:
    def __init__(self):
        self.type=config("EVENT_BUS_TYPE",default="generic")
       

class RabbitConfig:
    def  __init__(self):
        self.host=config("RABBITMQ_HOST",default="localhost")
        self.port=config("RABBITMQ_PORT",default=5672, cast=int)
        self.user=config("RABBITMQ_USER",default="guest")
        self.pwd=config("RABBITMQ_PWD",default="guest")

        self.exchange=config("RABBITMQ_EXCHANGE",default="sentimental")
        self.exchange_type=config("RABBITMQ_EXCHANGE_TYPE",default="topic")
        self.queue=config("RABBITMQ_QUEUE",default="")
        self.producer_master_routing_key=config("RABBITMQ_PRODUCER_MASTER_ROUTING",default="analyzer.text.analyzed")
        self.consumer_master_routing_key=config("RABBIT_CONSUMER_MASTER_ROUTING",default="master.text.created")
        self.auto_ack=config("RABBITMQ_AUTO_ACK",default=True, cast=bool)

class Config:
  def __init__(self):
      self.analyzer= config("PY_ANALYZER",default="transformers")
      self.event_bus= EventBus()
      self.rabbitmq = RabbitConfig()      
