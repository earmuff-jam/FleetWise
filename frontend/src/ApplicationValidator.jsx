import { Suspense, useEffect, useState } from 'react';

import { useSelector } from 'react-redux';
import { RouterProvider } from 'react-router-dom';

import { Dialog } from '@mui/material';
import { router } from '@common/router';
import LandingPage from '@features/LandingPage/LandingPage';

const ApplicationValidator = () => {
  const { loading } = useSelector((state) => state.auth);

  const [loggedInUser, setLoggedInUser] = useState(false);

  useEffect(() => {
    const userID = localStorage.getItem('userID');
    if (!userID) {
      setLoggedInUser(false);
      return;
    } else {
      setLoggedInUser(true);
    }
  }, [loading]);

  return loggedInUser ? (
    <Suspense fallback={<Dialog open={loading} title="Loading..." />}>
      <RouterProvider router={router} />
    </Suspense>
  ) : (
    <LandingPage />
  );
};

export default ApplicationValidator;
