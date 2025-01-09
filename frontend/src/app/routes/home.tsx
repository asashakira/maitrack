import {Box, Button, TextField} from '@mui/material';
import * as React from 'react';
import {useNavigate} from 'react-router-dom';

import {paths} from '@/config/paths';

const HomeRoute = () => {
  const navigate = useNavigate();
  const [userId, setUserId] = React.useState('');

  const handlePlayerIdChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setUserId(e.currentTarget.value);
  };

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    navigate(paths.users.getHref(userId));
  };

  return (
    <Box component="form" onSubmit={handleSubmit}>
      <TextField value={userId} onChange={handlePlayerIdChange} />
      <Button variant="outlined" type="submit">
        Submit
      </Button>
    </Box>
  );
};
export default HomeRoute;
