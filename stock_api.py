import requests
from dotenv import load_dotenv
import os

load_dotenv()

ALPHA_VANTAGE_API_KEY = os.getenv("ALPHA_VANTAGE_API_KEY")
ALPHA_VANTAGE_BASE_URL = "https://www.alphavantage.co/query"

def get_real_time_stock_data(symbol):
    parameters = {
        "function": "GLOBAL_QUOTE",
        "symbol": symbol,
        "apikey": ALPHA_VANTAGE_API_KEY
    }
    
    response = requests.get(ALPHA_VANTAGE_BASE_URL, params=parameters)
    data = response.json()
    
    if "Global Quote" in data:
        return {
            "symbol": symbol,
            "open": data["Global Quote"]["02. open"],
            "high": data["Global Quote"]["03. high"],
            "low": data["Global Quote"]["04. low"],
            "price": data["Global Quote"]["05. price"],
            "volume": data["Global Quote"]["06. volume"],
            "latest_trading_day": data["Global Quote"]["07. latest trading day"],
            "previous_close": data["Global Quote"]["08. previous close"],
            "change": data["Global Quote"]["09. change"],
            "change_percent": data["Global Quote"]["10. change percent"]
        }
    else:
        return {"error": "No data found for the symbol"}

def get_historical_stock_data(symbol, outputsize="compact"):
    parameters = {
        "function": "TIME_SERIES_DAILY",
        "symbol": symbol,
        "outputsize": outputsize,
        "apikey": ALPHA_VANTAGE_API_KEY
    }
    
    response = requests.get(ALPHA_VANTAGE_BASE_URL, params=parameters)
    data = response.json()
    
    if "Time Series (Daily)" in data:
        daily_data = data["Time Series (Daily)"]
        historical_data = []
        for date, daily_info in daily_data.items():
            parsed_data = {
                "date": date,
                "open": daily_info["1. open"],
                "high": daily_info["2. high"],
                "low": daily_info["3. low"],
                "close": daily_info["4. close"],
                "volume": daily_info["5. volume"]
            }
            historical_data.append(parsed_data)
        
        return historical_data
    else:
        return {"error": "No historical data found for the symbol"}