import {Box} from '@mui/material';

import {User} from '@/types/api';

import {useUser} from '../api/get-user';

export const UserData = ({
  gameName,
  tagLine,
}: {
  gameName: string;
  tagLine: string;
}) => {
  const userQuery = useUser({
    gameName,
    tagLine,
  });

  if (userQuery.isLoading) {
    return <Box>Error</Box>;
  }

  const user: User | undefined = userQuery?.data?.data;

  if (!user) return null;

  return (
    <>
      <Box>{'Name: ' + gameName}</Box>
      <Box>{'Rating: ' + user.rating}</Box>
      <Box>{'Play Count: ' + user.totalPlayCount}</Box>
    </>
  );
};
