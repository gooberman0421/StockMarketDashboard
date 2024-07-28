from flask import Flask, jsonify, request
import os
import datetime
import threading
import time
import random

app = Flask(__name__)

STOCK_DATA = {"AAPL": 150.0, "GOOGL": 2800.0, "MSFT": 300.0}
STOCK_HISTORY = {"AAPL": [], "GOOGL": [], "MSFT": []}
TRANSACTION_HISTORY = []

def update_stock_prices():
    while True:
        for symbol in STOCK_DATA:
            change_percent = random.uniform(-0.02, 0.02)
            STOCK_DATA[symbol] *= (1 + change_percent)
            STOCK_HISTORY[symbol].append({"timestamp": datetime.datetime.now().isoformat(), "price": STOCK_DATA[symbol]})
        time.sleep(10)

@app.route('/stock/<symbol>', methods=['GET'])
def get_stock(symbol):
    stock_price = STOCK_DATA.get(symbol.upper())
    if stock_price:
        return jsonify({"symbol": symbol.upper(), "price": stock_price}), 200
    else:
        return jsonify({"error": "Stock not found."}), 404

@app.route('/stock/history/<symbol>', methods=['GET'])
def get_stock_history(symbol):
    symbol = symbol.upper()
    if symbol in STOCK_HISTORY:
        return jsonify(STOCK_HISTORY[symbol]), 200
    else:
        return jsonify({"error": "Stock history not found."}), 404

@app.route('/transaction', methods=['POST'])
def submit_transaction():
    data = request.json
    if not data or 'symbol' not in data or 'quantity' not in data or 'type' not in data:
        return jsonify({"error": "Missing required fields."}), 400

    symbol = data['symbol'].upper()
    if symbol not in STOCK_DATA:
        return jsonify({"error": "Invalid stock symbol."}), 400
    
    if data['type'] not in ['buy', 'sell']:
        return jsonify({"error": "Invalid transaction type."}), 400
    
    try:
        quantity = int(data['quantity'])
        if quantity <= 0:
            raise ValueError
    except ValueError:
        return jsonify({"error": "Quantity must be a positive integer."}), 400

    transaction = {
        "id": len(TRANSACTION_HISTORY) + 1,
        "symbol": symbol,
        "quantity": quantity,
        "type": data['type'],
        "price": STOCK_DATA[symbol],
        "timestamp": datetime.datetime.now().isoformat()
    }
    TRANSACTION_HISTORY.append(transaction)
    return jsonify(transaction), 201

@app.route('/transactions', methods=['GET'])
def get_transactions():
    return jsonify(TRANSACTION_HISTORY), 200

if __name__ == '__main__':
    price_update_thread = threading.Thread(target=update_stock_prices, daemon=True)
    price_update_thread.start()
    app.run(debug=True)