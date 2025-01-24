import { IconButton, Stack, Typography } from '@mui/material';

import { ChevronLeftRounded, ChevronRightRounded } from '@mui/icons-material';

export default function MenuActionBarTitle({ theme, handleDrawerClose }) {
  return (
    <Stack
      sx={{
        flexDirection: 'row',
        justifyContent: 'space-between',
        alignItems: 'center',
        padding: '0.75rem',
        backgroundColor: 'accentColor.default',
      }}
    >
      <Typography variant="h4">Fleetwise</Typography>
      <IconButton onClick={handleDrawerClose}>
        {theme.direction === 'rtl' ? <ChevronRightRounded /> : <ChevronLeftRounded />}
      </IconButton>
    </Stack>
  );
}
