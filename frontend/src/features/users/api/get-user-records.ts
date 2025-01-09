import {useQuery, queryOptions} from '@tanstack/react-query';

import {api} from '@/lib/api-client';
import {QueryConfig} from '@/lib/react-query';
import {UserRecord} from '@/types/api';

export const getUserRecords = (): Promise<{data: UserRecord[]}> => {
  return api.get(`/records`);
};

export const getUserRecordsQueryOptions = () => {
  return queryOptions({
    queryKey: [],
    queryFn: () => getUserRecords(),
  });
};

type UseUserRecordsOptions = {
  queryConfig?: QueryConfig<typeof getUserRecordsQueryOptions>;
};

export const useUserRecords = ({queryConfig}: UseUserRecordsOptions) => {
  return useQuery({
    ...getUserRecordsQueryOptions(),
    ...queryConfig,
  });
};
