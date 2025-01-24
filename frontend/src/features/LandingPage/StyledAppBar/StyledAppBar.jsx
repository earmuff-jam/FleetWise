import { Paper, Stack, Typography } from '@mui/material';
import { BarChartRounded } from '@mui/icons-material';

const StyledAppBar = () => {
  return (
    <Paper elevation={0}>
      <Stack direction="row" spacing={1} alignItems="center" padding="1.125rem">
        <BarChartRounded color="primary" fontSize="small" />
        <Typography variant="h4" sx={{ fontFamily: 'Nunito', color: 'primary.main' }}>
          Fleetwise
        </Typography>
      </Stack>
    </Paper>
  );
};

export default StyledAppBar;
