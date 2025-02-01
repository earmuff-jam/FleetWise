import { Suspense, useEffect, useState } from 'react';

import { useSelector } from 'react-redux';

import { TourProvider } from '@reactour/tour';
import { RouterProvider } from 'react-router-dom';

import { Dialog } from '@mui/material';
import { router } from '@common/router';

import { PageContext } from '@src/PageContext';

import DEFAULT_TOUR_STEPS from '@utils/tour/steps';
import LandingPage from '@features/LandingPage/LandingPage';
import ResetPassword from '@features/LandingPage/Authentication/ResetPassword/ResetPassword';

const ApplicationValidator = () => {
  const { loading } = useSelector((state) => state.auth);

  const [page, setPage] = useState('reset');
  const [loggedInUser, setLoggedInUser] = useState(false);

  const navigateAuthComponents = () => {
    if (page === 'reset') {
      return <ResetPassword />;
    }
    return <LandingPage />;
  };

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
    <PageContext.Provider value={{ page, setPage }}>{navigateAuthComponents()}</PageContext.Provider>
  );
};

export default ApplicationValidator;
