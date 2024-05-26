from analyzer import Analyzer, AnalyzerResult
from pysentimiento import create_analyzer

class PySentimientoAnalyzer(Analyzer):
    def __init__(self):
        self.analyzer = create_analyzer(task="sentiment", lang="es")

    def Analyze(self,text:str)->AnalyzerResult:
        result = self.analyzer.predict(text)
        label=result.output
        score=result.probas[label]

        if label=="NEU":
            label= "NEUTRAL"

        if label=="NEG":
            label = "NEGATIVE"

        if label=="POS":
            label = "POSITIVE"

        response = AnalyzerResult(label, score)

        return response
