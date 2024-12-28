import { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';

import { Badge, Button, IconButton, Stack, Tooltip } from '@mui/material';

import {
  DarkModeRounded,
  LightModeOutlined,
  LogoutRounded,
  CircleNotifications,
  BlindRounded,
} from '@mui/icons-material';

import { authActions } from '@features/LandingPage/authSlice';
import { profileActions } from '@features/Profile/profileSlice';
import AppToolbarPopoverContent from '@features/Layout/AppToolbar/AppToolbarPopoverContent';
import { useTour } from '@reactour/tour';
import { useLocation, useParams } from 'react-router-dom';
import DEFAULT_TOUR_STEPS, { DEFAULT_STEP_MAPPER } from '@utils/tour/steps';

export default function AppToolbarActionButtons({ profileDetails }) {
  const { id } = useParams();
  const dispatch = useDispatch();
  const location = useLocation();

  const { setIsOpen, setCurrentStep, setSteps } = useTour();
  const { maintenanceNotifications = [], loading } = useSelector((state) => state.profile);

  const [anchorEl, setAnchorEl] = useState(null);

  const handleClick = (event) => setAnchorEl(event.currentTarget);
  const handleClose = () => setAnchorEl(null);

  const toggleReadOption = (id, selection) => {
    dispatch(profileActions.toggleMaintenanceNotificationReadOption({ maintenance_plan_id: id, is_read: !selection }));
  };

  const handleAppearance = () => {
    const draftData = { ...profileDetails, appearance: !profileDetails.appearance || false };
    dispatch(profileActions.updateProfileDetails({ draftData }));
  };

  const setTour = () => {
    const currentStep = id ? DEFAULT_STEP_MAPPER['/id'] : DEFAULT_STEP_MAPPER[location.pathname];
    const formattedSteps = DEFAULT_TOUR_STEPS.slice(currentStep.start, currentStep.end);
    setIsOpen(true);
    setCurrentStep(0);
    setSteps(formattedSteps);
  };

  const handleLogout = () => {
    dispatch(authActions.getLogout());
    localStorage.clear();
    window.location.href = '/';
  };

  useEffect(() => {
    dispatch(profileActions.getMaintenanceNotifications());
  }, []);

  return (
    <Stack direction="row" spacing="0.1rem">
      <IconButton size="small" onClick={handleClick} data-tour="overview-5">
        <Badge
          badgeContent={maintenanceNotifications?.filter((notification) => !notification?.is_read).length || 0}
          color="secondary"
        >
          <CircleNotifications />
        </Badge>
      </IconButton>
      <IconButton size="small" onClick={handleAppearance} data-tour="overview-6">
        {profileDetails?.appearance ? <LightModeOutlined fontSize="small" /> : <DarkModeRounded fontSize="small" />}
      </IconButton>
      <Tooltip title="log out">
        <IconButton size="small" onClick={handleLogout}>
          <LogoutRounded fontSize="small" />
        </IconButton>
      </Tooltip>
      <Button startIcon={<BlindRounded />} onClick={setTour}>
        Help with this page
      </Button>
      <AppToolbarPopoverContent
        loading={loading}
        anchorEl={anchorEl}
        handleClose={handleClose}
        toggleReadOption={toggleReadOption}
        options={maintenanceNotifications}
      />
    </Stack>
  );
}
