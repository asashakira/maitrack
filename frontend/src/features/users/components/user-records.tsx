import {Box} from '@mui/material';

import {UserRecord} from '@/types/api';

import {useUserRecords} from '../api/get-user-records';

export const UserRecords = ({
  gameName,
  tagLine,
}: {
  gameName: string;
  tagLine: string;
}) => {
  const userRecordsQuery = useUserRecords({});

  if (userRecordsQuery.isLoading) {
    return <Box>Error</Box>;
  }

  const records: UserRecord[] | undefined = userRecordsQuery.data?.data;

  if (!records) return null;

  return (
    <>
      {records.map((record: UserRecord) => (
        <Box key={record.id} sx={{margin: 2}}>
          <Box>{record.title}</Box>
          <Box>{record.percentage}</Box>
        </Box>
      ))}
    </>
  );
};
