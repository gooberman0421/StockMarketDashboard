from flask import Flask, jsonify, request
import os
import datetime

app = Flask(__name__)

STOCK_DATA = {"AAPL": 150.0, "GOOGL": 2800.0, "MSFT": 300.0}
TRANSACTION_HISTORY = []

@app.route('/stock/<symbol>', methods=['GET'])
def get_stock(symbol):
    stock_price = STOCK_DATA.get(symbol.upper())
    if stock_price:
        return jsonify({"symbol": symbol.upper(), "price": stock_price}), 200
    else:
        return jsonify({"error": "Stock not found."}), 404

@app.route('/transaction', methods=['POST'])
def submit_transaction():
    data = request.json
    if not data or 'symbol' not in data or 'quantity' not in data or 'type' not in data:
        return jsonify({"error": "Missing required fields."}), 400

    symbol = data['symbol'].upper()
    if symbol not in STOCK_DATA:
        return jsonify({"error": "Invalid stock symbol."}), 400
    
    transaction = {
        "id": len(TRANSACTION_HISTORY) + 1,
        "symbol": symbol,
        "quantity": data['quantity'],
        "type": data['type'],
        "price": STOCK_DATA[symbol],
        "timestamp": datetime.datetime.now().isoformat()
    }
    TRANSACTION_HISTORY.append(transaction)
    return jsonify(transaction), 201

@app.route('/transactions', methods=['GET'])
def get transactions():
    return jsonify(TRANSACTION_HISTORY), 200

if __name__ == '__main__':
    app.run(debug=True)