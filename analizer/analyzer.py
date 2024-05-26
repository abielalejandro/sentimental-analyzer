class AnalyzerResult:
    def __init__(self, label:str, score:float):
        self.label=label
        self.score = score

class Analyzer:
    def Analyze(self,text:str)->AnalyzerResult:
        pass

class GenericAnalyzer(Analyzer):
    def Analyze(self, text:str)->AnalyzerResult:
        return AnalyzerResult("POSITIVE",0.5)
