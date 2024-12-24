import { AppBar, Stack, Toolbar } from '@mui/material';

import AppToolbarTitle from '@features/Layout/AppToolbar/AppToolbarTitle';
import AppToolbarActionButtons from '@features/Layout/AppToolbar/AppToolbarActionButtons';
import MaintenanceNotificationButton from '@features/Layout/AppToolbar/MaintenanceNotificationButton';
import AppToolbarButtonGroup from '@features/Layout/AppToolbar/AppToolbarButtonGroup';

export default function AppToolbar({ profileDetails, handleDrawerOpen }) {
  return (
    <AppBar elevation={0}>
      <Toolbar sx={{ backgroundColor: 'accentColor.default' }}>
        <AppToolbarTitle onClick={handleDrawerOpen} />
        <Stack direction="row" spacing={1}>
          <MaintenanceNotificationButton profileDetails={profileDetails} />
          <AppToolbarActionButtons profileDetails={profileDetails} />
          <AppToolbarButtonGroup profileDetails={profileDetails} />
        </Stack>
      </Toolbar>
    </AppBar>
  );
}
