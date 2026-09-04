import { createConnectTransport } from "@connectrpc/connect-web";
import { TransportProvider } from "@connectrpc/connect-query";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import './App.css';
import EntitiesList from './modules/Entities/EntitiesList/EntitiesList.tsx';

const transport = createConnectTransport({
  baseUrl: "http://localhost:9002"
})

const queryClient = new QueryClient();

const App = () => {
  return (
    <TransportProvider transport={transport}>
      <QueryClientProvider client={queryClient}>
        <div className="content">
          <h1>Rsbuild with React</h1>
          <p>Start building amazing things with Rsbuild.</p>
          <EntitiesList />
        </div>
      </QueryClientProvider>
    </TransportProvider>
  );
};

export default App;
