from analyzer import Analyzer, AnalyzerResult
from transformers import pipeline

class TransformerAnalyzer(Analyzer):
    def __init__(self):
        #self.analyzer = pipeline("sentiment-analysis")
        self.analyzer = pipeline(model="finiteautomata/bertweet-base-sentiment-analysis")

    def Analyze(self,text:str)->AnalyzerResult:
        result = self.analyzer(text)
        label=result[0]["label"]
        score=result[0]["score"]

        if result[0]["label"]=="NEU":
            label= "NEUTRAL"

        if result[0]["label"]=="NEG":
            label = "NEGATIVE"

        if result[0]["label"]=="POS":
            label = "POSITIVE"
        
        response = AnalyzerResult(label, score)

        return response
