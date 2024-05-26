from config import Config
from analyzer_transformers import TransformerAnalyzer
from analyzer_pysentimiento import PySentimientoAnalyzer
from analyzer import Analyzer, GenericAnalyzer

def NewAnalyzer(config:Config)->Analyzer:
    match config.analyzer:
        case "transformers":
            return TransformerAnalyzer()  
        case "pysentimiento":
            return PySentimientoAnalyzer() 
        case _:
            return GenericAnalyzer() 

