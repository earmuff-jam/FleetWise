import { Stack, Typography } from '@mui/material';

export default function HeroContentCenterText() {
  return (
    <Stack textAlign="center">
      <Typography variant="h1" fontWeight="900">
        Fleetwise
      </Typography>
      <Typography variant="h4" fontWeight="100">
        Build and track your assets with ease
      </Typography>
    </Stack>
  );
}
