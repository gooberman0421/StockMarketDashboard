import sqlite3
from os import getenv
from dotenv import load_dotenv

load_dotenv()

def establish_db_connection():
    connection = sqlite3.connect(getenv('DATABASE_PATH', 'database.db'))
    connection.row_factory = sqlite3.Row
    return connection

def initialize_database():
    connection = establish_db_connection()
    cursor = connection.cursor()
    cursor.execute('''
        CREATE TABLE IF NOT EXISTS stocks (
            id INTEGER PRIMARY KEY,
            symbol TEXT NOT NULL,
            name TEXT NOT NULL,
            price REAL NOT NULL
        )
    ''')
    cursor.execute('''
        CREATE TABLE IF NOT EXISTS transactions (
            id INTEGER PRIMARY KEY,
            stock_id INTEGER,
            transaction_type TEXT NOT NULL,
            quantity INTEGER NOT NULL,
            price REAL NOT NULL,
            timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (stock_id) REFERENCES stocks (id)
        )
    ''')
    connection.commit()
    connection.close()

def insert_or_update_stock(symbol, name, current_price):
    connection = establish_db_connection()
    cursor = connection.cursor()
    cursor.execute('''
        INSERT INTO stocks (symbol, name, price)
        VALUES (?, ?, ?)
        ON CONFLICT(symbol) DO UPDATE SET
        name=excluded.name,
        price=excluded.price
    ''', (symbol, name, current_price))
    connection.commit()
    connection.close()

def fetch_stock_by_symbol(symbol):
    connection = establish_db_connection()
    cursor = connection.cursor()
    cursor.execute('SELECT id, symbol, name, price FROM stocks WHERE symbol = ?', (symbol,))
    stock_details = cursor.fetchone()
    connection.close()
    return stock_details

def log_stock_transaction(stock_symbol, transaction_type, amount, transaction_price):
    stock_detail = fetch_stock_by_symbol(stock_symbol)
    if stock_detail:
        connection = establish_db_connection()
        cursor = connection.cursor()
        cursor.execute('''
            INSERT INTO transactions (stock_id, transaction_type, quantity, price)
            VALUES (?, ?, ?, ?)
        ''', (stock_detail['id'], transaction_type, amount, transaction_price))
        connection.commit()
        connection.close()
    else:
        raise ValueError("Stock not found")

def fetch_transaction_history_by_symbol(stock_symbol):
    stock_detail = fetch_stock_by_symbol(stock_symbol)
    if stock_detail:
        connection = establish_db_connection()
        cursor = connection.cursor()
        cursor.execute('''
            SELECT transaction_type, quantity, price, timestamp
            FROM transactions
            WHERE stock_id = ?
            ORDER BY timestamp DESC
        ''', (stock_detail['id'],))
        transaction_history = cursor.fetchall()
        connection.close()
        return transaction_history
    else:
        raise ValueError("Stock not found")