import { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';

import { Badge, IconButton } from '@mui/material';
import { CircleNotifications } from '@mui/icons-material';

import { profileActions } from '@features/Profile/profileSlice';
import MaintenanceNotificationPopoverContent from '@features/Layout/AppToolbar/MaintenanceNotificationPopoverContent';

export default function MaintenanceNotificationButton({ profileDetails }) {
  const { maintenanceNotifications = [], loading } = useSelector((state) => state.profile);

  const dispatch = useDispatch();
  const [anchorEl, setAnchorEl] = useState(null);

  const handleClick = (event) => setAnchorEl(event.currentTarget);
  const handleClose = () => setAnchorEl(null);

  const toggleReadOption = (id, selection) => {
    dispatch(profileActions.toggleMaintenanceNotificationReadOption({ maintenance_plan_id: id, is_read: !selection }));
  };

  useEffect(() => {
    dispatch(profileActions.getMaintenanceNotifications());
  }, []);

  return (
    <>
      <IconButton size="small" onClick={handleClick}>
        <Badge
          badgeContent={maintenanceNotifications?.filter((notification) => !notification?.is_read).length || 0}
          color="secondary"
        >
          <CircleNotifications />
        </Badge>
      </IconButton>
      <MaintenanceNotificationPopoverContent
        loading={loading}
        anchorEl={anchorEl}
        handleClose={handleClose}
        toggleReadOption={toggleReadOption}
        options={maintenanceNotifications}
      />
    </>
  );
}
