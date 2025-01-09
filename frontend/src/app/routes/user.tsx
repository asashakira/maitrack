import {Container} from '@mui/material';
import {useParams} from 'react-router-dom';

import {UserData} from '@/features/users/components/user-data';
import {UserRecords} from '@/features/users/components/user-records';

const UserRoute = () => {
  const params = useParams();
  const gameName = params.gameName as string;
  const tagLine = params.tagLine as string;

  return (
    <Container
      sx={{
        marginTop: 2,
      }}
    >
      <UserData gameName={gameName} tagLine={tagLine} />
      {/* <UserRecords gameName={gameName} tagLine={tagLine} /> */}
    </Container>
  );
};

export default UserRoute;
