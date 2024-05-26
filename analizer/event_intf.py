from analyzer import Analyzer

class EventBus:
    def Listen(self):
        pass

class GenericEventBus(EventBus):
    def __init__(self, analyzer:Analyzer):
        self.analyzer = analyzer
    def Listen(self):
        pass

