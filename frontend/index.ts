import { useQuery } from 'react-query';

const fetchStockById = async (stockId: string) => {
  const response = await fetch(`/api/stocks/${stockId}`);
  if (!response.ok) {
    throw new Error('Failed to fetch stock details from the server');
  }
  return response.json();
};

const StockDetailsComponent = ({ stockId }: { stockId: string }) => {
  const { data: stockData, error: fetchError, isLoading: isStockLoading } = useQuery(['stockDetails', stockId], () => fetchStockById(stockId));

  if (isStockLoading) return <div>Loading stock information...</div>;
  if (fetchError instanceof Error) return <div>An error occurred while retrieving stock details: {fetchError.message}</div>;

  return (
    <div>
      <h1>Stock Details</h1>
    </div>
  );
};

export default StockDetailsComponent;