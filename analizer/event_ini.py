from config import Config
from event_intf import EventBus
from event_intf import GenericEventBus
from event_rabbitmq import RabbitEventBus

def NewEventBus(config:Config) ->EventBus:
    match config.event_bus.type:
        case "rabbitmq":
            return RabbitEventBus(config) 
        case _:
            return GenericEventBus() 