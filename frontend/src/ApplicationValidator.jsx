import steps from './tour/steps';
import { useSelector } from 'react-redux';
import { Suspense, useEffect, useState } from 'react';
import { SnackbarProvider } from 'notistack';
import { TourProvider } from '@reactour/tour';
import { RouterProvider } from 'react-router-dom';
import { Dialog } from '@mui/material';
import { router } from './common/router';
import LandingPage from './features/LandingPage/LandingPage';

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
    <TourProvider steps={steps}>
      <SnackbarProvider
        dense
        preventDuplicate
        maxSnack={3}
        anchorOrigin={{ horizontal: 'right', vertical: 'bottom' }}
        autoHideDuration={3000}
      >
        <Suspense fallback={<Dialog open={loading} title="Loading..." />}>
          <RouterProvider router={router} />
        </Suspense>
      </SnackbarProvider>
    </TourProvider>
  ) : (
    <LandingPage />
  );
};

export default ApplicationValidator;
