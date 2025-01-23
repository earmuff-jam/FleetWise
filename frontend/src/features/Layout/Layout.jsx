import { Suspense, useEffect, useState } from 'react';

import { Outlet } from 'react-router-dom';
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

import { useTheme } from '@emotion/react';
import { darkTheme, lightTheme } from '@utils/Theme';
import AppToolbar from '@features/Layout/AppToolbar/AppToolbar';
import { profileActions } from '@features/Profile/profileSlice';
import MenuActionBar from '@features/Layout/MenuActionBar/MenuActionBar';
import Banner from '@features/Layout/Banner';
import { authActions } from '@features/LandingPage/authSlice';

const Layout = () => {
  const theme = useTheme();
  const dispatch = useDispatch();
  const isVerified = localStorage.getItem('isVerified');
  const smScreenSizeAndHigher = useMediaQuery(theme.breakpoints.up('sm'));
  const lgScreenSizeAndHigher = useMediaQuery(theme.breakpoints.up('lg'));

  const { profileDetails, loading } = useSelector((state) => state.profile);

  const [openDrawer, setOpenDrawer] = useState(lgScreenSizeAndHigher ? true : false);

  const handleDrawerOpen = () => setOpenDrawer(true);
  const handleDrawerClose = () => setOpenDrawer(false);

  useEffect(() => {
    dispatch(profileActions.getProfileDetails());
    dispatch(profileActions.getFavItems({ limit: 10 }));
  }, []);

  if (loading) {
    return <Skeleton height="100vh" />;
  }

  return (
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
            <Banner
              isVerified={isVerified}
              revalidateEmail={() => dispatch(authActions.revalidateEmail({ email: profileDetails?.email_address }))}
            />
            <Outlet />
          </Container>
        </Stack>
      </Suspense>
    </ThemeProvider>
  );
};

export default Layout;
