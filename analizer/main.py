from config import Config
from event_ini import NewEventBus
from analyzer_ini import NewAnalyzer

configs = Config()
analyzer = NewAnalyzer(configs)
event = NewEventBus(configs, analyzer)
event.Listen()
