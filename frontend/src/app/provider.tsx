import {QueryClient, QueryClientProvider} from '@tanstack/react-query';
import * as React from 'react';

import {queryConfig} from '@/lib/react-query';

type AppProviderProps = {
  children: React.ReactNode;
};

export const AppProvider = ({children}: AppProviderProps) => {
  const [queryClient] = React.useState(
    () =>
      new QueryClient({
        defaultOptions: queryConfig,
      }),
  );
  return (
    <React.Suspense>
      <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    </React.Suspense>
  );
};
