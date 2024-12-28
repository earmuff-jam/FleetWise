import { Suspense, useEffect, useState } from 'react';

import { useSelector } from 'react-redux';

import { TourProvider } from '@reactour/tour';
import { RouterProvider } from 'react-router-dom';

import { Dialog } from '@mui/material';
import { router } from '@common/router';

import DEFAULT_TOUR_STEPS from '@utils/tour/steps';
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
    <TourProvider steps={DEFAULT_TOUR_STEPS}>
      <Suspense fallback={<Dialog open={loading} title="Loading..." />}>
        <RouterProvider router={router} />
      </Suspense>
    </TourProvider>
  ) : (
    <LandingPage />
  );
};

export default ApplicationValidator;
