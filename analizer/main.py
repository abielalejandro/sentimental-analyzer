from config import Config
from event_ini import NewEventBus

configs = Config()
event = NewEventBus(configs)
event.Listen()
