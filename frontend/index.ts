import express from 'express';
import path from 'path';
import React from 'react';
import ReactDOMServer from 'react-dom/server';
import { StaticRouter as Router, Route, Routes } from 'react-router-dom';

const Dashboard = () => <div>Dashboard</div>;
const StockDetails = () => <div>Stock Details</div>;
const NotFound = () => <div>404 - Not Found</div>;

const renderApp = (url: string) => {
  return ReactDOMServer.renderToString(
    <Router location={url}>
      <Routes>
        <Route path="/" element={<Dashboard />} />
        <Route path="/stock/:id" element={<StockDetails />} />
        <Route path="*" element={<NotFound />} />
      </Routes>
    </Router>
  );
};

const app = express();

app.use(express.static(path.resolve(__dirname, 'public')));

const PORT = process.env.PORT || 3000;

app.get('*', (req, res) => {
  const appHtml = renderApp(req.url);
  const htmlToSend = `<!DOCTYPE html>
<html>
<head>
<title>Stock Market Dashboard</title>
</head>
<body>
  <div id="app">${appHtml}</div>
  <script src="/bundle.js"></script>
</body>
</html>`;
  res.send(htmlToSend);
});

app.listen(PORT, () => {
  console.log(`Server is listening on port ${PORT}`);
});