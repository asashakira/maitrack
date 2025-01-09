import {useQuery, queryOptions} from '@tanstack/react-query';

import {api} from '@/lib/api-client';
import {QueryConfig} from '@/lib/react-query';
import {User} from '@/types/api';

export const getUser = ({
  gameName,
  tagLine,
}: {
  gameName: string;
  tagLine: string;
}): Promise<{data: User}> => {
  return api.get(`/users/by-mai-id/${gameName}/${tagLine}/data`);
};

export const getUserQueryOptions = (gameName: string, tagLine: string) => {
  return queryOptions({
    queryKey: ['user', gameName],
    queryFn: () => getUser({gameName, tagLine}),
  });
};

type UseUserOptions = {
  gameName: string;
  tagLine: string;
  queryConfig?: QueryConfig<typeof getUserQueryOptions>;
};

export const useUser = ({gameName, tagLine, queryConfig}: UseUserOptions) => {
  return useQuery({
    ...getUserQueryOptions(gameName, tagLine),
    ...queryConfig,
  });
};
