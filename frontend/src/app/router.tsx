import {useMemo} from 'react';
import {createBrowserRouter, RouterProvider} from 'react-router-dom';

import {paths} from '@/config/paths';

const convert = (m: any) => {
  return {
    Component: m.default,
  };
};

export const createAppRouter = () =>
  createBrowserRouter([
    {
      path: paths.home.path,
      lazy: () => import('./routes/home').then(convert),
    },
    {
      path: paths.users.path,
      lazy: () => import('./routes/user').then(convert),
    },
    {
      path: '*',
      lazy: () => import('./routes/not-found').then(convert),
    },
  ]);

export const AppRouter = () => {
  const router = useMemo(() => createAppRouter(), []);
  return <RouterProvider router={router} />;
};
