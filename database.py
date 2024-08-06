import time
import requests  # Assuming you use requests to make HTTP calls

STOCK_PRICE_TTL = 3600  # Time to live for cached prices in seconds (e.g., 1 hour)

def fetch_stock_price(symbol):
    now = time.time()
    connection = establish_db_connection()
    cursor = connection.cursor()
    
    # Check if the stock price is in the database and not outdated
    cursor.execute('''
        SELECT price, (strftime('%s', 'now') - strftime('%s', timestamp)) AS age
        FROM stocks
        WHERE symbol = ?
    ''', (symbol,))
    result = cursor.fetchone()
    
    if result and result['age'] < STOCK_PRICE_TTL:
        return result['price']
    else:
        # Assuming you have a function get_external_stock_price() that fetches
        # the price from an external API:
        new_price = get_external_stock_price(symbol)
        
        # Update the stock price in the database
        cursor.execute('''
            UPDATE stocks
            SET price = ?, timestamp = CURRENT_TIMESTAMP
            WHERE symbol = ?
        ''', (new_price, symbol))
        connection.commit()
        connection.close()
        
        return new_price

def get_external_stock_price(symbol):
    # Placeholder for API call to fetch new price
    # response = requests.get(f"https://api.example.com/stocks/price?symbol={symbol}")
    # new_price = response.json().get('price')
    new_price = 100.00  # Example static price
    return new_price