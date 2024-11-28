import { useQuery } from 'react-query';

const fetchStockDetails = async (id: string) => {
  const response = await fetch(`/api/stocks/${id}`);
  if (!response.ok) {
    throw new Error('Network response was not ok');
  }
  return response.json();
};

const StockDetails = ({ id }: { id: string }) => {
  const { data, error, isLoading } = useQuery(['stockDetails', id], () => fetchStockDetails(id));

  if (isLoading) return <div>Loading...</div>;
  if (error instanceof Error) return <div>An error occurred: {error.message}</div>;

  return (
    <div>
      <h1>Stock Details</h1>
      {/* Render your data here */}
    </div>
  );
};