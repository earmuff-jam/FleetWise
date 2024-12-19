import { AppBar, Toolbar } from '@mui/material';
import AppToolbarTitle from '@features/Layout/AppToolbar/AppToolbarTitle';
import AppToolbarActionButtons from '@features/Layout/AppToolbar/AppToolbarActionButtons';

export default function AppToolbar({ profileDetails, handleDrawerOpen }) {
  return (
    <AppBar elevation={0}>
      <Toolbar sx={{ backgroundColor: 'accentColor.default' }}>
        <AppToolbarTitle onClick={handleDrawerOpen} />
        <AppToolbarActionButtons profileDetails={profileDetails} />
      </Toolbar>
    </AppBar>
  );
}
