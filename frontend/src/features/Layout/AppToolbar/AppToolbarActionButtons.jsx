import { useDispatch } from 'react-redux';
import { IconButton, Tooltip } from '@mui/material';
import { LogoutRounded } from '@mui/icons-material';
import { authActions } from '@features/LandingPage/authSlice';

export default function AppToolbarActionButtons() {
  const dispatch = useDispatch();

  const handleLogout = () => {
    dispatch(authActions.getLogout());
    localStorage.clear();
    window.location.href = '/';
  };

  return (
    <Tooltip title="log out">
      <IconButton size="small" onClick={handleLogout}>
        <LogoutRounded fontSize="small" />
      </IconButton>
    </Tooltip>
  );
}
