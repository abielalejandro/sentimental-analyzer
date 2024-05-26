from config import Config
from event_intf import EventBus
from event_intf import GenericEventBus
from event_rabbitmq import RabbitEventBus
from analyzer import Analyzer

def NewEventBus(config:Config, analyzer:Analyzer) ->EventBus:
    match config.event_bus.type:
        case "rabbitmq":
            return RabbitEventBus(config, analyzer) 
        case _:
            return GenericEventBus(analyzer) 
