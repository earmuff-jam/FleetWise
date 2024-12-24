import { Suspense, useEffect, useState } from 'react';

import { Outlet, useLocation, useNavigate } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';

import {
  Box,
  CircularProgress,
  Container,
  CssBaseline,
  Skeleton,
  Stack,
  ThemeProvider,
  useMediaQuery,
} from '@mui/material';

import {
  ASSET_OVERVIEW_TOUR_STEPS,
  CATEGORIES_OVERVIEW_TOUR_STEPS,
  MAINTENANCE_PLAN_OVERVIEW_TOUR_STEPS,
  OVERVIEW_PAGE_TOUR_STEPS,
} from '@utils/tour/steps';
import { useTheme } from '@emotion/react';
import { TourProvider } from '@reactour/tour';

import { darkTheme, lightTheme } from '@utils/Theme';
import AppToolbar from '@features/Layout/AppToolbar/AppToolbar';
import { profileActions } from '@features/Profile/profileSlice';
import MenuActionBar from '@features/Layout/MenuActionBar/MenuActionBar';

const Layout = () => {
  const theme = useTheme();
  const location = useLocation();
  const dispatch = useDispatch();

  const smScreenSizeAndHigher = useMediaQuery(theme.breakpoints.up('sm'));
  const lgScreenSizeAndHigher = useMediaQuery(theme.breakpoints.up('lg'));

  const { profileDetails, loading } = useSelector((state) => state.profile);

  const [step, setStep] = useState(0);
  const [stepContent, setStepContent] = useState([
    ...OVERVIEW_PAGE_TOUR_STEPS,
    ...ASSET_OVERVIEW_TOUR_STEPS,
    ...CATEGORIES_OVERVIEW_TOUR_STEPS,
    ...MAINTENANCE_PLAN_OVERVIEW_TOUR_STEPS,
  ]);
  const [openDrawer, setOpenDrawer] = useState(lgScreenSizeAndHigher ? true : false);

  const handleDrawerOpen = () => setOpenDrawer(true);
  const handleDrawerClose = () => setOpenDrawer(false);

  const setCurrentStep = (step) => {
    setStep(step);
  };

  const buildTourSteps = (pathname = '') => {
    if (pathname != '') {
      const URL_PATH_BUILDER = {
        '/inventories/list': { steps: ASSET_OVERVIEW_TOUR_STEPS, startStepNumber: 6 },
        '/categories/list': { steps: CATEGORIES_OVERVIEW_TOUR_STEPS, startStepNumber: 13 },
        '/plan/list': { steps: MAINTENANCE_PLAN_OVERVIEW_TOUR_STEPS, startStepNumber: 20 },
      };
      setStepContent(URL_PATH_BUILDER[pathname]?.steps);
      setCurrentStep(URL_PATH_BUILDER[pathname]?.startStepNumber);
    }
  };

  useEffect(() => {
    dispatch(profileActions.getProfileDetails());
    dispatch(profileActions.getFavItems({ limit: 10 }));
    buildTourSteps(location.pathname || '');
  }, []);

  if (loading) {
    return <Skeleton height="100vh" />;
  }

  return (
    <TourProvider steps={stepContent} currentStep={step} setCurrentStep={setCurrentStep}>
      <ThemeProvider theme={profileDetails.appearance ? darkTheme : lightTheme}>
        <CssBaseline />
        <Suspense
          fallback={
            <Box display="flex" justifyContent="center" alignItems="center" minHeight="100vh">
              <CircularProgress color="inherit" />
            </Box>
          }
        >
          <AppToolbar profileDetails={profileDetails} handleDrawerOpen={handleDrawerOpen} />
          <Stack sx={{ marginTop: '5rem', marginBottom: '1rem' }}>
            <MenuActionBar
              openDrawer={openDrawer}
              handleDrawerClose={handleDrawerClose}
              smScreenSizeAndHigher={smScreenSizeAndHigher}
              lgScreenSizeAndHigher={lgScreenSizeAndHigher}
            />
            <Container maxWidth="md">
              <Outlet />
            </Container>
          </Stack>
        </Suspense>
      </ThemeProvider>
    </TourProvider>
  );
};

export default Layout;
