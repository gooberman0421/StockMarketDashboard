import fetch from 'node-fetch';

const API_BASE_URL = process.env.API_BASE_URL || 'http://localhost:3000';

export async function fetchStockData(tickerSymbol: string): Promise<any> {
    try {
        const response = await fetch(`${API_BASE_URL}/stocks/${tickerSymbol}`);
        
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        
        return response.json();
    } catch (error) {
        console.error('Error fetching stock data:', error);
        throw error;
    }
}

export async function submitTransaction(
    transactionData: { tickerSymbol: string; quantity: number; transactionType: 'buy' | 'sell' }
): Promise<any> {
    try {
        const response = await fetch(`${API_BASE_URL}/transactions`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(transactionData),
        });
        
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        
        return response.json();
    } catch (error) {
        console.error('Error submitting transaction:', error);
        throw error;
    }
}

export async function getTransactionHistory(): Promise<any> {
    try {
        const response = await fetch(`${API_BASE_URL}/transactions`);
        
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        
        return response.json();
    } catch (error) {
        console.error('Error fetching transaction history:', error);
        throw error;
    }
}