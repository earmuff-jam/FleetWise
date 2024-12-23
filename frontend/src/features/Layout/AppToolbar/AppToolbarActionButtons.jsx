import { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';

import { Badge, IconButton, Stack, Tooltip } from '@mui/material';

import { DarkModeRounded, LightModeOutlined, LogoutRounded, CircleNotifications } from '@mui/icons-material';

import { authActions } from '@features/LandingPage/authSlice';
import { profileActions } from '@features/Profile/profileSlice';
import AppToolbarPopoverContent from '@features/Layout/AppToolbar/AppToolbarPopoverContent';

export default function AppToolbarActionButtons({ profileDetails }) {
  const { maintenanceNotifications = [], loading } = useSelector((state) => state.profile);

  const dispatch = useDispatch();
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
      <IconButton size="small" onClick={handleClick}>
        <Badge
          badgeContent={maintenanceNotifications?.filter((notification) => !notification?.is_read).length || 0}
          color="secondary"
        >
          <CircleNotifications />
        </Badge>
      </IconButton>
      <IconButton size="small" onClick={handleAppearance}>
        {profileDetails?.appearance ? <LightModeOutlined fontSize="small" /> : <DarkModeRounded fontSize="small" />}
      </IconButton>
      <Tooltip title="log out">
        <IconButton size="small" onClick={handleLogout}>
          <LogoutRounded fontSize="small" />
        </IconButton>
      </Tooltip>
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
