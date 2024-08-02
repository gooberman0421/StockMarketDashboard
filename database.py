import sqlite3
from os import getenv
from dotenv import load_dotenv

load_dotenv()

def get_db_connection():
    conn = sqlite3.connect(getenv('DATABASE_PATH', 'database.db'))
    conn.row_factory = sqlite3.Row
    return conn

def create_tables():
    conn = get_db_connection()
    cursor = conn.cursor()
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
            type TEXT NOT NULL,
            quantity INTEGER NOT NULL,
            price REAL NOT NULL,
            timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (stock_id) REFERENCES stocks (id)
        )
    ''')
    conn.commit()
    conn.close()

def store_stock_data(symbol, name, price):
    conn = get_db_connection()
    cursor = conn.cursor()
    cursor.execute('''
        INSERT INTO stocks (symbol, name, price)
        VALUES (?, ?, ?)
        ON CONFLICT(symbol) DO UPDATE SET
        name=excluded.name,
        price=excluded.price
    ''', (symbol, name, price))
    conn.commit()
    conn.close()

def retrieve_stock_data(symbol):
    conn = get_db_connection()
    cursor = conn.cursor()
    cursor.execute('SELECT id, symbol, name, price FROM stocks WHERE symbol = ?', (symbol,))
    stock = cursor.fetchone()
    conn.close()
    return stock

def store_transaction(stock_symbol, transaction_type, quantity, price):
    stock = retrieve_stock_data(stock_symbol)
    if stock:
        conn = get_db_connection()
        cursor = conn.cursor()
        cursor.execute('''
            INSERT INTO transactions (stock_id, type, quantity, price)
            VALUES (?, ?, ?, ?)
        ''', (stock['id'], transaction_type, quantity, price))
        conn.commit()
        conn.close()
    else:
        raise ValueError("Stock not found")

def retrieve_transaction_history(stock_symbol):
    stock = retrieve_stock_data(stock_symbol)
    if stock:
        conn = get_db_connection()
        cursor = conn.cursor()
        cursor.execute('''
            SELECT type, quantity, price, timestamp
            FROM transactions
            WHERE stock_id = ?
            ORDER BY timestamp DESC
        ''', (stock['id'],))
        transactions = cursor.fetchall()
        conn.close()
        return transactions
    else:
        raise ValueError("Stock not found")